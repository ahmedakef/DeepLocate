package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// DirectoryPartition saves index of partitions
type DirectoryPartition map[string]int

func (d DirectoryPartition) getPathPartition(path string) int {
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

func getDirectoryPartition() DirectoryPartition {
	path := "indexFiles/directoryPartition.gob"

	var directoryPartition DirectoryPartition
	err := readGob(path, &directoryPartition)
	if err != nil {
		return make(map[string]int)
	}

	return directoryPartition
}

func (d *DirectoryPartition) saveAsGob() {

	path := "indexFiles/directoryPartition.gob"
	SaveAsJSON(d, "indexFiles/directoryPartition.json")
	err := saveGob(d, path)

	if err != nil {
		log.Errorf("Error while creating directoryPartitio file")
		os.Exit(1)
	}
}
