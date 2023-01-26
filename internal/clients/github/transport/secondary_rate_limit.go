package transport

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v49/github"
)

type secondaryRateLimitWaiter struct {
	base       http.RoundTripper
	sleepUntil *time.Time
	lock       sync.RWMutex
}

func NewRateLimitWaiter(base http.RoundTripper) *http.Client {
	return &http.Client{Transport: &secondaryRateLimitWaiter{
		base: newRateLimitInjecter(base),
	}}
}

// RoundTrip handles the secondary rate limit by waiting for it to finish before issuing new requests.
// If a request got a secondary rate limit error as a response, we retry the request after waiting.
// Issuing more requests during a secondary rate limit may cause a ban from the server side,
// so we want to prevent these requests, not just for the sake of cpu/network utilization.
// Nonetheless, there is no way to prevent subtle race conditions without completely serializing the requests,
// so we prefer to let some slip in case of a race condition, i.e.,
// after a retry-after response is received and before it it processed,
// a few other (parallel) requests may be issued.
func (t *secondaryRateLimitWaiter) RoundTrip(request *http.Request) (*http.Response, error) {
	t.waitForRateLimit()

	resp, err := t.base.RoundTrip(request)
	if err != nil {
		return resp, err
	}

	// primary limit is handled by the client as error (we don't want to sleep until end of hour anyway)
	secondaryLimit := parseSecondaryLimitTime(resp)
	if secondaryLimit == nil {
		return resp, nil
	}

	t.updateRateLimit(*secondaryLimit)

	return t.RoundTrip(request)
}

func (t *secondaryRateLimitWaiter) currentSleepTimeUnlocked() time.Duration {
	if t.sleepUntil == nil {
		return 0
	}
	return time.Until(*t.sleepUntil)
}

// waitForRateLimit waits for the cooldown time to finish if a secondary rate limit is active.
// it sleeps holding the RLock because no updates can happen during the sleep anyway.
func (t *secondaryRateLimitWaiter) waitForRateLimit() {
	t.lock.RLock()
	sleepTime := t.currentSleepTimeUnlocked()
	t.lock.RUnlock()

	time.Sleep(sleepTime)
}

// updateRateLimit updates the active rate limit and prints a message to the user.
// the rate limit is not updated if there's already an active rate limit.
// it never waits because the retry handles sleeping anyway.
func (t *secondaryRateLimitWaiter) updateRateLimit(secondaryLimit time.Time) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// check before update if there is already an active rate limit
	if t.currentSleepTimeUnlocked() > 0 {
		return
	}

	t.sleepUntil = &secondaryLimit
	newSleepTime := t.currentSleepTimeUnlocked()

	// check after updating because the secondary rate limit might have passed while we waited for the lock
	if newSleepTime <= 0 {
		return
	}
	log.Printf("Secondary rate limit reached! Sleeping for %.2f seconds [%v --> %v]", newSleepTime.Seconds(), time.Now(), t.sleepUntil)
}

func parseSecondaryLimitTime(resp *http.Response) *time.Time {
	err := github.CheckResponse(resp)
	if err == nil {
		return nil
	}
	abuse, ok := err.(*github.AbuseRateLimitError)
	if !ok {
		return nil
	}
	retryAfter := abuse.RetryAfter
	if retryAfter == nil {
		return nil
	}

	sleepUntil := time.Now().Add(*retryAfter)
	return &sleepUntil
}

type secondaryRateLimitInjecter struct {
	base          http.RoundTripper
	injectEvery   time.Duration
	sleepDuration time.Duration
	blockUntil    *time.Time
	lock          sync.Mutex
	AbuseAttempts int
}

const (
	RateLimitInjecterEvery    = "SECONDARY_RATE_LIMIT_INJECTER_EVERY"
	RateLimitInjecterDuration = "SECONDARY_RATE_LIMIT_INJECTER_DURATION"
)

func newRateLimitInjecter(base http.RoundTripper) http.RoundTripper {
	every := os.Getenv(RateLimitInjecterEvery)
	duration := os.Getenv(RateLimitInjecterDuration)
	if every == "" || every == "0" || duration == "" || duration == "0" {
		return base
	} else {
		injectEvery, err := strconv.ParseInt(every, 10, 0)
		if err != nil && injectEvery < 0 {
			log.Panicf("unexpected secondary rate limit injection every: %v / %v", injectEvery, err)
		}
		injectDuration, err := strconv.ParseInt(duration, 10, 0)
		if err != nil && injectDuration < 0 {
			log.Panicf("unexpected secondary rate limit injection duration: %v / %v", injectDuration, err)
		}

		everyDuration := time.Duration(int64(time.Second) * int64(injectEvery))
		sleepDuration := time.Duration(int64(time.Second) * int64(injectDuration))
		log.Printf("secondary rate limit injection is active. injecting %v sleep every %v", sleepDuration, everyDuration)

		injecter := &secondaryRateLimitInjecter{
			base:          base,
			injectEvery:   everyDuration,
			sleepDuration: sleepDuration,
		}
		return http.Client{Transport: injecter}.Transport
	}
}

func (t *secondaryRateLimitInjecter) toRetryResponse(resp *http.Response) *http.Response {
	resp.StatusCode = http.StatusForbidden
	timeUntil := time.Until(*t.blockUntil)
	if timeUntil.Nanoseconds()%int64(time.Second) > 0 {
		timeUntil += time.Second
	}
	resp.Header.Set("Retry-After", fmt.Sprintf("%v", int(timeUntil.Seconds())))
	doc_url := "https://docs.github.com/en/rest/guides/best-practices-for-integrators?apiVersion=2022-11-28#secondary-rate-limits"
	resp.Body = io.NopCloser(strings.NewReader(`{"documentation_url":"` + doc_url + `"}`))
	return resp
}

func (t *secondaryRateLimitInjecter) RoundTrip(request *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(request)
	if err != nil {
		return resp, err
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	now := time.Now()
	if t.blockUntil == nil {
		t.blockUntil = &now
	}

	// on-going rate limit
	if t.blockUntil.After(now) {
		t.AbuseAttempts++
		return t.toRetryResponse(resp), nil
	}

	nextStart := t.blockUntil.Add(t.injectEvery)
	if !now.Before(nextStart) {
		nextEnd := nextStart.Add(t.sleepDuration)
		t.blockUntil = &nextEnd
		log.Printf("inject sleep until: %v\n", t.blockUntil)
		return t.toRetryResponse(resp), nil
	}

	return resp, nil
}
