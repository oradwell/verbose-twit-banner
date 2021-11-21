package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func getJpegFilesInDirectory(directory string) []string {
	var files []string

	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		fileExt := filepath.Ext(path)

		if !info.IsDir() && (fileExt == ".jpg" || fileExt == ".jpeg") {
			files = append(files, path)
		}

		return nil
	})

	return files
}

func getRandIndexInArray(array []string) (int, error) {
	if len(array) == 0 {
		return -1, errors.New("Empty list")
	}

	rand.Seed(time.Now().Unix())

	return rand.Intn(len(array)), nil
}

func GetJpegPathInDirectory(directory string) (string, error) {
	files := getJpegFilesInDirectory(directory)

	rInd, err := getRandIndexInArray(files)
	if err != nil {
		return "", errEmptyDir
	}

	return files[rInd], nil
}
