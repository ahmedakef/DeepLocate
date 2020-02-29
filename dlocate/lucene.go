package main

import (
	"fmt"
	"os"

	std "github.com/balzaczyy/golucene/analysis/standard"
	_ "github.com/balzaczyy/golucene/core/codec/lucene410"
	"github.com/balzaczyy/golucene/core/document"
	"github.com/balzaczyy/golucene/core/index"
	"github.com/balzaczyy/golucene/core/search"
	"github.com/balzaczyy/golucene/core/store"
	"github.com/balzaczyy/golucene/core/util"

	log "github.com/Sirupsen/logrus"
)

func indexLucene(root string) {
	util.SetDefaultInfoStream(util.NewPrintStreamInfoStream(os.Stdout))

	index.DefaultSimilarity = func() index.Similarity {
		return search.NewDefaultSimilarity()
	}

	directory, _ := store.OpenFSDirectory("index")
	analyzer := std.NewStandardAnalyzer()
	conf := index.NewIndexWriterConfig(util.VERSION_LATEST, analyzer)
	writer, _ := index.NewIndexWriter(directory, conf)

	files := WalkSearch(root)
	for _, file := range files {
		log.WithFields(log.Fields{
			"fileName": file.FileName,
			"modTime":  file.ModTime,
		}).Info("indexed this file")
		//d := document.NewDocument()
		//d.Add(document.NewTextFieldFromString("name", file.FileName, document.STORE_YES))
		//d.Add(document.NewFieldFromString("path", file.LinuxPath, document.TEXT_FIELD_TYPE_STORED))
		//d.Add(document.NewFieldFromString("modified", file.ModTime.String(), document.STORE_YES))
		//d.Add(document.NewFieldFromString("type", string(file.Type), document.STORE_YES))
		//writer.AddDocument(d.Fields())
	}

	d := document.NewDocument()
	d.Add(document.NewTextFieldFromString("name", "new shape", document.STORE_YES))
	writer.AddDocument(d.Fields())

	writer.Close()
}

func findLucene(name string) {
	directory, _ := store.OpenFSDirectory("index")
	reader, _ := index.OpenDirectoryReader(directory)
	searcher := search.NewIndexSearcher(reader)

	q := search.NewTermQuery(index.NewTerm("name", name))
	res, _ := searcher.Search(q, nil, 1000)
	fmt.Printf("Found %v hit(s).\n", res.TotalHits)
	for _, hit := range res.ScoreDocs {
		fmt.Printf("Doc %v score: %v\n", hit.Doc, hit.Score)
		doc, _ := reader.Document(hit.Doc)
		fmt.Printf("name -> %v\n", doc.Get("name"))
	}
}
