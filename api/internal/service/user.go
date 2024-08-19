package service

import "github.com/Corray333/quiz/internal/types"

func (s *service) GetUserByTG(id int64) (*types.User, error) {
	return s.repo.GetUserByTG(id)
}

func (s *service) CreateUser(user *types.User) (int64, error) {
	return s.repo.CreateUser(user)
}

func (s *service) UpdateUser(user *types.User) error {
	return s.repo.UpdateUser(user)
}

func (s *service) CreateAdmin(username string) error {
	return s.repo.CreateAdmin(username)
}

func (s *service) IsAdminById(id int64) (bool, error) {
	return s.repo.IsAdminById(id)
}

func (s *service) GetAdmins() ([]types.Admin, error) {
	return s.repo.GetAdmins()
}

func (s *service) DeleteAdmin(id int64) error {
	return s.repo.DeleteAdmin(id)
}
