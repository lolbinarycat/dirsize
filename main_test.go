package main

import (
	"testing"
	"math/rand"
	"io/ioutil"

	dirwalk "github.com/karrick/godirwalk"
)

const MaxRand = 50000

func BenchmarkFormatFileSize(b *testing.B) {
	rGen := rand.New(rand.NewSource(86))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatFileSize(rGen.Int63n(MaxRand))
	}
}

func BenchmarkRandom(b *testing.B) {
	rGen := rand.New(rand.NewSource(86))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rGen.Int63n(MaxRand)
	}
}

func BenchmarkCalculateTotalSize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateTotalSize(".")
	}
}

func BenchmarkReadDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ioutil.ReadDir(".")
	}
}

func BenchmarkDirWalk(b *testing.B) {
	dirwalk.Walk
}
