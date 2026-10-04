package tutor

import (
	"strings"
	"testing"
)

func TestCleanLaTeXMath(t *testing.T) {
	input := `The accounting equation effect is:

\[
\text{Assets unchanged} = \text{Liabilities decrease } \$100 + \text{Equity increase } \$100
\]`

	cleaned := CleanLaTeXMath(input)
	if strings.Contains(cleaned, `\[`) || strings.Contains(cleaned, `\]`) {
		t.Errorf("equation delimiters not removed: %s", cleaned)
	}
	if strings.Contains(cleaned, `\text`) {
		t.Errorf("\\text not removed: %s", cleaned)
	}
	if strings.Contains(cleaned, `\$`) {
		t.Errorf("\\$ not replaced: %s", cleaned)
	}
	if !strings.Contains(cleaned, "Assets unchanged = Liabilities decrease  $100 + Equity increase  $100") {
		t.Errorf("unexpected cleaned text: %s", cleaned)
	}
}

func TestCleanLaTeXMathSymbolsAndFractions(t *testing.T) {
	input := `Formula: \(\frac{\text{Net Income}}{\text{Total Assets}} \times 100 \approx 15\%\)`
	cleaned := CleanLaTeXMath(input)
	if strings.Contains(cleaned, `\frac`) || strings.Contains(cleaned, `\times`) || strings.Contains(cleaned, `\approx`) {
		t.Errorf("failed to convert operators/fractions: %s", cleaned)
	}
	if !strings.Contains(cleaned, "(Net Income / Total Assets) × 100 ≈ 15%") {
		t.Errorf("unexpected output: %s", cleaned)
	}
}
