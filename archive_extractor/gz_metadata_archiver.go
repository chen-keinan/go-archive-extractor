package archive_extractor

import (
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/jfrog/go-archive-extractor/compression"
	"os"
	"time"
)

type GzMetadataArchiver struct {
}

func (ga GzMetadataArchiver) ExtractArchive(path string, processingFunc func(header *ArchiveHeader, params map[string]interface{}) error,
	params map[string]interface{}) error {
	archiveFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer archiveFile.Close()
	rc, err := compression.CreateCompression(path).GetReader(archiveFile)
	if err != nil {
		return archiver_errors.New(err)
	}
	archiveHeader := NewArchiveHeader(rc, "metadata", time.Now().Unix(), 0)
	err = processingFunc(archiveHeader, params)
	if err != nil {
		return err
	}
	return nil
}
