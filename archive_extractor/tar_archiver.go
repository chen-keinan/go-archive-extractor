package archive_extractor

import (
	"archive/tar"
	"github.com/go-archive-extractor/compression"
	"github.com/go-archive-extractor/utils"
	"io"
	"os"
)

type TarArchvier struct {
}

func (za TarArchvier) ExtractArchive(archivePath string, advanceProcessing func(header *ArchiveHeader, advanceProcessingParams map[string]interface{}) error,
	advanceProcessingParams map[string]interface{}) error { // recover handling in case of failure during indexing process
	archiveFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archiveFile.Close()
	fileReader, err := compression.CreateCompression(archivePath).GetReader(archiveFile)
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
			err = advanceProcessing(archiveHeader, advanceProcessingParams)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
