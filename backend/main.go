package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/bugii1995/backend/internal/quiz"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// ✅ Allow frontend requests
	r.Use(cors.Default())

	r.POST("/quiz/start", quiz.StartQuiz)
	r.POST("/quiz/answer", quiz.AnswerQuiz)

	r.Run(":8080")
}
