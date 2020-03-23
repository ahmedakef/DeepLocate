package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// DirectoryPartition saves index of partitions
type DirectoryPartition map[string]int

func (d DirectoryPartition) getDirectoryPartition(path string) int {
	var index = -1
	var ok = false
	for {
		index, ok = d[path]
		if !ok {
			lastslash := strings.LastIndex(path, "/")
			if lastslash == -1 {
				break
			}
			// search by the parent difrectory
			path = path[:lastslash]
		} else {
			// path found
			return index
		}
	}
	return -1
}

func readDirectoryPartitionGob() DirectoryPartition {
	path := "indexFiles/directoryPartition.gob"

	var directoryPartition DirectoryPartition
	err := readGob(path, &directoryPartition)
	if err != nil {
		log.Error("Error while reading directoryPartition")
		os.Exit(1)
	}

	return directoryPartition
}

func (d *DirectoryPartition) saveAsGob() {

	path := "indexFiles/directoryPartition.gob"
	err := saveGob(path, d)

	if err != nil {
		log.Errorf("Error while creating directoryPartitio file")
		os.Exit(1)
	}
}
