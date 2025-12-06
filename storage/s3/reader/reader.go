// Package reader wraps io.Reader and io.ReadCloser to provide better functionality for object storage
package reader

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type ObjectReaderCallback func(total int64) error

type ObjectReader struct {
	r        io.Reader
	ctx      context.Context
	ct       string
	size     int64
	callback ObjectReaderCallback
}

func (o *ObjectReader) ContentType() string { return o.ct }

func (o *ObjectReader) Size() int64 { return o.size }

func (o *ObjectReader) Read(p []byte) (int, error) {
	if err := o.ctx.Err(); err != nil {
		return 0, err
	}

	n, err := o.r.Read(p)
	if n > 0 {
		o.size += int64(n)

		if o.callback != nil {
			if cbErr := o.callback(o.size); cbErr != nil {
				return n, cbErr
			}
		}
	}

	if o.ctx.Err() != nil {
		return n, o.ctx.Err()
	}

	return n, err
}

func NewObjectReader(ctx context.Context, r io.Reader, cb ObjectReaderCallback) (*ObjectReader, error) {
	const sniffLen = 512
	buf := make([]byte, sniffLen)

	n, err := io.ReadAtLeast(io.LimitReader(r, sniffLen), buf, 1)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	ct := http.DetectContentType(buf[:n])

	rr := io.MultiReader(bytes.NewReader(buf[:n]), r)

	return &ObjectReader{
		r:        rr,
		ctx:      ctx,
		ct:       ct,
		size:     0,
		callback: cb,
	}, nil
}
