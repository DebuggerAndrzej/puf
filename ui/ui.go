package ui

import (
	"fmt"
	"os"

	"github.com/DebuggerAndrzej/puf/backend"
	tea "github.com/charmbracelet/bubbletea"
)

func InitTui(archivePath, searchedRegex, destination string) {
	p := tea.NewProgram(initialModel(archivePath, searchedRegex, destination))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel(archivePath, searchedRegex, destination string) model {
	return model{
		choices:     backend.GetAllFilesMatchingRegexInArchive(archivePath, searchedRegex),
		selected:    make(map[int]struct{}),
		archivePath: archivePath,
		destination: destination,
	}
}
