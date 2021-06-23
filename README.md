# go-archive-extractor

The archive-extractor is a library and set of tools
that can extract many various archive types with various tar compressions
and invoke advance processing function while iterating archive headers


* [Supported Archives](#supported-archives)
* [Supported Tar Compression](#supported-tar-compression)
* [Usage](#usage)




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
```
    zip := New(Zip)
    headers, err = zip.Extract("common.zip");
    fmt.Print(headers)
```

```
func main() {
	 zip := New(Zip)
    headers, err = zip.Extract("common.zip");
    if err != nil {
        fmt.Print(err.Error())
    }
    fmt.Print(headers)
}
```
