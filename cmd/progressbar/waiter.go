package progressbar

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/Legit-Labs/legitify/internal/common/group_waiter"
)

// pbWaiter waits for all required bars to be created before waiting for the progress bar to finish
type pbWaiter struct {
	realWait   group_waiter.Waitable
	minCount   int
	count      int
	reachedMin chan struct{}
}

func newPbWaiter(w group_waiter.Waitable) *pbWaiter {
	return &pbWaiter{
		realWait:   w,
		reachedMin: make(chan struct{}),
	}
}

func (w *pbWaiter) SetMinCount(min int) {
	w.minCount = min
	w.signal()
}

func (w *pbWaiter) ReportBarCreation() {
	w.count++
	w.signal()
}

func (w *pbWaiter) signal() {
	if w.count >= w.minCount && w.reachedMin != nil {
		close(w.reachedMin)
		w.reachedMin = nil
	}
}

func (w *pbWaiter) Wait() {
	var timeoutVerifier atomic.Bool
	go func() {
		// prevent the program from getting stucked if the progress bar does not reach expected number of bars
		const timeout = 1 * time.Minute
		time.Sleep(timeout)
		if !timeoutVerifier.Load() {
			log.Panicf("progress bar was not initialized within %v, quitting.", timeout)
		}
	}()

	<-w.reachedMin
	timeoutVerifier.Store(true)
	w.realWait.Wait()
}
