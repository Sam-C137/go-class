package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

	if hashes, err := searchTree(os.Args[1]); err != nil {
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
}

func searchTree(dir string) (result, error) {
	hashes := make(result)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		// ignore error for now
		if info.Mode().IsRegular() && info.Size() > 0 {
			h := hashFile(path)
			hashes[h.hash] = append(hashes[h.hash], h.path)
		}

		return nil
	})

	return hashes, err
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
