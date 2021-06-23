package extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"strings"
)

func TestTarUnexpectedEofArchiver(t *testing.T) {
	za := &TarArchvier{}
	var err error
	if _, err = za.Extract("./fixtures/test.deb"); err != nil {
		fmt.Print(err.Error() + "\n")
		assert.Equal(t, "archive/tar: invalid tar header", strings.Trim(err.Error(), ""))
	}
}

func TestTarArchiver(t *testing.T) {
	za := &TarArchvier{}
	var headers []*ArchiveHeader
	var err error
	if headers, err = za.Extract("./fixtures/test.tar.gz"); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	assert.Equal(t, headers[2].Name, "logRotator-1.0/log_rotator.go")
	assert.Equal(t, headers[2].ModTime, int64(1531307652))
	assert.Equal(t, headers[2].Size, int64(3685))
	assert.Equal(t, headers[2].Sha1, "f0509a5f7bf4ebbd8c69eee4e792391613ae3bf6")
	assert.Equal(t, headers[2].Sha2, "2b853b8a4021318cedabe7d3220d24786a48b936c17d7164ad56a3f8112c8d89")
}
