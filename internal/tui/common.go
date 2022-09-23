package tui

import (
	"golang.org/x/term"
	"os"
)

var (
	windowWidth  int
	windowHeight int
)

func init() {
	windowWidth, windowHeight, _ = term.GetSize(int(os.Stdout.Fd()))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
