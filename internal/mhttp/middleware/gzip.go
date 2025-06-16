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
		acceptTypes := c.Request.Header.Get("Accept")
		supportTypes := strings.Contains(acceptTypes, ContentTypeJSON) || strings.Contains(acceptTypes, ContentTypeHTML)
		supportsGzip := strings.Contains(acceptEncoding, AcceptEncodingGZIP)
		if supportsGzip {
			c.Header("Content-Encoding", ContentEncodingGZIP)
			gz := mgzip.NewGZIPWriter(c.Writer)
			c.Writer = gz
			defer func() {
				c.Header("Content-Length", "0")
				_ = gz.Close()
			}()
		}

		contentType := c.Request.Header.Get("Content-Type")
		contentEncoding := c.Request.Header.Get("Content-Encoding")
		supportTypes = strings.Contains(contentType, ContentTypeJSON) || strings.Contains(contentType, ContentTypeHTML)
		sendsGzip := strings.Contains(contentEncoding, ContentEncodingGZIP)
		if sendsGzip && supportTypes {
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
