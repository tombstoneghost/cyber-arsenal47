package dns

import (
	"context"
	"fmt"
	"net"
)

// DNS Cache Snooping module
func DNSSnooping_Init(domain string, dnsServer string) {
	// Set up DNS request with context
	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("udp", dnsServer+":53")
		},
	}

	// Query for the domain
	_, err := resolver.LookupHost(context.Background(), domain)
	if err != nil {
		fmt.Printf("[-] DNS Cache Snooping Failed: %v\n", err)
	} else {
		fmt.Printf("[+] %s is cached on the server %s\n", domain, dnsServer)
	}
}
