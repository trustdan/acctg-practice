package tui

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
	"github.com/trustdan/acctg-practice/internal/tutor"
)

// RenderMarkdown renders markdown text with ANSI styling suitable for the terminal,
// first passing it through tutor.CleanLaTeXMath to ensure equations render cleanly.
func RenderMarkdown(raw string, wrapWidth int) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}

	cleaned := tutor.CleanLaTeXMath(raw)

	if wrapWidth < 10 {
		wrapWidth = 10
	}

	customStyle := styles.DarkStyleConfig
	zeroMargin := uint(0)
	customStyle.Document.Margin = &zeroMargin
	customStyle.Document.Color = nil // Inherits box foreground color; avoids redundant per-word ANSI resets

	// Remove markdown '#' hash prefixes so headings render cleanly without raw markdown symbols
	empty := ""
	customStyle.H1.Prefix = empty
	customStyle.H1.Suffix = empty
	customStyle.H1.BackgroundColor = nil
	customStyle.H2.Prefix = empty
	customStyle.H3.Prefix = empty
	customStyle.H4.Prefix = empty
	customStyle.H5.Prefix = empty
	customStyle.H6.Prefix = empty

	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(customStyle),
		glamour.WithWordWrap(wrapWidth),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		return cleaned
	}

	rendered, err := r.Render(cleaned)
	if err != nil {
		return cleaned
	}

	return strings.TrimSpace(rendered)
}

// renderTutorContent renders and caches markdown explanations and hints for the model.
func (m *Model) renderTutorContent(raw string, wrapWidth int) string {
	if raw == "" {
		return ""
	}
	if raw == m.renderedHintRaw && wrapWidth == m.renderedHintWidth && m.renderedHint != "" {
		return m.renderedHint
	}
	m.renderedHintRaw = raw
	m.renderedHintWidth = wrapWidth
	m.renderedHint = RenderMarkdown(raw, wrapWidth)
	return m.renderedHint
}
