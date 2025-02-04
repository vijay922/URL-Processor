# URL Processor

## Overview
This Go script processes a list of URLs from an input file, extracts and modifies them based on specific rules, and saves the cleaned list to an output file. It ensures that directory paths and unique URLs are included while filtering out unnecessary versions or incomplete URLs.

## Features
- Reads URLs from an input file.
- Extracts and normalizes directory paths.
- Removes duplicate URLs.
- Excludes URLs with version parameters (`ver=` or `v=`) unless they are files.
- Saves the processed URLs to an output file in sorted order.
- Helps identify vulnerabilities such as:
  - Local File Inclusion (LFI)
  - Cross-Site Scripting (XSS)
  - SQL Injection (SQLi)
  - Path Traversal
  - DOM-based Cross-Site Scripting (DOM XSS)
  - HTML Injection
  - Command Injection

## Installation
1. Ensure you have Go installed on your system.
2. Clone this repository and navigate to the project directory.

## Usage
Run the script with the following flags:

```sh
 go run main.go -l input.txt -o output.txt
```

### Parameters
- `-l` : Path to the input file containing URLs.
- `-o` : Path to the output file where processed URLs will be saved.

## How It Works
1. The script reads URLs from the input file.
2. It checks if the URL has a valid scheme and host.
3. It extracts directory paths from URLs and adds them to the output list.
4. It filters URLs based on query parameters (ignores versioning parameters unless it's a file).
5. It sorts and writes the unique URLs to the output file.

## Example
### Input (input.txt)
```
https://archive.example.com/stable/dists/bionic/InRelease
https://chloe-ml.demo.example.com/comments/feed/
https://hiltongrandvacations.xcc.example.com/wp-content/plugins/xcloud/app/components/leadcapture/modal-thank-you.html
https://pitneybowes.cert.xcc.example.com/wp-content/plugins/xcloud/app/js/lib/angular/1.5.9/angular.js?ver=0.86-272&anim=katana
https://uhg.cert.xce.example.com/wp-content/plugins/xce/css/main.css?ver=0.1-55
https://chloe-ml.demo.example.com/wp-includes/js/jquery/ui/core.min.js?ver=1.11.4
https://chloe-ml.demo.example.com/wp-includes/js/mediaelement/wp-mediaelement.min.js?ver=f6347901ea6b2bb2d7a7251f50c6313d
https://www.test.com/hs-fs/hubfs/Google%20Drive%20Integration/blur%20for%20hubspot-Jul-22-2020-08-04-48-67-PM.png?width=133&name=blur%20for%20hubspot-Jul-22-2020-08-04-48-67-PM.png
https://www.test.com/helpers/loadCodeSequentially
https://www.test.com/newsroom/author/white-ops/page/1)
```

### Output (output.txt)
```
https://archive.example.com/stable/dists/
https://archive.example.com/stable/dists/bionic-updates/
https://archive.example.com/stable/dists/bionic/
https://chloe-ml.demo.example.com/comments/
https://chloe-ml.demo.example.com/comments/feed/
https://chloe-ml.demo.example.com/wp-includes/
https://chloe-ml.demo.example.com/wp-includes/js/
https://chloe-ml.demo.example.com/wp-includes/js/jquery/
https://chloe-ml.demo.example.com/wp-includes/js/jquery/ui/
https://chloe-ml.demo.example.com/wp-includes/js/mediaelement/
https://hiltongrandvacations.xcc.example.com/wp-content/
https://hiltongrandvacations.xcc.example.com/wp-content/plugins/
https://hiltongrandvacations.xcc.example.com/wp-content/plugins/xcloud/
https://hiltongrandvacations.xcc.example.com/wp-content/plugins/xcloud/app/
https://hiltongrandvacations.xcc.example.com/wp-content/plugins/xcloud/app/components/
https://hiltongrandvacations.xcc.example.com/wp-content/plugins/xcloud/app/components/leadcapture/
https://pitneybowes.cert.xcc.example.com/wp-content/
https://pitneybowes.cert.xcc.example.com/wp-content/plugins/
https://pitneybowes.cert.xcc.example.com/wp-content/plugins/xcloud/
https://pitneybowes.cert.xcc.example.com/wp-content/plugins/xcloud/app/
https://pitneybowes.cert.xcc.example.com/wp-content/plugins/xcloud/app/js/
https://pitneybowes.cert.xcc.example.com/wp-content/plugins/xcloud/app/js/lib/
https://pitneybowes.cert.xcc.example.com/wp-content/plugins/xcloud/app/js/lib/angular/
https://pitneybowes.cert.xcc.example.com/wp-content/plugins/xcloud/app/js/lib/angular/1.5.9/
https://uhg.cert.xce.example.com/wp-content/
https://uhg.cert.xce.example.com/wp-content/plugins/
https://uhg.cert.xce.example.com/wp-content/plugins/xce/
https://uhg.cert.xce.example.com/wp-content/plugins/xce/css/
https://www.test.com/helpers/
https://www.test.com/hs-fs/
https://www.test.com/hs-fs/hubfs/
https://www.test.com/hs-fs/hubfs/Google%20Drive%20Integration/
https://www.test.com/hs-fs/hubfs/Google%20Drive%20Integration/blur%20for%20hubspot-Jul-22-2020-08-04-48-67-PM.png?width=133&name=blur%20for%20hubspot-Jul-22-2020-08-04-48-67-PM.png
https://www.test.com/newsroom/
https://www.test.com/newsroom/author/
https://www.test.com/newsroom/author/white-ops/
https://www.test.com/newsroom/author/white-ops/page/)
```

## Functions Explained
### `main()`
- Parses command-line flags for input and output files.
- Reads URLs from the input file.
- Calls `processPath()` to extract directory paths.
- Filters and stores unique URLs.
- Sorts and writes URLs to the output file.

### `processPath(path string) []string`
- Extracts directory paths from a given URL path.
- Returns a list of parent directories.

### `readLines(path string) ([]string, error)`
- Reads lines from the specified file and returns them as a slice of strings.

### `writeLines(path string, lines []string) error`
- Writes a list of strings to the specified output file.

## License
This project is licensed under the MIT License.

