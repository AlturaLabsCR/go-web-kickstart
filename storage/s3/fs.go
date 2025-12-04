package s3

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"app/storage/kv"
	"app/utils"
)

type FileSystem struct {
	root       string
	bucketSize int64

	store kv.Store[Object]
	cache *kv.MemoryStore[Object]

	maxObjectSize int64
	maxBucketSize int64

	publicEndpoint string

	mu       sync.RWMutex
	keyLocks sync.Map
}

func NewFS(params S3Params) *FileSystem {
	cache := kv.NewMemoryStore[Object]()
	return &FileSystem{
		root:           params.Bucket,
		bucketSize:     0,
		store:          params.Store,
		cache:          cache,
		maxObjectSize:  params.MaxObjectSize,
		maxBucketSize:  params.MaxBucketSize,
		publicEndpoint: params.PublicEndpoint,
	}
}

func (fs *FileSystem) Put(ctx context.Context, params PutObjectParams) (key, publicURL string, err error) {
	mime, size, data, err := utils.InspectReader(
		params.Body,
		fs.maxObjectSize,
	)
	if err != nil {
		return "", "", err
	}

	if size > fs.maxObjectSize {
		return "", "", ErrTooLarge
	}

	parts := strings.Split(params.Path, "/")
	filename := parts[len(parts)-1]

	dir := strings.TrimSuffix(params.Path, filename)

	parts = strings.Split(params.Path, ".")
	ext := parts[len(parts)-1]

	if len(ext) < 31 && ext != "" {
		ext = "." + ext
	} else {
		ext = ""
	}

	md5sum := md5.Sum(data)
	md5 := hex.EncodeToString(md5sum[:])

	newKey := dir + md5 + ext

	keyLock := fs.getKeyLock(newKey)
	keyLock.Lock()
	defer keyLock.Unlock()

	var oldSize int64 = 0

	object, err := fs.cache.Get(ctx, newKey)
	if err != nil {
		if !errors.Is(err, kv.ErrNotFound) {
			return "", "", err
		}

		now := time.Now()
		object = Object{
			Bucket:   fs.root,
			Key:      newKey,
			Mime:     mime,
			MD5:      md5,
			Size:     size,
			Created:  now,
			Modified: now,
		}
	} else {
		if object.MD5 == md5 {
			return object.Key, fs.publicEndpoint + object.Key, nil
		}

		oldSize = object.Size

		object.Mime = mime
		object.MD5 = md5
		object.Size = size
		object.Modified = time.Now()
	}

	fs.mu.Lock()
	newSize := fs.bucketSize - oldSize + size
	if newSize > fs.maxBucketSize {
		fs.mu.Unlock()
		return "", "", ErrBucketTooLarge
	}
	fs.bucketSize = fs.bucketSize - oldSize + size
	fs.mu.Unlock()

	path, err := fs.objectPath(newKey)
	if err != nil {
		return "", "", err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", "", err
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", "", err
	}

	if err := fs.store.Set(ctx, newKey, object); err != nil {
		return "", "", err
	}

	if err := fs.cache.Set(ctx, newKey, object); err != nil {
		return "", "", err
	}

	return newKey, fs.publicEndpoint + newKey, err
}

func (fs *FileSystem) Get(ctx context.Context, key string) ([]byte, error) {
	keyLock := fs.getKeyLock(key)
	keyLock.Lock()
	defer keyLock.Unlock()

	path, err := fs.objectPath(key)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file %q not found: %w", key, err)
		}
		return nil, err
	}

	return data, nil
}

func (fs *FileSystem) Delete(ctx context.Context, key string) error {
	keyLock := fs.getKeyLock(key)
	keyLock.Lock()
	defer keyLock.Unlock()

	object, err := fs.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			fs.keyLocks.Delete(key)
			return nil
		}
		return err
	}

	path, err := fs.objectPath(key)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	fs.mu.Lock()
	fs.bucketSize -= object.Size
	fs.mu.Unlock()

	if err := fs.store.Delete(ctx, key); err != nil {
		return err
	}

	fs.keyLocks.Delete(key)

	return fs.cache.Delete(ctx, key)
}

func (fs *FileSystem) LoadCache(ctx context.Context) error {
	objects, err := fs.store.GetElems(ctx)
	if err != nil {
		return err
	}
	var size int64 = 0
	for key, object := range objects {
		if err := fs.cache.Set(ctx, key, object); err != nil {
			return err
		}
		size += object.Size
	}
	fs.bucketSize = size

	return nil
}

func (fs *FileSystem) objectPath(key string) (string, error) {
	p := filepath.Join(fs.root, key)
	p = filepath.Clean(p)

	root := filepath.Clean(fs.root)
	if !strings.HasPrefix(p, root+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid object key: path escapes bucket root")
	}
	return p, nil
}

func (fs *FileSystem) getKeyLock(key string) *sync.Mutex {
	lockIface, _ := fs.keyLocks.LoadOrStore(key, &sync.Mutex{})
	return lockIface.(*sync.Mutex)
}
