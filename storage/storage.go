// Package storage abstracts common storage solutions such as key-value or S3
package storage

import (
	"context"
	"io"
)

type ObjectStorage interface {
	PutObject(ctx context.Context, key string, body io.Reader) (publicURL string, err error)
	GetObject(ctx context.Context, key string) (body io.ReadCloser, err error)
	DeleteObject(ctx context.Context, key string) error
	LoadCache(ctx context.Context) error
}
