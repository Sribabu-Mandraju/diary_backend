package middleware

import (
	helper "backend/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := c.Request.Header.Get("Authorization")
		if userToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "network error",
			})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(userToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error  occoureded ",
			})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("user_id", claims.User_id)
		c.Next()
	}
}
