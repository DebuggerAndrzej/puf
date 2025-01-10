package ui

import (
	"maps"
	"slices"
	"sort"

	"github.com/DebuggerAndrzej/puf/backend"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	choices     []string
	items       map[string]string
	archivePath string
	cursor      int
	selected    map[int]struct{}
	destination string
	paginator   paginator.Model
}

func initialModel(archivePath, searchedRegex, destination string) model {
	items := backend.GetAllFilesMatchingRegexInArchive(archivePath, searchedRegex)
	choices := slices.Collect(maps.Keys(items))
	sort.Strings(choices)
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 15
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(items))
	return model{
		choices:     choices,
		items:       items,
		selected:    make(map[int]struct{}),
		archivePath: archivePath,
		destination: destination,
		paginator:   p,
	}
}
