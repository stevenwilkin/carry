package feed

import (
	"testing"
	"time"
)

func TestFeedFailing(t *testing.T) {
	f := &feed[any]{}

	if f.failing() {
		t.Fatal("Feed should not be failing")
	}

	f.errors = 1

	if !f.failing() {
		t.Fatal("Feed should be failing")
	}
}

func TestFeedFailed(t *testing.T) {
	f := &feed[any]{}

	if f.failed() {
		t.Fatal("Feed should not be failed")
	}

	f.errors = maxRetries + 1

	if !f.failed() {
		t.Fatal("Feed should be failed")
	}
}

func TestHandleInputOutput(t *testing.T) {
	var result int

	in := func() chan int {
		ch := make(chan int)
		go func() {
			ch <- 23
		}()
		return ch
	}

	out := func(x int) {
		result = x
	}

	f := NewFeed(in, out)

	f.handle()
	time.Sleep(time.Millisecond)

	if result != 23 {
		t.Fatalf("Expected: %d, got: %d", 23, result)
	}
}

func TestUpdateClearsErrorCount(t *testing.T) {
	in := func() chan int {
		ch := make(chan int)
		go func() {
			ch <- 23
		}()
		return ch
	}

	f := &feed[int]{
		inputF:  in,
		outputF: func(int) {},
		errors:  1}

	f.handle()
	time.Sleep(time.Millisecond)

	if f.errors != 0 {
		t.Fatal("Error count should be 0")
	}
}

func TestClosingChannel(t *testing.T) {
	in := func() chan int {
		ch := make(chan int)
		close(ch)
		return ch
	}

	f := NewFeed(in, func(int) {})

	f.handle()
	time.Sleep(time.Millisecond)

	if !f.failing() {
		t.Fatal("Feed should be failing")
	}
}

func TestClosingChannelRestart(t *testing.T) {
	delayBase = 0

	count := 0
	in := func() chan int {
		count += 1
		ch := make(chan int)
		close(ch)
		return ch
	}

	f := NewFeed(in, func(int) {})

	f.handle()
	time.Sleep(time.Millisecond)

	if count != maxRetries+1 {
		t.Fatalf("Expected: %d, got: %d", maxRetries+1, count)
	}

	if !f.failed() {
		t.Fatal("Feed should have failed")
	}
}
