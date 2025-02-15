package main

import "fmt"

func generate(limit int, ch chan<- int) {
	for n := range limit {
		if n < 2 {
			continue
		}
		ch <- n
	}

	close(ch)
}

func filter(src <-chan int, dst chan<- int, prime int) {
	for num := range src {
		if num%prime != 0 {
			dst <- num
		}
	}

	close(dst)
}

func sieve(limit int) {
	ch := make(chan int)
	go generate(limit, ch)

	for {
		prime, ok := <-ch

		if !ok {
			break
		}

		ch2 := make(chan int)
		go filter(ch, ch2, prime)
		ch = ch2
		fmt.Print(prime, " ")
	}
}

func main() {
	sieve(100)
}
