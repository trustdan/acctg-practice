package mastery

import (
	"sort"
	"time"

	"github.com/trustdan/acctg-practice/internal/bank"
	"github.com/trustdan/acctg-practice/internal/domain"
)

// InstanceConceptKey groups attempts by question instance and concept.
type InstanceConceptKey struct {
	InstanceID string
	ConceptID  string
}

// RebuildProjections aggregates immutable learner attempts into ConceptStats.
// Invariants enforced:
// 1. One question instance adds at most 1 count to a concept's independent evidence.
// 2. Unassisted first error stays an error even if a hinted retry succeeds.
// 3. Hinted/reference attempts are recorded as assisted evidence and do not inflate Alpha or reset decay.
// 4. Both scaffold reductions require two reviewed successful settings and delayed transfer.
// 5. Half-life stays fixed unless supported by >=2 spaced successes that change the reviewed setting.
// 6. Poor performance restores scaffolding to Level 0 (Full).
func RebuildProjections(attempts []domain.Attempt) map[string]*ConceptStats {
	// Sort attempts chronologically by AnsweredAt to ensure replay order
	sorted := make([]domain.Attempt, len(attempts))
	copy(sorted, attempts)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].AnsweredAt.Before(sorted[j].AnsweredAt)
	})

	projections := make(map[string]*ConceptStats)
	getOrCreate := func(cid string) *ConceptStats {
		s, ok := projections[cid]
		if !ok {
			s = NewConceptStats(cid)
			projections[cid] = s
		}
		return s
	}

	// Track instances already evaluated for each concept to avoid evidence multiplication
	seenInstanceConcepts := make(map[InstanceConceptKey]bool)

	// Track latest result per concept to detect recent poor performance
	latestHadError := make(map[string]bool)
	settings := make(map[string]map[string]bool)
	lastSetting := make(map[string]string)
	lastExposure := make(map[string]time.Time)
	seenAttempts := make(map[string]bool)

	for _, att := range sorted {
		if att.AttemptID != "" {
			if seenAttempts[att.AttemptID] {
				continue
			}
			seenAttempts[att.AttemptID] = true
		}
		if att.ConceptID == "" {
			continue
		}
		stats := getOrCreate(att.ConceptID)
		key := InstanceConceptKey{InstanceID: att.InstanceID, ConceptID: att.ConceptID}

		if stats.FirstAttemptAt == nil {
			firstTime := att.AnsweredAt
			stats.FirstAttemptAt = &firstTime
		}

		// Reference consultation tracking
		if att.ReferenceUsed || att.Assistance == domain.AssistanceReference {
			stats.ReferenceUses++
		}

		if !seenInstanceConcepts[key] {
			// First attempt for this concept in this instance
			seenInstanceConcepts[key] = true

			isAssisted := att.ReferenceUsed || att.Assistance != domain.AssistanceNone

			if !isAssisted {
				// Pure independent attempt
				stats.IndependentAttempts++
				answeredTime := att.AnsweredAt
				priorExposure, exposed := lastExposure[att.ConceptID]
				spaced := stats.LastIndependentAttemptAt != nil && answeredTime.Sub(*stats.LastIndependentAttemptAt) >= MinDelayedRetrievalGap && (!exposed || answeredTime.Sub(priorExposure) >= MinDelayedRetrievalGap)

				if att.IsCorrect {
					stats.Alpha += 1.0
					stats.IndependentSuccesses++
					latestHadError[att.ConceptID] = false

					groups := settings[att.ConceptID]
					if groups == nil {
						groups = make(map[string]bool)
						settings[att.ConceptID] = groups
					}
					// Only reviewed, snapshotted groups contribute transfer evidence.
					if att.PedagogyVersion > 0 && att.SettingGroup != "" {
						if spaced {
							if prior := lastSetting[att.ConceptID]; prior != "" && prior != att.SettingGroup {
								stats.TransferDelayedSuccesses++
							}
						}
						groups[att.SettingGroup] = true
						stats.SuccessfulSettings = len(groups)
					}
					lastSetting[att.ConceptID] = ""
					if att.PedagogyVersion > 0 {
						lastSetting[att.ConceptID] = att.SettingGroup
					}
					// Delayed retrieval verification:
					// Requires a prior learning attempt on this concept separated by at least MinDelayedRetrievalGap
					if spaced {
						stats.DelayedSuccesses++
						stats.HasDelayedRetrieval = true
					}

					succTime := answeredTime
					stats.LastSuccessAt = &succTime

					// Half-life bounded adjustment with delayed evidence
					if stats.TransferDelayedSuccesses >= 2 {
						hl := DefaultHalfLifeDays * (1.0 + 0.5*float64(stats.TransferDelayedSuccesses-1))
						if hl < MinHalfLifeDays {
							hl = MinHalfLifeDays
						}
						if hl > MaxHalfLifeDays {
							hl = MaxHalfLifeDays
						}
						stats.HalfLifeDays = hl
					} else {
						stats.HalfLifeDays = DefaultHalfLifeDays
					}

				} else {
					stats.Beta += 1.0
					latestHadError[att.ConceptID] = true
				}

				stats.LastIndependentAttemptAt = &answeredTime

			} else {
				// Assisted attempt before independent answer
				stats.AssistedAttempts++
				if !att.IsCorrect || att.Assistance == domain.AssistanceRevealed || att.Assistance == domain.AssistanceRetry {
					latestHadError[att.ConceptID] = true
				}
			}
		} else {
			// Subsequent attempt in the same instance (e.g. retry after initial error, or later stage with same concept)
			if att.Assistance == domain.AssistanceContrast || att.Assistance == domain.AssistanceRetry || att.Assistance == domain.AssistanceHinted || att.Assistance == domain.AssistanceRevealed || att.ReferenceUsed {
				stats.AssistedAttempts++
			}
			if !att.IsCorrect || att.Assistance == domain.AssistanceRevealed {
				latestHadError[att.ConceptID] = true
			}
		}
		lastExposure[att.ConceptID] = att.AnsweredAt

	}

	// Compute final per-concept scaffold level based on accumulated evidence
	for cid, stats := range projections {
		if stats.IsNew() || latestHadError[cid] || stats.BaseScore() < 0.60 {
			// Poor performance or unmastered: restore to full scaffolding
			stats.ScaffoldLevel = domain.ScaffoldFull
			continue
		}

		// Level 2 (Faded) requires high accuracy, multiple successes, AND delayed retrieval
		if stats.BaseScore() >= 0.80 && stats.IndependentSuccesses >= 3 && stats.TransferDelayedSuccesses >= 1 && stats.SuccessfulSettings >= 2 {
			stats.ScaffoldLevel = domain.ScaffoldFaded
		} else if stats.BaseScore() >= 0.60 && stats.IndependentSuccesses >= 2 && stats.SuccessfulSettings >= 2 && stats.TransferDelayedSuccesses >= 1 {
			// Level 1 (Intermediate)
			stats.ScaffoldLevel = domain.ScaffoldIntermediate
		} else {
			stats.ScaffoldLevel = domain.ScaffoldFull
		}
	}

	return projections
}

// ComputeQuestionScaffoldLevel determines the appropriate scaffold level for a question template
// based on the weakest concept tested and the active session intensity.
func ComputeQuestionScaffoldLevel(q bank.QuestionJSON, projections map[string]*ConceptStats, intensity SessionIntensity) domain.ScaffoldLevel {
	if len(q.Concepts) == 0 {
		return domain.ScaffoldFull
	}

	maxAllowed := domain.ScaffoldFaded
	if intensity == IntensityIntensive {
		// Under intensive drill, retain at least intermediate scaffolding
		maxAllowed = domain.ScaffoldIntermediate
	}

	minLevel := domain.ScaffoldFaded
	for _, cid := range q.Concepts {
		stats, ok := projections[cid]
		if !ok || stats.IsNew() {
			return domain.ScaffoldFull
		}
		if stats.ScaffoldLevel < minLevel {
			minLevel = stats.ScaffoldLevel
		}
	}

	if minLevel > maxAllowed {
		minLevel = maxAllowed
	}
	return minLevel
}
