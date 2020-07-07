package main


import (
	"os"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"flag"
	"strings"
)

// FileInfoOut is a struct with the strings to print.
// This is stored first so spacing can be correctly calculated.
type FileInfoOut struct {
	Name, Size string
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
	var dir string = flag.Arg(0)
	if dir == "" {
		dir = "."
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var totalSize int64 // only used if showTotal == true
	var outputArr = make([]FileInfoOut,len(files))
	for i, f := range files {
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
		fmtdSize := FormatFileSize(size)
		outputArr[i] = FileInfoOut{f.Name(),fmtdSize}
	}
	if showTotal {
		outputArr = append(outputArr, FileInfoOut{Name: "total:",Size: FormatFileSize(totalSize)})
	}


	// FmtOutput adds a newline, so we don't do it again
	fmt.Print(FmtOutput(outputArr,FmtOutputOptions{ExtraPadding: extraPadding}))
}

func CalculateSize(dirpath string) int64 {
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

type FmtOutputOptions struct {
	// ExtraPadding can be used to provide more distance between file names and sizes
	ExtraPadding string
}
func FmtOutput(fInfo []FileInfoOut,opt FmtOutputOptions) string {
	var maxNameLength int = 0

	for _, inf := range fInfo {
		if len(inf.Name) > maxNameLength {
			maxNameLength = len(inf.Name)
		}
	}

	bldr := strings.Builder{}

	// we want to be able to store len(fInfo) lines.
	// the length of a line is maxNameLength (all names will be padded up to this with spaces)
	// + 1 for newLine
	// + 7 for filesize (4 digits + 3 for suffix (unless suffix is B, but we dont check that)) and
	// another + 1 for padding
	var bldrLen = (maxNameLength + 1 + (4 + 3) + 1) * len(fInfo)
	if opt.ExtraPadding != "" {
		// increase length for ExtraPadding, if it is used
		bldrLen += len(opt.ExtraPadding) * len(fInfo)
	}
	bldr.Grow(bldrLen)
	for _, inf := range fInfo {
		// Filter out zero entries left by hidden files.
		if (inf == FileInfoOut{}) {
			continue
		}
		bldr.WriteString(inf.Name)
		for i := len(inf.Name); i < maxNameLength; i++ {
			// Pad name with spaces to make all names the same width (so sizes align)
			bldr.WriteRune(' ')
		}
		// One more space for padding
		bldr.WriteRune(' ')
		// Add ExtraPadding (no need for an if statement, it will be ignored if empty)
		bldr.WriteString(opt.ExtraPadding)
		bldr.WriteString(inf.Size)
		bldr.WriteRune('\n')
	}
	return bldr.String()
}
