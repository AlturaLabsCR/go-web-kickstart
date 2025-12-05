package s3

import (
	"sync"
	"time"

	"app/storage/kv"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	client         *s3.Client
	bucketName     string
	publicEndpoint string
	store          kv.Store[Object]
	cache          kv.Store[Object]
	maxObjectSize  int64
	maxBucketSize  int64

	sem chan struct{}

	objectLocks sync.Map

	mu         sync.RWMutex
	bucketSize int64
	inflight   int64
}

type Object struct {
	Bucket   string
	Key      string
	Mime     string
	Size     int64
	Created  time.Time
	Modified time.Time
}

type BucketParams struct {
	Client         *s3.Client
	BucketName     string
	PublicEndpoint string
	Store          kv.Store[Object]
	Cache          kv.Store[Object]
	MaxObjectSize  int64
	MaxBucketSize  int64
	MaxOps         int64
}
