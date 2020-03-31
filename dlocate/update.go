package main

import (
	"os"
	"path/filepath"

	utils "./osutils"
	log "github.com/Sirupsen/logrus"
)

func update(path string) bool {
	directoryPartition = getDirectoryPartition()
	indexInfo = getIndexInfo()
	indexInfo.loadRoots()

	if directoryPartition.getPathPartition(path) == -1 {
		log.Warn("The path hasn't been indexed, index it first")
		return false
	}

	filepath.Walk(path, updateIfChanged)

	indexInfo.savePartitions()
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
		partition := indexInfo.getPartition(partitionIndex)
		if partition.filePaths == nil {
			partition.filePaths = readPartitionFilesGob(partitionIndex)
		}

		lastChanged := utils.GetFileMetadata(path).CTime // real one from OS
		relativePath := partition.getRelativePath(path)
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
