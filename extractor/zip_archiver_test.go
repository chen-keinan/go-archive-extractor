package extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"strings"
)

func TestZipUnexpectedEofArchiver(t *testing.T) {
	za := &ZipArchvier{}
	if _, err := za.Extract("./fixtures/test.deb"); err != nil {
		fmt.Print(err.Error() + "\n")
		assert.Equal(t, "zip: not a valid zip file", strings.Trim(err.Error(), ""))
	}
}

func TestZipArchiver(t *testing.T) {
	za := &ZipArchvier{}
	var headers []*ArchiveHeader
	var err error
	if headers, err = za.Extract("./fixtures/test.zip"); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	assert.Equal(t, headers[0].Name, "test.txt")
	assert.Equal(t, headers[0].ModTime, int64(1534137067000000000))
	assert.Equal(t, headers[0].Size, int64(0))
	assert.Equal(t, headers[0].Sha1, "da39a3ee5e6b4b0d3255bfef95601890afd80709")
	assert.Equal(t, headers[0].Sha2, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
}
