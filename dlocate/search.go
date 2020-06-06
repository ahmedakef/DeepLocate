package main

import (
	"strconv"
	"strings"
	"time"

	structure "dlocate/dataStructures"
	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

// getPartitionFiles gets files of partition and its children
func getPartitionFiles(partitionIndex int, path string) []string {
	partition := indexInfo.getPartition(partitionIndex)

	// check that partition is related to the path
	if !partition.inSameDirection(path) {
		return []string{}
	}

	partitionFiles, ok := indexInfo.filesCache.Get(strconv.Itoa(partitionIndex))
	if !ok {
		partitionFiles = readPartitionFilesGob(partitionIndex)
		indexInfo.filesCache.Set(strconv.Itoa(partitionIndex), partitionFiles)
	}
	fileNames := make([]string, partition.FilesNumber)
	i := 0
	for path, files := range partitionFiles.(map[string][]string) {
		for _, fileName := range files {
			fileNames[i] = partition.Root + path + fileName
			i++
		}
	}
	for _, child := range partition.Children {
		fileNames = append(fileNames, getPartitionFiles(child, path)...)
	}

	return fileNames
}

// getPartitionClildren get the children of partioin that is related
// to the given path (either parent or child)
// and excludes the non relevant partitions
func getPartitionClildren(partitionIndex int, path string) []int {
	partition := indexInfo.getPartition(partitionIndex)

	// check that partition is related to the path
	if !partition.inSameDirection(path) {
		return []int{}
	}

	children := []int{partitionIndex}

	for _, child := range partition.Children {
		children = append(children, getPartitionClildren(child, path)...)
	}

	return children
}

// query: word to search
// path: directoy to search in
// searchContent : bool to indicate search content or not
func find(query, path string, searchContent bool) []string {
	var directoryPartition DirectoryPartition

	directoryPartition = getDirectoryPartition()
	partitionIndex := directoryPartition.getPathPartition(path)

	log.Info("Start searching file names ...")

	// get all files names in the partition and its children
	fileNames := getPartitionFiles(partitionIndex, path)

	var matchedFiles []string
	for _, fileName := range fileNames {
		if strings.Contains(fileName, query) {
			matchedFiles = append(matchedFiles, fileName)
			log.WithFields(log.Fields{
				"fileName": fileName,
			}).Info("found this file matched")
		}
	}

	if searchContent {
		log.Info("Start searching file content ...")

		invertedIndex.Load()
		children := getPartitionClildren(partitionIndex, path)
		contentResults := invertedIndex.Search(children, query, -1)
		matchedFiles = append(matchedFiles, contentResults...)

		log.WithFields(log.Fields{
			"fileNames": contentResults,
		}).Info("Result of content Search :")

	}
	return matchedFiles
}

func metaSearch() []string {
	//Size - file size in bytes
	//CTime - change time (last file name or path change)
	//MTime - modify time Max(last content change, CTime)
	//ATime - access time Max(last opened, MTime)

	start := utils.FileMetadata{ATime: time.Date(2019, 1, 1, 20, 34, 58, 651387237, time.UTC)}
	end := utils.FileMetadata{}

	//get parition index:
	partitionIndex := 0
	val, ok := indexInfo.metaCache.Get(strconv.Itoa(partitionIndex))
	var tree structure.KDTree
	if !ok {
		tree = readPartitionMetaGob(partitionIndex)
		indexInfo.metaCache.Set(strconv.Itoa(partitionIndex), tree)
	} else {
		tree = val.(structure.KDTree)
	}
	filesInfo := tree.SearchPartial(&start, &end)

	var files []string

	for _, file := range filesInfo {
		files = append(files, file.Path)
	}

	return files
}
