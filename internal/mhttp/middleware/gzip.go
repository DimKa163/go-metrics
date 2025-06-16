package middleware

import (
	"github.com/DimKa163/go-metrics/internal/mgzip"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	ContentTypeJSON    = "application/json"
	ContentTypeHTML    = "text/html"
	AcceptEncodingGZIP = "gzip"

	ContentEncodingGZIP = "gzip"
)

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptEncoding := c.Request.Header.Get("Accept-Encoding")
		contentType := c.Request.Header.Get("Content-Type")
		supportsGzip := strings.Contains(acceptEncoding, AcceptEncodingGZIP)
		supportTypes := strings.Contains(contentType, ContentTypeJSON) || strings.Contains(contentType, ContentTypeHTML)
		if supportsGzip && supportTypes {
			c.Header("Content-Encoding", ContentEncodingGZIP)
			gz := mgzip.NewGZIPWriter(c.Writer)
			c.Writer = gz
			defer func() {
				c.Header("Content-Length", "0")
				_ = gz.Close()
			}()
		}
		contentEncoding := c.Request.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, ContentEncodingGZIP)
		if sendsGzip {
			gz, err := mgzip.NewGZIPReader(c.Request.Body)
			if err != nil {
				return
			}
			c.Request.Body = gz
			defer gz.Close()
		}
		c.Next()
	}
}
