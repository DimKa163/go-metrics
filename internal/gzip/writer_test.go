package gzip

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

func TestWriter(t *testing.T) {
	data := []byte("hello world")
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	gw := NewWriter(c.Writer)
	_, _ = gw.Write(data)
	gw.Close()
	b := w.Body.Bytes()
	reader, err := NewReader(io.NopCloser(bytes.NewReader(b)))
	assert.NoError(t, err)

	unzipped, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, data, unzipped)
}
