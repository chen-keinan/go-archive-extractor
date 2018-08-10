package utils

import "strings"

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
