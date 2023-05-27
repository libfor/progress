package progress

import "sync"

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
