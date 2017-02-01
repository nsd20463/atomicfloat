/*
  atomic operations on float64

  The missing functions from the standard atomic package
  (along with atomic bools. I've wanted atomic bools many more times than I've wanted atomic floats,
   but atomic bools are tricker because most CPUs can't really do atomic ops on small data, so it becomes
   an op on the larger, enclosing 32-bit word)

  Copyright 2017 Nicolas Dade
*/
package atomicfloat

import (
	"math"
	"sync/atomic"
	"unsafe"
)

func LoadFloat64(p *float64) float64 {
	pi := (*uint64)(unsafe.Pointer(p))
	i := atomic.LoadUint64(pi)
	return math.Float64frombits(i)
}

func StoreFloat64(p *float64, v float64) {
	pi := (*uint64)(unsafe.Pointer(p))
	i := math.Float64bits(v)
	atomic.StoreUint64(pi, i)
}

func AddFloat64(p *float64, v float64) {
	pi := (*uint64)(unsafe.Pointer(p))
	for {
		i := atomic.LoadUint64(pi)
		x := math.Float64frombits(i)
		x += v
		n := math.Float64bits(x)
		if atomic.CompareAndSwapUint64(pi, i, n) {
			return
		}
	}
}
