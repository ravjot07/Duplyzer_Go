package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"duplyzer/sequential"
	// Import other models as they are created
	// "./fixedpool"
	// "./concurrentwalks"
	// "./limitedfs"
)

func printResults(hashes sequential.Results) {
	for hash, files := range hashes {
		if len(files) > 1 {
			fmt.Println(hash[len(hash)-7:], len(files))
			for _, file := range files {
				fmt.Println("  ", file)
			}
		}
	}
}

func main() {
	model := flag.String("model", "sequential", "Concurrency model to use: sequential, fixedpool, concurrentwalks, limitedfs")
	dir := flag.String("dir", ".", "Directory to scan for duplicate files")
	flag.Parse()

	if *dir == "" {
		log.Fatal("Missing parameter, provide dir name!")
	}

	var hashes sequential.Results
	var err error

	switch *model {
	case "sequential":
		fmt.Println("Running Sequential Program")
		hashes, err = sequential.SearchTreeSequential(*dir)
	case "fixedpool":
		fmt.Println("Running Fixed Pool of Worker Goroutines")
		// hashes, err = fixedpool.SearchTreeFixedPool(*dir)
	case "concurrentwalks":
		fmt.Println("Running Concurrent Directory Walks")
		// hashes, err = concurrentwalks.SearchTreeConcurrentWalks(*dir)
	case "limitedfs":
		fmt.Println("Running Limited Goroutines for File System Operations")
		// hashes, err = limitedfs.SearchTreeLimitedFS(*dir)
	default:
		fmt.Println("Unknown model:", *model)
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}

	printResults(hashes)
}
