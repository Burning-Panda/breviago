package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Burning-Panda/breviago/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	database *gorm.DB
)

func main() {
	// Initialize database
	database = db.GetGormDB()
	if database == nil {
		log.Fatal("Failed to initialize database")
	}
	defer func() {
		if err := db.CloseGormDB(); err != nil {
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

	api := r.Group("/api")
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
		response[i] = acronym.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func getAcronym(c *gin.Context) {
	id := c.Param("id")

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
		Where(&db.Acronym{UUID: validUUID.String()}).
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
		"data": acronym.ToResponse(),
	})
}

func createAcronym(c *gin.Context) {
	var acronym db.Acronym
	if err := c.ShouldBindJSON(&acronym); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

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
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
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
