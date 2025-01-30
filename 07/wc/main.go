package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var tlc, twc, tcc int
	for _, fname := range os.Args[1:] {
		var lc, wc, cc int

		file, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		scan := bufio.NewScanner(file)

		for scan.Scan() {
			s := scan.Text()

			wc += len(strings.Fields(s))
			cc += len(s)
			lc++
		}

		tlc, twc, tcc = tlc+lc, twc+wc, tcc+cc
		fmt.Printf("lines: %-7d words: %-7d characters: %-7d %s\n", lc, wc, cc, fname)
		file.Close()
	}

	if len(os.Args) > 2 {
		fmt.Printf("total lines: %-7d total words: %-7d total characters: %-7d\n", tlc, twc, tcc)
	}
}
