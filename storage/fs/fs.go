// Package fs implements 'ObjectStorage' interface with a filesystem backend
package fs

import (
	"context"
	"errors"
	"io"
	"time"

	"app/cache"
	"app/storage"
	"app/storage/reader"
)

func NewFS(params *FSParams) (*FS, error) {
	return &FS{
		root:          params.Root,
		cache:         params.Cache,
		l2cache:       params.L2Cache,
		maxObjectSize: params.MaxObjectSize,
		maxBucketSize: params.MaxBucketSize,
		bucketSize:    0,
	}, nil
}

func (fs *FS) PutObject(ctx context.Context, key string, body io.Reader) (*storage.Object, error) {
	objLock := fs.getObjectLock(key)
	objLock.Lock()
	defer objLock.Unlock()

	now := time.Now().Unix()
	created := now
	var oldSize int64 = 0

	if objStr, err := fs.cache.Get(ctx, key); err == nil {
		if o, err := storage.StringToObject(objStr); err == nil {
			created = o.Created
			oldSize = o.Size
		} else {
			return nil, err
		}
	} else {
		if !errors.Is(err, cache.ErrNotFound) {
			return nil, err
		}
	}

	rem := fs.remaining()
	if rem <= 0 {
		return nil, storage.ErrBucketTooLarge
	}

	r, ct, err := reader.NewObjectReader(ctx, body, rem)
	if err != nil {
		if err == reader.ErrOverflow {
			return nil, storage.ErrObjectTooLarge
		}
		return nil, err
	}

	newSize := r.Size()

	obj := &storage.Object{
		Bucket:   fs.root,
		Key:      key,
		Mime:     ct,
		Size:     newSize,
		Modified: now,
		Created:  created,
	}

	objStr, err := storage.ObjectToString(obj)
	if err != nil {
		return nil, err
	}

	fs.reserve(newSize)
	if err := fs.putObject(ctx, key, r); err != nil {
		fs.release(newSize)
		return nil, err
	}
	fs.commit(newSize, oldSize)

	if err := fs.cache.Set(ctx, key, objStr); err != nil {
		return nil, err
	}

	if fs.l2cache != nil {
		if err := fs.l2cache.Set(ctx, key, objStr); err != nil {
			return nil, err
		}
	}

	return obj, nil
}

func (fs *FS) GetObject(ctx context.Context, key string) (*storage.Object, io.ReadCloser, error) {
	objLock := fs.getObjectLock(key)
	objLock.Lock()
	defer objLock.Unlock()

	objStr, err := fs.cache.Get(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	obj, err := storage.StringToObject(objStr)
	if err != nil {
		return nil, nil, err
	}

	body, err := fs.getObject(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	return obj, body, nil
}

func (fs *FS) DeleteObject(ctx context.Context, key string) error {
	objLock := fs.getObjectLock(key)
	objLock.Lock()
	defer objLock.Unlock()

	objStr, err := fs.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return nil
		}
		return err
	}

	obj, err := storage.StringToObject(objStr)
	if err != nil {
		return err
	}

	if err := fs.deleteObject(ctx, key); err == nil {
		fs.commit(0, obj.Size)
	} else {
		return err
	}

	if fs.l2cache != nil {
		_ = fs.l2cache.Del(ctx, key)
	}

	return fs.cache.Del(ctx, key)
}

func (fs *FS) LoadCache(ctx context.Context) error {
	if fs.l2cache == nil {
		return storage.ErrBadCache
	}

	values, err := fs.l2cache.GetAll(ctx)
	if err != nil {
		return err
	}

	var total int64 = 0

	for key, val := range values {
		obj, err := storage.StringToObject(val)
		if err != nil {
			return err
		}

		if err := fs.cache.Set(ctx, key, val); err != nil {
			return err
		}

		total += obj.Size
	}

	fs.mu.Lock()
	fs.bucketSize = total
	fs.mu.Unlock()

	return nil
}
