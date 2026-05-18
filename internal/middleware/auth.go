package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireRole validates a Supabase JWT and checks that the user_metadata.role matches one of the allowed roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token", "code": "UNAUTHORIZED"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{"HS256"}))

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "code": "UNAUTHORIZED"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims", "code": "UNAUTHORIZED"})
			return
		}

		role := extractRole(claims)
		for _, allowed := range roles {
			if role == allowed {
				c.Set("role", role)
				c.Set("user_id", claims["sub"])
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions", "code": "FORBIDDEN"})
	}
}

func extractRole(claims jwt.MapClaims) string {
	if meta, ok := claims["user_metadata"].(map[string]interface{}); ok {
		if r, ok := meta["role"].(string); ok {
			return r
		}
	}
	if r, ok := claims["role"].(string); ok {
		return r
	}
	return ""
}
