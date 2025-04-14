// Auth middleware using OpenFGA

package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	openfgaClient "github.com/openfga/go-sdk/client"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte("your-secret-key") // In production, use environment variable

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UUID         string    `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	Username     string    `gorm:"unique" json:"username"`
	Email        string    `gorm:"unique" json:"email"`
	PasswordHash string    `gorm:"-" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AuthMiddleware is a middleware for Gin that checks if the user is authenticated
func AuthMiddleware(c *gin.Context) {
	// Skip authentication for login and public routes
	if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/" {
		c.Next()
		return
	}

	// Get token from Authorization header
	token := GetUserFromToken(c)
	if token == "" {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Validate JWT token
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !parsedToken.Valid {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// Set user in context
	c.Set("user", claims["user_id"])
	c.Next()
}

// LoginHandler handles user login
func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var user User
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// Compare password hash
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  user.UUID,
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"error": "Error generating token",
			})
			return
		}

		// Set token in cookie
		c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)
		c.Redirect(http.StatusFound, "/")
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

	return parts[1]
}
