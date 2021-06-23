package extractor

import (
	"strings"
	"testing"
)

func TestNewArchiveHeader(t *testing.T) {
	header, err := NewArchiveHeader(strings.NewReader(""), "aaa", int64(1), int64(2))
	if err != nil {
		t.Fatal("failed to create new archive header")
	}
	if header.Sha1 != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		t.Fatal("sha1 do not match")
	}
	if header.Sha2 != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Fatal("sha2 do not match")
	}
}
