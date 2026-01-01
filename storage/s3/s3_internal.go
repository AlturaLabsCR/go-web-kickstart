package s3

import (
	"context"
	"io"
	"sync"

	"app/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var _ storage.ObjectStorage = (*S3Bucket)(nil)

func (b *S3Bucket) putObject(ctx context.Context, key string, body io.Reader) error {
	_, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}

func (b *S3Bucket) getObject(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (b *S3Bucket) deleteObject(ctx context.Context, key string) (err error) {
	_, err = b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
	return err
}

func (b *S3Bucket) getObjectLock(key string) *sync.Mutex {
	lockIface, _ := b.objectLocks.LoadOrStore(key, &sync.Mutex{})
	return lockIface.(*sync.Mutex)
}

func (b *S3Bucket) remaining() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	rem := b.maxBucketSize - b.bucketSize - b.inflight
	if rem <= 0 {
		return 0
	}

	return min(b.maxObjectSize, rem)
}

func (b *S3Bucket) reserve(bytes int64) {
	b.mu.Lock()
	b.inflight += bytes
	b.mu.Unlock()
}

func (b *S3Bucket) release(bytes int64) {
	b.mu.Lock()
	b.inflight -= bytes
	b.mu.Unlock()
}

func (b *S3Bucket) commit(newBytes, oldBytes int64) {
	b.mu.Lock()
	b.inflight -= newBytes
	b.bucketSize -= oldBytes
	b.bucketSize += newBytes
	b.mu.Unlock()
}
