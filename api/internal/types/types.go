package types

import "encoding/json"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

const (
	QuestionTypeText        = "text"
	QuestionTypeSelect      = "select"
	QuestionTypeMultiSelect = "multi_select"
	QuestionTypeFileUpload  = "file_upload"
)

const (
	QuizTypePoll = "poll"
	QuizTypeQuiz = "quiz"
)

type User struct {
	ID              int64           `json:"id" db:"user_id"`
	TgID            int64           `json:"tg_id" db:"tg_id"`
	Username        string          `json:"username" db:"username"`
	Email           string          `json:"email" db:"email"`
	Phone           string          `json:"phone" db:"phone"`
	Password        string          `json:"password" db:"password"`
	Role            string          `json:"role" db:"role"`
	Data            json.RawMessage `json:"data" db:"data"`
	CurrentQuestion int64           `json:"currentQuestion" db:"current_question"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

type Quiz struct {
	ID          int64             `json:"id" db:"quiz_id"`
	Title       string            `json:"title" db:"title"`
	Description string            `json:"description,omitempty" db:"description"`
	CreatedAt   int64             `json:"createdAt" db:"created_at"`
	Cover       string            `json:"cover" db:"cover"`
	Type        string            `json:"type" db:"type"`
	NewAnswers  int64             `json:"newAnswers" db:"new_answers"`
	Questions   []json.RawMessage `json:"questions"`
}

type Question struct {
	Type   string          `json:"type" db:"type"`
	QuizID int64           `json:"quiz_id" db:"quiz_id"`
	ID     int64           `json:"id" db:"question_id"`
	Next   *int64          `json:"next"  db:"next_question_id"`
	Data   json.RawMessage `json:"question" db:"data"`
}

type IQuestion interface {
	GetType() string
	GetID() int64
	GetNext() *int64
	GetQuizID() int64
}

type QuestionBase struct {
	Type   string `json:"type"`
	QuizID int64  `json:"quiz_id"`
	ID     int64  `json:"id"`
	Next   *int64 `json:"next"`
	Image  string `json:"image"`
}

func (q *QuestionBase) GetType() string {
	return q.Type
}
func (q *QuestionBase) GetID() int64 {
	return q.ID
}
func (q *QuestionBase) GetNext() *int64 {
	return q.Next
}

func (q *QuestionBase) GetQuizID() int64 {
	return q.QuizID
}

type QuestionText struct {
	QuestionBase
	Question string   `json:"question"`
	Answer   []string `json:"answer"`
}

type QuestionSelect struct {
	QuestionBase
	Question string   `json:"question"`
	Answer   []string `json:"answer"`
	Options  []string `json:"options"`
}

type QuestionMultiSelect struct {
	QuestionBase
	Question string   `json:"question"`
	Answer   []string `json:"answer"`
	Options  []string `json:"options"`
}

type Answer struct {
	QuestionID int64           `json:"question_id" db:"question_id"`
	UserID     int64           `json:"user_id" db:"user_id"`
	Answer     []string        `json:"answer" db:"-"`
	AnswerRaw  json.RawMessage `json:"-" db:"answer"`
	Correct    []string        `json:"correct,omitempty"`
	Checked    bool            `json:"checked" db:"checked"`
	IsCorrect  bool            `json:"isCorrect" db:"is_correct"`
}
