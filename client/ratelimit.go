package client

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// request up to 10 times per second
const (
	reqInterval time.Duration = time.Second
	reqLimit    int           = 10
)

type rateLimitRoundTripper struct {
	transport http.RoundTripper
	l         *rate.Limiter
}

func newRateLimitRoundTripper() http.RoundTripper {
	return &rateLimitRoundTripper{
		l:         rate.NewLimiter(rate.Every(reqInterval/time.Duration(reqLimit)), reqLimit),
		transport: http.DefaultTransport,
	}
}

func (t *rateLimitRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if err := t.l.Wait(ctx); err != nil {
		return nil, err
	}
	return t.transport.RoundTrip(req)
}
