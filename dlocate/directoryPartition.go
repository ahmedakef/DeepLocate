package main

import "strings"

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
