package extractor

import (
	"archive/zip"
	"fmt"
	"path/filepath"
)

//ZipArchvier object
type ZipArchvier struct {
}

//Extract extract zip archive
//accept zip file path
//return file header metadata
func (za ZipArchvier) Extract(path string) ([]*ArchiveHeader, error) {
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
			errClose := rc.Close()
			if err != nil {
				fmt.Print(errClose.Error())
			}
			return nil, err
		}
		archiveHeader, err := NewArchiveHeader(rc, archiveEntry.Name, archiveEntry.ModTime().Unix(), archiveEntry.FileInfo().Size())
		if err != nil {
			errClose := rc.Close()
			if errClose != nil {

			}
			return nil, err
		}
		err = rc.Close()
		if err != nil {
			return nil, err
		}
		headers = append(headers, archiveHeader)
	}
	return headers, nil
}
