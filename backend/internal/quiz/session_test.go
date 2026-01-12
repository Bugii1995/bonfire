package quiz

import (
	"testing"
	"time"
)

func TestNewSessionCreatesProgressMap(t *testing.T) {
	questions := []Question{
		{ID: 1, TopicID: "present_simple", Difficulty: 2},
	}

	progress := []TopicProgress{
		{TopicID: "present_simple", Mastery: 40},
	}

	session := NewSession(questions, progress, nil)

	if session.ID == "" {
		t.Fatal("expected session ID to be set")
	}

	if len(session.Progress) != 1 {
		t.Fatalf("expected 1 topic in progress, got %d", len(session.Progress))
	}
}

func TestSessionNextQuestionReturnsQuestion(t *testing.T) {
	now := time.Now()

	questions := []Question{
		{ID: 1, TopicID: "articles", Difficulty: 2},
	}

	progress := []TopicProgress{
		{TopicID: "articles", Mastery: 50},
	}

	session := NewSession(questions, progress, nil)

	selected := session.NextQuestion(now)

	if selected == nil {
		t.Fatal("expected a question, got nil")
	}
}

func TestSubmitAnswerUpdatesMastery(t *testing.T) {
	now := time.Now()

	questions := []Question{
		{ID: 1, TopicID: "conditionals", Difficulty: 2},
		{ID: 2, TopicID: "conditionals", Difficulty: 2},
	}

	progress := []TopicProgress{
		{
			TopicID:  "conditionals",
			Mastery:  45,
			LastSeen: now,
		},
	}

	session := NewSession(questions, progress, nil)

	answer := Answer{
		QuestionID: 1,
		TopicID:    "conditionals",
		WasCorrect: true,
		Difficulty: 2,
	}

	_, update := session.SubmitAnswer(answer, now)

	if update.Mastery <= 45 {
		t.Fatalf("expected mastery to increase, got %v", update.Mastery)
	}

	if update.CorrectStreak != 1 {
		t.Fatalf("expected correct streak = 1, got %d", update.CorrectStreak)
	}
}

func TestSubmitAnswerTracksWrongTopics(t *testing.T) {
	now := time.Now()

	questions := []Question{
		{ID: 1, TopicID: "articles", Difficulty: 1},
	}

	progress := []TopicProgress{
		{
			TopicID:  "articles",
			Mastery:  60,
			LastSeen: now,
		},
	}

	session := NewSession(questions, progress, nil)

	answer := Answer{
		QuestionID: 1,
		TopicID:    "articles",
		WasCorrect: false,
		Difficulty: 1,
	}

	session.SubmitAnswer(answer, now)

	if !session.RecentWrongTopics["articles"] {
		t.Fatal("expected topic to be marked as recently wrong")
	}
}
