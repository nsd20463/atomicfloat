/*
  atomic operations on float64

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

func CompareAndSwapFloat64(p *float64, prev, next float64) bool {
	previ := math.Float64bits(prev) // note that doing the compare in uint64 conveniently makes NaN==NaN, which in this case is exactly what we want
	nexti := math.Float64bits(next) // and we ignore (or properly propagate) the many different values of NaN (all 2*(2^52-1) of them).
	pi := (*uint64)(unsafe.Pointer(p))
	return atomic.CompareAndSwapUint64(pi, previ, nexti)
}
