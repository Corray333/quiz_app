package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/Corray333/quiz/internal/types"
)

func (s *repository) UpdateUser(user *types.User) error {
	return nil
}
func (s *repository) CreateUser(user *types.User) (int64, error) {
	if user.Data == nil {
		user.Data = json.RawMessage("{}")
	}
	res := s.db.QueryRow("INSERT INTO users(username, email, password, tg_id, phone, data) VALUES($1, $2, $3, $4, $5, $6) RETURNING user_id", user.Username, user.Email, user.Password, user.TgID, user.Phone, user.Data)
	err := res.Scan(&user.ID)
	return user.ID, err
}
func (s *repository) GetAllUsers() ([]types.User, error) {
	return nil, nil
}
func (s *repository) GetUserByTG(tgID int64) (*types.User, error) {
	user := &types.User{}
	err := s.db.Get(user, "SELECT * FROM users WHERE tg_id = $1", tgID)
	return user, err
}

func (s *repository) SetCurrentQuestion(uid, qid int64) error {
	fmt.Println(uid, qid)
	_, err := s.db.Exec("UPDATE users SET current_question = $1 WHERE user_id = $2", qid, uid)
	return err
}

func (s *repository) GetCurrentQuestion(uid int64) (types.IQuestion, error) {
	req := &types.Question{}
	if err := s.db.Get(req, "SELECT * FROM questions WHERE question_id = (SELECT current_question FROM users WHERE user_id = $1)", uid); err != nil {
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

func (s *repository) SetAnswer(newAnswer *types.Answer) (*types.Answer, error) {
	ans := &types.Answer{}
	answerStr, err := json.Marshal(newAnswer.Answer)
	if err != nil {
		return nil, err
	}
	fmt.Println("New answer: ", newAnswer)
	err = s.db.Get(ans, "SELECT * FROM answers WHERE user_id = $1 AND question_id = $2", newAnswer.UserID, newAnswer.QuestionID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		_, err := s.db.Exec("INSERT INTO answers(user_id, question_id, answer) VALUES($1, $2, $3)", newAnswer.UserID, newAnswer.QuestionID, string(answerStr))
		if err != nil {
			return nil, err
		}
		return newAnswer, nil
	}

	if err := json.Unmarshal(ans.AnswerRaw, &ans.Answer); err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", ans)

	question, err := s.GetCurrentQuestion(newAnswer.UserID)
	if err != nil {
		return nil, err
	}
	if question.GetType() == types.QuestionTypeMultiSelect {
		for i := range newAnswer.Answer {
			fmt.Println(ans.Answer, newAnswer.Answer[i])
			if ind := slices.Index(ans.Answer, newAnswer.Answer[i]); ind != -1 {
				ans.Answer = append(ans.Answer[:ind], ans.Answer[ind+1:]...)
				fmt.Println(ans.Answer)
			} else {
				ans.Answer = append(ans.Answer, newAnswer.Answer[i])
			}
		}
	} else {
		ans.Answer = newAnswer.Answer
	}

	answerStr, err = json.Marshal(ans.Answer)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(answerStr))
	_, err = s.db.Exec("UPDATE answers SET answer = $3, checked = false WHERE user_id = $1 AND question_id = $2", newAnswer.UserID, newAnswer.QuestionID, string(answerStr))
	if err != nil {
		return nil, err
	}

	newAnswer.Answer = ans.Answer

	return newAnswer, nil
}

func (r repository) GetAnswer(uid int64, qid int64) (*types.Answer, error) {
	ans := &types.Answer{}
	if err := r.db.Get(ans, "SELECT * FROM answers WHERE user_id = $1 AND question_id = $2", uid, qid); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(ans.AnswerRaw, &ans.Answer); err != nil {
		return nil, err
	}

	return ans, nil
}

func (r repository) GetNextQuestion(uid int64) (types.IQuestion, error) {
	req := &types.Question{}
	if err := r.db.Get(req, `WITH current AS (
			SELECT quiz_id, question_number
			FROM questions
			WHERE question_id = (SELECT current_question FROM users WHERE user_id = $1)
		)
		SELECT q.quiz_id, q.question_id, q.type, q.data, q.question_number
		FROM questions q
		JOIN current c ON q.quiz_id = c.quiz_id AND q.question_number = c.question_number + 1`, uid); err != nil {
		if err == sql.ErrNoRows {
			return &types.QuestionBase{ID: 0}, nil
		}
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

func (r repository) CreateAdmin(username string) error {
	_, err := r.db.Exec("INSERT INTO admins(username) VALUES($1)", username)
	if err != nil {
		return fmt.Errorf("failed to create admin: %v", err)
	}

	return nil

}

func (r repository) IsAdminById(id int64) (bool, error) {
	count := 0
	r.db.Get(&count, "SELECT COUNT(*) FROM admins WHERE tg_id = $1", id)
	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (r repository) GetAdmins() ([]types.Admin, error) {
	admins := []types.Admin{}
	if err := r.db.Select(&admins, "SELECT * FROM admins"); err != nil {
		return nil, err
	}

	return admins, nil

}

func (r repository) DeleteAdmin(id int64) error {
	_, err := r.db.Exec("DELETE FROM admins WHERE tg_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete admin: %v", err)
	}

	return nil

}
