package main

import (
	"os"
	"path/filepath"
	"strconv"

	structure "./dataStructures"
	utils "./osutils"
	log "github.com/Sirupsen/logrus"
)

func savePartitionFilesGob(partitionIndex int, partitionFiles map[string][]string) {

	filePathsPath := filepath.FromSlash("indexFiles/filepaths/")
	if _, err := os.Stat(filePathsPath); os.IsNotExist(err) {
		os.MkdirAll(filePathsPath, os.ModePerm)
	}

	path := "indexFiles/filepaths/f" + strconv.Itoa(partitionIndex)

	err := utils.SaveGob(partitionFiles, path+".gob")

	if err != nil {
		log.Errorf("Error while creating partitionfiles file")
		os.Exit(1)
	}
}

func readPartitionFilesGob(partitionIndex int) map[string][]string {
	// TODO return error and results and check the error outside
	path := "indexFiles/filepaths/f" + strconv.Itoa(partitionIndex) + ".gob"

	var partitionFiles map[string][]string
	err := utils.ReadGob(path, &partitionFiles)
	if err != nil {
		log.Error("Error while reading partitionfiles")
		os.Exit(1)
	}

	return partitionFiles
}

func savePartitionMetaGob(partitionIndex int, tree structure.KDTree) {

	metadataPath := filepath.FromSlash("indexFiles/metadata/")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		os.MkdirAll(metadataPath, os.ModePerm)
	}

	path := "indexFiles/metadata/m" + strconv.Itoa(partitionIndex)

	err := utils.SaveGob(tree, path+".gob")

	if err != nil {
		log.Error("Error while creating files metadata tree")
		os.Exit(1)
	}
}

func readPartitionMetaGob(partitionIndex int) structure.KDTree {
	// TODO return error and results and check the error outside
	path := "indexFiles/metadata/m" + strconv.Itoa(partitionIndex) + ".gob"

	var tree structure.KDTree
	err := utils.ReadGob(path, &tree)
	if err != nil {
		log.Error("Error while reading files metadata tree")
		os.Exit(1)
	}
	return tree
}
