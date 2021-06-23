package archive_extractor

import (
	"errors"
	"fmt"
	"github.com/chen-keinan/go-archive-extractor/utils"
	archive "github.com/gen2brain/go-unarr"
	"io"
	"path/filepath"
)

type SevenZipArchvier struct {
}

//Extract extract 7zip archive
//accept 7zip file path
//return file header metadata
func (za SevenZipArchvier) Extract(path string) ([]*ArchiveHeader, error) {
	headers := make([]*ArchiveHeader, 0)
	r, err := archive.NewArchive(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	allFiles, err := r.List()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Close()
		if err != nil {
			fmt.Print(err.Error())
		}
	}()
	if err == nil {
		for _, archiveEntry := range allFiles {
			err := r.EntryFor(archiveEntry)
			if err != nil {
				return nil, err
			}
			if !utils.IsFolder(archiveEntry) {
				rc := &SevenZipReader{Archive: r, Size: r.Size()}
				archiveHeader, err := NewArchiveHeader(rc, r.Name(), r.ModTime().Unix(), int64(r.Size()))
				if err != nil {
					return nil, err
				}
				headers = append(headers, archiveHeader)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return headers, nil
}

type SevenZipReader struct {
	Archive *archive.Archive
	Size    int
}

func (a *SevenZipReader) Read(p []byte) (n int, err error) {
	if a.Size <= 0 {
		return 0, io.EOF
	}
	size := len(p)
	if len(p) > a.Size {
		size = a.Size
	}
	b := make([]byte, size)
	n, err = a.Archive.Read(b)
	if err != nil && err != io.EOF {
		return 0, err
	}
	copied := copy(p, b)
	if copied != n {
		return 0, errors.New(fmt.Sprintf("copy arrays failed, copied only %v from %v bytes", copied, n))
	}
	a.Size -= n
	return n, nil
}
