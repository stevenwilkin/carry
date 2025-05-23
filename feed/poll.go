package feed

import (
	"time"
)

func Poll[T any](f func() (T, error)) func() chan T {
	return func() chan T {
		ch := make(chan T)
		t := time.NewTicker(1 * time.Second)

		go func() {
			for {
				result, err := f()
				if err != nil {
					close(ch)
					return
				}

				ch <- result
				<-t.C
			}
		}()

		return ch
	}
}
