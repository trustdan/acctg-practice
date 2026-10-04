package tutor

import (
	"regexp"
	"strings"
)

var (
	reLaTeXText = regexp.MustCompile(`\\text\{([^}]*)\}`)
	reLaTeXBold = regexp.MustCompile(`\\textbf\{([^}]*)\}`)
	reLaTeXItal = regexp.MustCompile(`\\textit\{([^}]*)\}`)
	reLaTeXFrac = regexp.MustCompile(`\\frac\{([^}]*)\}\{([^}]*)\}`)
	reNewlines  = regexp.MustCompile(`\n{3,}`)
)

// CleanLaTeXMath converts simplified LaTeX math markup produced by LLMs into readable plain text.
func CleanLaTeXMath(s string) string {
	if s == "" {
		return ""
	}

	// 1. Unwrap \textbf{...} -> **$1**
	for reLaTeXBold.MatchString(s) {
		s = reLaTeXBold.ReplaceAllString(s, "**$1**")
	}

	// 2. Unwrap \textit{...} -> *$1*
	for reLaTeXItal.MatchString(s) {
		s = reLaTeXItal.ReplaceAllString(s, "*$1*")
	}

	// 3. Unwrap \text{...} -> $1
	for reLaTeXText.MatchString(s) {
		s = reLaTeXText.ReplaceAllString(s, "$1")
	}

	// 4. Unwrap \frac{a}{b} -> (a / b)
	for reLaTeXFrac.MatchString(s) {
		s = reLaTeXFrac.ReplaceAllString(s, "($1 / $2)")
	}

	// 5. Common math symbols & operators
	s = strings.ReplaceAll(s, `\times`, "×")
	s = strings.ReplaceAll(s, `\cdot`, "·")
	s = strings.ReplaceAll(s, `\approx`, "≈")
	s = strings.ReplaceAll(s, `\neq`, "≠")
	s = strings.ReplaceAll(s, `\le`, "≤")
	s = strings.ReplaceAll(s, `\ge`, "≥")
	s = strings.ReplaceAll(s, `\pm`, "±")
	s = strings.ReplaceAll(s, `\rightarrow`, "→")
	s = strings.ReplaceAll(s, `\to`, "→")
	s = strings.ReplaceAll(s, `\Rightarrow`, "⇒")

	// 6. Escaped characters: \$ -> $, \% -> %, \& -> &
	s = strings.ReplaceAll(s, `\$`, "$")
	s = strings.ReplaceAll(s, `\%`, "%")
	s = strings.ReplaceAll(s, `\&`, "&")

	// 7. LaTeX equation delimiters
	s = strings.ReplaceAll(s, `\[`, "")
	s = strings.ReplaceAll(s, `\]`, "")
	s = strings.ReplaceAll(s, `\(`, "")
	s = strings.ReplaceAll(s, `\)`, "")
	s = strings.ReplaceAll(s, "$$", "")

	// 8. Collapse excessive blank lines
	s = reNewlines.ReplaceAllString(s, "\n\n")

	return s
}
