package progress_test

import (
	"math"
	"testing"

	"github.com/libfor/progress"
)

func TestBasicTask(t *testing.T) {
	t.Parallel()

	track := progress.NewBasicTask()
	ch, fin := track.DoneChan()
	if fin {
		t.Fatalf("expected fresh task to be in progress")
	}

	t.Logf("adding more work than has been finished")
	track.AddWork(15)
	if fin, tot := track.Count(); fin != 0 || tot != 16 {
		t.Fatalf("expected finished %d == 0 and total %d == 16", fin, tot)
	}

	t.Logf("finishing more work than we have added")
	track.FinishWork(10)
	if fin, tot := track.Count(); fin != 10 || tot != 16 {
		t.Fatalf("expected finished %d == 10 and total %d == 16", fin, tot)
	}

	t.Logf("completion returns a full count")
	track.Complete()
	if fin, tot := track.Count(); fin != 16 || tot != 16 {
		t.Fatalf("expected finished %d == total %d == 16", fin, tot)
	}

	t.Log("consuming from channel that should be closed")
	if _, ok := <-ch; ok {
		t.Fatalf("channel was somehow open")
	}
}

func TestEdgeCases(t *testing.T) {
	t.Parallel()

	{
		track := progress.NewBasicTask()
		track.AddWork(15)
		track.FinishWork(15)
		if track.InProgress() {
			t.Fatalf("expected tracker to be finished")
		}
	}

	{
		track := progress.NewBasicTask()
		track.AddWork(100)
		track.AddWork(math.MaxUint64)
		if fin, tot := track.Count(); fin != 0 || tot != math.MaxUint64 {
			t.Fatalf("expected tracker to prevent overflow")
		}
	}
}
