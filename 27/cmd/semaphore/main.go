package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type pair struct {
	hash string
	path string
}

type fileList []string
type result map[string]fileList

const workers = 32

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, provide directory name!")
	}
	var wg sync.WaitGroup
	sem := make(chan any, workers)
	pairs := make(chan pair, workers)
	results := make(chan result)

	go collect(pairs, results)

	wg.Add(1)
	err := walkDir(os.Args[1], pairs, &wg, sem)

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	close(pairs)
	hashes := <-results

	for hash, files := range hashes {
		if len(files) > 1 {
			// use 7 characters like git
			fmt.Println(hash[len(hash)-7:], len(files))
			for _, file := range files {
				fmt.Println(" ", file)
			}
		}
	}
	close(results)
}

func collect(pairs <-chan pair, results chan<- result) {
	hashed := make(result)

	for p := range pairs {
		hashed[p.hash] = append(hashed[p.hash], p.path)
	}

	results <- hashed
}

func walkDir(dir string, pairs chan<- pair, wg *sync.WaitGroup, sem chan any) error {
	defer wg.Done()
	sem <- nil
	defer func() {
		<-sem
	}()

	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		// ignore the directory itself to avoid an infinite loop
		if info.Mode().IsDir() && path != dir {
			wg.Add(1)
			go walkDir(path, pairs, wg, sem)
			return filepath.SkipDir
		}

		if info.Mode().IsRegular() && info.Size() > 0 {
			wg.Add(1)
			go processFile(path, pairs, wg, sem)
		}

		return nil
	})
}

func processFile(path string, pairs chan<- pair, wg *sync.WaitGroup, sem chan any) {
	defer wg.Done()
	sem <- nil
	defer func() {
		<-sem
	}()

	pairs <- hashFile(path)
}

func hashFile(path string) pair {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Skipping %s: %s", path, err)
		return pair{"", ""}
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Printf("Error hashing %s: %s", path, err)
		return pair{"", ""}
	}

	return pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}
