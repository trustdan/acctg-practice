package domain

type QuestionRef struct {
	QuestionID      string `json:"question_id"`
	QuestionVersion int    `json:"question_version"`
}

// PedagogySnapshot is immutable review/evidence metadata, never an answer key.
type PedagogySnapshot struct {
	PolicyVersion int           `json:"policy_version"`
	SettingGroup  string        `json:"setting_group"`
	Contrasts     []QuestionRef `json:"contrasts,omitempty"`
	Remediation   bool          `json:"remediation,omitempty"`
}
