package extractor

import (
	"archive/tar"
	"fmt"
	"github.com/chen-keinan/go-archive-extractor/compression"
	"github.com/chen-keinan/go-archive-extractor/extractor/aerrors"
	"github.com/chen-keinan/go-archive-extractor/utils"
	"io"
	"os"
	"path/filepath"
)

//TarArchvier object
type TarArchvier struct {
}

//Extract extract tar archive
//accept tar file path
//return file header metadata
func (za TarArchvier) Extract(path string) ([]*ArchiveHeader, error) {
	headers := make([]*ArchiveHeader, 0)
	archiveFile, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = archiveFile.Close(); err != nil {
			fmt.Print(err.Error())
		}
	}()
	fileReader, err := compression.CreateCompression(path).GetReader(archiveFile)
	if err != nil {
		return nil, aerrors.New(err)
	}
	if fileReader == nil {
		fileReader = archiveFile
	}
	defer func() {
		err = fileReader.Close()
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	rc := tar.NewReader(fileReader)
	headers, archiveHeaders, err := za.extractHeaders(rc, headers)
	if err != nil {
		return archiveHeaders, err
	}
	return headers, nil
}

func (za TarArchvier) extractHeaders(rc *tar.Reader, headers []*ArchiveHeader) ([]*ArchiveHeader, []*ArchiveHeader, error) {
	for {
		archiveEntry, err := rc.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if !archiveEntry.FileInfo().IsDir() && !utils.PlaceHolderFolder(archiveEntry.FileInfo().Name()) {
			archiveHeader, err := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.FileInfo().Size())
			if err != nil {
				return nil, nil, err
			}
			headers = append(headers, archiveHeader)
		}
	}
	return headers, nil, nil
}
