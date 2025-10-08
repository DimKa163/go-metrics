package middleware

import (
	"bytes"
	"github.com/DimKa163/go-metrics/internal/crypto"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func CryptoMiddleware(decrypter *crypto.Decrypter) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		decrypted, err := decrypter.Decrypt(body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(decrypted))
	}
}
