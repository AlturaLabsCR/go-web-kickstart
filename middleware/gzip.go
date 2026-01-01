package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	acceptEncodingHeader  = "Accept-Encoding"
	contentEncodingHeader = "Content-Encoding"
	gzipHeaderValue       = "gzip"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(
			r.Header.Get(acceptEncodingHeader),
			gzipHeaderValue,
		) {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set(contentEncodingHeader, gzipHeaderValue)

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		next.ServeHTTP(gzw, r)
	})
}
