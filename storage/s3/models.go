package s3

import (
	"time"

	"app/storage/kv"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Object struct {
	Bucket    string
	Key       string
	PublicURL string
	Mime      string
	Size      int64
	Created   time.Time
	Modified  time.Time
}

type StorageParams struct {
	Client         *s3.Client
	BucketName     string
	PublicEndpoint string
	Store          kv.Store[Object]
	Cache          kv.Store[Object]
	MaxObjectSize  int64
	MaxBucketSize  int64
	MaxOps         int64
}
