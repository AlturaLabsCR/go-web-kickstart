package s3

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
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

	mu sync.RWMutex
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
	fs.mu.Lock()
	defer fs.mu.Unlock()

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

	parts := strings.Split(params.Key, "/")
	filename := parts[len(parts)-1]

	dir := strings.TrimSuffix(params.Key, filename)

	parts = strings.Split(params.Key, ".")
	ext := parts[len(parts)-1]

	if len(ext) < 31 && ext != "" {
		ext = "." + ext
	} else {
		ext = ""
	}

	md5sum := md5.Sum(data)
	md5 := hex.EncodeToString(md5sum[:])

	newKey := dir + md5 + ext

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
			return newKey, fs.publicEndpoint + newKey, err
		}

		oldSize = object.Size

		object.Mime = mime
		object.MD5 = md5
		object.Size = size
		object.Modified = time.Now()
	}

	newSize := fs.bucketSize - oldSize + size
	if newSize > fs.maxBucketSize {
		return "", "", ErrBucketTooLarge
	}

	path, err := fs.objectPath(newKey)
	if err != nil {
		return "", "", err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", "", err
	}
	fs.bucketSize = newSize

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

func (fs *FileSystem) Delete(ctx context.Context, key string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	object, err := fs.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, kv.ErrNotFound) {
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

	fs.bucketSize -= object.Size

	if err := fs.store.Delete(ctx, key); err != nil {
		return err
	}

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
	return filepath.Abs(filepath.Join(fs.root, key))
}
