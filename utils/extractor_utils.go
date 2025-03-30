package utils

import (
	"crypto/sha1"   // nolint:gosec
	"crypto/sha256" // nolint:gosec
	"fmt"
	"strings"
)

const (
	//FolderSuffix const test
	FolderSuffix string = "/"
)

//PlaceHolderFolder check if folder include special sign
//accept folder path
func PlaceHolderFolder(path string) bool {
	nameParts := strings.Split(path, "/")
	if len(nameParts) > 0 {
		return nameParts[len(nameParts)-1] == "-"
	}
	return false
}

//IsFolder check if path is folder
//accept file path
//return bool if path is folder
func IsFolder(path string) bool {
	return strings.HasSuffix(path, FolderSuffix)
}

// NewSHA2 NewSH2 calculate file sha256
// accept file byte
//return sha256 string
// nolint:gosec
func NewSHA2(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

// NewSHA1 NewSH2 calculate file sha1
// accept file byte
//return sha1 string
// nolint:gosec
func NewSHA1(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))

}
