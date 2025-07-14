package tripper

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

type HashTripper struct {
	rt  http.RoundTripper
	key string
}

func NewHashTripper(rt http.RoundTripper, key string) http.RoundTripper {
	return &HashTripper{
		rt:  rt,
		key: key,
	}
}

func (rt *HashTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	var err error
	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		hsh := hmac.New(sha256.New, []byte(rt.key))
		_, err := hsh.Write(body)
		if err != nil {
			return nil, err
		}
		h := hsh.Sum(nil)
		str := hex.EncodeToString(h)
		req.Header.Set("HashSHA256", str)
		req.Body = io.NopCloser(bytes.NewReader(body))
		hsh.Reset()
	}
	return rt.rt.RoundTrip(req)
}
