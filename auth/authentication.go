package auth

import (
	"net/http"
	"time"

	"slices"

	"github.com/Burning-Panda/acronyms-vault/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func IsAuthenticated(unprotectedRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allowed websites urls
		requestedUrl := c.Request.URL.Path
		if slices.Contains(unprotectedRoutes, requestedUrl) {
				c.Next()
				return
			}
		if !checkAuthStatus(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			// redirect to login page
			c.Redirect(http.StatusFound, "/login?redirectUrl=" + requestedUrl)
			c.Abort()
			return
		}
		c.Next()
	}
}

func checkAuthStatus(c *gin.Context) bool {
	// Check Authorization header first
	token := c.GetHeader("Authorization")
	if token == "" {
		// If no header, check cookie
		token, _ = c.Cookie("token")
		if token == "" {
			return false
		}
	}

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
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
			// Get redirect url from URI
			redirectUrl := c.Request.URL.Query().Get("redirectUrl")
			if redirectUrl == "" {
				redirectUrl = "/"
			}
			c.Redirect(http.StatusFound, redirectUrl)
		}
	}
}

func RegisterHandler(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.Request.Header.Get("Content-Type")
		
		var username, password, email string
		var user db.User
		
		if contentType == "application/json" {
			// Handle JSON request
			var registerRequest struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
				Email    string `json:"email" binding:"required"`
			}

			if err := c.ShouldBindJSON(&registerRequest); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			username = registerRequest.Username
			password = registerRequest.Password
			email = registerRequest.Email
		} else {
			// Handle form submission
			username = c.PostForm("username")
			password = c.PostForm("password")
			email = c.PostForm("email")
		}

		// Check if the user already exists
		if err := database.Where("username = ?", username).First(&user).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Create the user
		user = db.User{
			Username: username,
			Password: string(hashedPassword),
			Email:    email,
		}

		// Save the user to the database
		if err := database.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Redirect to the login page
		c.Redirect(http.StatusFound, "/login")
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