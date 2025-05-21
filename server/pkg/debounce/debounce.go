package debounce

import (
	"context"
	"sync"
	"time"
)

type TimerEntry struct {
	timer  *time.Timer
	cancel context.CancelFunc
}

type Debouncer struct {
	mu      sync.Mutex
	timers  map[string]TimerEntry
	timeout time.Duration
}

func NewDebouncer(timeout time.Duration) *Debouncer {
	return &Debouncer{
		timers:  make(map[string]TimerEntry),
		timeout: timeout,
	}
}

func (d *Debouncer) Debounce(key string, fn func(ctx context.Context)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if timerInfo, ok := d.timers[key]; ok {
		timerInfo.timer.Stop()
		timerInfo.cancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	newTimer := time.AfterFunc(d.timeout, func() {
		go func() {
			fn(ctx)
			d.mu.Lock()
			defer d.mu.Unlock()
			delete(d.timers, key)
		}()
	})
	d.timers[key] = TimerEntry{
		timer:  newTimer,
		cancel: cancel,
	}
}
