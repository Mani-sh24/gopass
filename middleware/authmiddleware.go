package middleware

import (
	"example/web-service-gin/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Read Authorization header

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "Missing Authorization header",
			})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "Invalid Authorization header",
			})
			return
		}

		tokenString := parts[1]
		claims, err := helpers.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "Invalid token",
			})
			return
		}
		c.Set("id", claims["id"])
		c.Set("email", claims["email"])

		c.Next()
	}
}
