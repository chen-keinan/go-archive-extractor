package archive_extractor

import (
	"github.com/cavaliercoder/go-cpio"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/jfrog/go-archive-extractor/compression"
	"github.com/jfrog/go-rpm/v2"
	"io"
	"math"
	"os"
)

type RpmArchvier struct {
}

func (za RpmArchvier) ExtractArchive(path string, processingFunc func(header *ArchiveHeader, params map[string]interface{}) error, params map[string]interface{}) error {
	rpmFile, err := rpm.OpenPackageFile(path)
	if err != nil {
		return archiver_errors.New(err)
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	//Read content of cpio archive which starts after headers
	headerEnd := za.getHeadersEnd(rpmFile.Headers)
	archiveHead := make([]byte, 6)

	_, err = file.ReadAt(archiveHead, int64(headerEnd))
	if err != nil {
		return err
	}

	//rewind to start of the file
	file.Seek(int64(headerEnd), 0)
	fileReader, err := compression.CreateCompressionFromBytes(archiveHead).GetReader(file)
	if err != nil || fileReader == nil {
		return archiver_errors.New(err)
	}
	defer fileReader.Close()
	err = za.readRpm(processingFunc, params, rpmFile, fileReader)
	if err != nil {
		return archiver_errors.New(err)
	}

	return nil
}

func (za RpmArchvier) getHeadersEnd(headers []rpm.Header) uint64 {
	var end uint64
	offset := 96
	for i := 0; i < 2; i++ {
		h := headers[i]
		// set start and end offsets
		start := offset
		end = uint64(start + 16 + (16 * h.IndexCount) + h.Length)
		offset = int(end)
		// calculate location of the end of the header by padding to a multiple of 8
		pad := 8 - int(math.Mod(float64(h.Length), 8))
		if pad < 8 {
			offset += pad
		}
	}
	return end
}

func (za RpmArchvier) readRpm(processingFunc func(header *ArchiveHeader, params map[string]interface{}) error,
	params map[string]interface{}, rpmFile *rpm.PackageFile, fileReader io.Reader) error {
	//create cpio reader
	cpioReader := cpio.NewReader(fileReader)
	// Parse the rpm
	var count = 0
	for {
		archiveEntry, err := cpioReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		count++
		if !archiveEntry.Mode.IsDir() {
			if archiveEntry != nil {
				archiveHeader := NewArchiveHeader(cpioReader, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.Size)
				err = processingFunc(archiveHeader, params)
				if _, ok := params["rpmPkg"]; !ok {
					params["rpmPkg"] = &RpmPkg{Name: rpmFile.Name(), Version: rpmFile.Version(), Release: rpmFile.Release(), Epoch: rpmFile.Epoch(), Licenses: []string{rpmFile.License()}, Vendor: rpmFile.Vendor()}
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type RpmPkg struct {
	Name     string
	Version  string
	Release  string
	Epoch    int
	Licenses []string
	Vendor   string
}
