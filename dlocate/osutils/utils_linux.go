package osutils

import (
	"os"
	"syscall"
)

//GetFileInfo return file metadata sturct using it's path
func GetFileInfo(fileinfo os.FileInfo, path string) FileMetadata {
	var file FileMetadata = NewFileMetadata(fileinfo)
	file.Path = path + "\\" + file.Name
	file.ATime = fileinfo.Sys().(*syscall.Stat_t).Atim
	file.CTime = fileinfo.Sys().(*syscall.Stat_t).Ctim
	file.MTime = fileinfo.Sys().(*syscall.Stat_t).Mtim
	return file
}
