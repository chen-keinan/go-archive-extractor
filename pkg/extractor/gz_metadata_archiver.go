package extractor

import (
	"fmt"
	compression2 "github.com/chen-keinan/go-archive-extractor/pkg/compression"
	aerrors2 "github.com/chen-keinan/go-archive-extractor/pkg/extractor/aerrors"
	"os"
	"path/filepath"
	"time"
)

//GzMetadataArchiver object
type GzMetadataArchiver struct {
}

//Extract extract gz metadata archive
//accept gz metadata file path
//return file header metadata
func (ga GzMetadataArchiver) Extract(path string) ([]*ArchiveHeader, error) {
	headers := make([]*ArchiveHeader, 0)
	archiveFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = archiveFile.Close()
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	rc, err := compression2.CreateCompression(path).GetReader(archiveFile)
	if err != nil {
		return nil, aerrors2.New(err)
	}
	archiveHeader, err := NewArchiveHeader(rc, "metadata", time.Now().Unix(), 0)
	if err != nil {
		return nil, err
	}
	headers = append(headers, archiveHeader)
	return headers, nil
}
