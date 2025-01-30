package main

import (
	"fmt"
	"io"
	"os"
)

// to run the program, `go run . *.txt > d.txt`

func main() {
	for _, fname := range os.Args[1:] {
		file, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if _, err := io.Copy(os.Stdout, file); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		file.Close()
	}
}
