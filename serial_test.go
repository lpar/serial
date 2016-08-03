package serial

import (
	"testing"
	"time"
)

func TestSerial(t *testing.T) {
	for i := 0; i < 100; i++ {
		n1 := Generate()
		n2 := Generate()
		if n1 == n2 {
			t.Error("Got the same value twice!")
		}
	}
}

func TestOneTime(t *testing.T) {
	n1 := Generate()
	SetSeen(n1)
	if !Seen(n1) {
		t.Error("Flagged value as seen, got 'not seen'")
	}
	n2 := Generate()
	if Seen(n2) {
		t.Error("Got 'seen' for unflagged value")
	}
	ExpireSeen(time.Duration(0))
	if Seen(n1) {
		t.Error("Emptied history but value was still 'seen'")
	}
}

func TestGC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping extended history test in short mode")
		return
	}
	vals := make([]Serial, 100)
	for i := 0; i < 100; i++ {
		v := Generate()
		vals = append(vals, v)
		SetSeen(v)
		time.Sleep(time.Second / 10)
	}
	before := len(history.seen)
	if before != 100 {
		t.Errorf("History wrong length, expected 100 got %d", before)
	}
	// 5050 = 5 seconds plus a little slop to make sure we don't occasionally
	// fail for no good reason
	ExpireSeen(time.Millisecond * 5050)
	after := len(history.seen)
	if after != 50 {
		t.Errorf("History wrong length after expire, expected 50 got %d", after)
	}
	count := 0
	for _, v := range vals {
		if Seen(v) {
			count++
		}
	}
	if count != len(history.seen) {
		t.Errorf("History had wrong number of values expected %d got %d", count, after)
	}
}
