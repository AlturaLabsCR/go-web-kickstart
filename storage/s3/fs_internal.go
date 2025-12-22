package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func (fs *FileSystem) putObject(_ context.Context, key string, body io.Reader) error {
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
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(f, body)
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

func (fs *FileSystem) remaining() int64 {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return min(fs.maxObjectSize, fs.maxBucketSize-fs.bucketSize-fs.inflight)
}

func (fs *FileSystem) reserve(bytes int64) {
	fs.mu.Lock()
	fs.inflight += bytes
	fs.mu.Unlock()
}

func (fs *FileSystem) release(bytes int64) {
	fs.mu.Lock()
	fs.inflight -= bytes
	fs.mu.Unlock()
}

func (fs *FileSystem) commit(newBytes, oldBytes int64) {
	fs.mu.Lock()
	fs.inflight -= newBytes
	fs.bucketSize -= oldBytes
	fs.bucketSize += newBytes
	fs.mu.Unlock()
}
