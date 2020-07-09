package main

//type SortOutputOptions struct {}
import "sort"

type SortMethod uint8

const {
	SortName SortMethod = iota
	SortSize
}

func SortFileInfo (info FileInfoList, method SortMethod) []FileInfoOut {
	sortingMethod = method
	sort.Sort(info)
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
	case SortSize:
		// use > because we want larger files first
		return l[i].Size > l[j].Size
	}
}
