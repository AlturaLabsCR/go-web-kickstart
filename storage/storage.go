// Package storage abstracts common storage solutions such as key-value or S3
package storage

import (
	"context"
	"io"

	"app/storage/s3"
)

type ObjectStorage interface {
	PutObject(ctx context.Context, key string, body io.Reader) (object *s3.Object, err error)
	GetObject(ctx context.Context, key string) (object *s3.Object, body io.ReadCloser, err error)
	DeleteObject(ctx context.Context, key string) error
	LoadCache(ctx context.Context) error
}
