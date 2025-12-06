package s3

import (
	"sync"

	"app/storage/kv"
)

type FileSystem struct {
	root           string
	publicEndpoint string
	store          kv.Store[Object]
	cache          kv.Store[Object]
	maxObjectSize  int64
	maxBucketSize  int64

	objectLocks sync.Map

	mu         sync.RWMutex
	bucketSize int64
	inflight   int64
}
