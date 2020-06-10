package main

import (
	"os"
	"path/filepath"
	"strconv"

	structure "dlocate/dataStructures"
	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

const filesLimit = 100

var invertedIndex structure.InvertedIndex

func startIndexing(path string) {
	//load index and check for repeated indexing
	if isRoot(path) != -1 {
		return
	}

	indexPath(path)

	indexInfo.clearPartitions()
	directoryPartition.saveAsGob()
	indexInfo.saveAsGob()
	invertedIndex.Save()
}

func indexPath(path string) {
	root := indexInfo.addRoot(path)
	directoryPartition[filepath.ToSlash(path)] = root.Index

	invertedIndex.Load()

	indexDir(path, &root)
	savePartition(&root)
}

func indexDir(path string, root *Partition) {
	files := utils.ListFiles(path)
	root.addDir(path)
	for _, file := range files {
		if file.IsDir {
			indexedUnder := isRoot(file.Path)
			if indexedUnder != -1 {
				parition := indexInfo.getPartition(indexedUnder)
				root.addChild(&parition)
				indexInfo.removeRoot(indexedUnder)
				continue
			}
			if root.FilesNumber >= filesLimit {
				child := NewPartition(indexInfo.getNextPartitionIndex(), file.Path)
				indexDir(file.Path, &child)
				indexInfo.addPartition(child)
				directoryPartition[filepath.ToSlash(file.Path)] = child.Index
				root.addChild(&child)
				savePartition(&child)
			} else {
				indexDir(file.Path, root)
			}
		}
	}
}

func savePartition(partition *Partition) {
	savePartitionGob(partition)
	savePartitionFilesGob(partition.Index, partition.filePaths)
	savePartitionMetaGob(partition.Index, partition.metadataTree)
	indexInfo.partitionsCache.Delete(strconv.Itoa(partition.Index))
	indexInfo.metaCache.Delete(strconv.Itoa(partition.Index))
	indexInfo.filesCache.Delete(strconv.Itoa(partition.Index))

	partition.printPartition()

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
	os.Remove("indexFiles/invertedIndex.json")
	err = os.Remove("indexFiles/invertedIndex.gob")
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
	err = utils.RemoveContents("indexFiles/content")
	if err != nil {
		log.Error(err)
	}

	log.Info("Index cleared successfully")
}
