package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func InitTui(archivePath, searchedRegex, destination string) {
	p := tea.NewProgram(initialModel(archivePath, searchedRegex, destination))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
