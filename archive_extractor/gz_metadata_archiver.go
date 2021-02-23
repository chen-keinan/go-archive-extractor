package archive_extractor

import (
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/jfrog/go-archive-extractor/compression"
	"time"
)

type GzMetadataArchiver struct{}

func (ga GzMetadataArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {

	cReader, err := compression.NewReader(path)
	if compression.IsGetReaderError(err) {
		return archiver_errors.New(err)
	}
	if err != nil {
		return err
	}
	defer cReader.Close()

	archiveHeader := NewArchiveHeader(cReader, "metadata", time.Now().Unix(), 0)
	err = processingFunc(archiveHeader, params)
	if err != nil {
		return err
	}
	return nil
}
