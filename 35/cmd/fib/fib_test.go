package fib

import "testing"

func BenchmarkFib20T(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(20, true)
	}
}

func BenchmarkFib20F(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(20, false)
	}
}
