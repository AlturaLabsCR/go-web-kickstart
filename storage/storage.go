// Package storage abstracts common storage solutions such as key-value or S3
package storage

import (
	"context"
	"encoding/json"
	"io"
)

const (
	ErrObjectTooLarge = errStr("object too large")
	ErrBucketTooLarge = errStr("reached bucket size limit")
	ErrBadCache       = errStr("nil Cache interface")
)

type ObjectStorage interface {
	PutObject(ctx context.Context, key string, body io.Reader) (object *Object, err error)
	GetObject(ctx context.Context, key string) (object *Object, body io.ReadCloser, err error)
	DeleteObject(ctx context.Context, key string) error
	LoadCache(ctx context.Context) error
}

type Object struct {
	Bucket   string
	Key      string
	Mime     string
	Size     int64
	Created  int64
	Modified int64
}

func ObjectToString(obj *Object) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func StringToObject(objStr string) (*Object, error) {
	var o Object
	if err := json.Unmarshal([]byte(objStr), &o); err != nil {
		return nil, err
	}
	return &o, nil
}
