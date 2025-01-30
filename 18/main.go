package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	input, err := os.Open("in.text")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	output, err := os.Create("out.txt")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	bytes, err := io.Copy(output, input)

	fmt.Printf("copied %d bytes", bytes)
}
