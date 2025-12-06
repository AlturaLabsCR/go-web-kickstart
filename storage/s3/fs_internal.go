package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"app/storage/s3/reader"
)

func (fs *FileSystem) putObject(key string, body io.Reader) error {
	path, err := fs.objectPath(key)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	rc, ok := body.(io.ReadCloser)
	if !ok {
		rc = io.NopCloser(body)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return err
	}

	return nil
}

func (fs *FileSystem) getObject(_ context.Context, key string) (io.ReadCloser, error) {
	path, err := fs.objectPath(key)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (fs *FileSystem) deleteObject(_ context.Context, key string) error {
	path, err := fs.objectPath(key)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (fs *FileSystem) makeSizeCallback() reader.ObjectReaderCallback {
	var lastSize int64

	return func(total int64) bool {
		delta := total - lastSize
		lastSize = total

		fs.mu.Lock()
		newInflight := fs.inflight + delta

		if fs.maxBucketSize > 0 && (fs.bucketSize+newInflight) > fs.maxBucketSize {
			fs.mu.Unlock()
			return false
		}

		fs.inflight = newInflight
		fs.mu.Unlock()

		if fs.maxObjectSize > 0 && total > fs.maxObjectSize {
			return false
		}

		return true
	}
}

func (fs *FileSystem) getObjectLock(key string) *sync.Mutex {
	lockIface, _ := fs.objectLocks.LoadOrStore(key, &sync.Mutex{})
	return lockIface.(*sync.Mutex)
}

func (fs *FileSystem) objectPath(key string) (string, error) {
	p := filepath.Join(fs.root, key)
	p = filepath.Clean(p)

	root := filepath.Clean(fs.root)
	if !strings.HasPrefix(p, root+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid object key: path escapes bucket root")
	}
	return p, nil
}
