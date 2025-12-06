package s3

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"app/storage/kv"
	"app/storage/s3/reader"
)

func (fs *FileSystem) PutObject(ctx context.Context, key string, body io.Reader) (object *Object, err error) {
	objLock := fs.getObjectLock(key)
	objLock.Lock()
	defer objLock.Unlock()

	now := time.Now()
	created := now
	var oldSize int64 = 0

	if obj, err := fs.cache.Get(ctx, key); err == nil {
		created = obj.Created
		oldSize = obj.Size
	} else {
		if !errors.Is(err, kv.ErrNotFound) {
			return nil, err
		}
	}

	r, err := reader.NewObjectReader(ctx, body, fs.makeSizeCallback())
	if err != nil {
		return nil, err
	}

	if err := fs.putObject(key, r); err != nil {
		return nil, err
	}

	if r.Size() > fs.maxObjectSize {
		fs.deleteObject(ctx, key)
		return nil, ErrObjectTooLarge
	}

	newSize := r.Size()

	fs.mu.Lock()

	fs.bucketSize -= oldSize
	fs.bucketSize += newSize
	fs.inflight -= newSize

	fs.mu.Unlock()

	object = &Object{
		Bucket:    fs.root,
		Key:       key,
		PublicURL: fs.publicEndpoint + key,
		Mime:      r.ContentType(),
		Size:      newSize,
		Modified:  now,
		Created:   created,
	}

	if err := fs.store.Set(ctx, key, object); err != nil {
		return nil, err
	}

	if err := fs.cache.Set(ctx, key, object); err != nil {
		return nil, err
	}

	return object, nil
}

func (fs *FileSystem) GetObject(ctx context.Context, key string) (*Object, io.ReadCloser, error) {
	objLock := fs.getObjectLock(key)
	objLock.Lock()
	defer objLock.Unlock()

	obj, err := fs.cache.Get(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	body, err := fs.getObject(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	return obj, body, nil
}

func (fs *FileSystem) DeleteObject(ctx context.Context, key string) error {
	keyLock := fs.getObjectLock(key)
	keyLock.Lock()
	defer keyLock.Unlock()

	obj, err := fs.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			return nil
		}
		return err
	}

	if err := fs.deleteObject(ctx, key); err == nil {
		fs.mu.Lock()
		fs.bucketSize -= obj.Size
		fs.mu.Unlock()
	} else {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

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
		if err := fs.cache.Set(ctx, key, &object); err != nil {
			return err
		}
		size += object.Size
	}
	fs.bucketSize = size

	return nil
}

func NewFS(params *StorageParams) (*FileSystem, error) {
	return &FileSystem{
		root:           params.BucketName,
		store:          params.Store,
		cache:          params.Cache,
		maxObjectSize:  params.MaxObjectSize,
		maxBucketSize:  params.MaxBucketSize,
		publicEndpoint: params.PublicEndpoint,

		bucketSize: 0,
	}, nil
}
