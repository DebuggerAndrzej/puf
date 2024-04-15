package ui

type model struct {
	choices     []string
	archivePath string
	cursor      int
	selected    map[int]struct{}
	destination string
}
