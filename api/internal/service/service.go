package service

import "github.com/Corray333/quiz/internal/types"

type Question interface {
	GetType() string
}

type Storage interface {
	CreateQuiz(quiz *types.Quiz) (int64, error)
	CreateQuestion(question *types.Question) (int64, error)
	GetQuestion(id int64) (*types.Question, error)
	GetUserByTG(user_id int64) (*types.User, error)
	UpdateUser(user *types.User) error
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
	return s.repo.CreateQuiz(quiz)
}

func (s *service) CreateQuestion(question *types.Question) (int64, error) {
	return s.repo.CreateQuestion(question)
}

func (s *service) GetQuestion(id int64) (*types.Question, error) {
	return s.repo.GetQuestion(id)
}
