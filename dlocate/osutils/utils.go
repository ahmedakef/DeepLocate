package osutils

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// GetFileMetadata recieves file path and returns FileMetadata
func GetFileMetadata(filePath string) FileMetadata {
	// convert to slash to be sure of how it lookslike
	filePath = filepath.ToSlash(filePath)

	fileInfo, err := os.Stat(filepath.FromSlash(filePath))
	if err != nil {
		log.Fatalf("Error while getting fileInfo : %v", err)
	}

	lastSlash := strings.LastIndex(filePath, "/")
	if lastSlash == -1 {
		log.Fatalf("Error while getting filename : %v", err)
	}
	path := filePath[:lastSlash]

	return GetFileInfo(fileInfo, path)

}

// ListFiles return a list of files and folders directly under the given dir
func ListFiles(path string) []FileMetadata {

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var filesData []FileMetadata
	for _, file := range files {
		filesData = append(filesData, GetFileInfo(file, path))
	}

	return filesData
}

//RemoveContents removes all files and directories under a given path
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
