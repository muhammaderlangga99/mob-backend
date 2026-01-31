package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	dto "mob-backend/internal/dto"
)

// AuthMiddleware validates JWT access tokens and injects userID into the context.
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized. Invalid or expired token", Code: "UNAUTHORIZED"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized. Invalid or expired token", Code: "UNAUTHORIZED"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !parsed.Valid {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized. Invalid or expired token", Code: "UNAUTHORIZED"})
			c.Abort()
			return
		}

		claims, ok := parsed.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized. Invalid or expired token", Code: "UNAUTHORIZED"})
			c.Abort()
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized. Invalid or expired token", Code: "UNAUTHORIZED"})
			c.Abort()
			return
		}

		c.Set("userID", sub)
		c.Next()
	}
}
