package transport

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
)

type nopServer struct {
}

func upTo1SecDelay() time.Duration {
	return time.Duration(int(time.Millisecond) * (rand.Int() % 1000))
}

func (n *nopServer) RoundTrip(r *http.Request) (*http.Response, error) {
	time.Sleep(upTo1SecDelay() / 100)
	return &http.Response{
		Body:   io.NopCloser(strings.NewReader("some response")),
		Header: http.Header{},
	}, nil
}

func TestSecondaryRateLimit(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	const requests = 5000
	os.Setenv(RateLimitInjecterEvery, "3")
	os.Setenv(RateLimitInjecterDuration, "1")

	c := NewRateLimitWaiter(&nopServer{})
	gw := group_waiter.New()
	for i := 0; i < requests; i++ {
		sleepTime := upTo1SecDelay() / 100
		if sleepTime.Milliseconds()%2 == 0 {
			sleepTime = 0 // bias towards no-sleep for high parallelism
		}
		time.Sleep(sleepTime)
		gw.Do(func() {
			time.Sleep(upTo1SecDelay() / 50)
			_, _ = c.Get("/bla")
		})
	}
	gw.Wait()
	log.Printf("%v abuse attempts out of %v requests", c.Transport.(*secondaryRateLimitWaiter).base.(*secondaryRateLimitInjecter).AbuseAttempts, requests)
}
