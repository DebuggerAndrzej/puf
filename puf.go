package main

import (
	"flag"
	"fmt"

	"github.com/DebuggerAndrzej/puf/backend"
)

func main() {
	archivePath := flag.String("f", "", "Path to archive (relative or absolute)")
	searchedRegex := flag.String("r", ".*", "Searched regex in files")
	flag.Parse()

	fmt.Println(backend.GetAllFilesMatchingRegexInArchive(*archivePath, *searchedRegex))
}
