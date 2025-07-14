package middleware

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"hash"
	"io"
	"net"
	"net/http"
)

const HashHeader = "HashSHA256"

func Hash(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(HashHeader)
		if header != "" && c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			expectedMAC := hmac.New(sha256.New, []byte(key))
			expectedMAC.Write(body)
			expectedSignature := expectedMAC.Sum(nil)
			sign, err := hex.DecodeString(header)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if !hmac.Equal(sign, expectedSignature) {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}
		c.Writer = NewHashWriter(c.Writer, key)
		c.Next()
	}
}

type HashWriter struct {
	gin.ResponseWriter
	hash.Hash
}

func NewHashWriter(writer gin.ResponseWriter, key string) *HashWriter {
	return &HashWriter{ResponseWriter: writer, Hash: hmac.New(sha256.New, []byte(key))}
}

func (h *HashWriter) Size() int {
	return h.ResponseWriter.Size()
}

func (h *HashWriter) Write(b []byte) (int, error) {
	_, err := h.Hash.Write(b)
	if err != nil {
		return 0, err
	}
	return h.ResponseWriter.Write(b)
}

func (h *HashWriter) WriteHeader(statusCode int) {
	sum := h.Sum(nil)
	hsh := hex.EncodeToString(sum)
	h.Header().Set("HashSHA256", hsh)
	h.ResponseWriter.WriteHeader(statusCode)
}

func (h *HashWriter) Flush() {
	h.ResponseWriter.Flush()
}

var _ http.Hijacker = (*HashWriter)(nil)

func (h *HashWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := h.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}
