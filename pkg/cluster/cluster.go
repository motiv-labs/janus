package cluster

import (
	"time"
)

func Watch(d time.Duration, fn func()) {
	t := time.NewTicker(d)

	go func(refreshTicker *time.Ticker) {
		defer refreshTicker.Stop()
		for {
			select {
			case <-refreshTicker.C:
				fn()
			}
		}
	}(t)
}
