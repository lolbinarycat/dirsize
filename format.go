package main

import 	(
	"fmt"
	"strings"
)

func FmtFileInfo(info FileInfo) FileInfoOut {
	name := info.Name
	if addSlashToDirs && info.IsDir {
		name += "/"
	}
	return FileInfoOut{name,FmtFileSize(info.Size)}
}

func FmtFileInfoList(list FileInfoList) FileInfoOutList {
	outList := make(FileInfoOutList,len(list))
	nilEnts := 0
	for i, m := range list {
		if m == nil {
			nilEnts++
			continue
		}
		out := FmtFileInfo(*m)
		outList[i-nilEnts] = out
	}
	return outList
}

func FmtFileSize(size int64) string {
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
func FmtOutput(fInfo FileInfoOutList,opt FmtOutputOptions) string {
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
