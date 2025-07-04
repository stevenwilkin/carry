package dispatcher

import (
	"errors"
	"testing"
	"time"
)

func TestRemaining(t *testing.T) {
	f := func(int) error {
		return nil
	}

	d := &Dispatcher{cb: f}
	d.Add(10)

	if d.Remaining() != 10 {
		t.Fatal("Should have 10 remaining")
	}

	d.Run()
	d.Wait()

	if d.Remaining() != 0 {
		t.Fatal("Should have 0 remaining")
	}
}

func TestFailingCallback(t *testing.T) {
	err := errors.New("fu")
	f := func(int) error {
		return err
	}

	d := NewDispatcher(f)
	d.Add(10)
	time.Sleep(time.Millisecond)

	if d.Remaining() != 10 {
		t.Fatal("Should have 10 remaining")
	}

	err = nil
	d.Wait()

	if d.Remaining() != 0 {
		t.Fatal("Should have 0 remaining")
	}
}

func TestNotMultiplesOfTen(t *testing.T) {
	f := func(int) error {
		return nil
	}

	d := NewDispatcher(f)
	d.Add(5)
	time.Sleep(20 * time.Millisecond)

	if d.Remaining() != 5 {
		t.Fatal("Should have 5 remaining")
	}

	d.Add(12)
	time.Sleep(20 * time.Millisecond)

	if d.Remaining() != 7 {
		t.Fatal("Should have 7 remaining")
	}

	d.Add(13)
	d.Wait()

	if d.Remaining() != 0 {
		t.Fatal("Should have 0 remaining")
	}
}

func TestCallCallback(t *testing.T) {
	var total, calls int
	f := func(x int) error {
		total += x
		calls += 1
		return nil
	}

	d := NewDispatcher(f)

	d.Add(10)
	d.Add(10)
	time.Sleep(20 * time.Millisecond)
	d.Add(10)
	d.Add(10)
	d.Wait()

	if total != 40 {
		t.Fatal("Should dispatch a total of 40")
	}

	if calls != 2 {
		t.Log(calls)
		t.Fatal("Should make 2 calls")
	}
}
