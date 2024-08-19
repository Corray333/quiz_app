package service

import (
	"github.com/Corray333/quiz/internal/types"
)

func (s *service) GetAnswers(userID int64, quizID int64) ([]types.Answer, error) {
	return s.repo.GetAnswers(userID, quizID)
}

func (s *service) GetUserAnswers(userID int64) ([]types.Answer, error) {
	return s.repo.GetUserAnswers(userID)
}

func (s *service) SetAnswer(answer *types.Answer) (*types.Answer, error) {
	return s.repo.SetAnswer(answer)
}

func (s *service) GetQuizAnswers(quiz_id int64, offset int) ([][]types.Answer, error) {
	answers, err := s.repo.GetQuizAnswers(quiz_id, offset)
	if err != nil {
		return nil, err
	}
	if len(answers) == 0 {
		return nil, nil
	}
	result := [][]types.Answer{
		{},
	}
	currentUser := answers[0].UserID
	for _, answer := range answers {
		if answer.UserID != currentUser {
			result = append(result, []types.Answer{})
			currentUser = answer.UserID
		}
		result[len(result)-1] = append(result[len(result)-1], answer)
	}
	return result, nil
}

func (s *service) GetAnswer(uid, qid int64) (*types.Answer, error) {
	return s.repo.GetAnswer(uid, qid)
}
