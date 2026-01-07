package s3

import (
	"sync"

	"app/cache"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Bucket struct {
	client        *s3.Client
	bucketName    string
	cache         cache.Cache
	l2cache       cache.Cache
	maxObjectSize int64
	maxBucketSize int64

	objectLocks sync.Map

	mu         sync.RWMutex
	bucketSize int64
	inflight   int64
}

type S3Params struct {
	Client        *s3.Client
	BucketName    string
	Cache         cache.Cache
	L2Cache       cache.Cache // optional
	MaxObjectSize int64
	MaxBucketSize int64
}
