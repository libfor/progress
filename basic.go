package progress

import (
	"math"
	"sync"
)

type LongRunningJob struct {
	l *sync.RWMutex

	total    *uint64
	finished *uint64

	done   chan struct{}
	frozen *bool

	ReaderWrapper
}

// NewLongRunningJob returns a basic task tracker.
// All methods are safe to call concurrently.
// It already uses Extend to provide utility methods.
func NewLongRunningJob() LongRunningJob {
	t := LongRunningJob{
		l:        new(sync.RWMutex),
		total:    new(uint64),
		finished: new(uint64),
		frozen:   new(bool),
		done:     make(chan struct{}),
	}
	t.ReaderWrapper = Extend(t)
	return t
}

// Count returns the finished work and total work.
// Note that the task itself counts as 1 piece of work.
func (t LongRunningJob) Count() (uint64, uint64) {
	t.l.RLock()
	defer t.l.RUnlock()

	total := *t.total
	finished := *t.finished

	if total < finished {
		total = finished
	}

	if *t.frozen {
		return total + 1, total + 1
	}

	return finished, total + 1
}

func (t LongRunningJob) Complete() {
	t.l.Lock()
	defer t.l.Unlock()
	if *t.frozen {
		return
	}
	*t.frozen = true
	close(t.done)
}

func (t LongRunningJob) DoneChan() (chan struct{}, bool) {
	t.l.RLock()
	defer t.l.RUnlock()
	return t.done, *t.frozen
}

const maxTasks = math.MaxUint64 - 1

// AddWork adds to the total work.
func (t LongRunningJob) AddWork(items uint64) {
	if items == 0 {
		return
	}
	t.l.Lock()
	defer t.l.Unlock()
	if *t.frozen {
		return
	}
	if maxTasks-*t.total < items {
		*t.total = maxTasks
		return
	}
	*t.total = *t.total + items
}

// FinishWork progresses the job. When the amount of finished work
// equals or exceed the total work, the job is marked as completed.
func (t LongRunningJob) FinishWork(items uint64) {
	if items == 0 {
		return
	}
	t.l.Lock()
	defer t.l.Unlock()
	if *t.frozen {
		return
	}
	if maxTasks-*t.finished < items {
		*t.finished = maxTasks
		return
	}
	*t.finished = *t.finished + items
	if *t.finished >= *t.total {
		*t.frozen = true
		close(t.done)
	}
}
