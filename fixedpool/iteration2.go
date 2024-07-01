// fixedpool/iteration2.go
package fixedpool

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

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

func Run(dir string) shared.Results {
	workers := 2 * runtime.GOMAXPROCS(0)
	fmt.Println("No of workers ie total double of total local CPUs", workers)
	paths := make(chan string)
	pairs := make(chan shared.Pair)
	done := make(chan bool)
	result := make(chan shared.Results)

	for i := 0; i < workers; i++ {
		go processFiles(paths, pairs, done)
	}
	go collectHashes(pairs, result)

	if err := searchTree(dir, paths); err != nil {
		return nil
	}
	close(paths)
	for i := 0; i < workers; i++ {
		<-done
	}
	close(pairs)

	return <-result
}
