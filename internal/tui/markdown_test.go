package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestHeadersHaveNoHashes(t *testing.T) {
	input := `## Header Two
Some text
### Header Three
More text`
	rendered := RenderMarkdown(input, 60)
	if strings.Contains(rendered, "##") || strings.Contains(rendered, "###") {
		t.Errorf("rendered output still contains markdown hashes: %s", rendered)
	}
	plain := ansi.Strip(rendered)
	if !strings.Contains(plain, "Header Two") || !strings.Contains(plain, "Header Three") {
		t.Errorf("missing header text in: %s", plain)
	}
}

func TestRenderMarkdownWithLaTeXAndFormatting(t *testing.T) {
	input := `The workshop is being earned **today**, but the cash was received **last month**. Accrual accounting separates those two events:

### Last month: cash received before work was performed
The company received $100 but still owed the customer a workshop:

- **Debit Cash $100** - Cash is an asset, and assets increase on the debit/left side.
- **Credit Unearned Revenue $100** - Unearned Revenue is a liability because the company had an obligation to provide the workshop.

At that point, the company had more cash, but it had not yet earned revenue.

### Today: the workshop is completed
The obligation is now satisfied:

- **Debit Unearned Revenue $100** - the liability decreases.
- **Credit Service Revenue $100** - revenue increases, which increases equity.

No cash changes hands today, so assets remain unchanged. The accounting equation effect is:

\[
\text{Assets unchanged} = \text{Liabilities decrease } \$100 + \text{Equity increase } \$100
\]`

	rendered := RenderMarkdown(input, 80)
	if strings.Contains(rendered, `\[`) || strings.Contains(rendered, `\]`) {
		t.Errorf("equation delimiters not removed: %s", rendered)
	}
	if strings.Contains(rendered, `\text`) {
		t.Errorf("\\text not removed: %s", rendered)
	}
	if strings.Contains(rendered, `\$`) {
		t.Errorf("\\$ not replaced: %s", rendered)
	}
	if strings.Contains(rendered, "###") {
		t.Errorf("heading still contains ### hashes: %s", rendered)
	}
	if !strings.Contains(rendered, "Debit Cash $100") {
		t.Errorf("expected bold account title preserved: %s", rendered)
	}
	if !strings.Contains(rendered, "Assets unchanged = Liabilities decrease  $100 + Equity increase  $100") {
		t.Errorf("expected clean equation in output: %s", rendered)
	}
}
