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
	inputFile := flag.String("l", "", "Path to the input file containing URLs.")
	outputFile := flag.String("o", "", "Path to the output file where processed URLs will be saved.")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Both -l and -o flags are required.")
		os.Exit(1)
	}

	urls, err := readLines(*inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	groupedURLs := groupURLs(urls)
	processedURLs := processGroupedURLs(groupedURLs)
	uniqueURLs := deduplicate(processedURLs)
	sort.Strings(uniqueURLs)

	if err := writeLines(*outputFile, uniqueURLs); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}
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
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func groupURLs(urls []string) []string {
	groups := make(map[string]string)

	for _, u := range urls {
		parsed, err := url.Parse(u)
		if err != nil {
			fmt.Printf("Error parsing URL %s: %v\n", u, err)
			continue
		}

		query := parsed.Query()
		paramKeys := make([]string, 0, len(query))
		for k := range query {
			paramKeys = append(paramKeys, k)
		}
		sort.Strings(paramKeys)
		key := parsed.Path + "?" + strings.Join(paramKeys, ",")

		if _, exists := groups[key]; !exists {
			groups[key] = u
		}
	}

	grouped := make([]string, 0, len(groups))
	for _, u := range groups {
		grouped = append(grouped, u)
	}

	return grouped
}

func processGroupedURLs(grouped []string) []string {
	var processed []string
	allowedExtensions := map[string]bool{
		".php":   true,
		".html":  true,
		".aspx":  true,
		".jsp":   true,
		".htm":   true,
		".asp":   true,
	}

	for _, u := range grouped {
		parsed, err := url.Parse(u)
		if err != nil {
			fmt.Printf("Error parsing URL %s: %v\n", u, err)
			continue
		}

		hasQuery := len(parsed.Query()) > 0
		extension := getExtension(parsed.Path)
		hasAllowedExt := allowedExtensions[extension]
		noExtension := isPathWithoutExtension(parsed.Path)

		if hasQuery && (hasAllowedExt || noExtension) {
			processed = append(processed, u)
			continue
		}

		pathComponents := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		var dirs []string

		for _, comp := range pathComponents {
			if strings.ContainsAny(comp, ".=") {
				break
			}
			dirs = append(dirs, comp)
			if len(dirs) >= 2 {
				break
			}
		}

		if len(dirs) == 0 {
			continue
		}

		newPath := "/" + strings.Join(dirs, "/") + "/"
		parsed.Path = newPath
		parsed.RawQuery = ""
		processedURL := parsed.String()

		processed = append(processed, processedURL)
	}

	return processed
}

func getExtension(path string) string {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return ""
	}
	lastComponent := path[lastSlash+1:]
	dotIndex := strings.LastIndex(lastComponent, ".")
	if dotIndex == -1 {
		return ""
	}
	return lastComponent[dotIndex:]
}

func isPathWithoutExtension(path string) bool {
	if path == "" {
		return true
	}
	lastSlash := strings.LastIndex(path, "/")
	lastComponent := path
	if lastSlash != -1 {
		lastComponent = path[lastSlash+1:]
	}
	return !strings.Contains(lastComponent, ".")
}

func deduplicate(urls []string) []string {
	seen := make(map[string]bool)
	unique := []string{}
	for _, u := range urls {
		if !seen[u] {
			seen[u] = true
			unique = append(unique, u)
		}
	}
	return unique
}

func writeLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}
