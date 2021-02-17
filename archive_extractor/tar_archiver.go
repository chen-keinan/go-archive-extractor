package archive_extractor

import (
	"archive/tar"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/jfrog/go-archive-extractor/compression"
	"github.com/jfrog/go-archive-extractor/utils"
	"io"
)

type TarArchiver struct{}

func (TarArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {

	cReader, err := compression.NewReader(path)
	if err != nil {
		return archiver_errors.New(err)
	}
	defer cReader.Close()

	rc := tar.NewReader(cReader)
	for {
		archiveEntry, err := rc.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if !archiveEntry.FileInfo().IsDir() && !utils.PlaceHolderFolder(archiveEntry.FileInfo().Name()) {
			archiveHeader := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.FileInfo().Size())
			err = processingFunc(archiveHeader, params)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
