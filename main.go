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
	IsDir bool
}

type FileInfoList []*FileInfo


var scratchBuffer []byte

func init() {
	scratchBuffer = make([]byte,godirwalk.MinimumScratchBufferSize)
}

var addSlashToDirs = true
//var sortBySize = true

var showHidden = false
var extraPadding = ""
var showTotal bool
func init() {
	flag.BoolVar(&showHidden, "a",false, "shows files starting with '.'")
	flag.StringVar(&extraPadding, "pad","", "a string of text to use as extra padding between file names and sizes")
	flag.BoolVar(&showTotal, "t", false, "show the total size of the directory")

	flag.String("sort","none","sorting method")
}


func main() {
	flag.Parse()
	var dir string = "."
	if flag.Arg(0) != "" {
		os.Chdir(flag.Arg(0))
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var totalSize int64 // only used if showTotal == true
	var infoArr = make(FileInfoList,len(files))
	skipped := 0 // how many entries have been skipped
	for i, f := range files {
		if err != nil {
			panic(err)
		}
		if !showHidden && f.Name()[0] == '.' {
			skipped++
			continue
		}
		var size int64
		var isDir bool
		if f.IsDir() {
			size = CalculateSize(filepath.Join(dir,f.Name()))
			isDir = true
		} else {
			size = f.Size()
			isDir = false
		}
		if showTotal {
			totalSize += size
		}

		infoArr[i-skipped] = &FileInfo{f.Name(),size,isDir}
	}
	infoArr = infoArr[:len(infoArr)-skipped]
	if showTotal {
		infoArr = append(infoArr, &FileInfo{Name: "total:",Size: totalSize})
	}
	if srtMethod := GetSortMethodFromFlags(); srtMethod != SortNone {
		SortFileInfo(infoArr,srtMethod)
	}

	// FmtOutput adds a newline, so we don't do it again
	fmt.Print(FmtOutput(FmtFileInfoList(infoArr),
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



