package archive_extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRpmArchiver(t *testing.T) {
	za := &RpmArchvier{}
	var headers []*ArchiveHeader
	var err error
	if headers, err = za.Extract("./fixtures/test.rpm"); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	assert.Equal(t, headers[18].Name, "./usr/share/doc/php-zstd-devel/tests/info.phpt")
	assert.Equal(t, headers[18].ModTime, int64(1517299253))
	assert.Equal(t, headers[18].Size, int64(183))
	assert.Equal(t, headers[18].Sha1, "5cd9bec6e67e4cc056abb3f507411dcef6133d09")
	assert.Equal(t, headers[18].Sha2, "7af6145ba6e6bc1da63404487e8564726f272ebb59b3a16f343234d552ecacce")
	assert.Equal(t, headers[18].PkgMeta["Release"].(string), "1.fc24.remi.7.0")
	assert.Equal(t, headers[18].PkgMeta["Version"].(string), "0.4.11")
	assert.Equal(t, headers[18].PkgMeta["Name"].(string), "php-zstd-devel")
}
