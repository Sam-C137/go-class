package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type kv struct {
	key string
	val int
}

func main() {
	scan := bufio.NewScanner(os.Stdin)
	words := make(map[string]int)
	scan.Split(bufio.ScanWords)

	for scan.Scan() {
		words[scan.Text()]++
	}

	fmt.Println(len(words), "unique words")

	var ss []kv

	for k, v := range words {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].val > ss[j].val
	})

	for _, s := range ss[:5] {
		fmt.Println(s.key, "appears", s.val, "times")
	}
}
