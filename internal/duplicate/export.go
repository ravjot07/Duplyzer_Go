package duplicate

import (
	"duplyzer/shared"
	"encoding/csv"
	"encoding/json"
	"os"
)

type FileHashInfo struct {
	Hash  string   `json:"hash"`
	Files []string `json:"files"`
}

func ExportToJSON(results shared.Results, outputPath string) error {
	fileHashInfos := []FileHashInfo{}
	for hash, files := range results {
		fileHashInfos = append(fileHashInfos, FileHashInfo{Hash: hash, Files: files})
	}
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	return encoder.Encode(fileHashInfos)
}

func ExportToCSV(results shared.Results, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for hash, files := range results {
		record := append([]string{hash}, files...)
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
