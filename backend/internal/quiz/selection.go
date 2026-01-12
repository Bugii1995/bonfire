package quiz

import "time"

//
// -------- Domain snapshots (DB-agnostic) --------
//

// TopicProgress represents the user's mastery state for a topic.
type TopicProgress struct {
	TopicID       string
	Mastery       float64 // 0–100
	CorrectStreak int
	WrongStreak   int
	IsMastered    bool
	LastSeen      time.Time
}

// ReviewItem represents a scheduled spaced-repetition review.
type ReviewItem struct {
	TopicID      string
	NextReviewAt time.Time
}

// Question is lightweight metadata used for selection only.
type Question struct {
	ID         int64
	TopicID    string
	Difficulty int // 1 = easy, 2 = medium, 3 = hard
}

//
// -------- Question purpose (semantic intent) --------
//

type QuestionPurpose string

const (
	PurposeReview    QuestionPurpose = "review"    // spaced repetition
	PurposeReinforce QuestionPurpose = "reinforce" // weak or mistake-prone topic
	PurposeProgress  QuestionPurpose = "progress"  // normal learning
	PurposeStretch   QuestionPurpose = "stretch"   // challenge confident users
)

// SelectedQuestion is the result of the selection algorithm.
type SelectedQuestion struct {
	QuestionID int64
	Purpose    QuestionPurpose
}

//
// -------- Core selection algorithm --------
//

// SelectNextQuestion decides which question should be served next.
//
// Priority order:
// 1. Due spaced-repetition reviews
// 2. Reinforcement of weak or recently mistaken topics
// 3. Normal progression
// 4. Stretch (hard questions for confident topics)
// 5. Safe fallback
//
// This function is deterministic and side-effect free.
func SelectNextQuestion(
	now time.Time,
	progress []TopicProgress,
	reviews []ReviewItem,
	questions []Question,
	recentWrongTopicIDs map[string]bool,
) *SelectedQuestion {

	// 1️⃣ Spaced repetition (highest priority)
	for _, r := range reviews {
		if !r.NextReviewAt.After(now) {
			for _, q := range questions {
				if q.TopicID == r.TopicID && q.Difficulty == 2 {
					return &SelectedQuestion{
						QuestionID: q.ID,
						Purpose:    PurposeReview,
					}
				}
			}
		}
	}

	// 2️⃣ Reinforce weak or mistake-prone topics
	for _, p := range progress {
		if p.Mastery < 40 || recentWrongTopicIDs[p.TopicID] {
			for _, q := range questions {
				if q.TopicID == p.TopicID && q.Difficulty == 1 {
					return &SelectedQuestion{
						QuestionID: q.ID,
						Purpose:    PurposeReinforce,
					}
				}
			}
		}
	}

	// 3️⃣ Normal progression
	for _, p := range progress {
		if p.Mastery >= 40 && p.Mastery < 80 {
			for _, q := range questions {
				if q.TopicID == p.TopicID && q.Difficulty == 2 {
					return &SelectedQuestion{
						QuestionID: q.ID,
						Purpose:    PurposeProgress,
					}
				}
			}
		}
	}

	// 4️⃣ Stretch confident users
	for _, p := range progress {
		if p.Mastery >= 80 {
			for _, q := range questions {
				if q.TopicID == p.TopicID && q.Difficulty == 3 {
					return &SelectedQuestion{
						QuestionID: q.ID,
						Purpose:    PurposeStretch,
					}
				}
			}
		}
	}

	// 5️⃣ Fallback (never block quiz flow)
	for _, q := range questions {
		if q.Difficulty == 2 {
			return &SelectedQuestion{
				QuestionID: q.ID,
				Purpose:    PurposeProgress,
			}
		}
	}

	return nil
}
