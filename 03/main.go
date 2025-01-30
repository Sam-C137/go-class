package main

import (
	"fmt"
	"os"
)

func main() {
	var sum float64
	var n int

	for {
		var val float64

		if _, err := fmt.Fscanln(os.Stdin, &val); err != nil {
			break
		}

		sum += val
		n++
	}

	if n == 0 {
		// best practice to write to the standard error output as it could be redirected somewhere
		// other than the standard output
		fmt.Fprintln(os.Stderr, "no values")
		os.Exit(-1)
	}

	fmt.Println("The average is", sum/float64(n))
}
