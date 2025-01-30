package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

const url = "https://jsonplaceholder.typicode.com/"

var form = `
<h1>Todo #{{.Id}}</h1>
<div>{{printf "User %d" .UserId}}</div>
<div>{{printf "%s (completed: %t)" .Title .Completed}}</div>
`

type todo struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(url + r.URL.Path[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var item todo

		if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.New("foo")
		tmpl.Parse(form)
		tmpl.Execute(w, item)
	}

	http.Error(w, "Unknown error", http.StatusInternalServerError)
}
