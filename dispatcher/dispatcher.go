package dispatcher

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type Dispatcher struct {
	remaining int
	cb        func(int) error
	m         sync.Mutex
	c         *sync.Cond
}

func (d *Dispatcher) Add(quantity int) {
	d.m.Lock()

	d.remaining += quantity

	d.m.Unlock()
	d.c.Broadcast()
}

func (d *Dispatcher) Remaining() int {
	d.m.Lock()
	defer d.m.Unlock()

	return d.remaining
}

func (d *Dispatcher) Run() {
	go func() {
		for {
			d.m.Lock()
			for d.remaining < 10 {
				d.c.Wait()
			}

			quantity := (d.remaining / 10) * 10
			d.m.Unlock()

			log.WithField("quantity", quantity).Debug("Market order")

			if err := d.cb(quantity); err != nil {
				log.Error(err)
				continue
			}

			d.m.Lock()
			d.remaining -= quantity
			d.m.Unlock()
			d.c.Broadcast()
		}
	}()
}

func (d *Dispatcher) Wait() {
	d.m.Lock()
	defer d.m.Unlock()

	for d.remaining != 0 {
		d.c.Wait()
	}
}

func NewDispatcher(cb func(int) error) *Dispatcher {
	d := &Dispatcher{cb: cb}
	d.c = sync.NewCond(&d.m)
	d.Run()

	return d
}
