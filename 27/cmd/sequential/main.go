package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type pair struct {
	hash string
	path string
}

type fileList []string
type result map[string]fileList

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, provide directory name!")
	}
	workers := 2 * runtime.GOMAXPROCS(0)
	paths := make(chan string)
	pairs := make(chan pair)
	done := make(chan bool)
	results := make(chan result)

	for range workers {
		go processFiles(paths, pairs, done)
	}

	// we need another go routine so we don't block here
	go collectHashes(pairs, results)

	hashes := searchTree(os.Args[1], workers,
		paths, pairs, results, done)

	for hash, files := range hashes {
		if len(files) > 1 {
			// use 7 characters like git
			fmt.Println(hash[len(hash)-7:], len(files))
			for _, file := range files {
				fmt.Println(" ", file)
			}
		}
	}
}

func searchTree(dir string, workers int,
	paths chan<- string, pairs chan<- pair,
	results <-chan result, done <-chan bool) result {

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.Mode().IsRegular() && info.Size() > 0 {
			paths <- path
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking %s %s", dir, err)
	}

	// close paths so that the workers stop
	close(paths)

	// wait for all workers to be done & close pairs
	for range workers {
		<-done
	}
	close(pairs)

	hashes := <-results
	return hashes
}

func collectHashes(pairs <-chan pair, results chan<- result) {
	hashes := make(result)

	for p := range pairs {
		hashes[p.hash] = append(hashes[p.hash], p.path)
	}

	results <- hashes
}

func processFiles(paths <-chan string, pairs chan<- pair, done chan<- bool) {
	for path := range paths {
		pairs <- hashFile(path)
	}

	done <- true
}

func hashFile(path string) pair {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening %s %s", path, err)
	}
	defer file.Close()

	hash := md5.New() // not secure but fast and good enough
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatalf("Error copying %s %s", path, err)
	}

	return pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}
