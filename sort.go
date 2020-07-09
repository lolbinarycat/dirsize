package main

//type SortOutputOptions struct {}
import (
	"sort"
	"flag"
)

type SortMethod int8

const (
	SortNone SortMethod = -1
	SortName SortMethod = iota
	SortSize
)

func SortFileInfo (info FileInfoList, method SortMethod) FileInfoList {
	sortingMethod = method
	sort.Sort(info)
	return info
}

func (l FileInfoList) Len() int {
	return len(l)
}

func (l FileInfoList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

var sortingMethod = SortSize
func (l FileInfoList) Less(i, j int) bool {
	switch sortingMethod {
	case SortNone:
		return i < j // if we don't want to sort, the correct position is where things already are
	case SortSize:
		// use > because we want larger files first
		return l[i].Size > l[j].Size
	default:
		panic("sorting method not implemented")
	}
}

func GetSortMethodFromFlags() SortMethod {
	switch flag.Lookup("sort").Value.String() {
	case "none":
		return SortNone
	case "name":
		return SortName
	case "size":
		return SortSize
	default:
		panic("invalid sort method")
	}
}
