package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database struct {
	mu sync.Mutex
	db map[string]dollars
}

func (d *database) list(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for item, price := range d.db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (d *database) add(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")
	d.mu.Lock()
	defer d.mu.Unlock()

	if item != "" && price != "" {
		if _, ok := d.db[item]; ok {
			http.Error(w, fmt.Sprintf("%s already exists in the record, skipping...\n", item), http.StatusBadRequest)
			return
		}

		p, err := strconv.ParseFloat(price, 32)
		if err != nil {
			http.Error(w, fmt.Sprintf("%q is not a valid digit\n", price), http.StatusBadRequest)
			return
		}

		d.db[item] = dollars(p)
		fmt.Fprintf(w, "Added %s with price %f to the record\n", item, p)
		return
	}

	http.Error(w, "Invalid request, item or price is missing", http.StatusBadRequest)
}

func (d *database) update(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")
	d.mu.Lock()
	defer d.mu.Unlock()

	if item != "" && price != "" {
		if _, ok := d.db[item]; !ok {
			http.Error(w, fmt.Sprintf("Cannot update non-existent item: %s\n", item), http.StatusNotFound)
			return
		}

		p, err := strconv.ParseFloat(price, 32)
		if err != nil {
			http.Error(w, fmt.Sprintf("%q is not a valid digit\n", price), http.StatusBadRequest)
			return
		}

		d.db[item] = dollars(p)
		fmt.Fprintf(w, "Updated %s with price %f in the record\n", item, p)
		return
	}
	http.Error(w, "Invalid request, item or price is missing", http.StatusBadRequest)
}

func (d *database) fetch(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	d.mu.Lock()
	defer d.mu.Unlock()

	if item != "" {
		price, ok := d.db[item]

		if !ok {
			http.Error(w, fmt.Sprintf("Cannot fetch non-existent item: %s\n", item), http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "Item: %s, Price: %s\n", item, price)
		return
	}

	http.Error(w, "Invalid request, item name is missing", http.StatusBadRequest)
}

func (d *database) delete(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	d.mu.Lock()
	defer d.mu.Unlock()

	if item != "" {
		if _, ok := d.db[item]; !ok {
			http.Error(w, fmt.Sprintf("Cannot delete non-existent item: %s\n", item), http.StatusNotFound)
			return
		}

		delete(d.db, item)
		fmt.Fprintf(w, "Deleted %s from the record\n", item)
		return
	}

	http.Error(w, "Invalid request, item name is missing", http.StatusBadRequest)
}

func main() {
	d := database{
		db: map[string]dollars{
			"shoes": 50,
			"socks": 5,
		},
	}

	// add some routes
	http.HandleFunc("/list", d.list)
	http.HandleFunc("/create", d.add)
	http.HandleFunc("/update", d.update)
	http.HandleFunc("/read", d.fetch)
	http.HandleFunc("/delete", d.delete)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
