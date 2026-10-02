package curriculum

import (
	_ "embed"
)

// AccountsJSON contains the canonical accounts definitions.
//
//go:embed accounts.json
var AccountsJSON []byte

// SeedQuestionsJSON contains the seed question bank templates.
//
//go:embed seed-questions.json
var SeedQuestionsJSON []byte
