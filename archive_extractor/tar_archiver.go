package archive_extractor

import (
	"archive/tar"
	"fmt"
	"github.com/chen-keinan/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/chen-keinan/go-archive-extractor/compression"
	"github.com/chen-keinan/go-archive-extractor/utils"
	"io"
	"os"
)

type TarArchvier struct {
}

func (za TarArchvier) Extract(path string) ([]*ArchiveHeader, error) {
	headers := make([]*ArchiveHeader, 0)
	archiveFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer archiveFile.Close()
	fileReader, err := compression.CreateCompression(path).GetReader(archiveFile)
	if err != nil {
		return nil, archiver_errors.New(err)
	}
	if fileReader == nil {
		fileReader = archiveFile
	}
	defer func() {
		err := fileReader.Close()
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	rc := tar.NewReader(fileReader)
	for {
		archiveEntry, err := rc.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !archiveEntry.FileInfo().IsDir() && !utils.PlaceHolderFolder(archiveEntry.FileInfo().Name()) {
			archiveHeader, err := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.FileInfo().Size())
			if err != nil {
				return nil, err
			}
			headers = append(headers, archiveHeader)
		}
	}
	return headers, nil
}
