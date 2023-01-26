package transport

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gregjones/httpcache"
)

func NewCacheTransport() http.RoundTripper {
	return httpcache.NewMemoryCacheTransport()
}

const CACHE_TRACKER_KEY = "CACHE_TRACKER_ENABLED"
const CACHE_TRACKER_ENABLED = "1"

type cacheTracker struct {
	lock   sync.Mutex
	base   http.RoundTripper
	total  int
	cached int
}

func NewCacheTracker(base *http.Client) *http.Client {
	if os.Getenv(CACHE_TRACKER_KEY) != CACHE_TRACKER_ENABLED {
		return base
	}

	return &http.Client{
		Transport: &cacheTracker{
			base: base.Transport,
		},
	}
}

func (c *cacheTracker) handleResp(resp *http.Response) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := resp.Header["X-From-Cache"]; ok {
		c.cached++
		log.Printf("cached %v", resp.Request.URL)
	}
	c.total++
	log.Printf("cached %v/%v", c.cached, c.total)
}

func (c *cacheTracker) RoundTrip(request *http.Request) (*http.Response, error) {
	resp, err := c.base.RoundTrip(request)
	if resp != nil {
		c.handleResp(resp)
	}
	return resp, err
}
