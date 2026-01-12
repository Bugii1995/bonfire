package quiz

import (
	"testing"
	"time"
)

func TestUpdateMastery_CorrectMedium(t *testing.T) {
	now := time.Now()

	current := TopicProgress{
		TopicID:       "present_simple",
		Mastery:       40,
		CorrectStreak: 0,
		WrongStreak:   0,
		LastSeen:      now.Add(-1 * time.Hour),
	}

	input := MasteryUpdateInput{
		WasCorrect:  true,
		Difficulty:  2,
		AnsweredAt:  now,
		CurrentTime: now,
	}

	result := UpdateMastery(current, input)

	if result.Mastery <= current.Mastery {
		t.Fatalf("expected mastery to increase, got %v", result.Mastery)
	}

	if result.CorrectStreak != 1 {
		t.Fatalf("expected correct streak = 1, got %d", result.CorrectStreak)
	}

	if result.WrongStreak != 0 {
		t.Fatalf("expected wrong streak reset to 0")
	}
}

func TestUpdateMastery_WrongHardPenalty(t *testing.T) {
	now := time.Now()

	current := TopicProgress{
		TopicID:       "articles",
		Mastery:       60,
		CorrectStreak: 2,
		WrongStreak:   1,
		LastSeen:      now,
	}

	input := MasteryUpdateInput{
		WasCorrect:  false,
		Difficulty:  3,
		AnsweredAt:  now,
		CurrentTime: now,
	}

	result := UpdateMastery(current, input)

	if result.Mastery >= current.Mastery {
		t.Fatalf("expected mastery to decrease, got %v", result.Mastery)
	}

	if result.CorrectStreak != 0 {
		t.Fatalf("expected correct streak reset")
	}

	if result.WrongStreak != 2 {
		t.Fatalf("expected wrong streak = 2")
	}
}

func TestUpdateMastery_DecayApplied(t *testing.T) {
	now := time.Now()

	current := TopicProgress{
		TopicID:  "conditionals",
		Mastery:  80,
		LastSeen: now.AddDate(0, 0, -20), // 20 days ago
	}

	input := MasteryUpdateInput{
		WasCorrect:  true,
		Difficulty:  1,
		AnsweredAt:  now,
		CurrentTime: now,
	}

	result := UpdateMastery(current, input)

	if result.Mastery >= 80 {
		t.Fatalf("expected decay before update, got %v", result.Mastery)
	}
}

func TestUpdateMastery_MasteredFlag(t *testing.T) {
	now := time.Now()

	current := TopicProgress{
		TopicID: "future_tense",
		Mastery: 98,
	}

	input := MasteryUpdateInput{
		WasCorrect:  true,
		Difficulty:  3,
		AnsweredAt:  now,
		CurrentTime: now,
	}

	result := UpdateMastery(current, input)

	if !result.IsMastered {
		t.Fatalf("expected topic to be marked mastered")
	}

	if result.Mastery != 100 {
		t.Fatalf("expected mastery clamped to 100, got %v", result.Mastery)
	}
}
