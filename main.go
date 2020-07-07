package main


import (
	"os"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"flag"
)

// FileInfoOut is a struct with the strings to print.
// This is stored first so spacing can be correctly calculated.
type FileInfoOut struct {
	Name, Size string
}

var ignoreDotfiles = true

func main() {
	flag.Parse()
	var dir string = flag.Arg(0)
	if dir == "" {
		dir = "."
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var outputArr = make([]FileInfoOut,len(files))
	for i, f := range files {
		if ignoreDotfiles && f.Name()[0] == '.' {
			continue
		}
		var size int64
		if f.IsDir() {
			size = CalculateTotalSize(filepath.Join(dir,f.Name()))
		} else {
			size = f.Size()
		}
		fmtdSize := FormatFileSize(size)
		outputArr[i] = FileInfoOut{f.Name(),fmtdSize}
	}
	// fmt.Println(FormatFileSize(CalculateTotalSize(dir)))
	fmt.Println(outputArr)
}

func CalculateTotalSize(dirpath string) int64 {
	var size int64 = 0 // size in bytes

	filepath.Walk(dirpath,func (path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		size += info.Size()
		return nil
	})
	return size
}

var MetricBinarySuffixes = [...]string{"B","kiB","MiB","GiB"}

func FormatFileSize(size int64) string {
	i := 0
	for size >= 1024 {
		i++
		size = size >> 10
	}
	return fmt.Sprintf("%d%v",size,MetricBinarySuffixes[i])
}
