package main

import (
	"flag"

	"github.com/DebuggerAndrzej/puf/ui"
)

func main() {
	archivePath := flag.String("f", "", "Path to archive (relative or absolute)")
	searchedRegex := flag.String("r", ".*", "Searched regex in files")
	destination := flag.String("d", "", "Path where to unzip files")
	flag.Parse()

	ui.InitTui(*archivePath, *searchedRegex, *destination)
}
