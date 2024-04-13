package ui

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}
