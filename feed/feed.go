package feed

import (
	"math"
	"time"
)

var (
	delayBase  = 2.0
	maxRetries = 6
)

type Feed interface {
	failing() bool
	failed() bool
	handle()
}

type feed[T any] struct {
	inputF  func() chan T
	outputF func(T)
	errors  int
}

func (f *feed[T]) failing() bool {
	return f.errors > 0
}

func (f *feed[T]) failed() bool {
	return f.errors > maxRetries
}

func (f *feed[T]) exponentialBackoff() {
	delaySeconds := math.Pow(delayBase, float64(f.errors))
	delay := time.Second * time.Duration(delaySeconds)
	time.Sleep(delay)
}

func (f *feed[T]) handle() {
	go func() {
		ch := f.inputF()
		for {
			item, ok := <-ch
			if ok {
				f.errors = 0
				f.outputF(item)
			} else {
				f.errors += 1
				if f.failed() {
					return
				} else {
					f.exponentialBackoff()
					ch = f.inputF()
				}
			}
		}
	}()
}

func NewFeed[T any](in func() chan T, out func(T)) Feed {
	return &feed[T]{
		inputF:  in,
		outputF: out}
}

var _ Feed = &feed[int]{}
