package tripper

import (
	"bytes"
	"github.com/DimKa163/go-metrics/internal/crypto"
	"io"
	"net/http"
)

type CryptoTripper struct {
	rt http.RoundTripper
	en *crypto.Encrypter
}

func NewCryptoTripper(rt http.RoundTripper, encrypter *crypto.Encrypter) *CryptoTripper {
	return &CryptoTripper{rt: rt, en: encrypter}
}

func (t *CryptoTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		cipherBody, err := t.en.Encrypt(body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(cipherBody))
	}
	return t.rt.RoundTrip(req)
}
