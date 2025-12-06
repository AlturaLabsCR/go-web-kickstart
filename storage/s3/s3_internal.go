package s3

import (
	"bytes"
	"context"
	"io"
	"sync"

	"app/storage/s3/reader"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (b *Bucket) putObject(ctx context.Context, key string, body io.Reader) error {
	// FIXME: pass a compatible reader, not the data

	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	_, err = b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
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

func (b *Bucket) makeSizeCallback() reader.ObjectReaderCallback {
	var lastSize int64

	return func(total int64) bool {
		delta := total - lastSize
		lastSize = total

		b.mu.Lock()
		newInflight := b.inflight + delta

		if b.maxBucketSize > 0 && (b.bucketSize+newInflight) > b.maxBucketSize {
			b.mu.Unlock()
			return false
		}

		b.inflight = newInflight
		b.mu.Unlock()

		if b.maxObjectSize > 0 && total > b.maxObjectSize {
			return false
		}

		return true
	}
}

func (b *Bucket) getObjectLock(key string) *sync.Mutex {
	lockIface, _ := b.objectLocks.LoadOrStore(key, &sync.Mutex{})
	return lockIface.(*sync.Mutex)
}
