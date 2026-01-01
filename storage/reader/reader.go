// Package reader wraps io.Reader and io.ReadCloser to provide better functionality for object storage
package reader

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

const (
	ErrOverflow = errStr("object too large")
	ErrCanceled = errStr("ctx closed")
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
