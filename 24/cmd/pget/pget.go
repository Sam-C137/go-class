package main

import (
	"io"
	"log"
	"net/http"
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
		ch <- result{
			url, err, 0,
		}
		return
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	end := time.Since(start).Round(time.Millisecond)

	ch <- result{
		url, err, end,
	}
}

func main() {
	stopper := time.After(time.Second * 3)
	results := make(chan result)
	sites := []string{
		"https://www.amazon.com",
		"https://www.google.com",
		"https://www.nytimes.com",
		"https://www.wsj.com",
		"http://localhost:8080/",
	}

	for i := range sites {
		go get(sites[i], results)
	}

	for range sites {
		select {
		case r := <-results:
			if r.err != nil {
				log.Printf("%-20s %s\n", r.url, r.err)
				continue
			}
			log.Printf("%-20s %s\n", r.url, r.latency)
		case <-stopper:
			log.Fatalf("Timeout after 3 seconds")
		}

	}
}
