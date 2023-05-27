package progress_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/libfor/progress"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	task := progress.NewBasicTask()
	task.AddWork(50)
	go progress.Logger(t.Logf, "building thing: ", task)
	for task.InProgress() {
		task.FinishWork(uint64(rand.Intn(10)))
		time.Sleep(10 * time.Millisecond)
	}
}
