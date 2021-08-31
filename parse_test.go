package main

import "testing"

func BenchmarkParseSVFile(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseSVFile("./gene-flow.csv", ',')
	}
}

func BenchmarkParseSVFileDirty(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseSVFileDirty("./gene-flow.csv", ',')
	}
}
