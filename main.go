package main


import (
	"os"
	"io/ioutil"
	"path/filepath"
	"fmt"
)

func main() {
	var dir string
	if len(os.Args) == 1 {
		dir = "."
	} else {
		dir = os.Args[1]
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}
	fmt.Println(FormatFileSize(CalculateTotalSize(dir)))
}

func CalculateTotalSize(path string) int64 {
	var size int64 = 0 // size in bytes

	filepath.Walk(path,func (_ string, info os.FileInfo, err error) error {
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
