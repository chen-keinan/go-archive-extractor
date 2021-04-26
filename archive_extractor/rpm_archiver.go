package archive_extractor

import (
	"github.com/cavaliercoder/go-cpio"
	"github.com/jfrog/go-archive-extractor/archive_extractor/archiver_errors"
	"github.com/jfrog/go-archive-extractor/compression"
	"github.com/jfrog/go-rpm/v2"
	"io"
	"math"
)

type RpmArchiver struct{}

func (ra RpmArchiver) ExtractArchive(path string,
	processingFunc func(*ArchiveHeader, map[string]interface{}) error, params map[string]interface{}) error {

	rpmFile, err := rpm.OpenPackageFile(path)
	if compression.IsGetReaderError(err) {
		return archiver_errors.New(err)
	}
	if err != nil {
		return err
	}

	headerEnd := ra.getHeadersEnd(rpmFile.Headers)
	cReader, err := compression.NewReaderSkipBytes(path, headerEnd)
	if err != nil {
		return archiver_errors.New(err)
	}
	defer cReader.Close()

	err = ra.readRpm(processingFunc, params, rpmFile, cReader)
	if err != nil {
		return archiver_errors.New(err)
	}
	return nil
}

func (RpmArchiver) getHeadersEnd(headers []rpm.Header) int64 {
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
	return int64(end)
}

func (RpmArchiver) readRpm(processingFunc func(*ArchiveHeader, map[string]interface{}) error,
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
		if archiveEntry != nil && !archiveEntry.Mode.IsDir() {
			archiveHeader := NewArchiveHeader(cpioReader, archiveEntry.Name, archiveEntry.ModTime.Unix(), archiveEntry.Size)
			err = processingFunc(archiveHeader, params)
			if _, ok := params["rpmPkg"]; !ok {
				modularityLabel := getModularityLabel(rpmFile)
				params["rpmPkg"] = &RpmPkg{Name: rpmFile.Name(), Version: rpmFile.Version(), Release: rpmFile.Release(),
					Epoch: rpmFile.Epoch(), Licenses: []string{rpmFile.License()}, Vendor: rpmFile.Vendor(), ModularityLabel: modularityLabel}
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

const (
	RpmTagModularityLabel = 5096
)

func getModularityLabel(rpmFile *rpm.PackageFile) string {
	return rpmFile.GetString(1, RpmTagModularityLabel)
}

type RpmPkg struct {
	Name            string
	Version         string
	Release         string
	Epoch           int
	Licenses        []string
	Vendor          string
	ModularityLabel string
}
