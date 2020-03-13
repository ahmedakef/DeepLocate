package main

import (
	"os"

	utils "./osutils"
	log "github.com/Sirupsen/logrus"
)

const filesLimit = 100

// ListFiles return a list of files and folders directly under the given dir
func ListFiles(path string) []utils.FileMetadata {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var filesData []utils.FileMetadata
	for _, file := range files {
		filesData = append(filesData, utils.GetFileInfo(file, path))
	}

	return filesData
}

var curIndex = 0
var partitions []*Partition
var curExtensionIndex = -1
var extensions map[string]int

func getExtensionIndex(name string) int {
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
		if file.IsDir {
			if root.FilesNumber >= filesLimit {
				child := NewPartition(getNextPartitionIndex(), file.Path)
				indexDir(file.Path, &child)
				partitions = append(partitions, &child)
				root.addChild(&child)
			} else {
				indexDir(file.Path, root)
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
		log.Debugf("start saving Partition %v\n", partition.Index)
		partition.printPartition()
		savePartitionGob(partition)
	}

	for _, partition := range partitions {
		log.Debugf("start reading Partition %v\n", partition.Index)
		p := readPartitionGob(partition.Index)
		p.printPartition()
	}
}
