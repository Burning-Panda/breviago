package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/Burning-Panda/breviago/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	database    *gorm.DB
	// Better-auth frontend URL - adjust this to your frontend URL
	frontendURL = "http://localhost:3000" // Change this to your actual frontend URL
	
	// Session cache to avoid repeated validation calls
	sessionCache = make(map[string]*CachedSession)
	sessionMutex = sync.RWMutex{}
	cacheExpiry  = 5 * time.Minute // Cache sessions for 5 minutes
)

// Cached session data
type CachedSession struct {
	Data      *BetterAuthSession
	ExpiresAt time.Time
}

// Better-auth session response structure
type BetterAuthSession struct {
	Session struct {
		ID        string `json:"id"`
		UserID    string `json:"userId"`
		ExpiresAt string `json:"expiresAt"`
	} `json:"session"`
	User struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}

// validateBetterAuthSession validates a session with the better-auth frontend
func validateBetterAuthSession(sessionToken string) (*BetterAuthSession, error) {
	// Check cache first
	sessionMutex.RLock()
	cached, exists := sessionCache[sessionToken]
	sessionMutex.RUnlock()
	
	if exists && time.Now().Before(cached.ExpiresAt) {
		return cached.Data, nil
	}
	
	// Create HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second, // Add timeout
	}
	
	// Create request to better-auth session endpoint
	req, err := http.NewRequest("GET", frontendURL+"/api/auth/session", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	// Add the session cookie
	req.Header.Set("Cookie", fmt.Sprintf("better-auth.session_token=%s", sessionToken))
	
	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %v", err)
	}
	defer resp.Body.Close()
	
	// Check if session is valid
	if resp.StatusCode != http.StatusOK {
		// Remove from cache if it exists
		sessionMutex.Lock()
		delete(sessionCache, sessionToken)
		sessionMutex.Unlock()
		return nil, fmt.Errorf("invalid session: status %d", resp.StatusCode)
	}
	
	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	var sessionData BetterAuthSession
	if err := json.Unmarshal(body, &sessionData); err != nil {
		return nil, fmt.Errorf("failed to parse session data: %v", err)
	}
	
	// Cache the validated session
	sessionMutex.Lock()
	sessionCache[sessionToken] = &CachedSession{
		Data:      &sessionData,
		ExpiresAt: time.Now().Add(cacheExpiry),
	}
	sessionMutex.Unlock()
	
	return &sessionData, nil
}

// cleanupExpiredSessions removes expired sessions from cache
func cleanupExpiredSessions() {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	
	now := time.Now()
	for token, cached := range sessionCache {
		if now.After(cached.ExpiresAt) {
			delete(sessionCache, token)
		}
	}
}

// invalidateSessionCache removes a specific session from cache
func invalidateSessionCache(sessionToken string) {
	sessionMutex.Lock()
	delete(sessionCache, sessionToken)
	sessionMutex.Unlock()
}

func AuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get session token from cookie
		sessionCookie, err := c.Cookie("better-auth.session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session cookie required"})
			c.Abort()
			return
		}
		
		// Validate session with better-auth
		sessionData, err := validateBetterAuthSession(sessionCookie)
		if err != nil {
			log.Printf("Session validation error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}
		
		// Store session data in context for UserMiddleware
		c.Set("sessionData", sessionData)
		c.Next()
	}
}

// TODO: Move this to a separate package
// Add user to the gin context for authentication
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionDataInterface, exists := c.Get("sessionData")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Session data not found in context"})
			c.Abort()
			return
		}

		sessionData, ok := sessionDataInterface.(*BetterAuthSession)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid session data type in context"})
			c.Abort()
			return
		}

		// Try to find user in our database using email from better-auth
		var user db.User
		err := database.
			Where("email = ?", sessionData.User.Email).
			First(&user).
			Error
		
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// User doesn't exist in our database, create them
				user = db.User{
					Name:  sessionData.User.Name,
					Email: sessionData.User.Email,
					// You might want to set other fields as needed
				}
				
				if createErr := database.Create(&user).Error; createErr != nil {
					log.Printf("Failed to create user: %v", createErr)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
					c.Abort()
					return
				}
				log.Printf("Created new user from better-auth session: %s", user.Email)
			} else {
				log.Printf("Database error when looking up user: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				c.Abort()
				return
			}
		} else {
			// User exists, check if we need to sync data from better-auth
			needsUpdate := false
			if user.Name != sessionData.User.Name {
				user.Name = sessionData.User.Name
				needsUpdate = true
			}
			// Add other fields you want to sync
			
			if needsUpdate {
				if updateErr := database.Save(&user).Error; updateErr != nil {
					log.Printf("Failed to sync user data: %v", updateErr)
					// Don't abort on sync errors, just log them
				} else {
					log.Printf("Synced user data for: %s", user.Email)
				}
			}
		}
		
		// Set user in context
		c.Set("user", user)
		c.Next()
	}
}

