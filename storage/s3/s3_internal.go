package s3

import (
	"context"
	"sync"
)

func (b *Bucket) getObjectLock(key string) *sync.Mutex {
	lockIface, _ := b.objectLocks.LoadOrStore(key, &sync.Mutex{})
	return lockIface.(*sync.Mutex)
}

func (b *Bucket) wait(ctx context.Context) error {
	select {
	case b.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *Bucket) done() {
	select {
	case <-b.sem:
	default:
	}
}
