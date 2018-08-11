package archive_extractor

import (
	"archive/tar"
	"github.com/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/go-archive-extractor/compression"
	"github.com/go-archive-extractor/utils"
	"io"
	"os"
)

type TarArchvier struct {
}

func (za TarArchvier) ExtractArchive(path string, processingFunc func(header *ArchiveHeader, params map[string]interface{}) error,
	params map[string]interface{}) error {
	archiveFile, err := os.Open(path)
	if err != nil {
		return archiver_errors.New(err)
	}
	defer archiveFile.Close()
	fileReader, err := compression.CreateCompression(path).GetReader(archiveFile)
	if err != nil {
		return nil
	}
	if fileReader == nil {
		fileReader = archiveFile
	}
	defer fileReader.Close()
	rc := tar.NewReader(fileReader)
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
