package tripper

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"net/http"
)

type HashTripper struct {
	rt  http.RoundTripper
	key string
	hash.Hash
}

func NewHashTripper(rt http.RoundTripper, key string) http.RoundTripper {
	return &HashTripper{
		rt:   rt,
		key:  key,
		Hash: hmac.New(sha256.New, []byte(key)),
	}
}

func (rt *HashTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		_, err := rt.Write(body)
		if err != nil {
			return nil, err
		}
		h := rt.Sum(nil)
		hsh := hex.EncodeToString(h)
		req.Header.Set("HashSHA256", hsh)
		req.Body = io.NopCloser(bytes.NewReader(body))
		rt.Reset()
	}
	return rt.rt.RoundTrip(req)
}
