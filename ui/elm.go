package ui

import (
	"fmt"

	"golang.org/x/exp/maps"

	"github.com/DebuggerAndrzej/puf/backend"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "tab", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			keys := maps.Keys(m.selected)
			var toUnzip []string
			for _, key := range keys {
				toUnzip = append(toUnzip, m.choices[key])
			}
			backend.UnzipRequestedFiles(m.archivePath, toUnzip)
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "What should we unpack?\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\nPress enter to unzip selected or q to quit.\n"
	return s
}
