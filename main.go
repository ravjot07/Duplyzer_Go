package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// "runtime/pprof"
	"time"

	"duplyzer/concurrentwalks"
	"duplyzer/fixedpool"
	"duplyzer/internal/duplicate"
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
	// Parse command-line flags
	model := flag.String("model", "fixedpool", "Concurrency model to use: sequential, fixedpool, concurrentwalks, limitedfs")
	dir := flag.String("dir", ".", "Directory to scan for duplicate files")
	outputFormat := flag.String("output-format", "text", "Output format: text, json, csv")
	flag.Parse()

	if *dir == "" {
		log.Fatal("Missing parameter, provide dir name!")
	}

	// // Create CPU profile
	// cpuProfile, err := os.Create("cpu_profile.prof")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// defer cpuProfile.Close()

	// if err := pprof.StartCPUProfile(cpuProfile); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }
	// defer pprof.StopCPUProfile()

	// // Create memory profile
	// memProfile, err := os.Create("mem_profile.prof")
	// if err != nil {
	// 	log.Fatal("could not create memory profile: ", err)
	// }
	// defer memProfile.Close()

	// Start timing
	start := time.Now()

	// Run selected model
	var hashes shared.Results
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

	// Stop timing
	elapsed := time.Since(start)
	// Export results based on the chosen output format
	outputPath := "output." + *outputFormat
	var err error
	switch *outputFormat {
	case "json":
		err = duplicate.ExportToJSON(hashes, outputPath)
	case "csv":
		err = duplicate.ExportToCSV(hashes, outputPath)
	default:
		fmt.Println("Unsupported output format")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error exporting results: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results saved to %s\n", outputPath)

	// Print results and performance metrics
	printResults(hashes)
	fmt.Printf("Execution time: %s\n", elapsed)

	// // Write memory profile
	// if err := pprof.WriteHeapProfile(memProfile); err != nil {
	// 	log.Fatal("could not write memory profile: ", err)
	// }
}
