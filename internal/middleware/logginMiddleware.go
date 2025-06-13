package middleware

import (
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func WithLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bytes []byte
		var body string
		_, err := c.Request.Body.Read(bytes)
		if err != nil {
			body = string(bytes)
		}
		logging.Log.Info(
			"got incoming HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("body", body),
		)
		startTime := time.Now()
		c.Next()
		elapsed := time.Since(startTime)
		logging.Log.Info("Processed HTTP request", zap.Int("status", c.Writer.Status()), zap.Int("size", c.Writer.Size()),
			zap.Duration("elapsed", elapsed))
	}
}
