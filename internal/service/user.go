package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"pvz/internal/logger"
	"pvz/internal/repository"
	"pvz/internal/repository/model"
)

const tokenTTL = time.Hour * 24

type UserService struct {
	repoUser repository.User
	logger   logger.Logger
}

func NewUserService(repoUser repository.User, log logger.Logger) *UserService {
	return &UserService{
		repoUser: repoUser,
		logger:   log,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	hashedPassword, err := GeneratePasswordHash(user.Password)
	if err != nil {
		s.logger.Errorw("Password hashing failed", "error", err)
		return model.User{}, fmt.Errorf("could not hash password: %w", err)
	}

	user.Password = hashedPassword

	id, err := s.repoUser.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	user.Id = id
	s.logger.Infow("User successfully created", "userID", user.Id, "email", user.Email)
	return user, nil
}

func GeneratePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repoUser.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Warnw("Incorrect password attempt", "email", email)
		return "", fmt.Errorf("invalid credentials")
	}

	claims := &model.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		s.logger.Errorw("Failed to sign JWT", "userID", user.Id, "error", err)
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	s.logger.Infow("User authentication successful", "userID", user.Id, "email", user.Email)
	return signedToken, nil
}

func (s *UserService) DummyLogin(ctx context.Context, role string) (string, error) {
	claims := &model.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
	if err != nil {
		s.logger.Errorw("Failed to sign dummy token", "role", role, "error", err)
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	s.logger.Infow("Dummy token created", "role", role)
	return signedToken, nil
}
