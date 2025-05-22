package app

import (
	"time"
)

func HandlePullRoom(f func() error, delay time.Duration) func() <-chan error {
	type req struct {
		resp chan error
	}

	requests := make(chan req)
	go func() {
		var (
			timer   *time.Timer
			pending *req
		)
		for {
			select {
			case r := <-requests:
				if timer != nil {
					if !timer.Stop() {
						<-timer.C
					}
				}
				r.resp <- nil
				timer = time.AfterFunc(delay, func() {
					requests <- req{resp: make(chan error)}
				})
				pending = &r
				timer = time.NewTimer(delay)
			case <-func() <-chan time.Time {
				if timer != nil {
					return timer.C
				}
				return make(chan time.Time)
			}():
				if pending != nil {
					r := pending
					pending = nil
					timer = nil
					r.resp <- f()
				}
				timer = nil
			}
		}
	}()
	return func() <-chan error {
		r := make(chan error)
		requests <- req{resp: r}
		return r
	}
}
