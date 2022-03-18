package model

import "golang.org/x/crypto/bcrypt"

//структура пользователя
type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

//очиска пароля из структуры
func (u *User) Snitized() {
	u.Password = ""
}

//проверка длинны, хеширование, присваивание пароля структуре
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encrypntString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
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
