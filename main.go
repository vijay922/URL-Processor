package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
)

func main() {
	inputFile := flag.String("l", "", "Path to the input file containing URLs")
	outputFile := flag.String("o", "", "Path to the output file")
	flag.Usage = func() {
		fmt.Println("Usage: go run main.go -l input.txt -o output.txt")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	urls, err := readLines(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	uniqueURLs := make(map[string]struct{})
	var mu sync.Mutex
	var wg sync.WaitGroup
	urlChan := make(chan string, len(urls))

	for _, u := range urls {
		wg.Add(1)
		go func(urlStr string) {
			defer wg.Done()
			processedURLs := processURL(urlStr)
			mu.Lock()
			for _, pURL := range processedURLs {
				uniqueURLs[pURL] = struct{}{}
			}
			mu.Unlock()
		}(u)
	}

	wg.Wait()
	close(urlChan)

	// Collect and sort the URLs
	sortedURLs := make([]string, 0, len(uniqueURLs))
	for u := range uniqueURLs {
		sortedURLs = append(sortedURLs, u)
	}
	sort.Strings(sortedURLs)

	// Write to output file
	if err := writeLines(*outputFile, sortedURLs); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}
}

func processURL(urlStr string) []string {
	var results []string
	u, err := url.Parse(strings.TrimSpace(urlStr))
	if err != nil || u.Host == "" || u.Scheme == "" {
		log.Printf("Skipping invalid URL: %s", urlStr)
		return results
	}

	path := u.EscapedPath()
	rawQuery := u.RawQuery
	hasVersion := strings.Contains(rawQuery, "ver=") || strings.Contains(rawQuery, "v=")
	splitDirPaths := processPath(path)

	for _, dirPath := range splitDirPaths {
		results = append(results, fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, dirPath))
	}

	isDir := strings.HasSuffix(path, "/")
	components := strings.Split(strings.Trim(path, "/"), "/")
	isFile := len(components) > 0 && !isDir && strings.Contains(components[len(components)-1], ".")

	if rawQuery != "" && !hasVersion && isFile {
		results = append(results, urlStr)
	}

	return results
}

func processPath(path string) []string {
	isDir := strings.HasSuffix(path, "/")
	components := strings.Split(strings.Trim(path, "/"), "/")
	splitLen := len(components)
	if !isDir && splitLen > 0 {
		splitLen--
	}

	var dirPaths []string
	currentPath := ""
	for i := 0; i < splitLen; i++ {
		currentPath += "/" + components[i]
		dirPaths = append(dirPaths, currentPath+"/")
	}

	return dirPaths
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func writeLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}
