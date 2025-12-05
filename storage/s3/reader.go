package s3

import (
	"bytes"
	"io"
	"net/http"
)

type ByteReserver interface {
	ReserveBytes(n int64) error
	ReleaseBytes(n int64)
}

type objectReader struct {
	r        io.Reader
	ct       string
	reserver ByteReserver
	reserved int64
}

func (b *Bucket) newStorageReader(r io.Reader) (*objectReader, error) {
	sniff := make([]byte, 512)
	n, err := io.ReadFull(r, sniff)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	sniff = sniff[:n]

	return &objectReader{
		r:        io.MultiReader(bytes.NewReader(sniff), r),
		ct:       http.DetectContentType(sniff),
		reserver: b,
	}, nil
}

func (r *objectReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	if n > 0 {
		if err2 := r.reserver.ReserveBytes(int64(n)); err2 != nil {
			return n, err2
		}
		r.reserved += int64(n)
	}
	return n, err
}

func (r *objectReader) Close() error {
	if r.reserver != nil && r.reserved > 0 {
		r.reserver.ReleaseBytes(r.reserved)
	}
	return nil
}
