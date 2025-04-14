package auth

import (
	"net/http"
	"time"

	"github.com/Burning-Panda/acronyms-vault/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func IsAuthenticated(c *gin.Context) bool {
	token := c.GetHeader("Authorization")
	if token == "" {
		return false
	}

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !parsedToken.Valid {
		return false
	}

	return true
}

// LoginHandler handles user login
func LoginHandler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request is a form submission or JSON
		contentType := c.Request.Header.Get("Content-Type")
		
		var username, password string
		var user db.User
		
		if contentType == "application/json" {
			// Handle JSON request
			var loginRequest struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}
			
			// Validate the request body
			if err := c.ShouldBindJSON(&loginRequest); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}
			
			username = loginRequest.Username
			password = loginRequest.Password
		} else {
			// Handle form submission
			username = c.PostForm("username")
			password = c.PostForm("password")
		}
		
		// Check if the user exists
		if err := database.Where("username = ?", username).First(&user).Error; err != nil {
			if contentType == "application/json" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			} else {
				c.HTML(http.StatusUnauthorized, "login.html", gin.H{
					"error": "Invalid credentials",
				})
			}
			return
		}
		
		// Compare password hash
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			if contentType == "application/json" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			} else {
				c.HTML(http.StatusUnauthorized, "login.html", gin.H{
					"error": "Invalid credentials",
				})
			}
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
			if contentType == "application/json" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
			} else {
				c.HTML(http.StatusInternalServerError, "login.html", gin.H{
					"error": "Error generating token",
				})
			}
			return
		}
		
		// Return response based on content type
		if contentType == "application/json" {
			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		} else {
			// Set token in cookie
			c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)
			c.Redirect(http.StatusFound, "/")
		}
	}
}

func RegisterHandler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
		
}

func LogoutHandler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func UserHandler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func RefreshHandler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}