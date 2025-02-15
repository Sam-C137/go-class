package main

import (
	"fmt"
	"log"
	"net/http"
)

type countChan chan int

func (ch countChan) handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Count is %d</h1>", <-ch)
}

func (ch countChan) count() {
	for i := 0; ; i++ {
		ch <- i
	}
}

func main() {
	ch := make(countChan)
	go ch.count()

	http.HandleFunc("/", ch.handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
