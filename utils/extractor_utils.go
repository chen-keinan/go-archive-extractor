package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"strings"
)

const (
	FolderSuffix string = "/"
)

func PlaceHolderFolder(path string) bool {
	nameParts := strings.Split(path, "/")
	if len(nameParts) > 0 {
		return nameParts[len(nameParts)-1] == "-"
	} else {
		return false
	}
}

func IsFolder(path string) bool {
	return strings.HasSuffix(path, FolderSuffix)
}

// NewSH2 ...
func NewSHA2(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

// NewSHA1 ...
func NewSHA1(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))

}
