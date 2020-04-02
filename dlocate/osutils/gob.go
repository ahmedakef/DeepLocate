package osutils

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//ReadGob read an object from a path in gob format
func ReadGob(path string, object interface{}) error {
	dataFile, err := os.Open(filepath.FromSlash(path))

	if err != nil {
		return err
	}
	// ensure to close the file after the fuction end
	defer dataFile.Close()

	var buf io.Reader = dataFile
	//buf, _ = gzip.NewReader(dataFile)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(object)

	if err != nil {
		return err
	}
	return nil
}

//SaveGob save an object from a path in gob format
func SaveGob(object interface{}, path string) error {
	saveAsJSON(object, path)

	// FromSlash converts / to the specific file system separator
	dataFile, err := os.Create(filepath.FromSlash(path))
	if err != nil {
		return err
	}
	defer dataFile.Close()

	var buf io.Writer = dataFile
	//buf = gzip.NewWriter(dataFile)
	enc := gob.NewEncoder(buf)
	enc.Encode(object)

	return nil
}

// saveAsJSON save aby datatype as json for better reading while debugging
func saveAsJSON(data interface{}, path string) {
	// create folder if not exits
	lastdot := strings.LastIndex(path, ".")
	filePath := path[:lastdot] + "json"
	lastSlash := strings.LastIndex(filePath, "/")
	directoryPath := filePath[:lastSlash]
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		os.MkdirAll(directoryPath, os.ModePerm)
	}

	b, _ := json.MarshalIndent(data, "", "\t")
	_ = ioutil.WriteFile(filepath.ToSlash(filePath), b, 0644)
}
