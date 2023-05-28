package progress

import (
	"sync"
	"time"
)

type ReaderWrapper struct {
	Reader

	lock      *sync.RWMutex
	startTime *time.Time
}

// Extend adds some nice utility methods to a Reader.
// Reader implementations can just embed this.
func Extend(t Reader) ReaderWrapper {
	if already, wrapped := t.(ReaderWrapper); wrapped {
		return already
	}
	now := time.Now()
	return ReaderWrapper{Reader: t,
		lock:      new(sync.RWMutex),
		startTime: &now,
	}
}

// PerSecond returns the throughput since this tracker began.
func (s ReaderWrapper) PerSecond() float64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	fin, _ := s.Count()
	return float64(fin) / time.Since(*s.startTime).Seconds()
}

// Remaining estimates the time remaining until completion.
func (s ReaderWrapper) Remaining() time.Duration {
	perSecond := s.PerSecond()
	fin, tot := s.Count()
	if perSecond == 0 {
		return time.Duration(0)
	}
	if !s.InProgress() {
		return 0
	}
	seconds := float64(tot-fin) / perSecond
	return time.Duration(seconds * float64(time.Second))
}

// Percentage returns the current progress as a percentage between 0 and 1.
func (s ReaderWrapper) Percentage() float64 {
	fin, tot := s.Count()
	return float64(fin) / float64(tot)
}

// InProgress will return false when the tracker is complete.
// It is debounced to try and return false for 200ms.
// After 200ms it will finally return true.
func (s ReaderWrapper) InProgress() bool {
	ch, closed := s.DoneChan()
	if closed {
		return false
	}
	t := time.NewTimer(200 * time.Millisecond)
	select {
	case <-t.C:
		return true
	case <-ch:
		return false
	}
}
