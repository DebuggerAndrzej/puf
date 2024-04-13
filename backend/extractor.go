package backend

import (
	"archive/zip"
	"bytes"
	"io"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func GetAllFilesMatchingRegexInArchive(archivePath, searchedRegex string) []string {
	zipFile, err := zip.OpenReader(archivePath)
	if err != nil {
		panic("Unable to open initial zip. Is it path correct?")
	}

	matcher, err := regexp.Compile(searchedRegex)
	if err != nil {
		panic("Regex couldn't compile. Please make sure it's correct regex expression")
	}

	var matchingFiles []string
	return handleFilesInZip(zipFile.File, matcher, matchingFiles)

}

func handleFilesInZip(files []*zip.File, matcher *regexp.Regexp, matchingFiles []string) []string {
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".zip") {
			innerZip, err := file.Open()
			if err != nil {
				panic("bad inner zip")
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
		panic("Couldn't read file from zip")
	}

	reader := bytes.NewReader(buffer)
	sizeToRead, _ := io.Copy(io.Discard, reader)
	zipFile, _ := zip.NewReader(reader, sizeToRead)

	return handleFilesInZip(zipFile.File, matcher, matchingFiles)
}
