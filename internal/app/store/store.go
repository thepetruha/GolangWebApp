package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store struct {
	Config         *Config
	db             *sql.DB
	userRepository *UserRepository
}

//функция возращающая новый store
func NewStore(config *Config) *Store {
	return &Store{
		Config: config,
	}
}

//открытие соединения с базой данных
func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.Config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db
	fmt.Println("Database connection")

	return nil
}

//закрытие соединения с базой данных
func (s *Store) Close() error {
	s.db.Close()
	return nil
}

func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
