package extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test7ZipArchiver(t *testing.T) {
	za := &SevenZipArchvier{}
	var headers []*ArchiveHeader
	var err error
	if headers, err = za.Extract("./fixtures/test.7z"); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	assert.Equal(t, headers[0].Name, "Interactive travel sample/media/Bali.jpg")
	assert.Equal(t, headers[0].ModTime, int64(1454061974))
	assert.Equal(t, headers[0].Size, int64(451984))
	assert.Equal(t, headers[0].Sha1, "9e4afe0b1e398b6d42fe4268ef7a18ac0762fe0d")
	assert.Equal(t, headers[0].Sha2, "6ca6536f591826dc99b1c64babddb10257bd5154def7be69eb3ad095a78d4c5f")
}
