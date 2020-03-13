package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strconv"
)

func savePartitionGob(partition *Partition) {
	dataFile, err := os.Create("indexFiles\\p" + strconv.Itoa(partition.Index) + ".gob")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var buf io.Writer = dataFile
	//buf = gzip.NewWriter(dataFile)
	enc := gob.NewEncoder(buf)
	enc.Encode(partition)

	dataFile.Close()
}

func readPartitionGob(index int) Partition {
	dataFile, err := os.Open("indexFiles\\p" + strconv.Itoa(index) + ".gob")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var buf io.Reader = dataFile
	//buf, _ = gzip.NewReader(dataFile)
	dec := gob.NewDecoder(buf)

	var partition Partition

	err = dec.Decode(&partition)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataFile.Close()

	return partition
}
