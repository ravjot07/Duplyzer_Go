#  Duplyzer_Go: Concurrent File Duplication Detector in Go

**Duplyzer_Go** is an tool developed in Go for detecting duplicate files in file systems based on their content using MD5 hashing. By leveraging Go's concurrency features, Duplyzer_Go can significantly speed up the process, making it an ideal solution for managing large data volumes efficiently.

## Features

-   **Content-Based Hashing**: Utilizes MD5 hashing to accurately detect duplicate files by their content.
-   **Concurrent Processing**: Implements four different concurrency models to optimize performance across various system configurations.
-   **Cross-Platform Compatibility**: Runs on multiple operating systems, thanks to Go's cross-platform support.
-   **Command-Line Interface**: Easy-to-use CLI for straightforward execution and minimal setup.
-   **Customizable Concurrency**: Allows users to choose the best concurrency model based on their system's performance.

## Concurrency Models

### 1. Sequential Program: iteration1.go

-   **Description**: This program performs a sequential file walk without using goroutines. It hashes each file's contents and checks for duplicates by comparing hash values. 
    
-   **Output**: Lists hashes and the number of times each hash appears, with corresponding file paths.
    
-   **Performance**: This method is the simplest and slowest, especially as the size of the filesystem increases. It serves as a baseline for measuring performance improvements in the following methods.
    

### 2.  Fixed Pool of Worker Goroutines: iteration2.go 

-   **Description**: Uses a fixed pool of worker goroutines that receive file paths from a sequential file walk. Each worker hashes the file and checks for duplicates.
      
-   **Output:** Similar to the sequential method but potentially faster due to concurrent processing of files.
    
-   **Performance:** This method introduces concurrency but still relies on a sequential walk to distribute work, which might still be a bottleneck if the walk itself is slow. 

## Installation

Ensure you have Go installed on your system.
Clone the repository:

`git clone https://github.com/yourusername/Duplyzer_Go.git`
`cd Duplyzer_Go`

## Usage

To start scanning for duplicate files in a directory, use the following command:

`go run duplyzer_modal.go path_to_directory` 

Replace `path_to_directory` with the path of the directory you want to scan.
### For example
`go run iteration2.go test` 
