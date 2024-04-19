package utils

import (
	"path/filepath"
)

func FileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func SplitExt(fileName string) (string, string) {
	return FileNameWithoutExt(fileName), filepath.Ext(fileName)
}