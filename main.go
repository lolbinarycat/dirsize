package main


import (
	"io/ioutil"
	"os"
	//"io/ioutil"
	"flag"
	"fmt"
	"path/filepath"


	//"time"
	"github.com/karrick/godirwalk"
	dirwalk "github.com/karrick/godirwalk"
)

// FileInfoOut is a struct with the info to print.
// This is stored first so spacing can be correctly calculated.
type FileInfoOut struct {
	Name, Size string
}
type FileInfoOutList []FileInfoOut

// FileInfo is the type used before Size is formatted into a string
// Mostly used to be sorted
type FileInfo struct {
	Name string
	Size int64
}

type FileInfoList []*FileInfo


var scratchBuffer []byte

func init() {
	scratchBuffer = make([]byte,godirwalk.MinimumScratchBufferSize)
}

var showHidden = false
var extraPadding = ""
var showTotal bool
func init() {
	flag.BoolVar(&showHidden, "a",false, "shows files starting with '.'")
	flag.StringVar(&extraPadding, "pad","", "a string of text to use as extra padding between file names and sizes")
	flag.BoolVar(&showTotal, "t", false, "show the total size of the directory")
}


func main() {
	flag.Parse()
	var dir string = "."
	if flag.Arg(0) != "" {
		os.Chdir(flag.Arg(0))
	}

	fmt.Println(os.Getwd())
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var totalSize int64 // only used if showTotal == true
	var outputArr = make(FileInfoList,len(files))
	for i, f := range files {
		if err != nil {
			panic(err)
		}
		if !showHidden && f.Name()[0] == '.' {
			continue
		}
		var size int64
		if f.IsDir() {
			size = CalculateSize(filepath.Join(dir,f.Name()))
		} else {
			size = f.Size()
		}
		if showTotal {
			totalSize += size
		}
		//fmtdSize := FormatFileSize(size)
		outputArr[i] = &FileInfo{f.Name(),size}
	}
	if showTotal {
		outputArr = append(outputArr, &FileInfo{Name: "total:",Size: totalSize})
	}


	// FmtOutput adds a newline, so we don't do it again
	fmt.Print(FmtOutput(FmtFileInfoList(outputArr),
			FmtOutputOptions{ExtraPadding: extraPadding} ))
}

func CalculateSize(dirpath string) int64 {
	var size int64 = 0 // size in bytes

	/*filepath.Walk(dirpath,func (path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		size += info.Size()
		return nil
	})*/

	dirwalk.Walk(dirpath, &dirwalk.Options{
		Callback: func (path string, entry *dirwalk.Dirent) error {
			isDirOrSym, err :=  entry.IsDirOrSymlinkToDir()
			if err != nil {
				return err
			}
			if isDirOrSym {
				return nil
			} else {
				stat, err := os.Lstat(path)
				if err != nil {
					return err
				}
				size += stat.Size()
				return nil
			}
		},
		Unsorted: true,
		ScratchBuffer: scratchBuffer,
	})
	return size
}

var MetricBinarySuffixes = [...]string{"B","kiB","MiB","GiB"}



