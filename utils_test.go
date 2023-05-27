package progress_test

import (
	"math"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/libfor/progress"
)

func TestWaitGroup(t *testing.T) {
	t.Parallel()

	wg := &progress.WaitGroup{}

	wg.Add(2)
	wg.Add(5)
	wg.Done()
	wg.Done()

	if fin, tot := wg.Count(); fin != 2 || tot != 8 {
		t.Fatalf(`expected fin %d = 2 and total %d = 8`, fin, tot)
	}

	go func() {
		wg.Done()
		wg.Done()
		wg.Done()
		wg.Done()
		wg.Done()
	}()

	wg.Wait()

	if wg.InProgress() {
		t.Fatalf("expected completion because waitgroup has completed")
	}
}

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

func TestLogger(t *testing.T) {
	t.Parallel()

	task := progress.NewLongRunningJob()
	task.AddWork(50)
	go progress.Logger(t.Logf, "building thing: ", task)
	for task.InProgress() {
		task.FinishWork(uint64(rand.Intn(10)))
		time.Sleep(10 * time.Millisecond)
	}
}
