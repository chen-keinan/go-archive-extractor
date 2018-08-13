package archive_extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"io/ioutil"
)

func TestRpmArchiver(t *testing.T) {
	za := &RpmArchvier{}
	funcParams := params()
	if err := za.ExtractArchive("./fixtures/test.rpm", processingFunc, funcParams); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	ad := funcParams["archiveData"].(*ArchiveData)
	assert.Equal(t, ad.Name, "./usr/share/doc/php-zstd-devel/tests/info.phpt")
	assert.Equal(t, ad.ModTime, int64(1517299253))
	assert.Equal(t, ad.IsFolder, false)
	assert.Equal(t, ad.Size, int64(183))
	b, err := ioutil.ReadAll(ad.ArchiveReader)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(b), "")
}
