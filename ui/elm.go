package ui

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/DebuggerAndrzej/puf/backend"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 && m.cursor%m.paginator.PerPage != 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 && m.cursor%m.paginator.PerPage != m.paginator.PerPage-1 {
				m.cursor++
			}
		case "right", "l":
			if m.cursor < len(m.choices)-m.paginator.PerPage {
				m.cursor += m.paginator.PerPage
			} else {
				m.cursor = len(m.choices) - 1
			}
		case "left", "h":
			if m.cursor > m.paginator.PerPage-1 {
				m.cursor -= m.paginator.PerPage
			} else {
				m.cursor = 0
			}
		case "tab", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "a":
			if len(m.selected) == len(m.choices) {
				for i := range len(m.choices) {
					delete(m.selected, i)
				}
			} else {
				for i := range len(m.choices) {
					m.selected[i] = struct{}{}
				}
			}
		case "enter":
			keys := maps.Keys(m.selected)
			var toUnzip []string
			for _, key := range keys {
				toUnzip = append(toUnzip, m.choices[key])
			}
			if len(toUnzip) > 0 {
				backend.UnzipRequestedFiles(m.archivePath, m.destination, toUnzip)
			} else {
				fmt.Println("No file selected. Nothing got unzipped")
			}
			return m, tea.Quit
		}
	}
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if len(m.choices) == 0 {
		return "No files found, press q to quit."
	}
	var sb strings.Builder
	sb.WriteString("\nWhat should we unpack?\n\n")
	start, end := m.paginator.GetSliceBounds(len(m.choices))
	for i, item := range m.choices[start:end] {
		cursor := " "
		absoluteIndex := i + start
		if m.cursor == absoluteIndex {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selected[absoluteIndex]; ok {
			checked = "x"
		}
		sb.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, item))
	}
	if len(m.choices) > m.paginator.PerPage {
		sb.WriteString("  " + m.paginator.View())
	}
	sb.WriteString("\n\n q: quit • tab/space: (de)select • a: (de)select all • enter: unzip selected \n")
	return sb.String()
}
