// Package serial generates unique serial numbers as 64 bit integers.
// It supports maintaining a blacklist to prevent their reuse and
// ensures thread safety. Generated numbers are based on nanosecond timestamps
// and so are most definitely not cryptographically random.
package serial

import (
	"sync"
	"time"
)

// Serial is a unique serial number.
type Serial int64

// Generator defines a generator of unique serial numbers. You can run any
// number of independent generators for different serial number problem
// domains, each with its own mutexes for thread safety.
type Generator struct {
	lastmutex  sync.RWMutex
	lastSerial Serial
	seenmutex  sync.RWMutex
	seen       map[Serial]struct{}
}

// NewGenerator creates and initializes a new serial number generator.
func NewGenerator() *Generator {
	gen := &Generator{}
	gen.seenmutex.Lock()
	gen.seen = make(map[Serial]struct{})
	gen.seenmutex.Unlock()
	return gen
}

// Seen returns a boolean to indicate whether the specified Serial value has
// been seen. Serial values are unseen until SetSeen is called. Once they have
// been set as seen, they remain seen until history is expired.
func (g *Generator) Seen(x Serial) bool {
	g.seenmutex.RLock()
	_, ok := g.seen[x]
	g.seenmutex.RUnlock()
	return ok
}

// SetSeen flags the specified Serial value as having been seen. This can
// then be interrogated using the Seen() method.
func (g *Generator) SetSeen(x Serial) {
	g.seenmutex.Lock()
	g.seen[x] = struct{}{}
	g.seenmutex.Unlock()
}

// ExpireSeen clears the history of seen Serial values, using an age limit
// provided as a time.Duration. All history data older than the specified
// duration is deleted.
//
// This function should be called periodically if you are using the Seen flag
// feature, or else eventually your memory will fill up.
func (g *Generator) ExpireSeen(agelimit time.Duration) {
	g.seenmutex.Lock()
	limit := time.Now().Add(-agelimit).UnixNano()
	for tok := range g.seen {
		if int64(tok) < limit {
			delete(g.seen, tok)
		}
	}
	g.seenmutex.Unlock()
}

// Generate generates a serial value based on Unix time in nanoseconds.
// You are guaranteed to get a different value each time you call the function.
// The value will be no earlier than the current Unix epoch time in nanoseconds.
func (g *Generator) Generate() Serial {
	g.lastmutex.Lock()
	id := Serial(time.Now().UnixNano())
	for id <= g.lastSerial {
		id = id + 1
	}
	g.lastSerial = id
	g.lastmutex.Unlock()
	return id
}
