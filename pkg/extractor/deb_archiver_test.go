package extractor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDebArchiver(t *testing.T) {
	za := New(Deb)
	var headers []*ArchiveHeader
	var err error
	if headers, err = za.Extract("./fixtures/test.deb"); err != nil {
		fmt.Print(err.Error())
		t.Fatal(err)
	}
	assert.Equal(t, headers[0].Name, "debian-binary")
	assert.Equal(t, headers[0].ModTime, int64(1485714631))
	assert.Equal(t, headers[0].Size, int64(4))
	assert.Equal(t, headers[0].Sha1, "7959c969e092f2a5a8604e2287807ac5b1b384ad")
	assert.Equal(t, headers[0].Sha2, "d526eb4e878a23ef26ae190031b4efd2d58ed66789ac049ea3dbaf74c9df7402")
}
