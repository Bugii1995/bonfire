package quiz

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ---------------- In-memory session store ----------------

var sessions = make(map[string]*Session)

// ---------------- DTOs ----------------

// Safe question sent to client (NO correct answer)
type QuestionResponse struct {
	ID      int64    `json:"id"`
	Prompt  string   `json:"prompt"`
	Options []string `json:"options"`
	Purpose string   `json:"purpose"`
}

type StartQuizResponse struct {
	SessionID string           `json:"session_id"`
	Question  QuestionResponse `json:"question"`
}

type AnswerQuizRequest struct {
	SessionID  string `json:"session_id"`
	QuestionID int64  `json:"question_id"`
	TopicID    string `json:"topic_id"`
	WasCorrect bool   `json:"was_correct"`
	Difficulty int    `json:"difficulty"`
}

type AnswerQuizResponse struct {
	Status       string              `json:"status"` // continue | finished
	NextQuestion *QuestionResponse   `json:"next_question,omitempty"`
	Mastery      MasteryUpdateResult `json:"mastery"`
	Explanation  string              `json:"explanation"`
}

// ---------------- Helpers ----------------

func toQuestionResponse(q Question, purpose QuestionPurpose) QuestionResponse {
	return QuestionResponse{
		ID:      q.ID,
		Prompt:  q.Prompt,
		Options: q.Options,
		Purpose: string(purpose),
	}
}

func findQuestionByID(questions []Question, id int64) Question {
	for _, q := range questions {
		if q.ID == id {
			return q
		}
	}
	panic("question not found")
}

// ---------------- Handlers ----------------

// POST /quiz/start
func StartQuiz(c *gin.Context) {
	now := time.Now()

	// ---- TEMP in-memory question set (Milestone 1) ----
	questions := []Question{
		{
			ID:             1,
			TopicID:        "articles",
			Difficulty:     1,
			Prompt:         "Choose the correct article: ___ apple",
			Options:        []string{"a", "an", "the"},
			CorrectAnswer:  "an",
			Explanation:    "We use 'an' before words that start with a vowel sound.",
		},
		{
			ID:             2,
			TopicID:        "articles",
			Difficulty:     2,
			Prompt:         "Choose the correct article: ___ university",
			Options:        []string{"a", "an", "the"},
			CorrectAnswer:  "a",
			Explanation:    "'University' starts with a 'you' sound, so we use 'a'.",
		},
	}
	// ---------------------------------------------------

	progress := []TopicProgress{
		{
			TopicID:  "articles",
			Mastery:  40,
			LastSeen: now,
		},
	}

	session := NewSession(questions, progress, nil)
	sessions[session.ID] = session

	selected := session.NextQuestion(now)
	if selected == nil {
		c.JSON(http.StatusOK, gin.H{"error": "no questions available"})
		return
	}

	q := findQuestionByID(session.Questions, selected.QuestionID)

	c.JSON(http.StatusOK, StartQuizResponse{
		SessionID: session.ID,
		Question:  toQuestionResponse(q, selected.Purpose),
	})
}

// POST /quiz/answer
func AnswerQuiz(c *gin.Context) {
	var req AnswerQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, ok := sessions[req.SessionID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	now := time.Now()

	next, update := session.SubmitAnswer(
		Answer{
			QuestionID: req.QuestionID,
			TopicID:    req.TopicID,
			WasCorrect: req.WasCorrect,
			Difficulty: req.Difficulty,
		},
		now,
	)

	answered := findQuestionByID(session.Questions, req.QuestionID)

	// ---- Finished ----
	if next == nil {
		c.JSON(http.StatusOK, AnswerQuizResponse{
			Status:      "finished",
			Mastery:     update,
			Explanation: answered.Explanation,
		})
		return
	}

	// ---- Continue ----
	nextQ := findQuestionByID(session.Questions, next.QuestionID)
	resp := toQuestionResponse(nextQ, next.Purpose)

	c.JSON(http.StatusOK, AnswerQuizResponse{
		Status:       "continue",
		NextQuestion: &resp,
		Mastery:      update,
		Explanation:  answered.Explanation,
	})
}
