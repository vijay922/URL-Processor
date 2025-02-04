package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
)

func main() {
	inputFile := flag.String("l", "", "Path to the input file containing URLs")
	outputFile := flag.String("o", "", "Path to the output file")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Both -l and -o flags are required")
		os.Exit(1)
	}

	urls, err := readLines(*inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	uniqueURLs := make(map[string]bool)

	for _, urlStr := range urls {
		u, err := url.Parse(urlStr)
		if err != nil {
			fmt.Printf("Error parsing URL %s: %v\n", urlStr, err)
			continue
		}

		if u.Host == "" || u.Scheme == "" {
			fmt.Printf("Invalid URL %s: missing scheme or host\n", urlStr)
			continue
		}

		path := u.EscapedPath()
		rawQuery := u.RawQuery
		hasVersion := strings.Contains(rawQuery, "ver=") || strings.Contains(rawQuery, "v=")

		splitDirPaths := processPath(path)

		// Add split directories as full URLs
		for _, dirPath := range splitDirPaths {
			fullURL := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, dirPath)
			uniqueURLs[fullURL] = true
		}

		// Check if we need to add the original URL
		isDir := strings.HasSuffix(path, "/")
		components := strings.Split(strings.Trim(path, "/"), "/")
		lastComponent := ""
		if len(components) > 0 && !isDir {
			lastComponent = components[len(components)-1]
		}
		isFile := !isDir && strings.Contains(lastComponent, ".")

		if rawQuery != "" && !hasVersion && isFile {
			uniqueURLs[urlStr] = true
		}
	}

	// Collect and sort the URLs
	var sortedURLs []string
	for u := range uniqueURLs {
		sortedURLs = append(sortedURLs, u)
	}
	sort.Strings(sortedURLs)

	// Write to output file
	err = writeLines(*outputFile, sortedURLs)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}
}

func processPath(path string) []string {
	isDir := strings.HasSuffix(path, "/")
	trimmedPath := strings.Trim(path, "/")
	components := strings.Split(trimmedPath, "/")

	splitLen := len(components)
	if !isDir {
		if len(components) == 0 {
			splitLen = 0
		} else {
			splitLen = len(components) - 1
		}
	}

	var dirPaths []string
	currentPath := ""
	for i := 0; i < splitLen; i++ {
		component := components[i]
		currentPath = currentPath + "/" + component
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func writeLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
