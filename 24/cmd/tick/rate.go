package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("start")
	tickRate := time.Second * 2
	ticker := time.NewTicker(tickRate).C
	stopper := time.After(tickRate * 5)

loop:
	for {
		select {
		case <-ticker:
			fmt.Println("tick")
		case <-stopper:
			break loop
		}
	}
	fmt.Println("finish")
}
