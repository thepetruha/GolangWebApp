package model

import (
	"errors"
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

//структура пользователя
type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

//валидация email адреса
func (u *User) ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		fmt.Println("Error validating email:", email)
		return false
	}

	return true
}

//очиска пароля из структуры
func (u *User) Snitized() {
	u.Password = ""
}

//проверка длинны, хеширование, присваивание пароля структуре
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 && u.ValidateEmail(u.Email) {
		enc, err := encrypntString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	} else {
		return errors.New("Error: incorrect email or password")
	}

	return nil
}

// проверка пароля на совместимость
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

//хеширование пароля
func encrypntString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
