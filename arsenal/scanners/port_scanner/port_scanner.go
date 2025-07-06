// Port Scanner Module
package scanners

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Ullaakut/nmap"
)

/*
Inputs:
- Target IPs (Command Seprated)
- Port Range (Default: 1 to 65535)
- Nmap Scan (Default: No)
- Host Discovery (Nmap -Pn option | Default: No)
- Nmap Flags (Default: -A)
- Protocol (Default: TCP)
*/

// Common Ports
var common_ports = []int{20, 21, 22, 23, 25, 53, 69, 80, 88, 110, 135, 137, 139, 143, 443, 587, 636, 1025, 1194, 2082, 2483, 2484, 3306, 3389, 27017}

// Helper Function
func contains(num int) bool {
	for _, val := range common_ports {
		if num == val {
			return true
		}
	}

	return false
}

func makeRange(start, end int) []int {
	r := []int{}

	for i := start; i <= end; i++ {
		if !contains(i) {
			r = append(r, i)
		}
	}

	return r
}

// Scanner
func scanner(target string, port int, protocol string, result_channel chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := target + ":" + strconv.Itoa(port)

	conn, err := net.Dial(protocol, addr)

	if err == nil {
		conn.Close()

		fmt.Println("[+] Port " + strconv.Itoa(port) + " is Open")
		result_channel <- port
	}
}

// Do Nmap Scan
func run_nmap(targets map[string][]int, host_discovery bool) map[int]string {
	nmap_flags := []string{"--version-all", "-v", "-A", "-T3", "-sV", "-O"}
	nmap_result := make(map[int]string)

	for host, ports := range targets {
		fmt.Printf("[!] Running Nmap for %s", host)

		// Convert []int ports to []string
		stringPorts := make([]string, len(ports))
		for i, v := range ports {
			stringPorts[i] = strconv.Itoa(v) // Convert int to string
		}

		if host_discovery {
			nmap_flags = append(nmap_flags, "-Pn")
		}

		scan, err := nmap.NewScanner(nmap.WithTargets(host), nmap.WithPorts(stringPorts...), nmap.WithServiceInfo(), nmap.WithAggressiveScan())

		if err != nil {
			fmt.Println("[-] Failed to initialize Nmap scanner:", err)
			continue
		}

		// Execute the scan and get the result
		result, warnings, err := scan.Run()
		if err != nil {
			fmt.Println("[-] Nmap scan failed:", err)
			continue
		}

		if len(warnings) > 0 {
			fmt.Println("Warnings:", warnings)
		}

		// Parse and display the scan results
		for _, host := range result.Hosts {
			if len(host.Ports) == 0 {
				fmt.Printf("[!] No open ports found on host %s\n", host.Addresses[0])
				continue
			}

			fmt.Printf("[+] Host %s - Open Ports:\n", host.Addresses[0])
			for _, port := range host.Ports {
				service := port.Service
				fmt.Printf("  - Port: %d/%s | Service: %s | Product: %s | Version: %s | Extra Info: %s\n",
					port.ID, port.Protocol, service.Name, service.Product, service.Version, service.ExtraInfo)
				nmap_result[int(port.ID)] = fmt.Sprintf("%s %s", service.Product, service.Version)
			}
		}
	}

	return nmap_result
}

// Initialize Scanner
func Scanner_init(target []string, port_start int, port_end int, nmap bool, host_discovery bool, protocol string) (map[string][]int, map[int]string) {
	fmt.Printf("[!] Provided Target(s): %v\n", target)

	scan_results := make(map[string][]int)

	result_channel := make(chan int)

	var wg sync.WaitGroup

	for _, ip := range target {
		open_ports := []int{}

		fmt.Println("[!] Scanning target: " + ip)

		for _, port := range common_ports {
			wg.Add(1)

			go scanner(ip, port, protocol, result_channel, &wg)
		}

		ports := makeRange(port_start, port_end)

		for _, port := range ports {
			wg.Add(1)
			go scanner(ip, port, protocol, result_channel, &wg)
		}

		go func() {
			wg.Wait()
			close(result_channel)
		}()

		for result := range result_channel {
			open_ports = append(open_ports, result)
		}

		sort.Ints(open_ports)

		scan_results[ip] = open_ports
	}

	for k, v := range scan_results {
		fmt.Printf("\n[!] Scan results for %s : %v\n", k, v)
	}

	nmap_results := make(map[int]string)

	if nmap {
		fmt.Println("[!] Initializing Nmap Scanner")
		time.Sleep(5 * time.Second)

		nmap_results = run_nmap(scan_results, host_discovery)
	}

	wg.Wait()

	return scan_results, nmap_results

}
