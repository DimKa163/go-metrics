package tripper

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
)

type GzipTripper struct {
	rt http.RoundTripper
}

func NewGzip(rt http.RoundTripper) http.RoundTripper {
	return &GzipTripper{rt: rt}
}

func (rt *GzipTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		buffer := bytes.NewBuffer(nil)
		writer := gzip.NewWriter(buffer)
		_, err = writer.Write(body)
		if err != nil {
			return nil, err
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}
		req.ContentLength = int64(buffer.Len())
		req.Header.Add("Content-Encoding", "gzip")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(buffer.Len()))
		req.Body = io.NopCloser(buffer)
	}
	return rt.rt.RoundTrip(req)
}
