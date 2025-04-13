package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx" // Используем sqlx вместо pgx
	"pvz/internal/logger"
	"pvz/internal/repository/model"
)

type UserPostgres struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewUserPostgres(db *sqlx.DB, log logger.Logger) *UserPostgres {
	return &UserPostgres{
		db:     db,
		logger: log,
	}
}

func (r *UserPostgres) CreateUser(ctx context.Context, user model.User) (uuid.UUID, error) {
	var id uuid.UUID

	query := `
		INSERT INTO users (email, role, password) 
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err := r.db.QueryRowxContext(ctx, query, user.Email, user.Role, user.Password).Scan(&id)
	if err != nil {
		r.logger.Errorw("Failed to insert user into database", "email", user.Email, "error", err)
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Infow("User created in database", "userID", id, "email", user.Email)
	return id, nil
}

func (r *UserPostgres) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	query := `SELECT id, email, role, password FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		r.logger.Warnw("User not found", "email", email, "error", err)
		return user, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}
