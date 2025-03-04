package utils

import (
	"path/filepath"
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

// In Windows, filepath.Clean operation will replace all slashes '/'
// to backslashes '\\'
// This can mess-up with the code that makes path comparisons - in indexer-app on Windows
func CleanPathKeepingUnixSlash(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

func JoinPathKeepingUnixSlash(elem ...string) string {
	return filepath.ToSlash(filepath.Join(elem...))
}
