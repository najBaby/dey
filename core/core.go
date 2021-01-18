package core

import "time"

// Run is
func Run(d time.Duration) {
	runAfter(d, func() {})
}

func runAfter(d time.Duration, f func()) {
	t := time.NewTicker(d)
	f()
	for {
		select {
		case <-t.C:
			go f()
		}
	}
}
