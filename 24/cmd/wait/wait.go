package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received, waiting...")
	time.Sleep(time.Second * 7)
	fmt.Fprintf(w, "<p>Hello foo</p>")
	fmt.Println("Response sent!")
}

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
