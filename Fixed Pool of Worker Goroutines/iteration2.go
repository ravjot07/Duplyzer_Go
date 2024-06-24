package main

import (
	"crypto/md5"
	"fmt"
	"io"
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
type results map[string]fileList

func hashFile(path string) pair {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	return pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}

// The processFiles function reads file paths from a channel, computes the hash for each file, and sends the results to another channel. It signals completion through a done channel.

func processFiles(paths <-chan string, pairs chan<- pair, done chan<- bool) {
	for path := range paths {
		pairs <- hashFile(path)
	}

	done <- true
}

// The collectHashes function reads pair structs from a channel, aggregates them into a results map, and sends the final result to another channel.

func collectHashes(pairs <-chan pair, result chan<- results) {
	hashes := make(results)

	for p := range pairs {
		hashes[p.hash] = append(hashes[p.hash], p.path)
	}
	result <- hashes
}

func searchTree(dir string, paths chan<- string) error {

	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}
		if fi.Mode().IsRegular() && fi.Size() > 0 {
			paths <- p
		}
		return nil
	}

	return filepath.Walk(dir, visit)
}

// The run function coordinates multiple goroutines to process files in a directory, compute their hashes, and aggregate the results.

func run(dir string) results {
	// The number of worker goroutines is set to twice the number of logical CPUs available.
	workers := 2 * runtime.GOMAXPROCS(0)
	fmt.Println("No of workers ie total double of total local CPUs", workers)
	paths := make(chan string)
	pairs := make(chan pair)
	done := make(chan bool)
	result := make(chan results)

	// Starting the worker go routine
	for i := 0; i < workers; i++ {
		go processFiles(paths, pairs, done)
	}
	// Starting the Hash Collection Goroutine, we need another goroutine so we don't block here
	go collectHashes(pairs, result)

	// Searching the Directory
	if err := searchTree(dir, paths); err != nil {
		return nil
	}
	// we must close the paths channel so the workers stop
	close(paths)
	// wait for all the workers to be done
	for i := 0; i < workers; i++ {
		<-done
	}
	// by closing pairs we signal that all the hashes
	// have been collected; we have to do it here AFTER
	// all the workers are done
	close(pairs)

	return <-result
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, provide dir name!")
	}
	if hashes := run(os.Args[1]); hashes != nil {
		for hash, files := range hashes {
			if len(files) > 1 {
				fmt.Println(hash[len(hash)-7:], len(files))

				for _, file := range files {
					fmt.Println("  ", file)
				}
			}
		}
	}
}
