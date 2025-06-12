package svl

import "testing"

func BenchmarkListA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := makeList(1200)
		sumList(l)
	}
}

func BenchmarkListB(b *testing.B) {
	l := makeList(1200)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sumList(l)
	}
}

func BenchmarkSliceA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := makeSlice(1200)
		sumSlice(l)
	}
}

func BenchmarkSliceB(b *testing.B) {
	l := makeSlice(1200)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sumSlice(l)
	}
}
