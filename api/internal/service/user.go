package service

import "github.com/Corray333/quiz/internal/types"

func (s *service) GetUserByTG(id int64) (*types.User, error) {
	return s.repo.GetUserByTG(id)
}

func (s *service) UpdateUser(user *types.User) error {
	return s.repo.UpdateUser(user)
}
