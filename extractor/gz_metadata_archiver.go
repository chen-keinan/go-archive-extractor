package extractor

import (
	"fmt"
	"github.com/chen-keinan/go-archive-extractor/compression"
	"github.com/chen-keinan/go-archive-extractor/extractor/archiver_errors"
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
func (ga GzMetadataArchiver) ExtractArchive(path string) ([]*ArchiveHeader, error) {
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
	rc, err := compression.CreateCompression(path).GetReader(archiveFile)
	if err != nil {
		return nil, archiver_errors.New(err)
	}
	archiveHeader, err := NewArchiveHeader(rc, "metadata", time.Now().Unix(), 0)
	if err != nil {
		return nil, err
	}
	headers = append(headers, archiveHeader)
	return headers, nil
}
