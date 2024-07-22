package main

import (
	"math/rand"
	"testing"
)

var (
	testRandom = rand.New(rand.NewSource(0))
	result     []*City
)

func BenchmarkFindCityBy(b *testing.B) {
	var err error
	dataFiles, err = findFileByExt(".", ".txt")
	if err != nil {
		b.Errorf("error while finding files: %v", err)
	}

	if len(dataFiles) == 0 {
		b.Skip("No data files found")
	}
	var (
		population int
		r          []*City
	)

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		population = testRandom.Intn(100000)
		r, _ = findCityWithPopulation(population)
	}
	result = r
}
