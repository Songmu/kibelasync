package client

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// request up to 10 times per second
// ref. https://github.com/kibela/kibela-api-v1-document/blob/master/README.md#1%E7%A7%92%E3%81%82%E3%81%9F%E3%82%8A%E3%81%AE%E3%83%AA%E3%82%AF%E3%82%A8%E3%82%B9%E3%83%88%E6%95%B0
const (
	reqInterval time.Duration = time.Second
	reqLimit    int           = 10
)

type rateLimitRoundTripper struct {
	transport http.RoundTripper
	l         *rate.Limiter

	mu           sync.RWMutex
	requestAfter time.Time
}

var _ http.RoundTripper = (*rateLimitRoundTripper)(nil)

// The cost budget is 300,000 per hour, and the cost limit for a single request is 10,000.
// So, when the remaining budget is under 10,000, the client waits next requesting until
// it recovers to 10,000. The cost budget will recover one per millisecond.
// ref. https://github.com/kibela/kibela-api-v1-document/blob/master/README.md#1%E6%99%82%E9%96%93%E3%81%94%E3%81%A8%E3%81%AB%E6%B6%88%E8%B2%BB%E3%81%A7%E3%81%8D%E3%82%8B%E3%82%B3%E3%82%B9%E3%83%88
func (rt *rateLimitRoundTripper) announceRemainingCost(remaining int) {
	const minimumCost = 10000
	if remaining < minimumCost {
		rt.setRequestAfter(
			time.Now().Add(time.Millisecond * time.Duration(remaining-minimumCost)),
		)
	}
}

func (rt *rateLimitRoundTripper) setRequestAfter(ti time.Time) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.requestAfter = ti
}

func (rt *rateLimitRoundTripper) canRequestAfter() time.Time {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	return rt.requestAfter
}

func newRateLimitRoundTripper() *rateLimitRoundTripper {
	return &rateLimitRoundTripper{
		l:         rate.NewLimiter(rate.Every(reqInterval/time.Duration(reqLimit)), reqLimit),
		transport: http.DefaultTransport,
	}
}

func (rt *rateLimitRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	if err := rt.l.Wait(ctx); err != nil {
		return nil, err
	}
	select {
	case <-time.After(time.Until(rt.canRequestAfter())):
		return rt.transport.RoundTrip(req)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
