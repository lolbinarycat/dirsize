package main

import (
	"testing"
	"math/rand"
	"io/ioutil"

	dirwalk "github.com/karrick/godirwalk"
	//"github.com/stretchr/testify/assert"
)

const MaxRand = 50000

var result interface{}

func BenchmarkFormatFileSize(b *testing.B) {
	rGen := rand.New(rand.NewSource(86))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FmtFileSize(rGen.Int63n(MaxRand))
	}
}

func BenchmarkRandom(b *testing.B) {
	rGen := rand.New(rand.NewSource(86))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rGen.Int63n(MaxRand)
	}
}

func BenchmarkCalculateSize(b *testing.B) {
	var r int64
	for i := 0; i < b.N; i++ {
		r = CalculateSize("~")
	}
	result = r
}

func BenchmarkReadDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ioutil.ReadDir(".")
	}
}

func BenchmarkGetFileInfoList(b *testing.B) {
	var r FileInfoList
	for i := 0; i < b.N; i++ {
		r = GetFileInfoList("/home/binarycat")
	}
	result = r
}

func BenchmarkReadDirnames(b *testing.B) {
	for i := 0; i<b.N;i++{
		dirwalk.ReadDirnames("/home/binarycat",nil)
	}
}

func BenchmarkReadDirnames_Buffed(b *testing.B) {
	scratchBuffer = make([]byte,dirwalk.MinimumScratchBufferSize)
	for i := 0; i<b.N;i++{
		dirwalk.ReadDirnames("/home/binarycat",scratchBuffer)
	}
}

//func TestGetFileInfoList(t *testing.T) {
	//infoList := GetFileInfoList("/home/binarycat")
	//rawInfoList, err := ioutil.ReadDir("/home/binarycat")
	//if err != nil {
	//	panic(err)
	//}
	//assert.Equal(t,len(rawInfoList),len(infoList),"lengths do not match")
	//}
