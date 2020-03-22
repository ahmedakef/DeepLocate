package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
)

// getPartitionFiles gets files of partition and its children
func getPartitionFiles(partitionIndex int) []string {
	fileNames := readPartitionFilesGob(partitionIndex)
	partition := readPartitionGob(partitionIndex)
	for i, fileName := range fileNames {
		fileNames[i] = partition.Root + fileName
	}
	for _, child := range partition.Children {
		fileNames = append(fileNames, getPartitionFiles(child)...)
	}

	return fileNames
}

// word: word to search
// root: directoy to search in
func find(word, root string) []string {
	var directoryPartition DirectoryPartition

	directoryPartition = readDirectoryPartitionGob()
	partitionIndex := directoryPartition.getDirectoryPartition(root)

	// get all files names in the partition and its children
	fileNames := getPartitionFiles(partitionIndex)

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
