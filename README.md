# go-archive-extractor

The archive-extractor is a library and set of tools
that can extract many various archive types with various tar compressions
and invoke advance processing function while iterating archive headers


* [Supported Archives](#supported-archives)
* [Supported Tar Compression](#supported-tar-compression)
* [Zip Usage](#zip-usage)
* [Tar Usage](#tar-usage)
* [Debian Usage](#debian-usage)
* [RPM Usage](#rpm-usage)
* [7z Usage](#7z-usage)




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

## Zip Usage
```
    zip := extractor.New(extractor.Zip)
    headers, err = zip.Extract("common.zip");
    fmt.Print(headers)
```
## Tar Usage
```
    tar := extractor.New(extractor.Tar)
    headers, err = tar.Extract("common.tar");
    fmt.Print(headers)
```
## Debian Usage
```
    deb := extractor.New(extractor.Deb)
    headers, err = deb.Extract("common.deb");
    fmt.Print(headers)
```
## RPM Usage
```
    rpm := extractor.New(extractor.Rpm)
    headers, err = rpm.Extract("common.rpm");
    fmt.Print(headers)
```
## 7z Usage
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
