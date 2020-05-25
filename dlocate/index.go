package main

import (
	"os"
	"path/filepath"

	structure "dlocate/dataStructures"
	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

const filesLimit = 100

var directoryPartition DirectoryPartition
var indexInfo IndexInfo
var invertedIndex structure.InvertedIndex

func startIndexing(path string) {
	//load index and check for repeated indexing
	indexInfo = getIndexInfo()
	indexInfo.loadRoots()
	if isRoot(path) != -1 {
		return
	}

	indexPath(path)

	indexInfo.savePartitions()
	directoryPartition.saveAsGob()
	indexInfo.saveAsGob()
	invertedIndex.Save()
}

func indexPath(path string) {
	root := indexInfo.addRoot(path)
	directoryPartition = getDirectoryPartition()
	directoryPartition[filepath.ToSlash(path)] = root.Index

	invertedIndex.Load()

	indexDir(path, root)
}

func indexDir(path string, root *Partition) {
	files := utils.ListFiles(path)
	root.addDir(path)
	for _, file := range files {
		if file.IsDir {
			indexedUnder := isRoot(file.Path)
			if indexedUnder != -1 {
				root.addChild(indexInfo.partitions[indexedUnder])
				indexInfo.removeRoot(indexedUnder)
				continue
			}
			if root.FilesNumber >= filesLimit {
				child := NewPartition(indexInfo.getNextPartitionIndex(), file.Path)
				indexDir(file.Path, &child)
				indexInfo.addPartition(&child)
				directoryPartition[filepath.ToSlash(file.Path)] = child.Index
				root.addChild(&child)
			} else {
				indexDir(file.Path, root)
			}
		}
	}
}

func isRoot(path string) int {
	for _, root := range indexInfo.Roots {
		if indexInfo.partitions[root].Root == path {
			return root
		}
	}
	return -1
}

func clearIndex() {
	os.Remove("indexFiles/directoryPartition.json")
	err := os.Remove("indexFiles/directoryPartition.gob")
	if err != nil {
		log.Error(err)
	}
	os.Remove("indexFiles/indexInfo.json")
	err = os.Remove("indexFiles/indexInfo.gob")
	if err != nil {
		log.Error(err)
	}
	err = utils.RemoveContents("indexFiles/filepaths")
	if err != nil {
		log.Error(err)
	}
	err = utils.RemoveContents("indexFiles/metadata")
	if err != nil {
		log.Error(err)
	}
	err = utils.RemoveContents("indexFiles/partitions")
	if err != nil {
		log.Error(err)
	}

	log.Info("Index cleared successfully")
}
