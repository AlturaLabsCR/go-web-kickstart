// Package reader wraps io.Reader and io.ReadCloser to provide better functionality for object storage
package reader

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type errStr string

func (e errStr) Error() string {
	return string(e)
}

const (
	ErrOverflow = errStr("object exceeds maximum allowed size")
	ErrCanceled = errStr("read canceled by context")
)

func NewObjectReader(ctx context.Context, r io.Reader, maxBytes int64) (*bytes.Reader, string, error) {
	limit := maxBytes + 1

	rd := &ctxReader{ctx: ctx, r: r}

	data, err := io.ReadAll(io.LimitReader(rd, limit))
	if err != nil {
		return nil, "", err
	}

	if int64(len(data)) > maxBytes {
		return nil, "", ErrOverflow
	}

	detectSize := min(len(data), 512)
	ct := http.DetectContentType(data[:detectSize])

	return bytes.NewReader(data), ct, nil
}

type ctxReader struct {
	ctx context.Context
	r   io.Reader
}

func (cr *ctxReader) Read(p []byte) (int, error) {
	if err := cr.ctx.Err(); err != nil {
		return 0, ErrCanceled
	}
	return cr.r.Read(p)
}
