package main

import (
	"strings"
	"time"

	utils "dlocate/osutils"
	log "github.com/sirupsen/logrus"
)

// getPartitionFiles gets files of partition and its children
func getPartitionFiles(partitionIndex int, root string) []string {
	partition := readPartitionGob(partitionIndex)

	// check that partition is related to the root
	if !partition.inSameDirection(root) {
		return []string{}
	}
	// check that partition have the root or its children
	// TODO : I think we don't need this check now as inSameDirection is enugh
	if !partition.containsDir(partition.getRelativePath(root)) {
		return []string{}
	}
	partitionFiles := readPartitionFilesGob(partitionIndex)
	fileNames := make([]string, partition.FilesNumber)
	i := 0
	for path, files := range partitionFiles {
		for _, fileName := range files {
			fileNames[i] = partition.Root + path + fileName
			i++
		}
	}
	for _, child := range partition.Children {
		fileNames = append(fileNames, getPartitionFiles(child, root)...)
	}

	return fileNames
}

// word: word to search
// root: directoy to search in
func findByFileName(word, root string) []string {
	var directoryPartition DirectoryPartition

	directoryPartition = getDirectoryPartition()
	partitionIndex := directoryPartition.getPathPartition(root)

	// get all files names in the partition and its children
	fileNames := getPartitionFiles(partitionIndex, root)

	var matchedFiles []string
	for _, fileName := range fileNames {
		if strings.Contains(fileName, word) {
			matchedFiles = append(matchedFiles, fileName)
			log.WithFields(log.Fields{
				"fileName": fileName,
			}).Info("found this file matched")
		}
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
	paritionIndex := 0
	tree := readPartitionMetaGob(paritionIndex)
	filesInfo := tree.SearchPartial(&start, &end)

	var files []string

	for _, file := range filesInfo {
		files = append(files, file.Path)
	}

	return files
}
