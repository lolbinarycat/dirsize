package main


import (
	"io/ioutil"
	"os"
	"path/filepath"

	//"io/ioutil"
	"flag"
	"fmt"

	//"path/filepath"

	//"time"
	"github.com/karrick/godirwalk"

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
var hideSize bool
func init() {
	flag.BoolVar(&showHidden, "a",false, "shows files starting with '.'")
	flag.StringVar(&extraPadding, "pad","", "a string of text to use as extra padding between file names and sizes")
	flag.BoolVar(&showTotal, "t", false, "show the total size of the directory")
	flag.BoolVar(&showTotal, "total", false, "show the total size of the directory")
	flag.BoolVar(&hideSize,"hide-size",false,"hide the sizes of files. ")
	flag.String("sort","none","sorting method")
}


func main() {
	flag.Parse()
	var dir string = "."
	if flag.Arg(0) != "" {
		os.Chdir(flag.Arg(0))
	}
	//showSize := !hideSize
	infoList := GetFileInfoList(dir,
		GetFileInfoListOpts{IgnoreHiddenFiles: !showHidden,
		GetSizes: !hideSize})

	if srtMethod := GetSortMethodFromFlags(); srtMethod != SortNone {
		SortFileInfo(infoList,srtMethod)
	}

	// FmtOutput adds a newline, so we don't do it again
	fmt.Print(
		FmtOutput(
			FmtFileInfoList(infoList,
				FmtFileInfoOpts{DirSuffix:"/",ShowSize:!hideSize}),
			FmtOutputOpts{ExtraPadding: extraPadding},
		),
	)
}

type GetFileInfoListOpts struct {
	IgnoreHiddenFiles bool
	GetSizes bool
}
func GetFileInfoList(dir string, opts GetFileInfoListOpts) FileInfoList {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var totalSize int64 // only used if showTotal == true
	var infoList = make(FileInfoList,len(files))
	skipped := 0 // how many entries have been skipped
	for i, f := range files {
		if opts.IgnoreHiddenFiles && f.Name()[0] == '.' {
			skipped++
			continue
		}
		var fInfo FileInfo
		//var size int64
		//var isDir bool
		if f.IsDir() {
			if opts.GetSizes {
				fInfo.Size = CalculateSize(filepath.Join(dir,f.Name()))
			}
			fInfo.IsDir = true
		} else {
			if opts.GetSizes {
				fInfo.Size = f.Size()
			}

			fInfo.IsDir = false
		}
		if showTotal {
			totalSize += fInfo.Size
		}
		fInfo.Name = f.Name()
		infoList[i-skipped] = &fInfo//&FileInfo{f.Name(),size,isDir}
	}
	if skipped != 0 {
		infoList = infoList[:len(infoList)-skipped]
	}
	
	if showTotal {
		infoList = append(infoList, &FileInfo{Name: "total:",Size: totalSize})
	}
	return infoList
}



var MetricBinarySuffixes = [...]string{"B","kiB","MiB","GiB"}
