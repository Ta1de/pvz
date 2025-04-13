package jwt

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"pvz/internal/logger"
	"pvz/internal/repository/model"
)

func AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Log.Warnw("Authorization header missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			logger.Log.Warnw("Token format is invalid", "token", tokenStr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SIGNING_KEY")), nil
		})
		if err != nil || !token.Valid {
			logger.Log.Warnw("Invalid or expired token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(*model.TokenClaims)
		if !ok {
			logger.Log.Warnw("Invalid token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		for _, role := range roles {
			if claims.Role == role {
				logger.Log.Infow("Token verified", "userId", claims.UserId, "role", claims.Role)
				c.Set("userClaims", claims)
				c.Next()
				return
			}
		}

		logger.Log.Warnw("Access forbidden", "allowedRoles", roles, "claims", claims)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
	}
}
