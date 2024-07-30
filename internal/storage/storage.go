package storage

import (
	"os"

	"github.com/Corray333/quiz_bot/internal/types"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	db, err := sqlx.Open("postgres", os.Getenv("DB_CONN_STR"))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &Storage{
		db: db,
	}
}

func (s *Storage) UpdateUser(user *types.User) error {
	return nil
}
func (s *Storage) CreateUser(user *types.User) error {
	return nil
}
func (s *Storage) GetAllUsers() ([]types.User, error) {
	return nil, nil
}
func (s *Storage) GetUserByID(user_id int64) (*types.User, error) {
	return nil, nil
}
