package config

import "os"

func InitDirs() {
	os.Mkdir(MediaPath, os.ModePerm)
	os.Mkdir(TempPath, os.ModePerm)
}