package reader

import (
	"context"
	"io"
)

type errStr string

func (e errStr) Error() string {
	return string(e)
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
