package main

import (
	"log"
	"net/http"
	"slices"

	"github.com/Burning-Panda/breviago/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	database *gorm.DB
)

func AuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sorry, you are not logged in"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// TODO: Move this to a separate package
// Add use to the gin context for authentication
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement a token validation for validating if user is logged in
		// For now, we will just set a user ID in the context

		var user db.User
		err := database.
			Where(&db.User{UUID: "00000000-0000-0000-0000-000000000000"}). // TODO: Get user UUID from token
			First(&user).
			Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
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

	r := gin.Default()

	// Middleware
	r.Use(AuthenticationMiddleware(database))

	// Routes
	r.GET("/", getDefault)

	api := r.Group("/api", AuthenticatedMiddleware(), UserMiddleware())
	api.GET("/", testAcronyms)

	api.GET("/acronyms", getAcronyms)
	api.POST("/acronyms", createAcronym)

	api.GET("/acronyms/:id", getAcronym)
	api.PUT("/acronyms/:id", updateAcronym)
	api.DELETE("/acronyms/:id", deleteAcronym)

	// User API
	user := api.Group("/users")
	user.GET("/", getUser)
	user.POST("/", createUser)
	user.GET("/me", getMe)
	user.GET("/:id", getUser)
	user.PUT("/:id", updateUser)
	user.DELETE("/:id", deleteUser)

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
	var acronyms []db.Acronym
	if err := database.
		Preload("Owner").
		Preload("Labels").
		Where(&db.Acronym{Visibility: db.VisibilityPrivate, OwnerID: 1}).
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


