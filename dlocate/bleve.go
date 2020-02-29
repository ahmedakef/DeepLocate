package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/blevesearch/bleve"
)

func createBleveIndex(indexName string) bleve.Index {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(indexName, mapping)
	if err != nil {
		log.Fatalln("Trouble making index!")
	}
	return index
}

func indexBleve(root string) {
	indexName := "index.bleve"
	index := createBleveIndex(indexName)

	files := WalkSearch(root)
	for _, file := range files {
		log.WithFields(log.Fields{
			"fileName": file.FileName,
			"modTime":  file.ModTime,
		}).Info("indexed this file")
		index.Index(file.FileName, file)
	}
	index.Close()
}

func findBleve(name string) {
	log.Info(name)
	index, _ := bleve.Open("index.bleve")
	query := bleve.NewQueryStringQuery(name)
	searchRequest := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Info(searchResults)
}
