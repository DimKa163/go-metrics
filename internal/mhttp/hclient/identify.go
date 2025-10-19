package hclient

import (
	"net"
	"net/http"
)

var ipAddr string

func init() {
	ipAddr, _ = getLocalIP()
}
func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
func UseIdentifyHandler() RequestHandler {
	return func(transport http.RoundTripper) http.RoundTripper {
		return NewIdentifyTripper(transport)
	}
}

type IdentifyTripper struct {
	rt http.RoundTripper
}

func NewIdentifyTripper(rt http.RoundTripper) *IdentifyTripper {
	return &IdentifyTripper{rt: rt}
}

func (t *IdentifyTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(
		"X-Real-IP",
		ipAddr,
	)
	return t.rt.RoundTrip(req)
}
