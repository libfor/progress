package progress_test

import (
	"testing"

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
