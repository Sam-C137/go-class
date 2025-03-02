package main

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(ctx context.Context, url string, ch chan<- result) {
	var r result
	start := time.Now()
	ticker := time.NewTicker(1 * time.Second).C
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r = result{
			url, err, 0,
		}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		r = result{
			url, err, t,
		}
		resp.Body.Close()
	}

	for {
		select {
		case ch <- r:
			return
		case <-ticker:
			log.Println("tick ", r)
		}
	}
}

// gets the first successful request and returns the response
func first(ctx context.Context, urls []string) (*result, error) {
	results := make(chan result, len(urls)) // buffer to avoid leaking
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel no matter what to free up resources

	for _, url := range urls {
		go get(ctx, url, results)
	}

	select {
	case r := <-results:
		return &r, nil
		// handle the case of the parent contex being cancelled from above
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func main() {
	sites := []string{
		"https://www.amazon.com",
		"https://www.google.com",
		"https://www.nytimes.com",
		"https://www.wsj.com",
		"http://localhost:8080/",
	}

	r, _ := first(context.Background(), sites)
	if r.err != nil {
		log.Printf("%-20s %s\n", r.url, r.err)
	} else {
		log.Printf("%-20s %s\n", r.url, r.latency)
	}

	time.Sleep(9 * time.Second)
	log.Println("quit anyway... ", runtime.NumGoroutine(), " still running")
}
