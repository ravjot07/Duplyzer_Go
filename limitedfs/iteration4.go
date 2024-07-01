package limitedfs

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"duplyzer/shared"
)

func hashFile(path string) shared.Pair {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	return shared.Pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}

func processFiles(path string, pairs chan<- shared.Pair, wg *sync.WaitGroup, limits chan bool) {
	defer wg.Done()
	limits <- true
	defer func() { <-limits }()
	pair := hashFile(path)
	fmt.Printf("Processed file: %s, Hash: %s\n", pair.Path, pair.Hash)
	pairs <- pair
}

func collectHashes(pairs <-chan shared.Pair, result chan<- shared.Results) {
	hashes := make(shared.Results)
	for p := range pairs {
		fmt.Printf("Collecting hash: %s for file: %s\n", p.Hash, p.Path)
		hashes[p.Hash] = append(hashes[p.Hash], p.Path)
	}
	result <- hashes
}

func searchTree(dir string, pairs chan<- shared.Pair, wg *sync.WaitGroup, limits chan bool) error {
	defer wg.Done()
	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}
		if fi.Mode().IsDir() && p != dir {
			wg.Add(1)
			go searchTree(p, pairs, wg, limits)
			return filepath.SkipDir
		}
		if fi.Mode().IsRegular() && fi.Size() > 0 {
			wg.Add(1)
			go processFiles(p, pairs, wg, limits)
		}
		return nil
	}
	limits <- true
	defer func() { <-limits }()
	return filepath.Walk(dir, visit)
}

func Run(dir string) shared.Results {
	workers := 2 * runtime.GOMAXPROCS(0)
	fmt.Println("Number of workers (double the number of logical CPUs):", workers)
	limits := make(chan bool, workers)
	pairs := make(chan shared.Pair)
	result := make(chan shared.Results)
	wg := new(sync.WaitGroup)

	go collectHashes(pairs, result)

	wg.Add(1)
	err := searchTree(dir, pairs, wg, limits)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	close(pairs)

	return <-result
}
