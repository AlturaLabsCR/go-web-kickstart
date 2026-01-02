package fs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"app/storage"
)

var _ storage.ObjectStorage = (*FS)(nil)

func (fs *FS) putObject(_ context.Context, key string, body io.Reader) error {
	path, err := fs.objectPath(key)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tmp := path + ".tmp"

	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(tmp)
	}()

	if _, err := io.Copy(f, body); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		return err
	}

	df, err := os.Open(dir)
	if err == nil {
		_ = df.Sync()
		_ = df.Close()
	}

	return nil
}

func (fs *FS) getObject(_ context.Context, key string) (io.ReadCloser, error) {
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

func (fs *FS) deleteObject(_ context.Context, key string) error {
	path, err := fs.objectPath(key)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (fs *FS) getObjectLock(key string) *sync.Mutex {
	lockIface, _ := fs.objectLocks.LoadOrStore(key, &sync.Mutex{})
	return lockIface.(*sync.Mutex)
}

func (fs *FS) objectPath(key string) (string, error) {
	p := filepath.Join(fs.root, key)
	p = filepath.Clean(p)

	root := filepath.Clean(fs.root)
	if !strings.HasPrefix(p, root+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid object key: path escapes bucket root")
	}
	return p, nil
}

func (fs *FS) remaining() int64 {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	rem := fs.maxBucketSize - fs.bucketSize - fs.inflight
	if rem <= 0 {
		return 0
	}

	return min(fs.maxObjectSize, rem)
}

func (fs *FS) reserve(bytes int64) {
	fs.mu.Lock()
	fs.inflight += bytes
	fs.mu.Unlock()
}

func (fs *FS) release(bytes int64) {
	fs.mu.Lock()
	fs.inflight -= bytes
	fs.mu.Unlock()
}

func (fs *FS) commit(newBytes, oldBytes int64) {
	fs.mu.Lock()
	fs.inflight -= newBytes
	fs.bucketSize -= oldBytes
	fs.bucketSize += newBytes
	fs.mu.Unlock()
}
