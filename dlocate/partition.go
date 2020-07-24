package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	structure "dlocate/dataStructures"
	utils "dlocate/osutils"

	log "github.com/sirupsen/logrus"
)

// map from directory path to its lastChanged Date
type directories map[string]time.Time

// Partition conatins basic info about partitions
type Partition struct {
	Index        int
	Root         string
	Directories  directories
	FilesNumber  int
	Children     []int // index for children partitions
	Extenstion   SignatureFile
	ExtenstionH  SignatureFile
	filePaths    map[string][]string
	toBeDeleted  map[string]bool
	metadataTree structure.KDTree
	//TODO implement versioning
}

// NewPartition creates new partition with 0 files and 0 children
func NewPartition(index int, root string) Partition {
	return Partition{
		Index: index, Root: root, Directories: make(map[string]time.Time),
		Children:  make([]int, 0),
		filePaths: make(map[string][]string), FilesNumber: 0,
		Extenstion: newSignatureFile(), ExtenstionH: newSignatureFileH(),
	}
}

func (p *Partition) addDir(path string) {
	lastChanged := utils.GetFileMetadata(path).CTime
	relativePath := p.getRelativePath(path)

	files := utils.ListFiles(path)
	cnt := 0
	for _, file := range files {
		if !file.IsDir {
			p.filePaths[relativePath] = append(p.filePaths[relativePath], file.Name)
			p.metadataTree.Insert(&file)
			cnt++
			p.addExtension(file.Extension)

			if deepScan {
				content := filesContent[file.Path]
				invertedIndex.Insert(p.Index, file.Path, content)
			}

		} else if _, ok := p.filePaths[relativePath]; !ok {
			p.filePaths[relativePath] = []string{}
		}
	}
	p.FilesNumber += cnt
	p.ExtenstionH.or(p.Extenstion)

	p.Directories[relativePath] = lastChanged
}

func readTxt(path string) map[string]float32 {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)
	var words = strings.Fields(text)
	fileContent := map[string]float32{}
	for _, word := range words {
		fileContent[word]++
	}

	return fileContent
}

func (p *Partition) getPartitionFiles() map[string][]string {
	partitionFiles, ok := indexInfo.filesCache.Get(strconv.Itoa(p.Index))
	if !ok {
		pF := readPartitionFilesGob(p.Index)
		indexInfo.filesCache.Set(strconv.Itoa(p.Index), pF)
		return pF
	}
	return partitionFiles.(map[string][]string)
}

func (p *Partition) clearDir(path string) {

	relativePath := p.getRelativePath(path)
	p.FilesNumber -= len(p.filePaths[relativePath])

	delete(p.Directories, relativePath)
	delete(p.filePaths, relativePath)

	// TODO : remove extensions from p.Extenstion
}

func (p *Partition) addExtension(extension string) {
	index := indexInfo.getExtensionIndex(extension)
	p.Extenstion.setBit(index)
}

func (p *Partition) hasExtension(extension string) bool {
	index := indexInfo.getExtensionIndex(extension)
	return p.Extenstion.getBit(index)
}

// dirInDirs ensures that root exits inside any of directories
func (p *Partition) containsDir(root string) bool {
	for dir := range p.Directories {
		if strings.HasPrefix(dir, root) {
			return true
		}
	}
	return false
}

// inSameDirection ensures that path is parent or child of the partition path
func (p *Partition) inSameDirection(path string) bool {
	// partition is the parent of the path
	if len(path) > len(p.Root) {
		return strings.HasPrefix(path, p.Root)
	}

	// partition is one of path's children
	return strings.HasPrefix(p.Root, path)
}

func (p *Partition) getRelativePath(path string) string {
	// if p.Root is longer then it is one of path's children
	if len(p.Root) > len(path) {
		return ""
	}
	return path[len(p.Root):] + "/"
}

func (p *Partition) addChild(c *Partition) {
	p.ExtenstionH.or(c.ExtenstionH)
	p.Children = append(p.Children, c.Index)
}

func (p *Partition) printPartition() {
	log.WithFields(log.Fields{
		"Index":       p.Index,
		"Root":        p.Root,
		"FilesNumber": p.FilesNumber,
		"Children":    p.Children,
		// "Directories": p.Directories,
		// "Extenstion":  p.Extenstion,
		// "ExtenstionH": p.ExtenstionH,
	}).Infof("partition %v info :", p.Index)

}
