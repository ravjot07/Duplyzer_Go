package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	// "runtime/pprof"
	"time"

	"duplyzer/concurrentwalks"
	"duplyzer/filemanagement"
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

func generateReportJSON(hashes shared.Results) ([]byte, error) {
	type ReportEntry struct {
		Hash  string   `json:"hash"`
		Files []string `json:"files"`
	}
	var report []ReportEntry
	for hash, files := range hashes {
		if len(files) > 1 {
			report = append(report, ReportEntry{Hash: hash, Files: files})
		}
	}
	return json.Marshal(report)
}

func reportHandler(hashes shared.Results) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		reportJSON, err := generateReportJSON(hashes)
		if err != nil {
			http.Error(w, "Unable to generate report", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(reportJSON)
	}
}

func main() {
	// Parse command-line flags
	model := flag.String("model", "fixedpool", "Concurrency model to use: sequential, fixedpool, concurrentwalks, limitedfs")
	dir := flag.String("dir", ".", "Directory to scan for duplicate files")
	outputFormat := flag.String("output-format", "text", "Output format: text, json, csv")
	action := flag.String("action", "", "Action to perform on duplicates: delete, move, hard-link")
	targetDir := flag.String("target-dir", "", "Target directory for move or hard-link actions")
	webReport := flag.Bool("web-report", false, "Generate a web-based report")
	flag.Parse()

	if *dir == "" {
		log.Fatal("Missing parameter, provide dir name!")
	}
	if (*action == "move" || *action == "hard-link") && *targetDir == "" {
		log.Fatal("Target directory must be specified for move or hard-link actions")
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

	// Manage duplicates based on specified action
	if *action != "" {
		err := filemanagement.ManageDuplicates(*action, *targetDir, hashes)
		if err != nil {
			fmt.Printf("Error managing duplicates: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Duplicate file management completed successfully")
	}

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

	// Export results based on the chosen output format
	if *webReport {
		http.HandleFunc("/report", reportHandler(hashes))
		fmt.Println("Starting web server at http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
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
}
