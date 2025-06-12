package fib

func Fib(n int, recursive bool) int {
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	default:
		if recursive {
			return Fib(n-1, true) + Fib(n-2, true)
		}

		a, b := 0, 1
		for i := 1; i < n; i++ {
			a, b = b, a+b
		}

		return b
	}
}
