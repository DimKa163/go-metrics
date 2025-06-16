package mgzip

import (
	"bufio"
	"compress/gzip"
	"errors"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

type GzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func NewGZIPWriter(writer gin.ResponseWriter) *GzipWriter {
	gz := gzip.NewWriter(writer)
	return &GzipWriter{
		ResponseWriter: writer,
		writer:         gz,
	}
}

func (g *GzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}
func (g *GzipWriter) Write(b []byte) (int, error) {
	return g.writer.Write(b)
}

func (g *GzipWriter) WriteHeader(statusCode int) {
	g.ResponseWriter.WriteHeader(statusCode)
}

func (g *GzipWriter) Flush() {
	_ = g.writer.Flush()
	g.ResponseWriter.Flush()
}

func (g *GzipWriter) Close() error {
	return g.writer.Close()
}

var _ http.Hijacker = (*GzipWriter)(nil)

func (g *GzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := g.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}
