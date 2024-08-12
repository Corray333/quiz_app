package service

import (
	"github.com/Corray333/quiz/internal/types"
)

func (s *service) GetAnswers(userID int64, quizID int64) ([]types.Answer, error) {
	return s.repo.GetAnswers(userID, quizID)
}

func (s *service) GetQuizAnswers(userID int64, quizID int64) ([]types.Answer, error) {
	return s.repo.GetQuizAnswers(userID, quizID)
}

func (s *service) SetAnswer(answer *types.Answer) (*types.Answer, error) {
	return s.repo.SetAnswer(answer)
}
