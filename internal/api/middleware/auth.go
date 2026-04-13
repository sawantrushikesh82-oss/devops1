package middleware

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(pub *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token", "code": "UNAUTHORIZED"})
			c.Abort()
			return
		}
		tok := strings.TrimPrefix(h, "Bearer ")
		claims := jwt.MapClaims{}
		parsed, err := jwt.ParseWithClaims(tok, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, jwt.ErrSignatureInvalid
			}
			return pub, nil
		})
		if err != nil || !parsed.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "code": "UNAUTHORIZED"})
			c.Abort()
			return
		}
		c.Set("user_claims", claims)
		c.Next()
	}
}
