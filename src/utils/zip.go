package utils

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(archivePath string, rootDirLikeArchiveName bool) error {
	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	destination, archiveName := filepath.Split(archivePath)
	rootFileName := archive.File[0].Name
	rootFileNameNew := FileNameWithoutExt(archiveName) + string(os.PathSeparator)

	for _, f := range archive.File {
		fileName := f.Name
		if rootDirLikeArchiveName {
			fileName = strings.Replace(f.Name, rootFileName, rootFileNameNew, 1)
		}
		filePath := filepath.Join(destination, fileName)

		if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return errors.New("invalid file path")
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		sourceFile, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(destinationFile, sourceFile); err != nil {
			return err
		}

		destinationFile.Close()
		sourceFile.Close()
	}

	return nil
}