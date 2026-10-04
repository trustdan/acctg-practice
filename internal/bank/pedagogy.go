package bank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/trustdan/acctg-practice/curriculum"
	"github.com/trustdan/acctg-practice/internal/domain"
)

type PedagogyQuestion struct {
	QuestionID        string               `json:"question_id"`
	QuestionVersion   int                  `json:"question_version"`
	SettingGroup      string               `json:"setting_group"`
	Contrasts         []domain.QuestionRef `json:"contrasts"`
	DistractorProfile string               `json:"distractor_profile"`
}
type PedagogyPolicy struct {
	SchemaVersion int                `json:"schema_version"`
	PolicyVersion int                `json:"policy_version"`
	Review        ReviewJSON         `json:"review"`
	Questions     []PedagogyQuestion `json:"questions"`
}

func LoadPedagogyPolicy(data []byte) (*PedagogyPolicy, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var p PedagogyPolicy
	if err := dec.Decode(&p); err != nil {
		return nil, err
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("unexpected data after pedagogy policy")
	}
	if p.SchemaVersion != 1 || p.PolicyVersion < 1 || len(p.Questions) == 0 || p.Review.Reviewer == nil || *p.Review.Reviewer == "" || p.Review.ApprovedAt == nil || *p.Review.ApprovedAt == "" {
		return nil, fmt.Errorf("pedagogy policy requires a supported version and review provenance")
	}
	if _, err := time.Parse(time.RFC3339, *p.Review.ApprovedAt); err != nil {
		return nil, fmt.Errorf("invalid pedagogy approval timestamp: %w", err)
	}
	if p.Review.Source == "" {
		return nil, fmt.Errorf("pedagogy policy requires review source")
	}
	seen := map[domain.QuestionRef]bool{}
	for _, q := range p.Questions {
		ref := domain.QuestionRef{QuestionID: q.QuestionID, QuestionVersion: q.QuestionVersion}
		if q.QuestionID == "" || q.QuestionVersion < 1 || q.SettingGroup == "" || seen[ref] {
			return nil, fmt.Errorf("invalid or duplicate pedagogy question %s", q.QuestionID)
		}
		seen[ref] = true
		switch q.DistractorProfile {
		case "", "service_vs_capital", "held_cash_vs_no_entry", "collection_vs_advance", "seller_receipt_vs_no_entry", "earned_advance_vs_receivable":
		default:
			if _, ok := ReviewedAmountDistractor(q.DistractorProfile); !ok {
				return nil, fmt.Errorf("unknown distractor profile %s", q.DistractorProfile)
			}
		}
		partners := map[domain.QuestionRef]bool{}
		for _, z := range q.Contrasts {
			if z == ref || z.QuestionID == "" || z.QuestionVersion < 1 || partners[z] {
				return nil, fmt.Errorf("invalid contrast for %s", q.QuestionID)
			}
			partners[z] = true
		}
	}
	for _, q := range p.Questions {
		for _, z := range q.Contrasts {
			if !seen[z] {
				return nil, fmt.Errorf("missing contrast target %s", z.QuestionID)
			}
		}
	}
	return &p, nil
}

func (p *PedagogyPolicy) ValidateReferences(questions []QuestionJSON) error {
	available := map[domain.QuestionRef]QuestionJSON{}
	for _, q := range questions {
		available[domain.QuestionRef{QuestionID: q.ID, QuestionVersion: q.Version}] = q
	}
	covered := make(map[domain.QuestionRef]bool)
	for _, q := range p.Questions {
		covered[domain.QuestionRef{QuestionID: q.QuestionID, QuestionVersion: q.QuestionVersion}] = true
		source, ok := available[domain.QuestionRef{QuestionID: q.QuestionID, QuestionVersion: q.QuestionVersion}]
		if !ok || !IsActiveForPractice(source.Status) {
			return fmt.Errorf("pedagogy source unavailable: %s", q.QuestionID)
		}
		for _, z := range q.Contrasts {
			target, ok := available[z]
			if !ok || !IsActiveForPractice(target.Status) || target.FamilyID == source.FamilyID {
				return fmt.Errorf("contrast must reference an active distinct event: %s", z.QuestionID)
			}
		}
		valid := q.DistractorProfile == "" || (q.DistractorProfile == "service_vs_capital" || q.DistractorProfile == "held_cash_vs_no_entry") && source.FamilyID == FamilyCashService || q.DistractorProfile == "collection_vs_advance" && source.FamilyID == FamilyCollectReceivable || q.DistractorProfile == "seller_receipt_vs_no_entry" && source.FamilyID == FamilyCustomerAdvance || q.DistractorProfile == "earned_advance_vs_receivable" && source.FamilyID == FamilyEarnAdvance
		if amount, ok := ReviewedAmountDistractor(q.DistractorProfile); ok {
			valid = amount.FamilyID == source.FamilyID && q.DistractorProfile == "amount_"+q.QuestionID
		}
		if !valid {
			return fmt.Errorf("distractor profile conflicts with family: %s", q.QuestionID)
		}
	}
	for ref, q := range available {
		if IsActiveForPractice(q.Status) && !covered[ref] {
			return fmt.Errorf("active seed missing reviewed setting: %s", q.ID)
		}
	}
	return nil
}

var pedagogyOnce sync.Once
var pedagogyPolicy *PedagogyPolicy
var pedagogyError error

func ReviewedPedagogy(q QuestionJSON) (domain.PedagogySnapshot, string, error) {
	pedagogyOnce.Do(func() {
		pedagogyPolicy, pedagogyError = LoadPedagogyPolicy(curriculum.PedagogyJSON)
		if pedagogyError != nil {
			return
		}
		catalog, _, err := LoadAccounts(bytes.NewReader(curriculum.AccountsJSON))
		if err != nil {
			pedagogyError = err
			return
		}
		seeds, err := LoadQuestionBank(bytes.NewReader(curriculum.SeedQuestionsJSON), catalog)
		if err != nil {
			pedagogyError = err
			return
		}
		pedagogyError = pedagogyPolicy.ValidateReferences(seeds.Questions)
	})
	if pedagogyError != nil {
		return domain.PedagogySnapshot{}, "", pedagogyError
	}
	if !RequiresReviewProvenance(q.Status) {
		return domain.PedagogySnapshot{}, "", nil
	}
	for _, p := range pedagogyPolicy.Questions {
		if p.QuestionID == q.ID && p.QuestionVersion == q.Version {
			return domain.PedagogySnapshot{PolicyVersion: pedagogyPolicy.PolicyVersion, SettingGroup: p.SettingGroup, Contrasts: append([]domain.QuestionRef(nil), p.Contrasts...)}, p.DistractorProfile, nil
		}
	}
	return domain.PedagogySnapshot{}, "", nil
}
