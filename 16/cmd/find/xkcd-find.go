package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type cartoon struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no file given")
		os.Exit(-1)
	}

	fileName := os.Args[1]

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "no search term")
		os.Exit(0)
	}

	var (
		items []cartoon
		terms []string
		input io.ReadCloser
		count int
		err   error
	)

	if input, err = os.Open(fileName); err != nil {
		fmt.Fprintf(os.Stderr, "bad file: %s\n", err)
		os.Exit(-1)
	}

	// decode file
	if err = json.NewDecoder(input).Decode(&items); err != nil {
		fmt.Fprintf(os.Stderr, "bad json: %s\n", err)
		os.Exit(-1)
	}

	fmt.Fprintf(os.Stderr, "read %d comics\n", len(items))

	// get search terms
	for _, t := range os.Args[2:] {
		terms = append(terms, strings.ToLower(t))
	}

	// search
cartoons:
	for _, item := range items {
		title := strings.ToLower(item.Title)
		transcript := strings.ToLower(item.Transcript)

		for _, term := range terms {
			if !strings.Contains(title, term) && !strings.Contains(transcript, term) {
				continue cartoons
			}
		}

		fmt.Printf(
			"https://xkcd.com/%d/ %s/%s/%s %q\n",
			item.Num, item.Month, item.Day, item.Year, item.Title,
		)
		count++
	}

	fmt.Fprintf(os.Stderr, "found %d comics\n", count)
}
