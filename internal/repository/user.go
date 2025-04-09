package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"pvz/internal/repository/model"
)

type UserPostgres struct {
	db *pgx.Conn
}

func NewUserPostgres(db *pgx.Conn) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(ctx context.Context, user model.User) (uuid.UUID, error) {
	var id uuid.UUID

	query := `
		INSERT INTO users (email, role, password) 
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err := r.db.QueryRow(ctx, query, user.Email, user.Role, user.Password).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка при добавлении пользователя: %w", err)
	}

	return id, nil
}

func (r *UserPostgres) GetUser(ctx context.Context, email, password string) (model.User, error) {
	var user model.User
	query := `SELECT id, email, role, password FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.Role, &user.Password)
	if err != nil {
		return user, fmt.Errorf("пользователь не найден: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, fmt.Errorf("неверный пароль: %w", err)
	}

	return user, nil
}
