package auxiliary

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hirochachacha/go-smb2"
)

const maxFileSize = 10 * 1024 * 1024 // 10MB

func downloadFilesUnder10MB(share *smb2.Share, dir string) {
	entries, err := share.ReadDir(dir)
	if err != nil {
		log.Fatalf("[-] Failed to read directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			downloadFilesUnder10MB(share, filepath.Join(dir, entry.Name()))
		} else if entry.Size() <= maxFileSize {
			downloadFile(share, filepath.Join(dir, entry.Name()))
		}
	}
}

func downloadFile(share *smb2.Share, path string) {
	fmt.Printf("[!] Downloading %s...\n", path)
	file, err := share.Open(path)
	if err != nil {
		log.Printf("[!] Failed to open file: %v", err)
		return
	}
	defer file.Close()

	localPath := filepath.Join("smb_downloads", path)
	os.MkdirAll(filepath.Dir(localPath), os.ModePerm)

	localFile, err := os.Create(localPath)
	if err != nil {
		log.Printf("[-] Failed to create local file: %v", err)
		return
	}
	defer localFile.Close()

	_, err = file.WriteTo(localFile)
	if err != nil {
		log.Printf("[-] Failed to download file: %v", err)
	} else {
		fmt.Printf("[+] Downloaded: %s\n", localPath)
	}
}

func SMBMiner_Init(server string, username string, password string) {
	server = server + ":445"

	conn, err := net.DialTimeout("tcp", server, 5*time.Second)
	if err != nil {
		log.Fatalf("[-] Failed to connect to server: %v", err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     username,
			Password: password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		log.Fatalf("[-] Failed to initiate SMB session: %v", err)
	}
	defer s.Logoff()

	shares, _ := s.ListSharenames()

	for _, shareName := range shares {
		share, err := s.Mount(shareName)
		fmt.Printf("\n[!] Checking Share: %s\n", shareName)
		if err != nil {
			fmt.Printf("\n[-] Failed to mount SMB share: %s", shareName)
			continue
		}
		defer share.Umount()

		// Start extracting files less than 10MB
		downloadFilesUnder10MB(share, ".")

		fmt.Println("[!] File extraction complete!")
	}
}
