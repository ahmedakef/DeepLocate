// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

//!+test
package main

import (
	"os"
	"testing"
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
