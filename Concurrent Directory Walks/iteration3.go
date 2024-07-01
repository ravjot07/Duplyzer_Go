package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
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

func searchTree(dir string, paths chan<- string, wg *sync.WaitGroup) error {

	defer wg.Done() // Decrement the wait group counter when the function returns.
	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}
		// If the file is a directory and it's not the root directory, recursively search this new directory in a new goroutine.
		if fi.Mode().IsDir() && p != dir {
			wg.Add(1) // Increment the wait group counter before starting a new goroutine.
			go searchTree(p, paths, wg)
			return filepath.SkipDir
		}
		if fi.Mode().IsRegular() && fi.Size() > 0 {
			paths <- p
		}
		return nil
	}

	return filepath.Walk(dir, visit)
}

// run sets up the environment and manages the flow of data between various components:
// it initiates file traversal, file processing, and hash aggregation.

func run(dir string) results {
	// The number of worker goroutines is set to twice the number of logical CPUs available.
	workers := 2 * runtime.GOMAXPROCS(0)
	fmt.Println("No of workers ie total double of total local CPUs", workers)
	paths := make(chan string)
	pairs := make(chan pair)
	done := make(chan bool)
	result := make(chan results)
	wg := new(sync.WaitGroup) // WaitGroup to synchronize goroutines.

	// Starting the worker go routine
	for i := 0; i < workers; i++ {
		go processFiles(paths, pairs, done)
	}
	// Starting the Hash Collection Goroutine, we need another goroutine so we don't block here
	go collectHashes(pairs, result)

	// multi-threaded walk of the directory tree; we need a
	// waitGroup because we don't know how many to wait for
	wg.Add(1)

	// Start the multi-threaded walk of the directory tree.
	err := searchTree(dir, paths, wg)

	if err != nil {
		log.Fatal(err)
	}
	// we must close the paths channel so the workers stop
	wg.Wait() // Wait for all file walking goroutines to complete.
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