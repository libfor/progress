package progress

import (
	"sync"
	"time"
)

// WaitGroup is a dropin replacement for sync.WaitGroup
// that provides progress tracking.
type WaitGroup struct {
	wg *sync.WaitGroup
	basicTask
}

func (wg *WaitGroup) init() {
	if wg.wg == nil {
		wg.wg = &sync.WaitGroup{}
		wg.basicTask = NewLongRunningJob()
	}
}

func (wg *WaitGroup) Add(delta int) {
	wg.init()
	if delta > 0 {
		wg.AddWork(uint64(delta))
	} else {
		wg.FinishWork(uint64(-delta))
	}
	wg.wg.Add(delta)
}
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}
func (wg *WaitGroup) Wait() {
	wg.init()
	wg.wg.Wait()
}

func Logger(to func(format string, v ...any), prefix string, r Reader) {
	sugar := Extend(r)
	to("%s  0%%", prefix)
	for sugar.InProgress() {
		to(`%s%3d%% (%.2f/s, %s remaining)`, prefix, int(100*sugar.Percentage()), sugar.PerSecond(), sugar.Remaining())
	}
	to("%s100%% (%.2f/s, %s remaining)", prefix, sugar.PerSecond(), sugar.Remaining())
}

// countPoller uses polling to create a Reader.
type countPoller struct {
	l        *sync.RWMutex
	c        chan struct{}
	finished *bool
	fin, tot *uint64
}

func (c countPoller) Count() (uint64, uint64) {
	c.l.RLock()
	defer c.l.RUnlock()

	if *c.finished {
		return *c.tot + 1, *c.tot + 1
	}
	return *c.fin, *c.tot + 1
}

func (c countPoller) DoneChan() (chan struct{}, bool) {
	c.l.RLock()
	defer c.l.RUnlock()

	if *c.finished {
		return c.c, true
	}
	return c.c, false
}

var _ Reader = countPoller{}

// NewReaderFromCount returns an implementation of Reader from just a count callback.
// It will regularly call the count function to determine the latest progress count.
// Note that `count` must be safe to call at any time.
// It will be considered complete when the callback returns equal numbers.
func NewReaderFromCount(count func() (uint64, uint64)) countPoller {
	fin, tot := count()
	c := countPoller{
		l:        &sync.RWMutex{},
		c:        make(chan struct{}),
		finished: new(bool),
		fin:      &fin,
		tot:      &tot,
	}
	if fin == tot {
		*c.finished = true
		close(c.c)
		return c
	}
	go func() {
		defer close(c.c)
		for {
			fin, tot := count()
			func() {
				c.l.Lock()
				defer c.l.Unlock()
				if *c.fin < fin {
					*c.fin = fin
				}
				if *c.tot < tot {
					*c.tot = tot
				}
				if *c.fin == *c.tot {
					*c.finished = true
				}
			}()
			if fin == tot {
				return
			}
			time.Sleep(30 * time.Millisecond)
		}
	}()
	return c
}
