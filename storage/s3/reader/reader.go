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

const ErrOverflow = errStr("reader overflow")

type ObjectReaderCallback func(total int64) bool

type ObjectReader struct {
	r        io.Reader
	ctx      context.Context
	ct       string
	size     int64
	callback ObjectReaderCallback
	stopped  bool
}

func (o *ObjectReader) ContentType() string { return o.ct }

func (o *ObjectReader) Size() int64 { return o.size }

func (o *ObjectReader) Read(p []byte) (int, error) {
	if o.stopped {
		return 0, ErrOverflow
	}

	if err := o.ctx.Err(); err != nil {
		o.stopped = true
		return 0, ErrOverflow
	}

	for {
		n, err := o.r.Read(p)

		if n > 0 {
			o.size += int64(n)

			if o.callback != nil && !o.callback(o.size) {
				o.stopped = true
				return n, ErrOverflow
			}

			return n, err
		}

		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}
			return 0, err
		}
	}
}

func NewObjectReader(ctx context.Context, r io.Reader, cb ObjectReaderCallback) (*ObjectReader, error) {
	const sniffLen = 512
	buf := make([]byte, sniffLen)

	n, err := r.Read(buf)
	if n < 0 {
		n = 0
	}
	if err != nil && err != io.EOF {
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
