package transport

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Legit-Labs/legitify/cmd/progressbar"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/gofri/go-github-ratelimit/github_ratelimit/github_ratelimit_test"
)

const (
	singleSleepLimit = 90 * time.Second
)

func NewRateLimitWaiter(ctx context.Context, base http.RoundTripper) (*http.Client, error) {
	sleepCB := github_ratelimit.WithLimitDetectedCallback(func(ctx *github_ratelimit.CallbackContext) {
		log.Printf("facing secondary rate limit with request: %v. sleeping until: %v", ctx.Request.URL, *ctx.SleepUntil)
		progressbar.Report(progressbar.NewTimedBar("secondary rate limit", *ctx.SleepUntil))
	})
	limitCB := github_ratelimit.WithSingleSleepLimit(singleSleepLimit, func(ctx *github_ratelimit.CallbackContext) {
		log.Printf("secondary rate limit sleep is too long with request: %v, failing the request (%v > %v)",
			ctx.Request.URL, time.Until(*ctx.SleepUntil), singleSleepLimit)
	})

	if context_utils.GetSimulateSecondaryRateLimit(ctx) {
		var err error
		base, err = github_ratelimit_test.NewRateLimitInjecter(base,
			&github_ratelimit_test.SecondaryRateLimitInjecterOptions{
				Every: 3 * time.Second,
				Sleep: 5 * time.Second,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("injection error: %v", err)
		}
	}

	return github_ratelimit.NewRateLimitWaiterClient(base, sleepCB, limitCB)
}
