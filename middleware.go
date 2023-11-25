package middleware

import (
	"net/http"
	"fmt"

	goFirebase "github.com/MyFitnessPro/firebase"
	"github.com/gin-gonic/gin"
)

// ProcessRequestMiddleware extracts and validates the query parameters and binds request body, 
// then sets these values in the gin.Context for use in handlers.
func ProcessRequestMiddleware(client *goFirebase.FirebaseClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, role, err := extractAndValidateQueryParams(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		var userData map[string]interface{}
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			if err := c.BindJSON(&userData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				c.Abort()
				return
			}
		}

		// Set extracted values in the context
		c.Set("uid", uid)
		c.Set("role", role)
		c.Set("userData", userData)
		c.Next()
	}
}

// extractAndValidateQueryParams is a helper function to extract and validate query parameters.
func extractAndValidateQueryParams(c *gin.Context) (string, string, error) {
	uid := c.Query("uid")
	role := c.Query("role")

	if uid == "" || role == "" {
		return "", "", fmt.Errorf("missing uid or role")
	}
	return uid, role, nil
}