package main

import "time"

// FileInfo conatins basic info about explored files
type FileInfo struct {
	FileName  string
	LinuxPath string
	ModTime   time.Time
	Type      FileType
}

// FileType like folder, audio, book ... etc
type FileType string
