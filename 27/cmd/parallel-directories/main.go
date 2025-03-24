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

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, provide directory name!")
	}
	paths := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	err := walkDir(os.Args[1], paths, &wg)

	if err != nil {
		log.Fatalf("Error walking directory %s %s", os.Args[1], err)
	}

	var swg sync.WaitGroup
	var hashes result
	swg.Add(1)
	go func() {
		hashes = processFile(paths)
		swg.Done()
	}()
	wg.Wait()
	close(paths)
	swg.Wait()

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

func walkDir(dir string, paths chan<- string, wg *sync.WaitGroup) error {
	defer wg.Done()

	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		// ignore the directory itself to avoid an infinite loop
		if info.Mode().IsDir() && path != dir {
			wg.Add(1)
			go walkDir(path, paths, wg)
			return filepath.SkipDir
		}

		if info.Mode().IsRegular() && info.Size() > 0 {
			paths <- path
		}

		return nil
	})
}

func processFile(paths <-chan string) result {
	hashed := make(result)

	for path := range paths {
		p := hashFile(path)
		hashed[p.hash] = append(hashed[p.hash], p.path)
	}

	return hashed
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
