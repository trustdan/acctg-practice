package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func TestTitleSlidesThenAutomaticallyRestoresApplication(t *testing.T) {
	m := &Model{State: StateIntro, Intro: &IntroState{}, IntroReturn: StateHelp, Width: 80, Height: 24}
	_, cmd := m.exitIntro()
	m.Intro = nil
	if m.State != StateTitle || cmd == nil {
		t.Fatal("missing animated transition")
	}
	first := m.View()
	for i := 0; i < 30; i++ {
		m.Update(titleTickMsg{})
	}
	centered := ansi.Strip(m.View())
	if first == m.View() || !strings.Contains(centered, "AccounTutor 9000") {
		t.Fatal("title failed to arrive")
	}
	if !strings.Contains(centered, "/__\\") {
		t.Fatal("standard terminal did not show ASCII lettering")
	}
	for _, size := range [][2]int{{80, 24}, {20, 6}, {1, 1}} {
		m.Update(tea.WindowSizeMsg{Width: size[0], Height: size[1]})
		lines := strings.Split(ansi.Strip(m.View()), "\n")
		if len(lines) > size[1] {
			t.Fatal("title exceeds terminal height")
		}
		for _, line := range lines {
			if len(line) > size[0] {
				t.Fatal("title exceeds terminal width")
			}
		}
	}
	for i := 30; i < titleFrames; i++ {
		m.Update(titleTickMsg{})
	}
	if m.State != StateHelp {
		t.Fatal("title did not restore original app screen")
	}
	m.Update(titleTickMsg{})
	if m.State != StateHelp {
		t.Fatal("stale tick changed app state")
	}
}
