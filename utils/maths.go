package utils

import (
	"math/rand"
	"sync"
	"time"
)

// MaxInt returns the larger one of a and b.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt returns the smaller one of a and b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

type Proba struct {
	// rand.New(...) returns a non thread safe object
	r    *rand.Rand
	lock sync.Mutex
}

// NewProba returns a Proba.
func NewProba() *Proba {
	return &Proba{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// TrueOnProba checks if true on given probability.
func (p *Proba) TrueOnProba(proba float64) (truth bool) {
	p.lock.Lock()
	truth = p.r.Float64() < proba
	p.lock.Unlock()
	return
}
