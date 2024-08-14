package service

import (
	"github.com/Corray333/quiz/internal/types"
)

func (s *service) GetAnswers(userID int64, quizID int64) ([]types.Answer, error) {
	return s.repo.GetAnswers(userID, quizID)
}

func (s *service) GetQuizAnswers(userID int64) ([]types.Answer, error) {
	return s.repo.GetQuizAnswers(userID)
}

func (s *service) SetAnswer(answer *types.Answer) (*types.Answer, error) {
	return s.repo.SetAnswer(answer)
}

func (s *service) GetAllAnswers(quiz_id int64, offset int) ([][]types.Answer, error) {
	answers, err := s.repo.GetAllAnswers(quiz_id, offset)
	if err != nil {
		return nil, err
	}
	result := [][]types.Answer{
		{answers[0]},
	}
	currentUser := answers[0].UserID
	for _, answer := range answers {
		if answer.UserID != currentUser {
			result = append(result, []types.Answer{})
			result[len(result)-1] = append(result[len(result)-1], answer)
			currentUser = answer.UserID
		}
		result[len(result)-1] = append(result[len(result)-1], answer)
	}
	return result, nil
}
