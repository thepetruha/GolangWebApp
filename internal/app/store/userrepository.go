package store

import (
	"errors"
	"fmt"
	"webapp/internal/app/model"
)

//структура хранилища
type UserRepository struct {
	store *Store
}

//занесение пользователя в БД
func (r *UserRepository) Create(u *model.User) error {
	if err := u.BeforeCreate(); err != nil {
		return err
	}

	var row = r.store.db.QueryRow(`INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id`, u.Email, u.EncryptedPassword)
	if row == nil {
		return errors.New("Error querying")
	}

	if err := row.Scan(&u.ID); err != nil {
		fmt.Println("Error querying insert user")
		return err
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		`SELECT id, email, encrypted_password FROM users WHERE email = $1`,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword); err != nil {
		return nil, err
	}

	return u, nil
}

//Поиск пользователя по id
func (r *UserRepository) FindUserId(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		`SELECT id, email, encrypted_password FROM users WHERE id = $1`,
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword); err != nil {
		return nil, err
	}

	return u, nil
}
