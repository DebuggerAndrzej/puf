package ui

import (
	"fmt"
	"os"

	"github.com/DebuggerAndrzej/puf/backend"
	tea "github.com/charmbracelet/bubbletea"
)

func InitTui(archivePath, searchedRegex string) {
	p := tea.NewProgram(initialModel(archivePath, searchedRegex))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel(archivePath, searchedRegex string) model {
	return model{
		choices:  backend.GetAllFilesMatchingRegexInArchive(archivePath, searchedRegex),
		selected: make(map[int]struct{}),
	}
}
