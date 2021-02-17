package archive_extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"strings"
)

func TestTarUnexpectedEofArchiver(t *testing.T) {
	za := &TarArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.deb", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error() + "\n")
		assert.Equal(t, "archive/tar: invalid tar header", strings.Trim(err.Error(), ""))
	}
}

func TestTarArchiver(t *testing.T) {
	za := &TarArchiver{}
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
}

func TestTarArchiver_Lzma(t *testing.T) {
	za := &TarArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/junit.tar.lzma", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, "junit-4.12.jar", ad.Name)
	assert.Equal(t, int64(1534397548), ad.ModTime)
	assert.False(t, ad.IsFolder)
	assert.Equal(t, int64(314932), ad.Size)
}
