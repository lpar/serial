// Generates unique serial numbers as 64 bit integers.
// Supports maintaining a blacklist to prevent their reuse.
// Ensures thread safety.
package serial

import (
	"sync"
	"time"
)

type Serial int64

var ratchet struct {
	sync.RWMutex
	lastSerial Serial
}

var history struct {
	sync.RWMutex
	// Map with no values, just keys
	seen map[Serial]struct{}
}

func init() {
	history.Lock()
	history.seen = make(map[Serial]struct{})
	history.Unlock()
}

// Seen returns a boolean to indicate whether the specified Serial value has
// been seen. Serial values are unseen until SetSeen is called. Once they have
// been set as seen, they remain seen until history is expired.
func Seen(x Serial) bool {
	history.RLock()
	_, ok := history.seen[x]
	history.RUnlock()
	return ok
}

// SetSeen flags the specified Serial value as having been seen. This can
// then be interrogated using the Seen() method.
func SetSeen(x Serial) {
	history.Lock()
	history.seen[x] = struct{}{}
	history.Unlock()
}

// ExpireSeen clears the history of seen Serial values, using an age limit
// provided as a time.Duration. All history data older than the specified
// duration is deleted.
//
// This function should be called periodically if you are using the Seen flag
// feature, or else eventually your memory will fill up.
func ExpireSeen(agelimit time.Duration) {
	history.Lock()
	limit := time.Now().Add(-agelimit).UnixNano()
	for tok := range history.seen {
		if int64(tok) < limit {
			delete(history.seen, tok)
		}
	}
	history.Unlock()
}

// Generate generates a serial value based on Unix time in nanoseconds.
// You are guaranteed to get a different value each time you call the function.
// The value will be no earlier than the current Unix epoch time in nanoseconds.
func Generate() Serial {
	ratchet.Lock()
	id := Serial(time.Now().UnixNano())
	for id <= ratchet.lastSerial {
		id = id + 1
	}
	ratchet.lastSerial = id
	ratchet.Unlock()
	return id
}
