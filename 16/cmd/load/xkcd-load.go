package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

// returns the metadata for one comic by number
func getOne(i int) []byte {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)
	response, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "can't read: %s\n", err)
		os.Exit(-1)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "skipping %d: got %d\n", i, response.StatusCode)
		return nil
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid body: %s\n", err)
		os.Exit(-1)
	}

	return body
}

func main() {
	var (
		output io.WriteCloser = os.Stdout
		err    error
		count  int
		fails  int
		data   []byte
	)

	if len(os.Args) > 1 {
		output, err = os.Create(os.Args[1])

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		defer output.Close()
	}

	fmt.Fprintln(output, "[")
	defer fmt.Fprintln(output, "]")

	// stop if we get 2 404s in a row
	for i := 1; fails < 2; i++ {
		if data = getOne(i); data == nil {
			fails++
			continue
		}

		if count > 0 {
			fmt.Fprint(output, ",")
		}

		_, err = io.Copy(output, bytes.NewBuffer(data))

		if err != nil {
			fmt.Fprintf(os.Stderr, "stopped: %s\n", err)
			os.Exit(-1)
		}

		fails = 0
		count++
	}

	fmt.Fprintf(os.Stderr, "read %d commics\n", count)
}
