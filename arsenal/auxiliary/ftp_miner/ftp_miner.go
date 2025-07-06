package auxiliary

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

func LoginAndListFiles(conn *ftp.ServerConn) {
	// List files and directories
	entries, err := conn.List("/")
	if err != nil {
		fmt.Println("[X] Error listing files:", err)
		return
	}

	fmt.Println("[>] Files and directories:")
	for _, entry := range entries {
		fmt.Println(" -", entry.Name)
	}

	fmt.Print("[?] Downloading all files")

	// Create a folder to store downloaded files
	folderName := "ftp_downloads"
	err = os.Mkdir(folderName, 0755)
	if err != nil {
		fmt.Println("[X] Error creating directory:", err)
		return
	} else {
		fmt.Println("[!] Created directory: ftp_downloads")
	}

	// Download all files
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFile {
			fmt.Printf("[*] Downloading file: %s\n", entry.Name)
			err = downloadFile(conn, entry.Name, folderName)
			if err != nil {
				fmt.Println("[X] Error downloading file:", err)
			}
		}
	}
	fmt.Println("[!] Download complete")
}

func downloadFile(conn *ftp.ServerConn, fileName, folderName string) error {
	// Retrieve the file from FTP
	resp, err := conn.Retr(fileName)
	if err != nil {
		return fmt.Errorf("error retrieving file: %s", err)
	}
	defer resp.Close()

	// Read the file data
	fileData, err := io.ReadAll(resp)
	if err != nil {
		return fmt.Errorf("error reading file data: %s", err)
	}

	// Create the file locally
	localFilePath := filepath.Join(folderName, fileName)
	err = os.WriteFile(localFilePath, fileData, 0644)
	if err != nil {
		return fmt.Errorf("error saving file locally: %s", err)
	}

	return nil
}

func FTPMiner_Init(server string, username string, password string) {
	// Connect to the FTP server
	server = server + ":21"

	conn, err := ftp.Dial(server, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		fmt.Println("[X] Error connecting to FTP server:", err)
		return
	}
	defer conn.Quit()

	// Login to FTP
	err = conn.Login(username, password)
	if err != nil {
		fmt.Println("[X] Error logging into FTP server:", err)
		return
	}
	fmt.Println("[!] Successfully logged into FTP server")

	LoginAndListFiles(conn)

	fmt.Println("[!] FTP Miner ran successfully")
}
