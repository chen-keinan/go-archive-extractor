package archive_extractor

import (
	"archive/zip"
	"fmt"
	"path/filepath"
)

type ZipArchvier struct {
}

func (za ZipArchvier) ExtractArchive(path string) ([]*ArchiveHeader, error) {
	headers := make([]*ArchiveHeader, 0)
	r, err := zip.OpenReader(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Close()
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	for _, archiveEntry := range r.File {
		rc, err := archiveEntry.Open()
		if err != nil {
			rc.Close()
			return nil, err
		}
		archiveHeader, err := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime().Unix(), archiveEntry.FileInfo().Size())
		if err != nil {
			rc.Close()
			return nil, err
		}
		rc.Close()
		headers = append(headers, archiveHeader)
	}
	return headers, nil
}
