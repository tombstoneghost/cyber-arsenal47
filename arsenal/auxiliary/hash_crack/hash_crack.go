// Hash Cracker
package auxiliary

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf16"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"
)

// Convert string to UTF-16LE encoded bytes
func utf16le(data string) []byte {
	encoded := utf16.Encode([]rune(data))
	buf := make([]byte, len(encoded)*2)
	for i, v := range encoded {
		binary.LittleEndian.PutUint16(buf[i*2:], v)
	}
	return buf
}

func MD5Hash(data string) []byte {
	hash := md5.Sum([]byte(data))
	return hash[:]
}

func SHA1Hash(data string) []byte {
	hash := sha1.Sum([]byte(data))
	return hash[:]
}

func SHA256Hash(data string) []byte {
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

func SHA512Hash(data string) []byte {
	hash := sha512.Sum512([]byte(data))
	return hash[:]
}

func BcryptHash(data string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
}

func Blake2bHash(data string) ([]byte, error) {
	hash, err := blake2b.New256(nil)
	if err != nil {
		return nil, err
	}
	hash.Write([]byte(data))
	return hash.Sum(nil), nil
}

func Blake2sHash(data string) ([]byte, error) {
	hash, err := blake2s.New256(nil)
	if err != nil {
		return nil, err
	}
	hash.Write([]byte(data))
	return hash.Sum(nil), nil
}

func NTLMHash(data string) []byte {
	md4 := md4.New()
	md4.Write(utf16le(data))
	return md4.Sum(nil)
}

// Identify hash based on length and pattern
func identifyHashType(hash string) string {
	switch len(hash) {
	case 32:
		if regexp.MustCompile("^[a-f0-9]{32}$").MatchString(hash) {
			return "md5"
		}
	case 40:
		if regexp.MustCompile("^[a-f0-9]{40}$").MatchString(hash) {
			return "sha1"
		}
	case 64:
		if regexp.MustCompile("^[a-f0-9]{64}$").MatchString(hash) {
			return "sha256"
		}
	case 128:
		if regexp.MustCompile("^[a-f0-9]{128}$").MatchString(hash) {
			return "sha512"
		}
	case 60: // BCrypt always has 60 chars
		if strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") {
			return "bcrypt"
		}
	default:
		if strings.HasPrefix(hash, "$blake2b$") {
			return "blake2b"
		}
		if strings.HasPrefix(hash, "$blake2s$") {
			return "blake2s"
		}
		if len(hash) == 32 && regexp.MustCompile("^[A-F0-9]{32}$").MatchString(hash) {
			return "ntlm"
		}
	}
	return "unknown"
}

func HashCracker_init(hash string, wordlist string) {
	fmt.Println("[!] Hash Provided: " + hash)

	h_type := identifyHashType(hash)

	fmt.Println("[!] Hash Type Identified: " + h_type)

	if strings.Compare(h_type, "unknown") == 0 {
		fmt.Println("[X] Unable to identify supported hash")
		return
	}

	h_type = strings.ToUpper(h_type)

	fmt.Println("[!] Wordlist: " + wordlist)

	file, err := os.Open(wordlist)
	if err != nil {
		log.Fatal("[X] Unable to Read Wordlist")
	}
	defer file.Close()

	file_read := bufio.NewScanner(file)

	var curr_hash_raw []byte
	var curr_hash string

	for file_read.Scan() {
		word := strings.TrimSpace(file_read.Text())
		switch h_type {
		case "MD5":
			curr_hash_raw = MD5Hash(word)
		case "SHA1":
			curr_hash_raw = SHA1Hash(word)
		case "SHA256":
			curr_hash_raw = SHA256Hash(word)
		case "SHA512":
			curr_hash_raw = SHA512Hash(word)
		case "BCrypt":
			curr_hash_raw, _ = BcryptHash(word)
		case "Blake2b":
			curr_hash_raw, _ = Blake2bHash(word)
		case "Blake2s":
			curr_hash_raw, _ = Blake2sHash(word)
		case "NTLM":
			curr_hash_raw = NTLMHash(word)
		}

		curr_hash = hex.EncodeToString(curr_hash_raw)

		if curr_hash != hash {
			continue
		} else {
			fmt.Println("[+] Hash Cracked! Plaintext: " + word)
		}
	}
}
