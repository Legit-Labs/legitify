package transport

import (
	"log"
	"net/http"
	"time"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
)

const (
	singleSleepLimit = 90 * time.Second
)

func NewRateLimitWaiter(base http.RoundTripper) (*http.Client, error) {
	sleepCB := github_ratelimit.WithLimitDetectedCallback(func(ctx *github_ratelimit.CallbackContext) {
		log.Printf("secondary rate limit detected. Sleeping for %v until: %v", time.Until(*ctx.SleepUntil), *ctx.SleepUntil)
	})
	limitCB := github_ratelimit.WithSingleSleepLimit(singleSleepLimit, func(ctx *github_ratelimit.CallbackContext) {
		log.Printf("secondary rate limit sleep is too long, failing the request (%v > %v)", time.Until(*ctx.SleepUntil), singleSleepLimit)
	})
	return github_ratelimit.NewRateLimitWaiterClient(base, sleepCB, limitCB)
}
