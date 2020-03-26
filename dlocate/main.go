package main

import (
	"flag"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

var operation = flag.String("o", "index", "the operation to do (index or search or clear")
var destination = flag.String("d", "./", "the search directory")
var searchWord = flag.String("s", "", "the search word")

func main() {

	log.SetLevel(log.DebugLevel)

	flag.Parse()
	root := *destination
	op := *operation

	// root = "/home/ahmed/Downloads/cloud computing/"
	// op = "search"
	// remove trailling backslash
	if filepath.ToSlash(root)[len(root)-1] == '/' {
		root = root[:len(root)-1]
	}

	if op == "index" {
		startIndexing(root)
	} else if op == "clear" {
		clearIndex()
	} else if op == "search" {
		word := *searchWord
		word = "run"
		find(word, root)

	}

}
