package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func downloadFile(url, outputDir string, wg *sync.WaitGroup, verbose bool) {
	defer wg.Done()

	// Extract the file name from the URL
	filename := filepath.Base(url)
	outFile := filepath.Join(outputDir, filename)

	// Handle file name conflicts
	for i := 1; fileExists(outFile); i++ {
		newFilename := fmt.Sprintf("%s-%d%s", strings.TrimSuffix(filename, filepath.Ext(filename)), i, filepath.Ext(filename))
		outFile = filepath.Join(outputDir, newFilename)
	}

	if verbose {
		fmt.Printf("Starting download: %s\n", url)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory %s: %v\n", outputDir, err)
		return
	}

	// Create the file
	out, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", outFile, err)
		return
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Handle non-live URLs (HTTP status codes other than 200)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-live URL %s: HTTP %d\n", url, resp.StatusCode)
		return
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", outFile, err)
		return
	}

	// Print success message
	fmt.Printf("Downloaded %s\n", filename)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func usage() {
	fmt.Println("Usage: downloader [-h] [-l inputfile] [-o outputdir] [-v]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h              Display this help message.")
	fmt.Println("  -l inputfile    Specify the input file containing URLs (default: js-urls.txt).")
	fmt.Println("  -o outputdir    Specify the output directory (default: jsoutput-files).")
	fmt.Println("  -v              Enable verbose output.")
}

func main() {
	// Define command line flags
	showHelp := flag.Bool("h", false, "Display this help message")
	inputFile := flag.String("l", "js-urls.txt", "Input file containing URLs")
	outputDir := flag.String("o", "jsoutput-files", "Output directory")
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	// Show usage information if -h flag is provided
	if *showHelp {
		usage()
		return
	}

	// Open the input file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	// Read URLs and download in parallel
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		wg.Add(1)
		go downloadFile(url, *outputDir, &wg, *verbose)
	}
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
	}
}