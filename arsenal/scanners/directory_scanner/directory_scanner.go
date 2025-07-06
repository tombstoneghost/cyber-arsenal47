// Directory Scanner Module
package scanners

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
- Target IP or Hostname including "http or https"
- Wordlist to use (Default: /usr/share/wordlists/dirbuster/directory-list-2.3-medium.txt)
- HTTP Codes (Default: 2XX, 3XX)
- Search Files (Default: true)
- File Extensions (Default: txt,php,html)
- Search Files Only (Default: false)
*/

var (
	mu           sync.Mutex // Mutex to coordinate printing
	total        int        // Total number of URLs to scan
	completed    uint32     // Atomic counter for completed scans
	printedLines int        // Track the number of lines printed for results
)

func scanner(url string, errorCodes []int, client *http.Client) {
	resp, err := client.Get(url)
	if err != nil {
		printResult(fmt.Sprintf("[-] Error scanning %s: %v", url, err))
	} else {
		defer resp.Body.Close()
		statusCode := resp.StatusCode
		if statusCode == 200 {
			printResult(fmt.Sprintf("[+] %s - [%d]", url, statusCode))
		} else if !slices.Contains(errorCodes, statusCode) {
			printResult(fmt.Sprintf("[!] %s - [%d]", url, statusCode))
		}
	}
	atomic.AddUint32(&completed, 1)
}

func printResult(result string) {
	mu.Lock()
	defer mu.Unlock()
	// Move cursor up by one line to print the result above the progress
	fmt.Printf("\033\r\033[K%s\n", result)
	printedLines++  // Increment the number of lines printed
	printProgress() // Reprint the progress on the last line
}

func printProgress() {
	fmt.Printf("\033[%dB\r\033[K[?] Progress: %.2f%%", printedLines+1, getProgress()) // Move cursor to the last line, clear it, and print progress
}

func getProgress() float64 {
	if total == 0 {
		return 0.0
	}
	completedCount := atomic.LoadUint32(&completed)
	return (float64(completedCount) / float64(total)) * 100
}

func worker(jobs <-chan string, errorCodes []int, wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()

	for url := range jobs {
		scanner(url, errorCodes, client)
	}
}

func Scanner_init(target string, wordlist string, error_codes []int, extensions []string, file_only bool, workerCount int) {
	fmt.Println("[!] Scanning target: " + target)
	fmt.Println("[!] Wordlist: " + wordlist)
	fmt.Printf("[!] Error Codes: %v\n", error_codes)
	fmt.Printf("[!] Extensions: %v\n", extensions)

	if file_only {
		if len(extensions) == 0 {
			fmt.Println("[!] Extenions not set")
			os.Exit(1)
		}
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	jobs := make(chan string)
	wg := sync.WaitGroup{}

	// Open the wordlist file
	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Printf("Error opening wordlist file: %v\n", err)
		return
	}

	defer file.Close()

	// Count total lines in the wordlist
	file_read := bufio.NewScanner(file)
	for file_read.Scan() {
		line := strings.TrimSpace(file_read.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			total++
		}
	}

	// Reset file reader to the beginning
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Printf("Error resetting file reader: %v\n", err)
		return
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(jobs, error_codes, &wg, client)
	}

	fmt.Println("\n[!] Scanning the target")

	// Start ticker to print progress every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			printProgress()
		}
	}()

	file_read = bufio.NewScanner(file)

	go func() {
		for file_read.Scan() {
			word := strings.TrimSpace(file_read.Text())
			if word == "" || strings.HasPrefix(word, "#") {
				continue // Skip empty lines or comments
			}

			escapedWord := url.PathEscape(word)

			// If file_only is false and extensions are provided, scan both directories and files
			if !file_only && len(extensions) > 0 {
				// First, scan the directory itself
				fullURL := target + escapedWord
				jobs <- fullURL

				// Then, append each extension and scan the files
				for _, ext := range extensions {
					if !strings.HasPrefix(ext, ".") {
						ext = "." + ext // Ensure the extension starts with a dot
					}
					fileURL := target + escapedWord + ext
					jobs <- fileURL
				}
			} else {
				// If file_only is true, just scan files with extensions
				if file_only && len(extensions) > 0 {
					for _, ext := range extensions {
						if !strings.HasPrefix(ext, ".") {
							ext = "." + ext // Ensure the extension starts with a dot
						}
						fileURL := target + escapedWord + ext
						jobs <- fileURL
					}
				} else {
					// If no extensions or file_only is false, just scan the directories
					fullURL := target + escapedWord
					jobs <- fullURL
				}
			}
		}

		close(jobs)
	}()

	if err := file_read.Err(); err != nil {
		fmt.Printf("Error reading wordlist file: %v\n", err)
	}

	wg.Wait()

	fmt.Printf("\r\033[K[?] Progress: 100.00%%\n[!] Scanning complete\n") // Clear line before printing final message
}
