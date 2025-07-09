package client

import (
	"bytes"
	"github.com/cenkalti/backoff/v5"
	"io"
	"net/http"
)

type RetryRoundTripper struct {
	rt http.RoundTripper
}

func NewRetryRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &RetryRoundTripper{rt: rt}
}

func (rt *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	times := [3]int{1, 3, 5}
	attempt := 0
	return backoff.Retry(req.Context(), func() (*http.Response, error) {

		response, err := rt.rt.RoundTrip(req)
		if err != nil {
			return nil, backoff.Permanent(err)
		}
		if err = rt.drain(response); err != nil {
			return nil, backoff.Permanent(err)
		}

		if rt.shouldRetry(response) && attempt < len(times) {
			if req.Body != nil {
				req.Body = io.NopCloser(bytes.NewBuffer(body))
			}
			at := attempt
			attempt++
			return nil, backoff.RetryAfter(times[at])
		}
		return response, nil
	}, backoff.WithBackOff(backoff.NewExponentialBackOff()))
}

func (rt *RetryRoundTripper) shouldRetry(resp *http.Response) bool {
	switch resp.StatusCode {
	case http.StatusRequestTimeout:
	case http.StatusTooManyRequests:
	case http.StatusBadGateway:
	case http.StatusGatewayTimeout:
	case http.StatusServiceUnavailable:
	case http.StatusInternalServerError:
		return true
	default:
		return false
	}
	return false
}

func (rt *RetryRoundTripper) drain(response *http.Response) error {
	if response.Body != nil {
		_, err := io.Copy(io.Discard, response.Body)
		if err != nil {
			return err
		}
		response.Body.Close()
	}
	return nil
}
