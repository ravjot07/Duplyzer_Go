// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"duplyzer/concurrentwalks"
	"duplyzer/fixedpool"
	"duplyzer/limitedfs"
	"duplyzer/sequential"
	"duplyzer/shared"
)

func printResults(hashes shared.Results) {
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
	start := time.Now()
	model := flag.String("model", "fixedpool", "Concurrency model to use: sequential, fixedpool, concurrentwalks, limitedfs")
	dir := flag.String("dir", ".", "Directory to scan for duplicate files")
	flag.Parse()

	if *dir == "" {
		log.Fatal("Missing parameter, provide dir name!")
	}

	var hashes shared.Results
	var err error

	switch *model {
	case "fixedpool":
		fmt.Println("Running Fixed Pool of Worker Goroutines")
		hashes = fixedpool.Run(*dir)
	case "sequential":
		fmt.Println("Running Sequential Program")
		hashes = sequential.Run(*dir)
	case "concurrentwalks":
		fmt.Println("Running Concurrent Directory Walks")
		hashes = concurrentwalks.Run(*dir)
	case "limitedfs":
		fmt.Println("Running Limited Goroutines for File System Operations")
		hashes = limitedfs.Run(*dir)
	default:
		fmt.Println("Unknown model:", *model)
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}

	printResults(hashes)

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}
