package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"pvz/internal/logger"
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
		logger.SugaredLogger.Errorw("Failed to insert user into database", "email", user.Email, "error", err)
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.SugaredLogger.Infow("User created in database", "userID", id, "email", user.Email)
	return id, nil
}

func (r *UserPostgres) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	query := `SELECT id, email, role, password FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Role,
		&user.Password,
	)
	
	if err != nil {
		logger.SugaredLogger.Warnw("User not found", "email", email, "error", err)
		return user, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}
