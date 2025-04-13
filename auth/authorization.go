// Auth middleware using OpenFGA

package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	openfgaClient "github.com/openfga/go-sdk/client"
)

// AuthMiddleware is a middleware for Gin that checks if the user has the required permissions

// AuthorizationMiddleware creates a middleware that checks permissions using OpenFGA
func AuthorizationMiddleware(fgaClient *openfgaClient.OpenFgaClient, objectType string, relation string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (assuming it's set by authentication middleware)
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		// Get object ID from path parameters
		objectID := c.Param("id")
		if objectID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "object ID is required"})
			c.Abort()
			return
		}

		// Construct the object string
		object := objectType + ":" + objectID

		// Check permission using OpenFGA
		body := openfgaClient.ClientCheckRequest{
			User:     user.(string),
			Relation: relation,
			Object:   object,
		}

		response, err := fgaClient.Check(c.Request.Context()).Body(body).Execute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check authorization"})
			c.Abort()
			return
		}

		if response == nil || response.Allowed == nil || !*response.Allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserFromToken extracts user information from the Authorization header
func GetUserFromToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Assuming Bearer token format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	// In a real implementation, you would validate the token and extract the user
	// For now, we'll just return the token as the user ID
	return parts[1]
}
