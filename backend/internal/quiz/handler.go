package quiz

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// In-memory session store (TEMPORARY)
var sessions = make(map[string]*Session)

// ---- Request / Response DTOs ----

type StartQuizResponse struct {
	SessionID string             `json:"session_id"`
	Question  *SelectedQuestion  `json:"question"`
}

// ---- Handlers ----

// StartQuiz creates a new quiz session and returns the first question
func StartQuiz(c *gin.Context) {
	now := time.Now()

	// --- TEMP DATA (replace with DB later) ---
	questions := []Question{
		{ID: 1, TopicID: "articles", Difficulty: 1},
		{ID: 2, TopicID: "articles", Difficulty: 2},
		{ID: 3, TopicID: "articles", Difficulty: 3},
	}

	progress := []TopicProgress{
		{
			TopicID:  "articles",
			Mastery:  40,
			LastSeen: now,
		},
	}

	// ----------------------------------------

	session := NewSession(questions, progress, nil)
	sessions[session.ID] = session

	first := session.NextQuestion(now)

	c.JSON(http.StatusOK, StartQuizResponse{
		SessionID: session.ID,
		Question:  first,
	})
}

// ---- Answer DTOs ----

type AnswerQuizRequest struct {
	SessionID  string `json:"session_id"`
	QuestionID int64  `json:"question_id"`
	TopicID    string `json:"topic_id"`
	WasCorrect bool   `json:"was_correct"`
	Difficulty int    `json:"difficulty"`
}

type AnswerQuizResponse struct {
	Status       string                `json:"status"` // continue | finished
	NextQuestion *SelectedQuestion     `json:"next_question,omitempty"`
	Mastery      MasteryUpdateResult   `json:"mastery"`
}
// AnswerQuiz processes an answer and returns the next question
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

	// ---- END OF QUIZ ----
	if next == nil {
		c.JSON(http.StatusOK, AnswerQuizResponse{
			Status:  "finished",
			Mastery: update,
		})
		return
	}

	// ---- CONTINUE ----
	c.JSON(http.StatusOK, AnswerQuizResponse{
		Status:       "continue",
		NextQuestion: next,
		Mastery:      update,
	})
}
