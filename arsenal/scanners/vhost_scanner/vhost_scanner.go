// VHost Scanner
package scanners

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
- Target IP or Hostname including "http or https"
- Wordlist to use (Default: /usr/share/seclists/Discovery/DNS/subdomains-top1million-20000.txt)
- HTTP Codes (Default: 2XX, 3XX)
*/

var (
	mu           sync.Mutex // Mutex to coordinate printing
	total        int        // Total number of domains to scan
	completed    uint32     // Atomic counter for completed scans
	printedLines int        // Track the number of lines printed for results
)

func scanner(target string, basename string, ssl bool, ignore_length int, client *http.Client) {
	var host string

	if ssl {
		host = basename + "." + target[8:]
	} else {
		host = basename + "." + target[7:]
	}

	host = host[:len(host)-1]

	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		printResult(fmt.Sprintf("[-] Error creating request for %s: %v", target, err))
		return
	}

	req.Host = host

	resp, err := client.Do(req)
	if err != nil {
		printResult(fmt.Sprintf("[-] Error scanning %s: %v", host, err))
		return
	}
	defer resp.Body.Close()

	status_code := resp.StatusCode
	body, _ := io.ReadAll(resp.Body)

	if status_code == 200 && len(body) != ignore_length {
		printResult(fmt.Sprintf("[+] %s - Status: [%d] Body: [%d]", host, status_code, len(body)))
	}

	atomic.AddUint32(&completed, 1)
}

func printResult(result string) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Printf("\033[%d\r\033[K%s\n", printedLines+1, result)
	printedLines++
	printProgress()
}

func printProgress() {
	fmt.Printf("\033[%dB\r\033[K[?] Progress: %.2f%%", printedLines+1, getProgress())
}

func getProgress() float64 {
	if total == 0 {
		return 0.0
	}
	completedCount := atomic.LoadUint32(&completed)
	return (float64(completedCount) / float64(total)) * 100
}

func worker(ctx context.Context, jobs <-chan string, target string, ssl bool, ignore_length int, wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case domain, ok := <-jobs:
			if !ok {
				return
			}
			scanner(target, domain, ssl, ignore_length, client)
		}
	}
}

func Scanner_init(target string, wordlist string, ssl bool, ignore_length int, workerCount int) {
	fmt.Println("[!] Scanning target: " + target)
	fmt.Println("[!] Wordlist: " + wordlist)
	fmt.Printf("[!] SSL: %s\n", strconv.FormatBool(ssl))
	fmt.Printf("[!] Ignore Content Length: %d\n", ignore_length)

	file, err := os.Open(wordlist)
	if err != nil {
		log.Fatal("[X] Unable to Read Wordlist")
	}
	defer file.Close()

	// Set up channel to handle interrupt signals
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file_read := bufio.NewScanner(file)
	jobs := make(chan string)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if ssl {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	// Count total lines in the wordlist
	for file_read.Scan() {
		word := strings.TrimSpace(file_read.Text())
		if word != "" && !strings.HasPrefix(word, "#") {
			total++
		}
	}

	fmt.Printf("[!] Total: %d\n", total)

	// Reset file reader to the beginning
	if _, err := file.Seek(0, 0); err != nil {
		log.Fatal("Error resetting file reader:", err)
	}

	// Reinitialize the file scanner after resetting the file reader
	file_read = bufio.NewScanner(file) // Reset the scanner

	wg := sync.WaitGroup{}
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(ctx, jobs, target, ssl, ignore_length, &wg, client)
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

	go func() {
		for file_read.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				word := strings.TrimSpace(file_read.Text())
				if word != "" && !strings.HasPrefix(word, "#") {
					jobs <- word
				}
			}
		}
		close(jobs)
	}()

	if err := file_read.Err(); err != nil {
		log.Fatal("Error reading wordlist file:", err)
	}

	// Wait for interrupt signal
	go func() {
		<-interruptChan
		fmt.Println("\n[!] Interrupt received, stopping scan...")
		// Stop the scanner gracefully (e.g., by closing channels or using a flag)
		cancel() // Singla cancellation to workers
	}()

	wg.Wait()
	fmt.Printf("\r\033[K[?] Progress: 100.00%%\n[!] Scanning complete\n")
}
