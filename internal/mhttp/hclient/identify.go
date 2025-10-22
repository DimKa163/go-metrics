package hclient

import (
	"net/http"
)

func UseIdentifyHandler(ip string) RequestHandlerFactory {
	return func(transport http.RoundTripper) http.RoundTripper {
		return NewIdentifyTripper(transport, ip)
	}
}

type IdentifyTripper struct {
	rt http.RoundTripper
	ip string
}

func NewIdentifyTripper(rt http.RoundTripper, ip string) *IdentifyTripper {
	return &IdentifyTripper{rt: rt, ip: ip}
}

func (t *IdentifyTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(
		"X-Real-IP",
		t.ip,
	)
	return t.rt.RoundTrip(req)
}
