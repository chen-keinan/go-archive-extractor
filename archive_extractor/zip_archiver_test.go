package archive_extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"strings"
)

func TestZipUnexpectedEofArchiver(t *testing.T) {
	za := &ZipArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.deb", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error() + "\n")
		assert.Equal(t, "zip: not a valid zip file", strings.Trim(err.Error(), ""))
	}
}

func TestZipArchiver(t *testing.T) {
	za := &ZipArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.zip", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "test.txt")
	assert.Equal(t, ad.ModTime, int64(1534137067))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(0))
}
