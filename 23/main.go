package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(url string, ch chan<- result) {
	start := time.Now()
	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get from %s, %s", url, err)
		ch <- result{
			url, err, 0,
		}
		return
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	t := time.Since(start).Round(time.Millisecond)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid body %s", err)
		os.Exit(-1)
	}

	ch <- result{
		url, nil, t,
	}
}

func main() {
	results := make(chan result)
	list := []string{
		"https://www.amazon.com",
		"https://www.google.com",
		"https://www.nytimes.com",
		"https://www.wsj.com",
	}

	for _, url := range list {
		go get(url, results)
	}

	for range list {
		r := <-results
		if r.err != nil {
			log.Printf("%-20s %s\n", r.url, r.err)
			continue
		}

		log.Printf("%-20s %s\n", r.url, r.latency)
	}
}
