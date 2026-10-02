package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const titleFrames = 90

type titleTickMsg struct{}

func titleTick() tea.Cmd {
	return tea.Tick(time.Second/30, func(time.Time) tea.Msg { return titleTickMsg{} })
}

// Small ASCII glyphs keep the full title legible on a standard terminal.
var titleGlyphs = map[rune][]string{
	'A': {" /\\ ", "/__\\", "|  |", "|  |"},
	'C': {" ___", "|   ", "|   ", "|___"},
	'O': {" __ ", "|  |", "|  |", "|__|"},
	'U': {"    ", "|  |", "|  |", "|__|"},
	'N': {"    ", "|\\ |", "| \\|", "|  |"},
	'T': {"____", " || ", " || ", " || "},
	'R': {" __ ", "|__|", "|\\  ", "| \\ "},
	'9': {" __ ", "|__|", "   |", " __|"},
	'0': {" __ ", "| /|", "|/ |", "|__|"},
	' ': {"  ", "  ", "  ", "  "},
}

func titleArt() []string {
	rows := make([]string, 4)
	for _, ch := range "ACCOUNTUTOR 9000" {
		for i, glyph := range titleGlyphs[ch] {
			rows[i] += glyph + " "
		}
	}
	return rows
}

func (m *Model) renderTitle() string {
	width, height := m.Width, m.Height
	if width < 1 {
		width = 80
	}
	if height < 1 {
		height = 24
	}
	art := titleArt()
	if width < len(art[0]) {
		art = []string{"AccounTutor 9000", "=== DEBITS / CREDITS ==="}
	}
	target := (width - len(art[0])) / 2
	if target < 0 {
		target = 0
	}
	frame := m.TitleFrame
	if frame > 30 {
		frame = 30
	}
	// Ease out as the title settles into the center in the first second.
	remaining := 30 - frame
	x := target + (width-target)*remaining*remaining/900
	lines := make([]string, height)
	top := (height - len(art) - 4) / 2
	if top < 0 {
		top = 0
	}
	for i, row := range art {
		if top+i >= height {
			break
		}
		line := strings.Repeat(" ", x) + row
		if len(line) > width {
			line = line[:width]
		}
		lines[top+i] = lipgloss.NewStyle().Foreground(lipgloss.Color("#22D3EE")).Bold(true).Render(line)
	}
	if m.TitleFrame >= 30 {
		for i, text := range []string{"AccounTutor 9000", "AI creates. Rules validate. Machine drills. History adapts."} {
			y := top + len(art) + 1 + i
			if y >= height {
				break
			}
			if len(text) > width {
				text = text[:width]
			}
			lines[y] = strings.Repeat(" ", (width-len(text))/2) + text
		}
	}
	return strings.Join(lines, "\n")
}
