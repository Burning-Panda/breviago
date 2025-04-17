// Auth middleware using OpenFGA

package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	openfgaClient "github.com/openfga/go-sdk/client"
)

// TODO: Use environment variable
var jwtSecret = []byte("secret-key") // In production, use environment variable

// AuthenticationMiddleware is a middleware for Gin that checks if the user is authenticated
func AuthenticationMiddleware(unprotectedRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, route := range unprotectedRoutes {
			if strings.HasPrefix(c.Request.URL.Path, route) {
				c.Next()
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !parsedToken.Valid {
			fmt.Printf("Auth failure from %s: %v\n", c.ClientIP(), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
				c.Abort()
				return
			}
		}

		c.Set("user", claims["user_id"])
		c.Next()
	}
}


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