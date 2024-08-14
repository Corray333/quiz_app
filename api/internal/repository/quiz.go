package repository

import (
	"encoding/json"
	"fmt"
	"slices"

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
	fmt.Println(quiz)

	qid := 0
	for i := 0; i < len(quiz.Questions); i++ {
		qtype := struct {
			Type string `json:"type"`
		}{}
		err := json.Unmarshal([]byte(quiz.Questions[i]), &qtype)
		if err != nil {
			return 0, err
		}

		res := tx.QueryRow("INSERT INTO questions (quiz_id, data, type, question_number) VALUES ($1, $2, $3, $4) RETURNING question_id", quiz.ID, quiz.Questions[i], qtype.Type, i+1)
		if err := res.Scan(&qid); err != nil {
			return 0, err
		}
	}

	return quiz.ID, tx.Commit()
}

func (s *repository) CreateQuestion(question *types.Question) (int64, error) {
	res := s.db.QueryRow("INSERT INTO questions (quiz_id, data, type) VALUES ($1, $2, $3) RETURNING question_id", question.QuizID, question.Data, question.Type)
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
	WHERE qst.question_number = 1
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
	WHERE qst.question_number = 1 AND q.quiz_id = $1
	GROUP BY q.quiz_id
	`, id)
	if err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *repository) GetFirstQuestion(quizID int64) (types.IQuestion, error) {
	req := &types.Question{}
	err := s.db.Get(req, `
	SELECT *
	FROM questions q1
	WHERE quiz_id = $1 AND q1.question_number = 1
	`, quizID)
	if err != nil {
		return nil, err
	}

	var result types.IQuestion

	switch req.Type {
	case types.QuestionTypeText:
		question := &types.QuestionText{}
		if err := json.Unmarshal(req.Data, question); err != nil {
			return nil, err
		}
		question.QuizID = req.QuizID
		question.Type = req.Type
		question.ID = req.ID
		result = question
	case types.QuestionTypeSelect:
		question := &types.QuestionSelect{}
		if err := json.Unmarshal(req.Data, question); err != nil {
			return nil, err
		}
		question.QuizID = req.QuizID
		question.Type = req.Type
		question.ID = req.ID
		result = question
	case types.QuestionTypeMultiSelect:
		question := &types.QuestionMultiSelect{}
		if err := json.Unmarshal(req.Data, question); err != nil {
			return nil, err
		}
		question.QuizID = req.QuizID
		question.Type = req.Type
		question.ID = req.ID
		result = question
	}

	return result, nil
}

func (s *repository) GetAnswers(userID int64, quizID int64) ([]types.Answer, error) {
	answers := []types.Answer{}
	if err := s.db.Select(&answers, "SELECT question_id, user_id, answer, checked FROM answers NATURAL JOIN questions WHERE user_id = $1 AND quiz_id = $2", userID, quizID); err != nil {
		return nil, fmt.Errorf("error getting answers: %v", err)
	}

	for i := range answers {
		if err := json.Unmarshal(answers[i].AnswerRaw, &answers[i].Answer); err != nil {
			return nil, fmt.Errorf("error unmarshalling answer: %v", err)
		}
		answers[i].AnswerRaw = nil
		answers[i].Correct = nil
	}
	return answers, nil
}

func (s *repository) GetAllAnswers(quizID int64, offset int) ([]types.Answer, error) {
	answers := []types.Answer{}
	if err := s.db.Select(&answers, `SELECT questions.question_id, user_id, answer, checked FROM answers
	JOIN questions ON answers.question_id = questions.question_id
	WHERE questions.quiz_id = $1
	ORDER BY checked, answers.user_id, questions.question_number LIMIT 50 OFFSET $2`, quizID, offset); err != nil {
		return nil, fmt.Errorf("error getting answers: %v", err)
	}

	questions := []types.Question{}
	if err := s.db.Select(&questions, "SELECT question_id, data FROM questions WHERE quiz_id = $1", quizID); err != nil {
		return nil, fmt.Errorf("error getting correct answers: %v", err)
	}

	users := map[int64]struct{}{}

	for i := range answers {
		if err := json.Unmarshal(answers[i].AnswerRaw, &answers[i].Answer); err != nil {
			return nil, fmt.Errorf("error unmarshalling answer: %v", err)
		}
		users[answers[i].UserID] = struct{}{}
		for _, question := range questions {
			if question.ID == answers[i].QuestionID {
				correct := struct {
					Answer []string `json:"answer"`
				}{}
				if err := json.Unmarshal(question.Data, &correct); err != nil {
					return nil, fmt.Errorf("error unmarshalling correct answer: %v", err)
				}
				answers[i].Correct = correct.Answer
				answers[i].IsCorrect = answerIsCorrect(answers[i].Answer, correct.Answer)
				answers[i].AnswerRaw = nil
			}
		}
		answers[i].AnswerRaw = nil
	}
	tx, err := s.db.Beginx()
	defer tx.Rollback()

	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	for i := range users {
		if _, err := tx.Exec(`UPDATE answers
		SET checked = true
		FROM questions
		WHERE answers.question_id = questions.question_id
		AND answers.user_id = $1
		AND questions.quiz_id = $2`, i, quizID); err != nil {
			return nil, fmt.Errorf("error updating answers: %v", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return answers, nil
}

func answerIsCorrect(answer, correct []string) bool {
	if len(answer) != len(correct) {
		return false
	}
	for i := range answer {
		if !slices.Contains(correct, answer[i]) {
			return false
		}
	}
	return true
}

func (s *repository) GetQuizAnswers(userID int64) ([]types.Answer, error) {

	quizID := 0

	if err := s.db.QueryRow("SELECT quiz_id FROM users JOIN questions ON users.current_question = questions.question_id WHERE user_id = $1", userID).Scan(&quizID); err != nil {
		return nil, fmt.Errorf("error getting quiz id: %v", err)
	}

	answers := []types.Answer{}
	if err := s.db.Select(&answers, "SELECT question_id, user_id, answer, checked FROM answers NATURAL JOIN questions WHERE user_id = $1 AND quiz_id = $2 ORDER BY question_number", userID, quizID); err != nil {
		return nil, fmt.Errorf("error getting answers: %v", err)
	}
	questions := []types.Question{}
	if err := s.db.Select(&questions, "SELECT question_id, data FROM questions WHERE quiz_id = $1", quizID); err != nil {
		return nil, fmt.Errorf("error getting correct answers: %v", err)
	}

	for i := range answers {
		if err := json.Unmarshal(answers[i].AnswerRaw, &answers[i].Answer); err != nil {
			return nil, fmt.Errorf("error unmarshalling answer: %v", err)
		}

		for _, question := range questions {
			if question.ID == answers[i].QuestionID {
				correct := struct {
					Answer []string `json:"answer"`
				}{}
				if err := json.Unmarshal(question.Data, &correct); err != nil {
					return nil, fmt.Errorf("error unmarshalling correct answer: %v", err)
				}
				answers[i].Correct = correct.Answer
				answers[i].IsCorrect = answerIsCorrect(answers[i].Answer, correct.Answer)
				answers[i].AnswerRaw = nil
			}
		}
	}
	return answers, nil
}
