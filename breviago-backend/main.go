package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Burning-Panda/acronyms-vault/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	database *gorm.DB
)

type Owner struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

// AcronymResponse represents the public view of an Acronym
type AcronymResponse struct {
	UUID    string    `json:"uuid"`
	Acronym string    `json:"acronym"`
	Meaning string    `json:"meaning"`
	Owner   Owner   	`json:"owner"`
	Labels  []db.Label `json:"labels"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
	r.Use(func(ctx *gin.Context) {
		// TODO: Add Authentication middleware here
		ctx.Next()
	})

	r.GET("/", getDefault)

	api := r.Group("/api")
	api.GET("/", testAcronyms)

	api.GET("/acronyms", getAcronyms)
	api.POST("/acronyms", createAcronym)

	api.GET("/acronyms/:id", getAcronym)
	api.PUT("/acronyms/:id", updateAcronym)
	api.DELETE("/acronyms/:id", deleteAcronym)

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
		Find(&acronyms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch acronyms"})
		return
	}

	// Convert to response format
	response := make([]AcronymResponse, len(acronyms))
	for i, a := range acronyms {
		response[i] = AcronymResponse{
			UUID:      a.UUID,
			Acronym:   a.Acronym,
			Meaning:   a.Meaning,
			Owner:     Owner{
				Name: a.Owner.Name,
				UUID: a.Owner.UUID,
			},
			Labels:    a.Labels,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func getAcronym(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
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
		"data": acronym,
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







