package quiz

import (
	"time"

	"github.com/google/uuid"
)

// Session represents one quiz run (in-memory for now).
type Session struct {
	ID         string
	StartedAt time.Time

	Progress map[string]TopicProgress // topic_id -> progress
	Reviews  []ReviewItem
	Questions []Question

	RecentWrongTopics map[string]bool
	AskedQuestions    map[int64]bool
}

// Answer represents a user answer submission.
type Answer struct {
	QuestionID int64
	TopicID    string
	WasCorrect bool
	Difficulty int
}

func NewSession(
	questions []Question,
	initialProgress []TopicProgress,
	reviews []ReviewItem,
) *Session {

	progressMap := make(map[string]TopicProgress)
	for _, p := range initialProgress {
		progressMap[p.TopicID] = p
	}

	return &Session{
		ID:                uuid.NewString(),
		StartedAt:         time.Now(),
		Progress:          progressMap,
		Reviews:           reviews,
		Questions:         questions,
		RecentWrongTopics: make(map[string]bool),
		AskedQuestions:    make(map[int64]bool),
	}
}

func (s *Session) NextQuestion(now time.Time) *SelectedQuestion {
	progressList := make([]TopicProgress, 0, len(s.Progress))
	for _, p := range s.Progress {
		progressList = append(progressList, p)
	}

	// Filter out already-asked questions
	available := make([]Question, 0)
	for _, q := range s.Questions {
		if !s.AskedQuestions[q.ID] {
			available = append(available, q)
		}
	}

	return SelectNextQuestion(
		now,
		progressList,
		s.Reviews,
		available,
		s.RecentWrongTopics,
	)
}

func (s *Session) SubmitAnswer(
	answer Answer,
	now time.Time,
) (*SelectedQuestion, MasteryUpdateResult) {

	// Mark question as asked
	s.AskedQuestions[answer.QuestionID] = true

	current := s.Progress[answer.TopicID]

	update := UpdateMastery(
		current,
		MasteryUpdateInput{
			WasCorrect:  answer.WasCorrect,
			Difficulty:  answer.Difficulty,
			AnsweredAt:  now,
			CurrentTime: now,
		},
	)

	s.Progress[answer.TopicID] = TopicProgress{
		TopicID:       answer.TopicID,
		Mastery:       update.Mastery,
		CorrectStreak: update.CorrectStreak,
		WrongStreak:   update.WrongStreak,
		IsMastered:    update.IsMastered,
		LastSeen:      update.LastSeen,
	}

	if !answer.WasCorrect {
		s.RecentWrongTopics[answer.TopicID] = true
	}

	next := s.NextQuestion(now)
	return next, update
}
