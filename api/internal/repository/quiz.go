package repository

import (
	"encoding/json"
	"fmt"

	"github.com/Corray333/quiz/internal/types"
)

func (s *repository) CreateQuiz(quiz *types.Quiz) (int64, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res := tx.QueryRow("INSERT INTO quizzes (title, description, cover, type) VALUES ($1, $2, $3, $4) RETURNING quiz_id", quiz.Title, quiz.Description, quiz.Cover, quiz.Type)
	if err := res.Scan(&quiz.ID); err != nil {
		return 0, err
	}

	qid := 0
	for i := len(quiz.Questions) - 1; i >= 0; i-- {
		fmt.Println(quiz.Questions[i])
		qtype := struct {
			Type string `json:"type"`
		}{}
		err := json.Unmarshal([]byte(quiz.Questions[i]), &qtype)
		if err != nil {
			return 0, err
		}

		res := tx.QueryRow("INSERT INTO questions (quiz_id, data, type, next_question_id) VALUES ($1, $2, $3, $4) RETURNING question_id", quiz.ID, quiz.Questions[i], qtype.Type, func() any {
			if qid == 0 {
				return nil
			}
			return qid
		}())
		if err := res.Scan(&qid); err != nil {
			return 0, err
		}
	}

	return quiz.ID, tx.Commit()
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

func (s *repository) ListQuizzes(offset int) ([]types.Quiz, error) {
	quizzes := []types.Quiz{}
	err := s.db.Select(&quizzes, `SELECT q.quiz_id, q.title, q.description, q.created_at, q.cover, q.type, COUNT(a.question_id) AS new_answers
	FROM quizzes q
	LEFT JOIN questions qst ON q.quiz_id = qst.quiz_id
	LEFT JOIN answers a ON qst.question_id = a.question_id AND a.checked = false
	WHERE qst.next_question_id IS NULL
	GROUP BY q.quiz_id`)
	if err != nil {
		return nil, err
	}
	return quizzes, nil
}

func (s *repository) GetQuiz(id int64) (*types.Quiz, error) {
	quiz := &types.Quiz{}
	err := s.db.Get(quiz, `
	SELECT q.quiz_id, q.title, q.description, q.created_at, q.cover, q.type, COUNT(a.question_id) AS new_answers
	FROM quizzes q
	LEFT JOIN questions qst ON q.quiz_id = qst.quiz_id
	LEFT JOIN answers a ON qst.question_id = a.question_id AND a.checked = false
	WHERE qst.next_question_id IS NULL AND q.quiz_id = $1
	GROUP BY q.quiz_id
	`, id)
	if err != nil {
		return nil, err
	}
	return quiz, nil
}
