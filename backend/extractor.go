package backend

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var recoverableErrorsEncountered []string

func GetAllFilesMatchingRegexInArchive(archivePath, searchedRegex string) map[string]string {
	zipFile, err := zip.OpenReader(archivePath)
	if err != nil {
		printMessageAndExit("Unable to open initial zip. Is it's path correct?")
	}

	matcher, err := regexp.Compile(searchedRegex)
	if err != nil {
		printMessageAndExit("Regex couldn't compile. Please make sure it's correct regex expression")
	}

	matchingFiles := make(map[string]string)
	handleFilesInZip(zipFile.File, matcher, matchingFiles, "")
	return matchingFiles
}

func UnzipRequestedFiles(archivePath, destination string, filenames []string) {
	zipFile, err := zip.OpenReader(archivePath)
	if err != nil {
		printMessageAndExit("Unable to open initial zip. Is it's path correct?")
	}

	unzipFilesInZip(zipFile.File, filenames, destination, "")
	if len(recoverableErrorsEncountered) > 0 {
		fmt.Printf("\u001b[33m\nWARNING: %s\u001b[0m\n", strings.Join(recoverableErrorsEncountered, "\nWARNING: "))
	}
	fmt.Println("\nUnzipped requested files")
}

func updateMatchingFiles(matchingFiles map[string]string, key, filename, dirPrefix string) {
	if dirPrefix != "" {
		matchingFiles[key] = filepath.Join(dirPrefix, filename)
	} else {
		matchingFiles[key] = filename
	}
}

func handleFilesInZip(files []*zip.File, matcher *regexp.Regexp, matchingFiles map[string]string, dirPrefix string) {
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".zip") {
			innerZip, err := file.Open()
			if err != nil {
				recoverableErrorsEncountered = append(
					recoverableErrorsEncountered,
					fmt.Sprintf("Ommiting file %s. File couldn't be opened", file.Name),
				)
				continue
			}
			getMatchingFilesFromZip(innerZip, matcher, matchingFiles, file.Name)
		}
		if !file.FileInfo().IsDir() {
			if matcher.MatchString(file.Name) {
				if _, contains := matchingFiles[filepath.Base(file.Name)]; contains {
					for i := 1; ; i++ {
						newFilename := fmt.Sprintf("%s(%d)", filepath.Base(file.Name), i)
						if _, contains := matchingFiles[newFilename]; !contains {
							updateMatchingFiles(matchingFiles, newFilename, file.Name, dirPrefix)
							break
						}
					}
				} else {
					updateMatchingFiles(matchingFiles, filepath.Base(file.Name), file.Name, dirPrefix)

				}
			}
		}
	}
}

func getMatchingFilesFromZip(
	openedZip io.ReadCloser,
	matcher *regexp.Regexp,
	matchingFiles map[string]string,
	filename string,
) {
	defer openedZip.Close()
	buffer, err := io.ReadAll(openedZip)
	if err != nil {
		recoverableErrorsEncountered = append(
			recoverableErrorsEncountered,
			fmt.Sprintf("Ommiting file %s. File couldn't be read", filename),
		)
		return
	}

	reader := bytes.NewReader(buffer)
	zipFile, err := zip.NewReader(reader, int64(len(buffer)))
	if err != nil {
		recoverableErrorsEncountered = append(
			recoverableErrorsEncountered,
			fmt.Sprintf("Ommiting file %s. File couldn't be reopened from buffer", filename),
		)
		return
	}

	handleFilesInZip(zipFile.File, matcher, matchingFiles, filename)
}

func unzipFilesInZip(files []*zip.File, filenames []string, destination, dirPrefix string) {
	for _, file := range files {
		if strings.HasSuffix(file.Name, ".zip") {
			innerZip, err := file.Open()
			if err != nil {
				continue
			}
			getZipFile(innerZip, filenames, destination, file.Name)
		}
		if !file.FileInfo().IsDir() {
			if dirPrefix != "" {
				if slices.Contains(filenames, filepath.Join(dirPrefix, file.Name)) {
					unzipFile(*file, destination)
				}
			} else {
				if slices.Contains(filenames, file.Name) {
					unzipFile(*file, destination)
				}
			}
		}
	}
}

func getZipFile(openedZip io.ReadCloser, filenames []string, destination, fileName string) {
	defer openedZip.Close()
	buffer, err := io.ReadAll(openedZip)
	if err != nil {
		return
	}

	reader := bytes.NewReader(buffer)
	zipFile, err := zip.NewReader(reader, int64(len(buffer)))

	if err != nil {
		return
	}

	unzipFilesInZip(zipFile.File, filenames, destination, fileName)
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

	var dstFile *os.File
	if _, err := os.Stat(filePath); err == nil {
		for i := 1; ; i++ {
			newFilePath := fmt.Sprintf(
				"%s(%d)%s",
				strings.TrimSuffix(filePath, filepath.Ext(filePath)),
				i,
				filepath.Ext(filePath),
			)
			if _, err := os.Stat(newFilePath); errors.Is(err, os.ErrNotExist) {
				dstFile, err = os.OpenFile(newFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
				break
			}
		}
	} else {
		dstFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			printMessageAndExit("Couldn't open file to save")
		}
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
	if len(recoverableErrorsEncountered) > 0 {
		fmt.Printf("\u001b[33m\nWARNING: %s\u001b[0m\n", strings.Join(recoverableErrorsEncountered, "\nWARNING: "))
	}
	fmt.Printf("\u001b[31m%s\u001b[0m\n", msg)
	os.Exit(1)
}
