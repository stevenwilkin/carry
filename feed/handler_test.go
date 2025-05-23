package feed

import (
	"testing"
	"time"
)

func testFeed() *feed[int] {
	return &feed[int]{
		inputF:  func() chan int { return make(chan int) },
		outputF: func(int) {}}
}

func failingFeed() *feed[int] {
	f := testFeed()
	f.errors = 1
	return f
}

func failedFeed() *feed[int] {
	f := testFeed()
	f.errors = maxRetries + 1
	return f
}

func TestAddfeed(t *testing.T) {
	h := NewHandler()

	f1 := testFeed()
	f2 := testFeed()

	h.Add(f1)
	h.Add(f2)

	if len(h.feeds) != 2 {
		t.Fatal("Incorrect number of feeds")
	}

	for i, f := range []*feed[int]{f1, f2} {
		if h.feeds[i] != f {
			t.Fatal("Missing expected feed")
		}
	}
}

func TestFailing(t *testing.T) {
	h := NewHandler()
	h.Add(testFeed())

	if h.Failing() {
		t.Fatal("Handler should not be failing")
	}

	h.Add(failingFeed())

	if !h.Failing() {
		t.Fatal("Handler should be failing")
	}
}

func TestFailed(t *testing.T) {
	h := NewHandler()
	h.Add(testFeed())
	h.Add(failingFeed())

	if h.Failed() {
		t.Fatal("Handler should not be failed")
	}

	h.Add(failedFeed())

	if !h.Failed() {
		t.Fatal("Handler should be failed")
	}
}

func TestActivatesfeed(t *testing.T) {
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

	h := NewHandler()
	h.Add(f)
	time.Sleep(time.Millisecond)

	if result != 23 {
		t.Fatalf("Expected: %d, got: %d", 23, result)
	}
}
