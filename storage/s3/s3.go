// Package s3 has abstraction methods to upload objects to an s3-like store
package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"app/storage/kv"
	"app/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const (
	ErrTooLarge       = errStr("object size exceeds maximum allowed by this interface")
	ErrBucketTooLarge = errStr("bucket size exceeds maximum allowed by this interface")
)

type Storage interface {
	// require read/write from/to a database
	Put(ctx context.Context, parms PutObjectParams) error
	Delete(ctx context.Context, key string) error

	// avoids excessive API queries
	LoadCache(ctx context.Context) error
}

type Object struct {
	ObjectBucket   string
	ObjectKey      string
	ObjectMime     string
	ObjectSize     int64
	ObjectCreated  time.Time
	ObjectModified time.Time
}

type S3 struct {
	client     *s3.Client
	bucket     string
	bucketSize int64

	store kv.Store[Object]
	cache *kv.MemoryStore[Object]

	maxObjectSize int64
	maxBucketSize int64

	mu sync.RWMutex
}

type PutObjectParams struct {
	Key  string
	Body io.Reader
}

type S3Params struct {
	Bucket        string
	Store         kv.Store[Object]
	MaxObjectSize int64
	MaxBucketSize int64
}

func New(params S3Params) (*S3, error) {
	s3c, err := s3config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	s3client := s3.NewFromConfig(s3c)

	cache := kv.NewMemoryStore[Object]()

	return &S3{
		client:        s3client,
		bucket:        params.Bucket,
		bucketSize:    0,
		store:         params.Store,
		cache:         cache,
		maxObjectSize: params.MaxObjectSize,
		maxBucketSize: params.MaxBucketSize,
	}, nil
}

func (s *S3) Put(ctx context.Context, params PutObjectParams) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	mime, size, data, err := utils.InspectReader(
		params.Body,
		s.maxObjectSize,
	)
	if err != nil {
		return err
	}

	if size > s.maxObjectSize {
		return ErrTooLarge
	}

	var oldSize int64 = 0

	object, err := s.cache.Get(ctx, params.Key)
	if err != nil {
		if !errors.Is(err, kv.ErrNotFound) {
			return err
		}

		now := time.Now()
		object = Object{
			ObjectBucket:   s.bucket,
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

	newSize := s.bucketSize - oldSize + size
	if newSize > s.maxBucketSize {
		return ErrBucketTooLarge
	}

	if _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(params.Key),
		Body:   bytes.NewReader(data),
	}); err != nil {
		return err
	}
	s.bucketSize = newSize

	if err := s.store.Set(ctx, params.Key, object); err != nil {
		return err
	}

	return s.cache.Set(ctx, params.Key, object)
}

func (s *S3) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	object, err := s.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			return nil
		}
		return err
	}

	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}); err != nil {
		return err
	}

	s.bucketSize -= object.ObjectSize

	if err := s.store.Delete(ctx, key); err != nil {
		return err
	}

	return s.cache.Delete(ctx, key)
}

func (s *S3) LoadCache(ctx context.Context) error {
	buckets, err := s.store.GetElems(ctx)
	if err != nil {
		return err
	}
	var size int64 = 0
	for key, object := range buckets {
		if err := s.cache.Set(ctx, key, object); err != nil {
			return err
		}
		size += object.ObjectSize
	}
	s.bucketSize = size
	return nil
}
