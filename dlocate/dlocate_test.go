// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

//!+test
package main

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestCleared(t *testing.T) {
	clearIndex()

	// list of files to check that its cleared
	fileNames := []string{
		"indexFiles/indexInfo.gob",
		"indexFiles/indexInfo.json",
		"indexFiles/invertedIndex.json",
		"indexFiles/invertedIndex.gob",
	}
	for _, fileName := range fileNames {
		if _, err := os.Stat(fileName); err == nil {
			t.Errorf(`file : %q still exits after clear`, fileNames)
		}
	}

}

// command to run this: go test -benchmem -run=^$ -bench=.
// to test find in linux run : time find . -name "*run*"
func BenchmarkSearch(b *testing.B) {

	log.SetLevel(log.ErrorLevel)

	indexInfo = getIndexInfo()
	directoryPartition = getDirectoryPartition()

	//	change this
	destination := "/home/ahmed/Downloads/cloud computing/"
	deepScan := false
	searchWord := "run"

	if deepScan {
		invertedIndex.Load()
	}

	// to test that function work
	// matchedFiles, _ := find(searchWord, destination, deepScan)
	// log.Info(matchedFiles)

	for i := 0; i < b.N; i++ {
		find(searchWord, destination, deepScan)
	}
}
