package transport

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Legit-Labs/legitify/cmd/progressbar"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/gofri/go-github-ratelimit/github_ratelimit/github_ratelimit_test"
)

const (
	singleSleepLimit = 90 * time.Second
	injecterEnvKey   = "LEGITIFY_SECONDARY_RATE_LIMIT_INJECTION"
)

func NewRateLimitWaiter(base http.RoundTripper) (*http.Client, error) {
	sleepCB := github_ratelimit.WithLimitDetectedCallback(func(ctx *github_ratelimit.CallbackContext) {
		log.Printf("facing secondary rate limit with request: %v. sleeping until: %v", ctx.Request.URL, *ctx.SleepUntil)
		progressbar.Report(progressbar.NewTimedBar("Secondary rate limit", *ctx.SleepUntil))
	})
	limitCB := github_ratelimit.WithSingleSleepLimit(singleSleepLimit, func(ctx *github_ratelimit.CallbackContext) {
		log.Printf("secondary rate limit sleep is too long, failing the request (%v > %v)", time.Until(*ctx.SleepUntil), singleSleepLimit)
	})

	if os.Getenv(injecterEnvKey) == "1" {
		var err error
		base, err = github_ratelimit_test.NewRateLimitInjecter(base,
			&github_ratelimit_test.SecondaryRateLimitInjecterOptions{
				Every: time.Second,
				Sleep: time.Second,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("injection error: %v", err)
		}
	}

	return github_ratelimit.NewRateLimitWaiterClient(base, sleepCB, limitCB)
}
