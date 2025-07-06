package main

import (
	"C"
	auto_pentest "arsenal/automate/auto_pentest"
	db_miner "arsenal/auxiliary/db_miner"
	dns_cache_snooping "arsenal/auxiliary/dns_cache_snooping"
	ftp_miner "arsenal/auxiliary/ftp_miner"
	hash_cracker "arsenal/auxiliary/hash_crack"
	smb_miner "arsenal/auxiliary/smb_miner"
	directory_scanner "arsenal/scanners/directory_scanner"
	ftp_login "arsenal/scanners/ftp_login"
	port_scanner "arsenal/scanners/port_scanner"
	smb_login "arsenal/scanners/smb_login"
	vhost_scanner "arsenal/scanners/vhost_scanner"
	exploit_db "arsenal/search/exploit_db_search"
)
import (
	"log"
	"strconv"
	"unsafe"
)

//export PortScanner
func PortScanner(targetCount C.int, targetPtr **C.char, portStartPtr *C.char, portEndPtr *C.char, nmapPtr *C.char, hostDiscoveryPtr *C.char, protocolPtr *C.char) {
	target := []string{}

	if targetPtr == nil {
		log.Fatal("targetPtr is nil")
	}

	for i := 0; i < int(targetCount); i++ {
		t := C.GoString(*(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(targetPtr)) + uintptr(i)*unsafe.Sizeof(targetPtr))))
		target = append(target, t)
	}

	port_start, _ := strconv.Atoi(C.GoString(portStartPtr))
	port_end, _ := strconv.Atoi(C.GoString(portEndPtr))
	nmap, _ := strconv.ParseBool(C.GoString(nmapPtr))

	host_discovery, _ := strconv.ParseBool(C.GoString(hostDiscoveryPtr))
	protocol := C.GoString(protocolPtr)

	port_scanner.Scanner_init(target, port_start, port_end, nmap, host_discovery, protocol)
}

//export DirectoryScanner
func DirectoryScanner(targetPtr *C.char, wordlistPtr *C.char, errorCodeCount C.int, errorCodesPtr **C.char, extensionCount C.int, extensionsPtr **C.char, fileOnlyPtr *C.char, workerCount C.int) {
	target := C.GoString(targetPtr)
	wordlist := C.GoString(wordlistPtr)

	error_codes_str := []string{}

	for i := 0; i < int(errorCodeCount); i++ {
		t := C.GoString(*(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(errorCodesPtr)) + uintptr(i)*unsafe.Sizeof(errorCodesPtr))))
		error_codes_str = append(error_codes_str, t)
	}

	error_codes := []int{}

	for _, str := range error_codes_str {
		val, _ := strconv.Atoi(str)

		error_codes = append(error_codes, val)
	}

	extensions := []string{}

	for i := 0; i < int(extensionCount); i++ {
		t := C.GoString(*(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(extensionsPtr)) + uintptr(i)*unsafe.Sizeof(extensionsPtr))))
		extensions = append(extensions, t)
	}

	file_only, _ := strconv.ParseBool(C.GoString(fileOnlyPtr))

	directory_scanner.Scanner_init(target, wordlist, error_codes, extensions, file_only, int(workerCount))
}

//export VhostScanner
func VhostScanner(targetPtr *C.char, wordlistPtr *C.char, sslPtr *C.char, ignoreBodyLengthPtr *C.char, workerCount C.int) {
	target := C.GoString(targetPtr)
	wordlist := C.GoString(wordlistPtr)

	ssl, _ := strconv.ParseBool(C.GoString(sslPtr))
	ignore_body_length, _ := strconv.Atoi(C.GoString(ignoreBodyLengthPtr))

	vhost_scanner.Scanner_init(target, wordlist, ssl, ignore_body_length, int(workerCount))
}

//export HashCracker
func HashCracker(hashPtr *C.char, wordlistPtr *C.char) {
	hash := C.GoString(hashPtr)
	wordlist := C.GoString(wordlistPtr)

	hash_cracker.HashCracker_init(hash, wordlist)
}

//export DBMiner
func DBMiner(dbTypePtr *C.char, hostPtr *C.char, portPtr C.int, usernamePtr *C.char, passwordPtr *C.char, dbNamePtr *C.char) {
	dbType := C.GoString(dbTypePtr)
	host := C.GoString(hostPtr)
	port := int(portPtr)
	username := C.GoString(usernamePtr)
	password := C.GoString(passwordPtr)
	dbName := C.GoString(dbNamePtr)

	db_miner.DBMiner_Init(dbType, host, port, username, password, dbName)
}

