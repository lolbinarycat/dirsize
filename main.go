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

func init() {
	flag.BoolVar(&showHidden, "a",false, "shows files starting with '.'")
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
	var outputArr = make([]FileInfoOut,len(files))
	for i, f := range files {
		if !showHidden && f.Name()[0] == '.' {
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
	// FmtOutput adds a newline, so we don't do it again
	fmt.Print(FmtOutput(outputArr))
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

func FmtOutput(fInfo []FileInfoOut) string {
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
	bldr.Grow((maxNameLength + 1 + (4 + 3) + 1) * len(fInfo))
	for _, inf := range fInfo {
		if (inf == FileInfoOut{}) {
			continue
		}
		bldr.WriteString(inf.Name)
		for i := len(inf.Name); i < maxNameLength; i++ {
			// Pad name with spaces to make all names the same width (so sizes align)
			bldr.WriteRune(' ')
		}
		// One more space for extra padding
		bldr.WriteRune(' ')
		bldr.WriteString(inf.Size)
		bldr.WriteRune('\n')
	}
	return bldr.String()
}
