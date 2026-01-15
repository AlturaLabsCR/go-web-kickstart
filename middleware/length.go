package middleware

import (
	"bytes"
	"net/http"
	"strconv"
)

func ContentLength(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := &bytes.Buffer{}
		bw := &bufferedResponseWriter{
			ResponseWriter: w,
			buf:            buf,
			statusCode:     http.StatusOK,
			headers:        w.Header(),
		}

		next.ServeHTTP(bw, r)

		for key, values := range bw.headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

		w.WriteHeader(bw.statusCode)
		_, _ = buf.WriteTo(w)
	})
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
	headers    http.Header
}

func (bw *bufferedResponseWriter) Header() http.Header {
	return bw.headers
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(code int) {
	bw.statusCode = code
}
