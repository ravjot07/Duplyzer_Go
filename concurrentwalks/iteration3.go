// concurrentwalks/iteration3.go
package concurrentwalks

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

func processFiles(paths <-chan string, pairs chan<- shared.Pair, done chan<- bool) {
	for path := range paths {
		pairs <- hashFile(path)
	}
	done <- true
}

func collectHashes(pairs <-chan shared.Pair, result chan<- shared.Results) {
	hashes := make(shared.Results)
	for p := range pairs {
		hashes[p.Hash] = append(hashes[p.Hash], p.Path)
	}
	result <- hashes
}

func searchTree(dir string, paths chan<- string, wg *sync.WaitGroup) error {
	defer wg.Done()
	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}
		if fi.Mode().IsDir() && p != dir {
			wg.Add(1)
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

func Run(dir string) shared.Results {
	workers := 2 * runtime.GOMAXPROCS(0)
	fmt.Println("Number of workers (double the number of logical CPUs):", workers)
	paths := make(chan string)
	pairs := make(chan shared.Pair)
	done := make(chan bool)
	result := make(chan shared.Results)
	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		go processFiles(paths, pairs, done)
	}
	go collectHashes(pairs, result)

	wg.Add(1)
	err := searchTree(dir, paths, wg)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	close(paths)
	for i := 0; i < workers; i++ {
		<-done
	}
	close(pairs)

	return <-result
}
