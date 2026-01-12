package quiz

import (
	"testing"
	"time"
)

func TestSelectNextQuestion_ReviewPriority(t *testing.T) {
	now := time.Now()

	reviews := []ReviewItem{
		{
			TopicID:      "past_simple",
			NextReviewAt: now.Add(-1 * time.Hour),
		},
	}

	questions := []Question{
		{ID: 1, TopicID: "past_simple", Difficulty: 2},
	}

	result := SelectNextQuestion(
		now,
		nil,
		reviews,
		questions,
		nil,
	)

	if result == nil {
		t.Fatal("expected a question")
	}

	if result.Purpose != PurposeReview {
		t.Fatalf("expected review purpose, got %s", result.Purpose)
	}
}

func TestSelectNextQuestion_ReinforceWeakTopic(t *testing.T) {
	now := time.Now()

	progress := []TopicProgress{
		{
			TopicID: "articles",
			Mastery: 25,
		},
	}

	questions := []Question{
		{ID: 2, TopicID: "articles", Difficulty: 1},
	}

	result := SelectNextQuestion(
		now,
		progress,
		nil,
		questions,
		nil,
	)

	if result.Purpose != PurposeReinforce {
		t.Fatalf("expected reinforce purpose")
	}
}

func TestSelectNextQuestion_Progression(t *testing.T) {
	now := time.Now()

	progress := []TopicProgress{
		{
			TopicID: "present_continuous",
			Mastery: 55,
		},
	}

	questions := []Question{
		{ID: 3, TopicID: "present_continuous", Difficulty: 2},
	}

	result := SelectNextQuestion(
		now,
		progress,
		nil,
		questions,
		nil,
	)

	if result.Purpose != PurposeProgress {
		t.Fatalf("expected progress purpose")
	}
}

func TestSelectNextQuestion_Stretch(t *testing.T) {
	now := time.Now()

	progress := []TopicProgress{
		{
			TopicID: "relative_clauses",
			Mastery: 90,
		},
	}

	questions := []Question{
		{ID: 4, TopicID: "relative_clauses", Difficulty: 3},
	}

	result := SelectNextQuestion(
		now,
		progress,
		nil,
		questions,
		nil,
	)

	if result.Purpose != PurposeStretch {
		t.Fatalf("expected stretch purpose")
	}
}

func TestSelectNextQuestion_Fallback(t *testing.T) {
	now := time.Now()

	questions := []Question{
		{ID: 5, TopicID: "misc", Difficulty: 2},
	}

	result := SelectNextQuestion(
		now,
		nil,
		nil,
		questions,
		nil,
	)

	if result == nil {
		t.Fatal("expected fallback question")
	}
}
