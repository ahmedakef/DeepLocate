package main

import "time"

// FileInfo conatins basic info about explored files
type FileInfo struct {
	FileName string
	ModTime  time.Time
}
