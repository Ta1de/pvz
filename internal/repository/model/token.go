package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID
	Role   string
}
