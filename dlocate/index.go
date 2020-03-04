package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const filesLimit = 100

// ListFiles return a list of files and folders directly under the given dir
func ListFiles(path string) []os.FileInfo {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func getDirName(path string) string {
	index := strings.LastIndex(path, "\\")
	dir := path[index+1:]
	return dir
}

func getFileExtenstion(name string) string {
	dot := strings.LastIndex(name, ".")
	extention := name[dot+1:]
	return extention
}

var curIndex = 0
var partitions []*Partition
var curExtensionIndex = -1
var extensions map[string]int

func getExtenstionIndex(name string) int {
	index, ok := extensions[name]
	if !ok {
		index = getNextExtensionIndex()
		extensions[name] = index
	}
	return index
}

func getNextExtensionIndex() int {
	curExtensionIndex++
	return curExtensionIndex
}

func getNextPartitionIndex() int {
	curIndex++
	return curIndex
}

func indexDir(path string, root *Partition) {
	files := ListFiles(path)
	root.addDir(path)
	for _, file := range files {
		if file.IsDir() {
			dirPath := path + "\\" + file.Name()
			if root.FilesNumber >= filesLimit {
				child := NewPartition(getNextPartitionIndex(), dirPath)
				indexDir(dirPath, &child)
				partitions = append(partitions, &child)
				root.addChild(&child)
			} else {
				indexDir(dirPath, root)
			}
		}
	}
}

func startIndexing(path string) {
	root := NewPartition(0, path)
	partitions = append(partitions, &root)
	extensions = make(map[string]int)
	indexDir(path, &root)

	for _, partition := range partitions {
		partition.printPartition()
		savePartitionGob(partition)
	}

	for _, partition := range partitions {
		p := readPartitionGob(partition.Index)
		p.printPartition()
	}
}
