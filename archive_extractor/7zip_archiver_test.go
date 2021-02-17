package archive_extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test7ZipArchiver(t *testing.T) {
	za := &SevenZipArchiver{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.7z", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "Interactive travel sample/.spxproperties")
	assert.Equal(t, ad.ModTime, int64(-11644473600))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(44))
}
