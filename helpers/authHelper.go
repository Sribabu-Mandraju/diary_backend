package helper

import (
	"github.com/gin-gonic/gin" // Importing Gin web framework package
)

// MatchUserTypeToUid function matches the user type to the user ID in the context.
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	// Retrieving the user type from the Gin context.
	userType := c.GetString("user_type")

	// Ignoring the user type for now (to be used for future implementation).
	_ = userType

	return nil // Returning nil error indicating successful execution.
}
