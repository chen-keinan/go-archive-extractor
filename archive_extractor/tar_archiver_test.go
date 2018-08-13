package archive_extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"io/ioutil"
	"strings"
)

func TestTarUnexpectedEofArchiver(t *testing.T) {
	za := &TarArchvier{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.deb", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error() + "\n")
		assert.Equal(t, "archive/tar: invalid tar header", strings.Trim(err.Error(), ""))
	}
}

func TestTarArchiver(t *testing.T) {
	za := &TarArchvier{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.tar.gz", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "logRotator-1.0/log_rotator.go")
	assert.Equal(t, ad.ModTime, int64(1531307652))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(3685))
	b, err := ioutil.ReadAll(ad.ArchiveReader)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(b), "")
}
