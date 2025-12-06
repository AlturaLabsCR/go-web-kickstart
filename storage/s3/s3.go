// Package s3 has abstraction methods to upload objects to an s3-like store
package s3

import (
	"context"
	"errors"
	"io"
	"time"

	"app/storage/kv"
	"app/storage/s3/reader"
)

func (b *Bucket) PutObject(ctx context.Context, key string, body io.Reader) (object *Object, err error) {
	objLock := b.getObjectLock(key)
	objLock.Lock()
	defer objLock.Unlock()

	now := time.Now()
	created := now
	var oldSize int64 = 0

	if obj, err := b.cache.Get(ctx, key); err == nil {
		created = obj.Created
		oldSize = obj.Size
	} else {
		if !errors.Is(err, kv.ErrNotFound) {
			return nil, err
		}
	}

	r, err := reader.NewObjectReader(ctx, body, b.makeSizeCallback())
	if err != nil {
		return nil, err
	}

	if err := b.putObject(ctx, key, r); err != nil {
		return nil, err
	}

	if r.Size() > b.maxObjectSize {
		b.deleteObject(ctx, key)
		return nil, ErrObjectTooLarge
	}

	newSize := r.Size()

	b.mu.Lock()

	b.bucketSize -= oldSize
	b.bucketSize += newSize
	b.inflight -= newSize

	b.mu.Unlock()

	object = &Object{
		Bucket:    b.bucketName,
		Key:       key,
		PublicURL: b.publicEndpoint + key,
		Mime:      r.ContentType(),
		Size:      newSize,
		Modified:  now,
		Created:   created,
	}

	if err := b.store.Set(ctx, key, object); err != nil {
		return nil, err
	}

	if err := b.cache.Set(ctx, key, object); err != nil {
		return nil, err
	}

	return object, nil
}

func (b *Bucket) GetObject(ctx context.Context, key string) (*Object, io.ReadCloser, error) {
	objLock := b.getObjectLock(key)
	objLock.Lock()
	defer objLock.Unlock()

	obj, err := b.cache.Get(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	body, err := b.getObject(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	return obj, body, nil
}

func (b *Bucket) DeleteObject(ctx context.Context, key string) error {
	objLock := b.getObjectLock(key)
	objLock.Lock()
	defer objLock.Unlock()

	obj, err := b.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			return nil
		} else {
			return err
		}
	}

	if err := b.deleteObject(ctx, key); err == nil {
		b.mu.Lock()
		b.bucketSize -= obj.Size
		b.mu.Unlock()
	} else {
		return err
	}

	if err := b.store.Delete(ctx, key); err != nil {
		return err
	}

	return b.cache.Delete(ctx, key)
}

func (b *Bucket) LoadCache(ctx context.Context) error {
	buckets, err := b.store.GetElems(ctx)
	if err != nil {
		return err
	}

	var size int64 = 0
	for key, object := range buckets {
		if err := b.cache.Set(ctx, key, &object); err != nil {
			return err
		}
		size += object.Size
	}

	b.mu.Lock()
	b.bucketSize = size
	b.mu.Unlock()

	return nil
}

func NewS3(params *StorageParams) (*Bucket, error) {
	return &Bucket{
		client:         params.Client,
		bucketName:     params.BucketName,
		store:          params.Store,
		cache:          params.Cache,
		maxObjectSize:  params.MaxObjectSize,
		maxBucketSize:  params.MaxBucketSize,
		publicEndpoint: params.PublicEndpoint,

		bucketSize: 0,
	}, nil
}
