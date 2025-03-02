package main

import (
	"context"
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

func get(ctx context.Context, url string, ch chan<- result) {
	start := time.Now()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := http.DefaultClient.Do(req)
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
	results := make(chan result)
	sites := []string{
		"https://www.amazon.com",
		"https://www.google.com",
		"https://www.nytimes.com",
		"https://www.wsj.com",
		"http://localhost:8080/",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := range sites {
		go get(ctx, sites[i], results)
	}

	for range sites {
		r := <-results
		if r.err != nil {
			log.Printf("%-20s %s\n", r.url, r.err)
			continue
		}
		log.Printf("%-20s %s\n", r.url, r.latency)
	}
}
