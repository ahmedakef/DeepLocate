package main

import (
	"flag"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

var operation = flag.String("o", "index", "the operation to do (index or search")
var destination = flag.String("d", "./", "the search directory")
var searchWord = flag.String("s", "", "the search word")

func main() {

	log.SetLevel(log.DebugLevel)

	flag.Parse()
	root := *destination
	// remove trailling backslash
	if filepath.ToSlash(root)[len(root)-1] == '/' {
		root = root[:len(root)-1]
	}

	op := *operation
	// root = "/home/ahmed/Downloads/cloud computing/"
	// op = "search"

	if op == "index" {
		startIndexing(root)
	} else if op == "search" {
		word := *searchWord
		// word = "hello"
		find(word, root)

	}

}
