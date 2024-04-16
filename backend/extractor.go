package backend

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func GetAllFilesMatchingRegexInArchive(archivePath, searchedRegex string) []string {
	zipFile, err := zip.OpenReader(archivePath)
	if err != nil {
		printMessageAndExit("Unable to open initial zip. Is it's path correct?")
	}

	matcher, err := regexp.Compile(searchedRegex)
	if err != nil {
		printMessageAndExit("Regex couldn't compile. Please make sure it's correct regex expression")
	}

	var matchingFiles []string
	return handleFilesInZip(zipFile.File, matcher, matchingFiles)
}

func UnzipRequestedFiles(archivePath, destination string, filenames []string) {
	zipFile, err := zip.OpenReader(archivePath)
	if err != nil {
		printMessageAndExit("Unable to open initial zip. Is it's path correct?")
	}

	unzipFilesInZip(zipFile.File, filenames, destination)
	fmt.Println("\nUnzipped requested files")
}

func handleFilesInZip(files []*zip.File, matcher *regexp.Regexp, matchingFiles []string) []string {
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".zip") {
			innerZip, err := file.Open()
			if err != nil {
				printMessageAndExit("bad inner zip")
			}
			matchingFiles = getMatchingFilesFromZip(innerZip, matcher, matchingFiles)
		}
		if !file.FileInfo().IsDir() {
			if matcher.MatchString(file.Name) && !slices.Contains(matchingFiles, filepath.Base(file.Name)) {
				matchingFiles = append(matchingFiles, filepath.Base(file.Name))
			}
		}
	}
	return matchingFiles
}

func getMatchingFilesFromZip(openedZip io.ReadCloser, matcher *regexp.Regexp, matchingFiles []string) []string {
	defer openedZip.Close()
	buffer, err := io.ReadAll(openedZip)
	if err != nil {
		printMessageAndExit("Couldn't read file from zip")
	}

	reader := bytes.NewReader(buffer)
	sizeToRead, _ := io.Copy(io.Discard, reader)
	zipFile, _ := zip.NewReader(reader, sizeToRead)

	return handleFilesInZip(zipFile.File, matcher, matchingFiles)
}

func unzipFilesInZip(files []*zip.File, filenames []string, destination string) {
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".zip") {
			innerZip, err := file.Open()
			if err != nil {
				panic("bad inner zip")
			}
			getZipFile(innerZip, filenames, destination)
		}
		if !file.FileInfo().IsDir() {
			if slices.Contains(filenames, filepath.Base(file.Name)) {
				unzipFile(*file, destination)
			}
		}
	}
}

func getZipFile(openedZip io.ReadCloser, filenames []string, destination string) {
	defer openedZip.Close()
	buffer, err := io.ReadAll(openedZip)
	if err != nil {
		printMessageAndExit("Couldn't read file from zip")
	}

	reader := bytes.NewReader(buffer)
	sizeToRead, _ := io.Copy(io.Discard, reader)
	zipFile, _ := zip.NewReader(reader, sizeToRead)

	unzipFilesInZip(zipFile.File, filenames, destination)
}

func unzipFile(file zip.File, destination string) {
	var filePath string
	if destination == "" {
		filePath = filepath.Base(file.Name)
	} else {
		filePath = filepath.Join(destination, filepath.Base(file.Name))
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			printMessageAndExit("Couldn't create requested destination dir")
		}
		filePath = filepath.Join(destination, filepath.Base(file.Name))
	}

	dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		printMessageAndExit("Couldn't open file to save")
	}
	fileInArchive, err := file.Open()
	if err != nil {
		printMessageAndExit("Couldn't open archive file requested to save")
	}
	if _, err := io.Copy(dstFile, fileInArchive); err != nil {
		printMessageAndExit("Couldn't copy archive data to file")
	}
	dstFile.Close()
	fileInArchive.Close()
}

func printMessageAndExit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
