package s3

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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

	mu sync.RWMutex
}

func NewFS(params S3Params) *FileSystem {
	cache := kv.NewMemoryStore[Object]()
	return &FileSystem{
		root:          params.Bucket,
		bucketSize:    0,
		store:         params.Store,
		cache:         cache,
		maxObjectSize: params.MaxObjectSize,
		maxBucketSize: params.MaxBucketSize,
	}
}

func (fs *FileSystem) Put(ctx context.Context, params PutObjectParams) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	mime, size, data, err := utils.InspectReader(
		params.Body,
		fs.maxObjectSize,
	)
	if err != nil {
		return err
	}

	if size > fs.maxObjectSize {
		return ErrTooLarge
	}

	var oldSize int64 = 0

	object, err := fs.cache.Get(ctx, params.Key)
	if err != nil {
		if !errors.Is(err, kv.ErrNotFound) {
			return err
		}

		now := time.Now()
		object = Object{
			ObjectBucket:   fs.root,
			ObjectKey:      params.Key,
			ObjectMime:     mime,
			ObjectSize:     size,
			ObjectCreated:  now,
			ObjectModified: now,
		}
	} else {
		oldSize = object.ObjectSize

		object.ObjectMime = mime
		object.ObjectSize = size
		object.ObjectModified = time.Now()
	}

	newSize := fs.bucketSize - oldSize + size
	if newSize > fs.maxBucketSize {
		return ErrBucketTooLarge
	}

	path := fs.objectPath(params.Key)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	fs.bucketSize = newSize

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}

	if err := fs.store.Set(ctx, params.Key, object); err != nil {
		return err
	}

	return fs.cache.Set(ctx, params.Key, object)
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

	path := fs.objectPath(key)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	fs.bucketSize -= object.ObjectSize

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
		size += object.ObjectSize
	}
	fs.bucketSize = size
	return nil
}

func (fs *FileSystem) objectPath(key string) string {
	return filepath.Join(fs.root, key)
}
