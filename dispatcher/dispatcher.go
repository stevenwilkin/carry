package dispatcher

import (
	"sync"
	"time"
)

type Dispatcher struct {
	remaining int
	cb        func(int) error
	m         sync.Mutex
}

func (d *Dispatcher) Add(quantity int) {
	d.m.Lock()
	defer d.m.Unlock()

	d.remaining += quantity
}

func (d *Dispatcher) Remaining() int {
	return d.remaining
}

func (d *Dispatcher) Run() {
	go func() {
		t := time.NewTicker(10 * time.Millisecond)

		for {
			quantity := (d.remaining / 10) * 10

			if quantity == 0 {
				<-t.C
				continue
			}

			if err := d.cb(quantity); err != nil {
				continue
			}

			d.m.Lock()
			d.remaining -= quantity
			d.m.Unlock()
		}
	}()
}

func (d *Dispatcher) Wait() {
	t := time.NewTicker(10 * time.Millisecond)

	for {
		if d.remaining == 0 {
			return
		}

		<-t.C
	}
}

func NewDispatcher(cb func(int) error) *Dispatcher {
	d := &Dispatcher{cb: cb}
	d.Run()

	return d
}
