package model

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Role     string    `db:"role"`
	Password string    `db:"password"`
}