//export DNSCacheSnooping
func DNSCacheSnooping(domainPtr *C.char, dnsServerPtr *C.char) {
	domain := C.GoString(domainPtr)
	dnsServer := C.GoString(dnsServerPtr)

	dns_cache_snooping.DNSSnooping_Init(domain, dnsServer)
}

//export FTPBruteForcer
func FTPBruteForcer(targetPtr *C.char, userFilePtr *C.char, passFilePtr *C.char, workerCountPtr C.int) {
	target := C.GoString(targetPtr)
	userFile := C.GoString(userFilePtr)
	passFile := C.GoString(passFilePtr)
	workerCount := int(workerCountPtr)

	ftp_login.FTPBruteForcer_Init(target, userFile, passFile, workerCount)
}

//export FTPMiner
func FTPMiner(serverPtr *C.char, usernamePtr *C.char, passwordPtr *C.char) {
	server := C.GoString(serverPtr)
	username := C.GoString(usernamePtr)
	password := C.GoString(passwordPtr)

	ftp_miner.FTPMiner_Init(server, username, password)
}

//export SMBBruteForce
func SMBBruteForce(targetPtr *C.char, userFilePtr *C.char, passFilePtr *C.char, workerCountPtr C.int) {
	target := C.GoString(targetPtr)
	userFile := C.GoString(userFilePtr)
	passFile := C.GoString(passFilePtr)
	workerCount := int(workerCountPtr)

	smb_login.SMBBruteForce_Init(target, userFile, passFile, workerCount)
}

//export SMBMiner
func SMBMiner(serverPtr *C.char, usernamePtr *C.char, passwordPtr *C.char) {
	server := C.GoString(serverPtr)
	username := C.GoString(usernamePtr)
	password := C.GoString(passwordPtr)

	smb_miner.SMBMiner_Init(server, username, password)
}

//export ExploitDBSearch
func ExploitDBSearch(queryPtr *C.char) {
	query := C.GoString(queryPtr)

	exploit_db.SearchExploitDB(query)
}

//export AutoPentest
func AutoPentest(targetPtr *C.char, portStartPtr *C.char, portEndPtr *C.char, nmapPtr *C.char, hostDiscoveryPtr *C.char, protocolPtr *C.char,
	dirScanWordlistPtr *C.char, errorCodeCount C.int, errorCodesPtr **C.char, extensionCount C.int, extensionsPtr **C.char, fileOnlyPtr *C.char,
	vhostScanWordlistPtr *C.char, sslPtr *C.char, ignoreBodyLengthPtr *C.char, userFilePtr *C.char, passFilePtr *C.char, workerCountPtr C.int) {
	target := C.GoString(targetPtr)
	portStart, _ := strconv.Atoi(C.GoString(portStartPtr))
	portEnd, _ := strconv.Atoi(C.GoString(portEndPtr))
	nmap, _ := strconv.ParseBool(C.GoString(nmapPtr))
	host_discovery, _ := strconv.ParseBool(C.GoString(hostDiscoveryPtr))
	protocol := C.GoString(protocolPtr)
	dirScanWordlist := C.GoString(dirScanWordlistPtr)

	error_codes_str := []string{}

	for i := 0; i < int(errorCodeCount); i++ {
		t := C.GoString(*(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(errorCodesPtr)) + uintptr(i)*unsafe.Sizeof(errorCodesPtr))))
		error_codes_str = append(error_codes_str, t)
	}

	error_codes := []int{}

	for _, str := range error_codes_str {
		val, _ := strconv.Atoi(str)

		error_codes = append(error_codes, val)
	}

	extensions := []string{}

	for i := 0; i < int(extensionCount); i++ {
		t := C.GoString(*(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(extensionsPtr)) + uintptr(i)*unsafe.Sizeof(extensionsPtr))))
		extensions = append(extensions, t)
	}

	file_only, _ := strconv.ParseBool(C.GoString(fileOnlyPtr))

	vhostScanWordlist := C.GoString(vhostScanWordlistPtr)
	ssl, _ := strconv.ParseBool(C.GoString(sslPtr))
	ignore_body_length, _ := strconv.Atoi(C.GoString(ignoreBodyLengthPtr))
	userFile := C.GoString(userFilePtr)
	passFile := C.GoString(passFilePtr)
	workerCount := int(workerCountPtr)

	auto_pentest.AutoPentest_Init(target, portStart, portEnd, nmap, host_discovery, protocol,
		dirScanWordlist, int(errorCodeCount), error_codes, int(extensionCount), extensions, file_only,
		vhostScanWordlist, ssl, ignore_body_length, userFile, passFile, workerCount)
}

func main() {

}
