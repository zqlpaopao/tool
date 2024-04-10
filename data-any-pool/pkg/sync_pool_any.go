package pkg

// Package syncpool provides a generic wrapper around sync.Pool

import (
	"sync"
)

// A SyncPool is a generic wrapper around a sync.Pool.
type SyncPool[T any] struct {
	pool sync.Pool
}

// New creates a new Pool with the provided new function.
//
// The equivalent sync.Pool construct is "sync.Pool{New: fn}"
func New[T any](fn func() T) SyncPool[T] {
	return SyncPool[T]{
		pool: sync.Pool{New: func() interface{} { return fn() }},
	}
}

// Get is a generic wrapper around sync.Pool's Get method.
func (p *SyncPool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put is a generic wrapper around sync.Pool's Put method.
func (p *SyncPool[T]) Put(x T) {
	p.pool.Put(x)
}
