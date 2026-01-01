package fs

import (
	"sync"

	"app/cache"
)

type FS struct {
	root          string
	cache         cache.Cache
	l2cache       cache.Cache
	maxObjectSize int64
	maxBucketSize int64

	objectLocks sync.Map

	mu         sync.RWMutex
	bucketSize int64
	inflight   int64
}

type FSParams struct {
	Root          string
	Cache         cache.Cache
	MaxObjectSize int64
	MaxBucketSize int64

	// optional params
	L2Cache cache.Cache
}
