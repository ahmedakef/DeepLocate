package main

import (
	"os"
	"path/filepath"
	"strconv"

	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

func update(path string) bool {
	directoryPartition = getDirectoryPartition()
	indexInfo = getIndexInfo()

	if directoryPartition.getPathPartition(path) == -1 {
		log.Warn("The path hasn't been indexed, index it first")
		return false
	}

	filepath.Walk(path, updateIfChanged)

	//delete directories info that was deleted

	for index := 1; index <= indexInfo.CurIndex; index++ {
		partition := indexInfo.getPartition(index)
		for directory, toBeDeleted := range partition.toBeDeleted {
			if toBeDeleted {
				log.Warnf("Directory %v has been deleted "+
					"and will be removed from index", directory)
				partition.clearDir(partition.Root + directory[:len(directory)-1])
			}
		}
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
		partition := indexInfo.getPartition(partitionIndex)
		if partition.filePaths == nil {
			partitionFiles, _ := indexInfo.filesCache.Get(strconv.Itoa(partitionIndex))
			partition.filePaths = partitionFiles.(map[string][]string)
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
			indexDir(path, &partition) // index the directoy and its subdirectories

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
