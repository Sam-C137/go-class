package main

import (
	"log"
	"time"
)

func main() {
	chans := []chan int{
		make(chan int),
		make(chan int),
	}

	for i := range chans {
		go func(idx int, ch chan<- int) {
			for {
				time.Sleep(time.Second * time.Duration(idx))
				ch <- idx
			}
		}(i+1, chans[i])
	}

	for i := 0; i < 12; i++ {
		select {
		case v1 := <-chans[0]:
			log.Println("received ", v1)
		case v2 := <-chans[1]:
			log.Println("received ", v2)
		}
	}
}
