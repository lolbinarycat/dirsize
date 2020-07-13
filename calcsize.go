package main

import (
	dirwalk "github.com/karrick/godirwalk"
	"os"
)

func CalculateSize(dirpath string) int64 {

	/*filepath.Walk(dirpath,func (path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		size += info.Size()
		return nil
	})*/
	var size int64 = 0 // size in bytes
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
