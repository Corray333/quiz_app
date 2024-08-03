package repository

import (
	"os"

	"github.com/Corray333/quiz/internal/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type repository struct {
	db *sqlx.DB
}

func New() *repository {
	db, err := sqlx.Open("postgres", os.Getenv("DB_CONN_STR"))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &repository{
		db: db,
	}
}

func (s *repository) UpdateUser(user *types.User) error {
	return nil
}
func (s *repository) CreateUser(user *types.User) error {
	return nil
}
func (s *repository) GetAllUsers() ([]types.User, error) {
	return nil, nil
}
func (s *repository) GetUserByTG(user_id int64) (*types.User, error) {
	return nil, nil
}