func main() {
	// Initialize database
	database = db.GetGormDB(database)
	if database == nil {
		log.Fatal("Failed to initialize database")
	}
	defer func() {
		if err := db.CloseGormDB(database); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Initialize database schema and default data
	db.InitDB(database)

	// Start background cleanup routine for expired sessions
	go func() {
		ticker := time.NewTicker(10 * time.Minute) // Cleanup every 10 minutes
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cleanupExpiredSessions()
				log.Println("Cleaned up expired session cache entries")
			}
		}
	}()

	r := gin.Default()

	// CORS middleware - adjust the origin to match your frontend URL
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", frontendURL)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Routes
	r.GET("/", getDefault)
	r.GET("/health", getHealth) // Health check endpoint

	// Session management endpoints (for frontend to notify backend)
	r.POST("/api/session/invalidate", invalidateSession) // Endpoint for frontend to invalidate sessions
	r.POST("/api/session/refresh", refreshSession)       // Endpoint to force session refresh

	// API group with authentication
	// AuthenticatedMiddleware will run first, then UserMiddleware
	api := r.Group("/api", AuthenticatedMiddleware(), UserMiddleware())
	api.GET("/", testAcronyms) // This will now be protected

	api.GET("/acronyms", getAcronyms)
	api.POST("/acronyms", createAcronym)

	api.GET("/acronyms/:id", getAcronym)
	api.PUT("/acronyms/:id", updateAcronym)
	api.DELETE("/acronyms/:id", deleteAcronym)

	// User API
	user := api.Group("/users")
	user.GET("/", getUser)       // Protected
	user.POST("/", createUser)   // This createUser is for backend admin, registration is separate
	user.GET("/me", getMe)       // Protected
	user.GET("/:id", getUser)    // Protected
	user.PUT("/:id", updateUser) // Protected
	user.DELETE("/:id", deleteUser) // Protected

	r.Run(":8080")
}

func getDefault(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func testAcronyms(c *gin.Context) {
	var acronyms []db.Acronym
	if err := database.
		Preload("Owner").
		Preload("Labels").
		Find(&acronyms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch acronyms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": acronyms,
	})
}

func getAcronyms(c *gin.Context) {
	usr := c.MustGet("user").(db.User)
	
	var acronyms []db.Acronym
	if err := database.
		Preload("Owner").
		Preload("Labels").
		Where(&db.Acronym{OwnerID: usr.ID}).
		Find(&acronyms).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch acronyms"})
		return
	}

	response := make([]db.AcronymResponse, 0, len(acronyms))
	for i, acronym := range acronyms {
		response = slices.Insert(response, i, acronym.ToJson())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func getAcronym(c *gin.Context) {
	id := c.Param("id")

	var usr db.User
	usr = c.MustGet("user").(db.User)

	// Validate UUID format
	validUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No acronym found with this ID"})
		return
	}

	var acronym db.Acronym

	if err := database.
		Preload("Owner").
		Preload("Labels").
		Where(&db.Acronym{
			UUID:    validUUID.String(),
			OwnerID: usr.ID,
		}).
		First(&acronym).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Acronym not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch acronym"})
		return
	}
	// Check if the user is the owner of the acronym
	// Assuming the user ID is stored in the context
	// TODO: Implement user ID retrieval from context
	// TODO: Implement ownership check

	c.JSON(http.StatusOK, gin.H{
		"data": acronym.ToJson(),
	})
}

func createAcronym(c *gin.Context) {
	usr := c.MustGet("user").(db.User)

	var acronym db.Acronym
	if err := c.ShouldBindJSON(&acronym); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set the owner ID to the current user's ID
	acronym.OwnerID = usr.ID
	acronym.Owner = usr

	if err := database.Create(&acronym).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create acronym"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Acronym created successfully",
		"data":    acronym,
	})
}

func updateAcronym(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func deleteAcronym(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

/* User API */

func getUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func getMe(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
		"user":    user,
	})
}

func createUser(c *gin.Context) {
	var userInput struct {
		Name      string `json:"name" binding:"required"`
		LegalName string `json:"legal_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := db.User{
		Name:      userInput.Name,
		LegalName: userInput.LegalName,
		Email:     userInput.Email,
		Password:  string(hashedPassword),
	}

	if err := database.Create(&user).Error; err != nil {
		// TODO: Handle database errors, e.g., unique constraint violations
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

func updateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func deleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func getAcronymsWithGrants(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func getHealth(c *gin.Context) {
	// Try to get session token from cookie to check auth status
	sessionCookie, err := c.Cookie("better-auth.session_token")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"authenticated": false,
		})
		return
	}
	
	// Validate session with better-auth
	sessionData, err := validateBetterAuthSession(sessionCookie)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"authenticated": false,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"authenticated": true,
		"user": gin.H{
			"email": sessionData.User.Email,
			"name": sessionData.User.Name,
		},
	})
}

// Session management endpoints

// invalidateSession allows the frontend to invalidate a session in the backend cache
func invalidateSession(c *gin.Context) {
	var request struct {
		SessionToken string `json:"sessionToken"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	
	if request.SessionToken == "" {
		// Try to get from cookie if not in body
		sessionCookie, err := c.Cookie("better-auth.session_token")
		if err == nil {
			request.SessionToken = sessionCookie
		}
	}
	
	if request.SessionToken != "" {
		invalidateSessionCache(request.SessionToken)
		log.Printf("Invalidated session cache for token: %s", request.SessionToken[:8]+"...")
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Session invalidated",
	})
}

// refreshSession forces a session refresh by removing it from cache
func refreshSession(c *gin.Context) {
	sessionCookie, err := c.Cookie("better-auth.session_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No session cookie found"})
		return
	}
	
	// Remove from cache to force fresh validation
	invalidateSessionCache(sessionCookie)
	
	// Validate fresh session
	sessionData, err := validateBetterAuthSession(sessionCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session refresh failed"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Session refreshed",
		"user": gin.H{
			"email": sessionData.User.Email,
			"name":  sessionData.User.Name,
		},
	})
}
