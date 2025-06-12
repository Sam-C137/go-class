package share

import "testing"

func BenchmarkShare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run()
	}
}
