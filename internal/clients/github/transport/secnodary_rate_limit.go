package transport

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Legit-Labs/legitify/internal/screen"
)

const (
	secondaryRateLimitResetHeader = "Retry-After"
)

type secondaryRateLimitWaiter struct {
	base       http.RoundTripper
	sleepUntil *time.Time
	lock       sync.RWMutex
}

func NewRateLimitWaiter(base http.RoundTripper) *http.Client {
	return &http.Client{Transport: &secondaryRateLimitWaiter{
		base: base,
	}}
}

func (t *secondaryRateLimitWaiter) RoundTrip(request *http.Request) (*http.Response, error) {
	t.waitForRateLimit()

	resp, err := t.base.RoundTrip(request)
	if err != nil {
		return resp, err
	}

	// primary limit is handled by the client as error (we don't want to sleep until end of hour anyway)
	secondaryLimit, err := getSecondaryLimitTime(resp)
	if err != nil {
		return resp, err
	}

	t.updateRateLimit(secondaryLimit)

	return resp, nil
}

func (t *secondaryRateLimitWaiter) waitForRateLimit() {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if t.sleepUntil == nil {
		return
	}
	sleepTime := time.Until(*t.sleepUntil)
	if sleepTime < 0 {
		return
	}

	screen.Printf("Secondary rate limit reached! Sleeping for %v seconds", sleepTime.Seconds())
	time.Sleep(sleepTime)
}

func (t *secondaryRateLimitWaiter) updateRateLimit(secondaryLimit *time.Time) {
	if secondaryLimit == nil {
		return
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	t.sleepUntil = secondaryLimit
}

func getSecondaryLimitTime(resp *http.Response) (*time.Time, error) {
	// insipred by go-github code
	v := resp.Header[secondaryRateLimitResetHeader]
	if len(v) == 0 {
		return nil, nil
	}

	// According to GitHub support, the "Retry-After" header value will be
	// an integer which represents the number of seconds that one should
	// wait before resuming making requests.
	retryAfterSeconds, err := strconv.ParseInt(v[0], 10, 64) // Error handling is noop.
	if err != nil {
		return nil, err
	}
	retryAfter := time.Duration(retryAfterSeconds) * time.Second
	sleepUntil := time.Now().Add(retryAfter)

	return &sleepUntil, nil
}
