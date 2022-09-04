package models

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        int
	UUID      uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserDTO - A representation shared with consumers of the API
type UserDTO struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (user *User) AsDTO() UserDTO {
	return UserDTO{
		UUID:      user.UUID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) GetByEmail(email string) (User, error) {
	stmt, err := m.DB.Prepare("SELECT * FROM users WHERE email = $1")

	if err != nil {
		log.Fatal(err)
	}

	var user User
	err = stmt.QueryRow(email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m UserModel) GetByUUID(uuid string) (User, error) {
	stmt, err := m.DB.Prepare("SELECT * FROM users WHERE uuid = $1")

	if err != nil {
		log.Fatal(err)
	}

	var user User
	err = stmt.QueryRow(uuid).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
