// Package s3 has abstraction methods to upload objects to an s3-like store
package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const (
	ErrObjectTooLarge = errStr("object size exceeds maximum allowed by this interface")
	ErrBucketTooLarge = errStr("bucket size exceeds maximum allowed by this interface")
)

func New(params BucketParams) (*Bucket, error) {
	return &Bucket{
		client:         params.Client,
		bucketName:     params.BucketName,
		store:          params.Store,
		cache:          params.Cache,
		maxObjectSize:  params.MaxObjectSize,
		maxBucketSize:  params.MaxBucketSize,
		publicEndpoint: params.PublicEndpoint,

		sem: make(chan struct{}, params.MaxOps),

		bucketSize: 0,
	}, nil
}

func (b *Bucket) PutObject(ctx context.Context, key string, body io.Reader) (publicURL string, err error) {
	if err := b.wait(ctx); err != nil {
		return "", err
	}
	defer b.done()

	objLock := b.getObjectLock(key)
	defer objLock.Unlock()

	r, err := b.newStorageReader(body)
	if err != nil {
		return "", err
	}
	defer r.Close()

	_, err = b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(b.bucketName),
		Key:         aws.String(key),
		Body:        r,
		ContentType: aws.String(r.ct),
	})
	if err != nil {
		return "", err
	}

	return b.publicEndpoint + key, nil
}

func (b *Bucket) LoadCache(ctx context.Context) error {
	buckets, err := b.store.GetElems(ctx)
	if err != nil {
		return err
	}
	var size int64 = 0
	for key, object := range buckets {
		if err := b.cache.Set(ctx, key, object); err != nil {
			return err
		}
		size += object.Size
	}
	b.bucketSize = size
	return nil
}

func (b *Bucket) ReserveBytes(n int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.maxBucketSize > 0 && b.bucketSize+b.inflight+n > b.maxBucketSize {
		return ErrBucketTooLarge
	}

	b.inflight += n
	return nil
}

func (b *Bucket) ReleaseBytes(n int64) {
	b.mu.Lock()
	b.inflight -= n
	if b.inflight < 0 {
		b.inflight = 0
	}
	b.mu.Unlock()
}
