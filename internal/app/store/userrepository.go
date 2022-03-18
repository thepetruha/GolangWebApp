package store

import (
	"errors"
	"fmt"
	"webapp/internal/app/model"
)

type UserRepository struct {
	store *Store
}

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
