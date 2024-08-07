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
	ID       int64           `json:"id" db:"user_id"`
	TgID     int64           `json:"tg_id" db:"tg_id"`
	Username string          `json:"username" db:"username"`
	Email    string          `json:"email" db:"email"`
	Phone    string          `json:"phone" db:"phone"`
	Password string          `json:"password" db:"password"`
	Role     string          `json:"role" db:"role"`
	Data     json.RawMessage `json:"data" db:"data"`
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
	NewAnswers  int64             `json:"newAnswers" db:"-"`
	Questions   []json.RawMessage `json:"questions"`
}

type Question struct {
	Type     string          `json:"type"`
	QuizID   int64           `json:"quiz_id"`
	ID       int64           `json:"id"`
	Next     int64           `json:"next"`
	Question json.RawMessage `json:"question"`
}

type QuestionBase struct {
	Type   string `json:"type"`
	QuizID int64  `json:"quiz_id"`
	ID     int64  `json:"id"`
	Next   int    `json:"next"`
}

func (q *QuestionBase) GetType() string {
	return q.Type
}

type QuestionText struct {
	QuestionBase
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type QuestionSelect struct {
	QuestionBase
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Options  []string `json:"options"`
}

type QuestionMultiSelect struct {
	QuestionBase
	Question string   `json:"question"`
	Answer   []string `json:"answer"`
	Options  []string `json:"options"`
}
