[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/go-archive-extractor)](https://goreportcard.com/report/github.com/chen-keinan/go-archive-extractor)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/go-archive-extractor/blob/master/LICENSE)
<img src="./pkg/img/coverage_badge.png" alt="test coverage badge">
[![Gitter](https://badges.gitter.im/beacon-sec/community.svg)](https://gitter.im/beacon-sec/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
# go-archive-extractor

The archive-extractor is an open-source library for extracting various archive types.

it returns archive headers metadata (name,size,timestamp,sha1 and sha256).

it also support  diffrent types of tar compressions.


* [Installation](#installation)
* [Supported Archives](#supported-archives)
* [Supported Tar Compression](#supported-tar-compression)
* [Usage](#usage)


## Installation
``
go get github.com/chen-keinan/go-archive-extractor
``

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
```
    zip := extractor.New(extractor.Zip)
    headers, err = zip.Extract("common.zip");
    fmt.Print(headers)
```
### Tar Usage
```
    tar := extractor.New(extractor.Tar)
    headers, err = tar.Extract("common.tar");
    fmt.Print(headers)
```
### Debian Usage
```
    deb := extractor.New(extractor.Deb)
    headers, err = deb.Extract("common.deb");
    fmt.Print(headers)
```
### RPM Usage
```
    rpm := extractor.New(extractor.Rpm)
    headers, err = rpm.Extract("common.rpm");
    fmt.Print(headers)
```
### 7z Usage
```
    sevenZip := extractor.New(extractor.SevenZip)
    headers, err = sevenZip.Extract("common.7z");
    fmt.Print(headers)
```

```
func main() {
    zip := extractor.New(extractor.Zip)
    headers, err = zip.Extract("common.zip");
    if err != nil {
        fmt.Print(err.Error())
    }
    fmt.Print(headers)
}
```
