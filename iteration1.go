package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Define a struct type 'pair' to hold the hash and the path of a file
type pair struct {
	hash string
	path string
}
type fileList []string
type results map[string]fileList

// Function to hash a file given its path and return a 'pair' struct
func hashFile(path string) pair {
	// Open the file specified by the path
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	// Ensure the file is closed when the function returns
	defer file.Close()
	// Create a new MD5 hash
	hash := md5.New()

	// Copy the contents of the file into the hash
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	// Return a pair struct with the computed hash and file path
	return pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}

// Function to traverse a directory tree and compute file hashes

func searchTree(dir string) (results, error) {
	hashes := make(results)

	// Define a visit function to be called for each file and directory found by Walk
	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil && err != os.ErrNotExist {
			return err
		}
		// Check if the file is regular and has a size greater than 0
		if fi.Mode().IsRegular() && fi.Size() > 0 {
			h := hashFile(p) // Compute the hash of the file
			// Append the file path to the list of files with the same hash
			hashes[h.hash] = append(hashes[h.hash], h.path)
		}
		return nil // Continue walking the directory tree
	}
	// Walk the directory tree rooted at dir, calling visit for each file and directory
	err := filepath.Walk(dir, visit)

	return hashes, err
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, provide dir name!")
	}
	// Call searchTree with the directory name provided as a command-line argument
	if hashes, err := searchTree(os.Args[1]); err == nil {
		for hash, files := range hashes {
			if len(files) > 1 {
				// Print the last 7 characters of the hash and the number of files like git
				fmt.Println(hash[len(hash)-7:], len(files))

				// Print each file path indented for readability
				for _, file := range files {
					fmt.Println("  ", file)
				}
			}
		}
	}
}
