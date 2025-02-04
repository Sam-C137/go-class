package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, r *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) add(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")

	if item != "" && price != "" {
		if _, ok := db[item]; ok {
			http.Error(w, fmt.Sprintf("%s already exists in the record, skipping...\n", item), http.StatusBadRequest)
			return
		}

		p, err := strconv.ParseFloat(price, 32)
		if err != nil {
			http.Error(w, fmt.Sprintf("%q is not a valid digit\n", price), http.StatusBadRequest)
			return
		}

		db[item] = dollars(p)
		fmt.Fprintf(w, "Added %s with price %f to the record\n", item, p)
		return
	}

	http.Error(w, "Invalid request, item or price is missing", http.StatusBadRequest)
}

func (db database) update(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")

	if item != "" && price != "" {
		if _, ok := db[item]; !ok {
			http.Error(w, fmt.Sprintf("Cannot update non-existent item: %s\n", item), http.StatusNotFound)
			return
		}

		p, err := strconv.ParseFloat(price, 32)
		if err != nil {
			http.Error(w, fmt.Sprintf("%q is not a valid digit\n", price), http.StatusBadRequest)
			return
		}

		db[item] = dollars(p)
		fmt.Fprintf(w, "Updated %s with price %f in the record\n", item, p)
		return
	}
	http.Error(w, "Invalid request, item or price is missing", http.StatusBadRequest)
}

func (db database) fetch(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")

	if item != "" {
		price, ok := db[item]

		if !ok {
			http.Error(w, fmt.Sprintf("Cannot fetch non-existent item: %s\n", item), http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "Item: %s, Price: %s\n", item, price)
		return
	}

	http.Error(w, "Invalid request, item name is missing", http.StatusBadRequest)
}

func (db database) delete(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")

	if item != "" {
		if _, ok := db[item]; !ok {
			http.Error(w, fmt.Sprintf("Cannot delete non-existent item: %s\n", item), http.StatusNotFound)
			return
		}

		delete(db, item)
		fmt.Fprintf(w, "Deleted %s from the record\n", item)
		return
	}

	http.Error(w, "Invalid request, item name is missing", http.StatusBadRequest)
}

func main() {
	db := database{
		"shoes": 50,
		"socks": 5,
	}

	// add some routes
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/create", db.add)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/read", db.fetch)
	http.HandleFunc("/delete", db.delete)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
