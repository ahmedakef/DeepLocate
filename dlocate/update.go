package main

import (
	"os"
	"path/filepath"
	"strings"

	utils "dlocate/osutils"
	python "dlocate/python"

	log "github.com/sirupsen/logrus"
)

func update(path string) bool {

	if deepScan {
		log.Info("get all files content from the machine learning model")
		log.Info("This should take some minutes ...")
		python.ExecuteScript("Extract.py", path, &filesContent)
		log.Info("Finished reading all files content in the given path")
	}

	if directoryPartition.getPathPartition(path) == -1 {
		log.Warn("The path hasn't been indexed, index it first")
		return false
	}

	filepath.Walk(path, updateIfChanged)

	//delete directories info that was deleted
	for _, partition := range indexInfo.updatedPartitions {
		// partition = indexInfo.getPartition(partitionIndex)
		for directory, toBeDeleted := range partition.toBeDeleted {
			directorInSearch := strings.HasPrefix(partition.Root+directory, path+"/")
			if toBeDeleted && directorInSearch {
				log.Warnf("Directory %v has been deleted "+
					"and will be removed from index", directory)
				partition.clearDir(partition.Root + directory[:len(directory)-1])
			}
		}
		savePartition(partition)

	}

	indexInfo.clearPartitions()
	directoryPartition.saveAsGob()
	indexInfo.saveAsGob()

	return true
}

func updateIfChanged(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Errorf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}

	if info.IsDir() {
		partitionIndex := directoryPartition.getPathPartition(path)

		// load partition if not loaded
		partition, ok := indexInfo.updatedPartitions[partitionIndex]
		if !ok {
			p := indexInfo.getPartition(partitionIndex)
			partition = &p
			indexInfo.updatedPartitions[partition.Index] = partition
		}

		if partition.filePaths == nil {
			partition.filePaths = partition.getPartitionFiles()
		}
		if partition.toBeDeleted == nil {
			partition.toBeDeleted = make(map[string]bool)
			// initialize all directories to be deleted
			for directory := range partition.Directories {
				partition.toBeDeleted[directory] = true
			}
		}

		lastChanged := utils.GetFileMetadata(path).CTime // real one from OS
		relativePath := partition.getRelativePath(path)
		partition.toBeDeleted[relativePath] = false
		savedLastChanged, ok := partition.Directories[relativePath]
		// if new directory, index it and its subdirectories
		if !ok {
			log.WithFields(log.Fields{
				"Path": path,
			}).Infof("New Directory :")
			indexDir(path, partition) // index the directoy and its subdirectories

			// update lastchanged as folder is already indexed
			partition.Directories[relativePath] = lastChanged
			return nil
		}

		if lastChanged.After(savedLastChanged) {
			// redindex all files and only the changed folders
			log.Warnf("Directory %v has changed and will index files again", path)

			partition.clearDir(path)
			partition.addDir(path) // index files only (1 level)

			partition.Directories[relativePath] = lastChanged
		}
	}

	return nil
}
