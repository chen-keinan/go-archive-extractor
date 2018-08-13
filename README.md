# go-archive-extractor

The archive-extractor is a library and set of tools
that can extract many archive types (tar , zip , rpm ,deb, 7zip) with supported compressions (bz2,gz,Z,infl,xp3,xz) on tar files
and invoke advance processing function while iterating archive headers

Example:

- Define advance params to be uses in advance processing method :
 ```
type ArchiveData struct {
	ArchiveReader io.Reader
	IsFolder      bool
	Name          string
	ModTime       int64
	Size          int64
}
```
```
func params() map[string]interface{} {
	return map[string]interface{}{
		"archveData": &ArchiveData{},
	}
}
```
- Define advance processing function to be invoke during archive extraction :
```
func processingFunc(header *ArchiveHeader, params map[string]interface{}) error {
	if len(params) == 0 {
		return errors.New("Advance processing params are missing")
	}
	var ok bool
	var archiveData *ArchiveData
	if archiveData, ok = params["archiveData"].(*ArchiveData); !ok {
		return errors.New("Advance processing archveData param is missing")
	}
	archiveData.Name = header.Name
	archiveData.ModTime=header.ModTime
	archiveData.Size=header.Size
	archiveData.IsFolder=header.IsFolder
 	fmt.Print(archiveData)
	return nil
}
```
- create archive extractor type and pass advance processing function and params :
```
func main() {
	za := &ZipArchvier{}
	if err:=za.ExtractArchive("/User/Name/file.zip",processingFunc,params()); err != nil{
 		fmt.Print(err)
 		}
}
```
