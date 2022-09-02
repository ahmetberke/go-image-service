package handlers

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

type gzipHandler struct {
}

func NewGzipHandler() *gzipHandler {
	return &gzipHandler{}
}

func (g *gzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Printf("ACCEPT ENCODING : %s", request.Header.Get("Accept-Encoding"))
		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			// create a gziped response
			wrw := NewWrappedResponseWriter(writer)
			wrw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(wrw, request)
			defer wrw.Flush()

			return
		}

		// handle normal
		next.ServeHTTP(writer, request)
	})
}

type WrappedReponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedReponseWriter {
	gw := gzip.NewWriter(rw)

	return &WrappedReponseWriter{rw: rw, gw: gw}
}

func (wr *WrappedReponseWriter) Header() http.Header {
	return wr.rw.Header()
}

func (wr *WrappedReponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

func (wr *WrappedReponseWriter) WriteHeader(statuscode int) {
	wr.rw.WriteHeader(statuscode)
}

func (wr *WrappedReponseWriter) Flush() {
	wr.gw.Flush()
	wr.gw.Close()
}
