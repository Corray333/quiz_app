package repository

import "github.com/Corray333/quiz/internal/types"

func (s *repository) UpdateUser(user *types.User) error {
	return nil
}
func (s *repository) CreateUser(user *types.User) (int64, error) {
	res := s.db.QueryRow("INSERT INTO users(username, email, password, tg_id, phone, data) VALUES($1, $2, $3, $4, $5) RETURNING user_id", user.Username, user.Email, user.Password, user.TgID, user.Phone, user.Data)
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
