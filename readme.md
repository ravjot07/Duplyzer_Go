#  Duplyzer_Go: Concurrent File Duplication Detector in Go

**Duplyzer_Go** is an tool developed in Go for detecting duplicate files in file systems based on their content using MD5 hashing. By leveraging Go's concurrency features, Duplyzer_Go can significantly speed up the process, making it an ideal solution for managing large data volumes efficiently.

## Features

-   **Content-Based Hashing**: Utilizes MD5 hashing to accurately detect duplicate files by their content.
-   **Concurrent Processing**: Implements four different concurrency models to optimize performance across various system configurations.
-   **Cross-Platform Compatibility**: Runs on multiple operating systems, thanks to Go's cross-platform support.
-   **Command-Line Interface**: Easy-to-use CLI for straightforward execution and minimal setup.
-   **Customizable Concurrency**: Allows users to choose the best concurrency model based on their system's performance.
- **Local file input**:Implemented volume mounting in Docker run command to allow access to local directories - Users can now specify a directory on their local machine to be scanned for duplicate files
- - **Support Report Generation**:Implemented feature to generate output of duplicate files as json or csv file.

## Concurrency Models

### 1. Sequential Program: iteration1.go

-   **Description**: This program performs a sequential file walk without using goroutines. It hashes each file's contents and checks for duplicates by comparing hash values. 
    
-   **Output**: Lists hashes and the number of times each hash appears, with corresponding file paths.
    
-   **Performance**: This method is the simplest and slowest, especially as the size of the filesystem increases. It serves as a baseline for measuring performance improvements in the following methods.
    

### 2.  Fixed Pool of Worker Goroutines: iteration2.go 

-   **Description**: Uses a fixed pool of worker goroutines that receive file paths from a sequential file walk. Each worker hashes the file and checks for duplicates.
      
-   **Output:** Similar to the sequential method but potentially faster due to concurrent processing of files.
    
-   **Performance:** This method introduces concurrency but still relies on a sequential walk to distribute work, which might still be a bottleneck if the walk itself is slow. 

### 3. Concurrent Directory Walks: iteration3.go 

  

- **Description:** Enhances the concurrency by initiating a new file walk goroutine for each directory encountered. This spreads the hashing workload more efficiently across multiple cores.

  

- **Output:** Same structured output, demonstrating efficiency.

- **Performance:** Expected to show noticeable performance improvements over Method 1, especially in environments with deeply nested directory structures.

  

### 4. Limited Goroutines for File System Operations: iteration4.go 

  

- **Description:** Creates a goroutine for each file and directory but limits the number of goroutines that can perform file system operations concurrently to reduce contention.

  

- **Output:** Continues to output duplicate hashes and file paths.

- **Performance:** Likely the most sophisticated and efficient version, balancing concurrency and system resource utilization to minimize bottlenecks caused by excessive thread creation and I/O operations.

## Installation

Ensure you have Go installed on your system.
Clone the repository:

`git clone https://github.com/ravjot07/Duplyzer_Go.git`
`cd Duplyzer_Go`

## Usage
### Without Docker
To start scanning for duplicate files in a directory, use the following command:

`go run main.go --model=<model> --dir=<path_to_directory> --output-format=<format>` 

Replace `<model>` with the concurrency model you want to use (`sequential`, `fixedpool`, `concurrentwalks`, `limitedfs`) and `<path_to_directory>` with the path of the directory you want to scan.
Replace `<format>` with the type of report you want to generate  (`csv`, `json`)
###  Using Build 
Build Using:

`go build -o duplizer-cli .`

Run the CLI Tool:

`./duplizer-cli --model=<model> --dir=<path_to_directory> --output-format=<format>`


### With Docker 
Build the Docker Image:

`docker build -t duplizer-cli .`

Replace `/path/to/local/directory` with the actual path to the directory you want to scan for duplicate files.

```
docker run --rm -v "/path/to/local/directory:/app/scan-dir" duplizer-cli --model concurrentwalks --dir /app/scan-dir
```
Example:
```
docker run --rm -v "/home/ravjot/Desktop/Go Lang Projects/Duplyzer_Go/test:/app/scan-dir" duplizer-cli --model concurrentwalks --dir /app/scan-dir
```

### For Example: Limited Goroutines for File System Operations
`go run main.go --model=limitedfs --dir=./test` 

Output:
```
Running Limited Goroutines for File System Operations
Number of workers (double the number of logical CPUs): 16
Processed file: test/rc.txt, Hash: c1c7f5decb9ff01edf1af096ebb8f4a4
Collecting hash: c1c7f5decb9ff01edf1af096ebb8f4a4 for file: test/rc.txt
Processed file: test/rb.txt, Hash: 6f11ac20bf1d3c85c586fa793fa03186
Collecting hash: 6f11ac20bf1d3c85c586fa793fa03186 for file: test/rb.txt
Processed file: test/r1/ra.txt, Hash: 6f11ac20bf1d3c85c586fa793fa03186
Collecting hash: 6f11ac20bf1d3c85c586fa793fa03186 for file: test/r1/ra.txt
fa03186 2
   test/rb.txt
   test/r1/ra.txt
Execution time: 503.904µs
  ```
  

