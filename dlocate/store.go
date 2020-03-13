package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func savePartitionGob(partition *Partition) {
	partitionsPath := "indexFiles/partitions/"
	if _, err := os.Stat(partitionsPath); os.IsNotExist(err) {
		os.MkdirAll(partitionsPath, os.ModePerm)
	}

	// FromSlash converts / to the specific file system separator
	dataFile, err := os.Create(filepath.FromSlash("indexFiles/partitions/p" + strconv.Itoa(partition.Index) + ".gob"))

	if err != nil {
		log.Errorf("Error while storing index for partition %v: %v\n", partition.Index, err)
		os.Exit(1)
	}

	var buf io.Writer = dataFile
	//buf = gzip.NewWriter(dataFile)
	enc := gob.NewEncoder(buf)
	enc.Encode(partition)

	dataFile.Close()
}

func readPartitionGob(index int) Partition {
	dataFile, err := os.Open(filepath.FromSlash("indexFiles/partitions/p" + strconv.Itoa(index) + ".gob"))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// ensure to close the file after the fuction end
	defer dataFile.Close()

	var buf io.Reader = dataFile
	//buf, _ = gzip.NewReader(dataFile)
	dec := gob.NewDecoder(buf)

	var partition Partition

	err = dec.Decode(&partition)

	if err != nil {
		log.Errorf("Error while reading index for partition %q: %v\n", index, err)
		os.Exit(1)
	}

	return partition
}
