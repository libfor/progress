package progress_test

import (
	"math"
	"sync/atomic"
	"testing"

	"github.com/libfor/progress"
)

// TestPolling tests the polling wrapper.
// It's important to note that this test relies on the debounce mechanism
// to ensure that enough time has passed for the polling to catch new values.
func TestPolling(t *testing.T) {
	t.Parallel()

	totalWork := atomic.Uint64{}
	currentWork := atomic.Uint64{}

	totalWork.Store(20)

	counter := func() (uint64, uint64) {
		return currentWork.Load(), totalWork.Load()
	}

	tracker := progress.Extend(progress.NewReaderFromCount(counter))
	if !tracker.InProgress() {
		t.Fatalf("should not be already done")
	}
	if p := tracker.Percentage(); p != 0 {
		t.Fatalf("expected 0 percentage, got %f", p)
	}

	currentWork.Add(10)
	if !tracker.InProgress() {
		t.Fatalf("should not be already done")
	}
	if p := tracker.Percentage(); math.Abs(p-.5) > .1 {
		t.Fatalf("expected ~.5 percentage, got %f", p)
	}

	currentWork.Add(10)
	if tracker.InProgress() {
		t.Fatalf("should be done")
	}
	if p := tracker.Percentage(); p != 1 {
		t.Fatalf("expected 1 percentage, got %f", p)
	}
}
