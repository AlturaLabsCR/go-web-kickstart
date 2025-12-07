package s3

import (
	"context"
	"io"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (b *Bucket) putObject(ctx context.Context, key string, body io.Reader) error {
	_, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}

func (b *Bucket) getObject(ctx context.Context, key string) (body io.ReadCloser, err error) {
	out, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
	return out.Body, err
}

func (b *Bucket) deleteObject(ctx context.Context, key string) (err error) {
	_, err = b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
	return err
}

func (b *Bucket) getObjectLock(key string) *sync.Mutex {
	lockIface, _ := b.objectLocks.LoadOrStore(key, &sync.Mutex{})
	return lockIface.(*sync.Mutex)
}

func (b *Bucket) remaining() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return min(b.maxObjectSize, b.maxBucketSize-b.bucketSize-b.inflight)
}

func (b *Bucket) reserve(bytes int64) {
	b.mu.Lock()
	b.inflight += bytes
	b.mu.Unlock()
}

func (b *Bucket) release(bytes int64) {
	b.mu.Lock()
	b.inflight -= bytes
	b.mu.Unlock()
}

func (b *Bucket) commit(newBytes, oldBytes int64) {
	b.mu.Lock()
	b.inflight -= newBytes
	b.bucketSize -= oldBytes
	b.bucketSize += newBytes
	b.mu.Unlock()
}
