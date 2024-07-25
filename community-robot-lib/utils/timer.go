package utils

import "time"

type Timer interface {
	Start(f func(), interval, delay time.Duration)
	Stop()
}

func NewTimer() Timer {
	return timer{
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

type timer struct {
	stop    chan struct{}
	stopped chan struct{}
}

// Start triggers the work by the interval until recieving a stop signal.
func (t timer) Start(f func(), interval, delay time.Duration) {
	go func(f func(), interval, delay time.Duration) {
		if delay > 0 {
			select {
			case <-time.After(delay):
				f()
			case <-t.stop:
				close(t.stopped)
				return
			}
		}

		// the ticker will fire after interval,
		// which means the f will run for the first time after interval.
		ticker := time.Tick(interval)

		for {
			select {
			case <-ticker:
				f()
			case <-t.stop:
				close(t.stopped)
				return
			}
		}
	}(f, interval, delay)
}

func (t timer) Stop() {
	close(t.stop)

	<-t.stopped
}
