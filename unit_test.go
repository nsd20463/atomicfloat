package atomicfloat

import (
	"sync"
	"testing"
)

func TestSimple(t *testing.T) {
	// there isn't much that can go wrong, but we might as well exercise the code

	var x, y float64

	if LoadFloat64(&x) != 0 {
		t.Error("zero-val not 0.0")
	}

	StoreFloat64(&x, 1.25) // NOTE WEL these constants are chosen b/c they are not repeating fractions in binary
	StoreFloat64(&y, 3.5)

	if LoadFloat64(&x) != 1.25 || LoadFloat64(&y) != 3.5 {
		t.Error("Load(Set(x)) != x")
	}

	AddFloat64(&x, 1.75)
	if LoadFloat64(&x) != 3 {
		t.Error("1.25+1.75 != 3")
	}
}

func TestAtomic(t *testing.T) {
	var x float64
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 10000; j++ {
				AddFloat64(&x, float64(j))
			}
			wg.Done()
		}()
	}

	// do the same math serially
	var y float64
	for i := 0; i < 100; i++ {
		for j := 0; j < 10000; j++ {
			y += float64(j)
		}
	}

	wg.Wait()
	if x != y { // in general float!=float is wrong, but here all inputs are ints and all intermediately results are within the 53 bit matissa, so we should get a perfect match
		t.Errorf("100*sum(0,1,...99) %v != %v", y, x)
	}
}

// just to check that it does indeed fail if we aren't atomic
func TestNonAtomic(t *testing.T) {
	var x float64
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 10000; j++ {
				nonatomicAddFloat64(&x, float64(j))
			}
			wg.Done()
		}()
	}

	// do the same math serially
	var y float64
	for i := 0; i < 100; i++ {
		for j := 0; j < 10000; j++ {
			y += float64(j)
		}
	}

	wg.Wait()
	t.Logf("concurrent non-atomic adds results in a mess %v != %v\n", x, y)
}

func BenchmarkAtomicAdd(b *testing.B) {
	var x float64
	for i := 0; i < b.N; i++ {
		AddFloat64(&x, float64(i))
	}
}

// out of curiosity, see what the CAS() and the non-pipelining/unrolling costs us
func BenchmarkRegularAdd(b *testing.B) {
	var x float64
	for i := 0; i < b.N; i++ {
		nonatomicAddFloat64(&x, float64(i))
	}
}

func nonatomicAddFloat64(x *float64, y float64) {
	*x += y
}

func TestCAS(t *testing.T) {
	var x float64
	if CompareAndSwapFloat64(&x, 1, 2) {
		t.Error("should have failed")
	}
	if !CompareAndSwapFloat64(&x, 0, 4) {
		t.Error("should have succeeded")
	}
	if x != 4 {
		t.Error("should have updated")
	}
}

func BenchmarkCAS(b *testing.B) {
	var x float64
	for i := 0; i < b.N; i++ {
		CompareAndSwapFloat64(&x, x, x+1)
	}
}
