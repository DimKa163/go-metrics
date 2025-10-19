package middleware

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

func IdentifyMiddleware(ipNet *net.IPNet) gin.HandlerFunc {
	return func(c *gin.Context) {
		ipStr := c.Request.Header.Get("X-Real-IP")
		if ipStr == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		ip := net.ParseIP(ipStr)

		if !ipNet.Contains(ip) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
