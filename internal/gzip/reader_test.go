package gzip

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

func TestReader(t *testing.T) {
	data := []byte("hello world")
	respReader := prepareBody(data)
	reader, err := NewReader(io.NopCloser(bytes.NewReader(respReader)))
	assert.NoError(t, err)

	unzipped, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, data, unzipped)
}

func prepareBody(data []byte) []byte {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	gw := NewWriter(c.Writer)
	_, _ = gw.Write(data)
	gw.Close()

	return w.Body.Bytes()
}
