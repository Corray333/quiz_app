package repository

import "github.com/Corray333/quiz/internal/types"

func (s *repository) CreateQuiz(quiz *types.Quiz) (int64, error) {
	res := s.db.QueryRow("INSERT INTO quizzes (title, description, cover, type) VALUES ($1, $2, $3, $4) RETURNING quiz_id", quiz.Title, quiz.Description, quiz.Cover, quiz.Type)
	if err := res.Scan(&quiz.ID); err != nil {
		return 0, err
	}
	return quiz.ID, nil
}

func (s *repository) CreateQuestion(question *types.Question) (int64, error) {
	res := s.db.QueryRow("INSERT INTO questions (quiz_id, data, type) VALUES ($1, $2, $3) RETURNING question_id", question.QuizID, question.Question, question.Type)
	if err := res.Scan(&question.ID); err != nil {
		return 0, err
	}
	return question.ID, nil

}

func (s *repository) GetQuestion(id int64) (*types.Question, error) {
	return nil, nil
}
