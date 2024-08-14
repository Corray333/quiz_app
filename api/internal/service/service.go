package service

import (
	"errors"

	"github.com/Corray333/quiz/internal/types"
)

var (
	ErrWrongImageUrl = errors.New("wrong image url")
)

type Question interface {
	GetType() string
}

type Storage interface {
	CreateQuiz(quiz *types.Quiz) (int64, error)
	CreateQuestion(question *types.Question) (int64, error)
	GetQuestion(id int64) (*types.Question, error)
	GetUserByTG(user_id int64) (*types.User, error)
	UpdateUser(user *types.User) error
	ListQuizzes(offset int) ([]types.Quiz, error)
	GetQuiz(id int64) (*types.Quiz, error)
	SetAnswer(newAnswer *types.Answer) (*types.Answer, error)
	GetFirstQuestion(quizID int64) (types.IQuestion, error)

	CreateUser(user *types.User) (int64, error)
	SetCurrentQuestion(uid, qid int64) error
	GetCurrentQuestion(uid int64) (types.IQuestion, error)
	GetNextQuestion(uid int64) (types.IQuestion, error)
	GetAnswers(userID int64, quizID int64) ([]types.Answer, error)
	GetQuizAnswers(userID int64) ([]types.Answer, error)
	GetAllAnswers(quizID int64, offset int) ([]types.Answer, error)
}

type service struct {
	repo Storage
}

func NewService(store Storage) *service {
	return &service{
		repo: store,
	}
}

func (s *service) CreateQuiz(quiz *types.Quiz) (int64, error) {
	// if !strings.HasPrefix(quiz.Cover, viper.GetString("image_url")) {
	// 	return 0, ErrWrongImageUrl
	// }
	// quiz.Cover = strings.TrimPrefix(quiz.Cover, viper.GetString("image_url"))

	return s.repo.CreateQuiz(quiz)
}

func (s *service) CreateQuestion(question *types.Question) (int64, error) {
	return s.repo.CreateQuestion(question)
}

func (s *service) GetQuestion(id int64) (*types.Question, error) {
	return s.repo.GetQuestion(id)
}

func (s *service) ListQuizzes(offset int) ([]types.Quiz, error) {
	return s.repo.ListQuizzes(offset)
}

func (s *service) GetQuiz(id int64) (*types.Quiz, error) {
	return s.repo.GetQuiz(id)
}

func (s *service) GetFirstQuestion(quizID int64) (types.IQuestion, error) {
	return s.repo.GetFirstQuestion(quizID)
}

func (s *service) SetCurrentQuestion(uid, qid int64) error {
	return s.repo.SetCurrentQuestion(uid, qid)
}

func (s *service) GetCurrentQuestion(uid int64) (types.IQuestion, error) {
	return s.repo.GetCurrentQuestion(uid)
}

func (s *service) GetNextQuestion(uid int64) (types.IQuestion, error) {
	return s.repo.GetNextQuestion(uid)
}
