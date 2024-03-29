[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/go-archive-extractor)](https://goreportcard.com/report/github.com/chen-keinan/go-archive-extractor)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/go-archive-extractor/blob/master/LICENSE)
<img src="./pkg/img/coverage_badge.png" alt="test coverage badge">
[![Gitter](https://badges.gitter.im/beacon-sec/community.svg)](https://gitter.im/beacon-sec/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

<br><img src="./pkg/img/arc-extract.png" width="300" alt="arc-extract logo"><br>

# go-archive-extractor

The archive-extractor is an open-source library for extracting various archive types.

it returns archive headers metadata (name,size,timestamp,sha1 and sha256)

it also supports different types of tar compressions

* [Installation](#installation)
* [Supported Archives](#supported-archives)
* [Supported Tar Compression](#supported-tar-compression)
* [Usage](#usage)
* [Archive Header](#archive-header)

## Installation

```shell
go get github.com/chen-keinan/go-archive-extractor
```

## Supported Archives

- tar
- zip
- rpm
- deb
- 7zip

## Supported Tar Compression

- bz2
- gz
- Z
- infla
- xp3
- xz

## Usage

### Zip Usage

```go
zip := extractor.New(extractor.Zip)
headers, err = zip.Extract("common.zip");
fmt.Print(headers)
```

### Tar Usage

```go
tar := extractor.New(extractor.Tar)
headers, err = tar.Extract("common.tar");
if err != nil {
   fmt.Print(err.Error())
}
fmt.Print(headers)
```

### Debian Usage

```go
deb := extractor.New(extractor.Deb)
headers, err = deb.Extract("common.deb");
if err != nil {
  fmt.Print(err.Error())
}
fmt.Print(headers)
```

### RPM Usage

```go
rpm := extractor.New(extractor.Rpm)
headers, err = rpm.Extract("common.rpm");
if err != nil {
  fmt.Print(err.Error())
}
fmt.Print(headers)
```

### 7z Usagey
```go
sevenZip := extractor.New(extractor.SevenZip)
headers, err = sevenZip.Extract("common.7z");
if err != nil {
  fmt.Print(err.Error())
}
fmt.Print(headers)
```

```go
package main

import (
	"fmt"
	"github.com/chen-keinan/go-archive-extractor/pkg/extractor"
)

func main() {
	zip := extractor.New(extractor.Zip)
	headers, err := zip.Extract("common.zip")
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print(headers)
}
```

### Archive Header

```go
type ArchiveHeader struct {
Name          string
ModTime       int64
Size          int64
Sha1          string
Sha2          string
PkgMeta       map[string]interface{}
ArchiveReader io.Reader
}
```
