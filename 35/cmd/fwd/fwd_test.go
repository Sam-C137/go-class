package fwd

import (
	"math/rand"
	"testing"
)

func BenchmarkDirect(b *testing.B) {
	r := randString(rand.Intn(24), defaultChars)
	h := make([]int, b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h[i] = length(r)
	}
}

func BenchmarkForward(b *testing.B) {
	r := randString(rand.Intn(24), defaultChars)
	h := make([]int, b.N)

	var t3 forwarder = &thing3{}
	var t2 forwarder = &thing2{t3}
	var t1 forwarder = &thing1{t2}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h[i] = t1.forward(r)
	}
}
