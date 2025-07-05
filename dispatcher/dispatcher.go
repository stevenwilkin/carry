package dispatcher

import (
	"sync"
	"time"

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
	d.c.Signal()
}

func (d *Dispatcher) Remaining() int {
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
	d.c = sync.NewCond(&d.m)
	d.Run()

	return d
}
