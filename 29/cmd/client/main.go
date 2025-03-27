package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type sku struct {
	item, price string
}

var items = []sku{
	{"shoes", "46"},
	{"socks", "6"},
	{"sandals", "27"},
	{"clogs", "36"},
	{"pants", "30"},
	{"shorts", "20"},
}

func main() {
	go runCreates()
	go runDeletes()
	go runUpdates()

	time.Sleep(5 * time.Second)
}

func query(cmd, params string) error {
	resp, err := http.Get("http://localhost:8080/" + cmd + "?" + params)
	if err != nil {
		log.Printf("err %s = %v\n", params, err)
		return nil
	}
	defer resp.Body.Close()

	log.Printf("got %s = %d (no err)\n", params, resp.StatusCode)
	return nil
}

func runCreates() {
	for {
		for _, s := range items {
			if err := query("create", fmt.Sprintf("item=%s&price=%s", s.item, s.price)); err != nil {
				return
			}
		}
	}
}

func runUpdates() {
	for {
		for _, s := range items {
			if err := query("update", fmt.Sprintf("item=%s&price=%s", s.item, s.price)); err != nil {
				return
			}
		}
	}
}

func runDeletes() {
	for {
		for _, s := range items {
			if err := query("delete", fmt.Sprintf("item=%s", s.item)); err != nil {
				return
			}
		}
	}
}
