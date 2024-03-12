package middleware

import (
	helper "backend/helpers" // Importing helper package for token validation
	"net/http"               // Importing net/http package for HTTP status codes

	"github.com/gin-gonic/gin" // Importing Gin web framework package
)

// Authenticate middleware verifies the authentication token in the request header.
func Authenticate() gin.HandlerFunc {
	// Returns a Gin middleware handler function.
	return func(c *gin.Context) {
		// Extracting the authentication token from the request header.
		userToken := c.Request.Header.Get("Authorization")

		// Checking if the authentication token is empty.
		if userToken == "" {
			// Responding with an internal server error if the token is missing.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "network error",
			})
			c.Abort() // Aborting the middleware chain.
			return
		}

		// Validating the authentication token.
		claims, err := helper.ValidateToken(userToken)
		if err != "" {
			// Responding with an internal server error if token validation fails.
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error  occoureded",
			})
			c.Abort() // Aborting the middleware chain.
			return
		}

		// Setting user information extracted from the token into Gin context for further request processing.
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("user_id", claims.User_id)

		c.Next() // Proceeding to the next middleware or route handler.
	}
}
