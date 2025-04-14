package main

import (
	"log"
	"net/http"

	"github.com/Burning-Panda/acronyms-vault/auth"
	"github.com/Burning-Panda/acronyms-vault/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	// Initialize database
	database := db.GetGormDB()
	if database == nil {
		log.Fatal("Failed to initialize database")
	}
	defer func() {
		if err := db.CloseGormDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Authentication routes
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.POST("/login", auth.LoginHandler(database))

	// Middleware
	r.Use(auth.AuthMiddleware)

	r.GET("/", getDefault)

	// API v1
	v1 := r.Group("/api/v1")

	// Public routes
	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	v1.GET("/acronyms", getAcronyms)
	v1.POST("/acronyms", createAcronym)
	v1.POST("/acronyms/batch", createAcronyms)
	v1.PUT("/acronyms/:id", updateAcronym)
	v1.DELETE("/acronyms/:id", deleteAcronym)
	v1.GET("/acronyms/:id", getAcronym)
	v1.GET("/acronyms/search", searchAcronyms)

	authGroup := v1.Group("/auth")
	authGroup.POST("/login", auth.LoginHandler(database))
	authGroup.POST("/register", auth.RegisterHandler(database))
	authGroup.POST("/logout", auth.LogoutHandler(database))
	authGroup.GET("/user", auth.UserHandler(database))
	authGroup.POST("/refresh", auth.RefreshHandler(database))

	r.Run(":8080")
}

func getDefault(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func getAcronyms(c *gin.Context) {
	var acronyms []db.Acronym
	if err := db.GetGormDB().Find(&acronyms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch acronyms"})
		return
	}
	c.JSON(http.StatusOK, acronyms)
}

func createAcronym(c *gin.Context) {
	var acronym db.Acronym
	if err := c.ShouldBindJSON(&acronym); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.GetGormDB().Create(&acronym).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create acronym"})
		return
	}

	c.JSON(http.StatusCreated, acronym)
}

func createAcronyms(c *gin.Context) {
	var acronyms []db.Acronym
	if err := c.ShouldBindJSON(&acronyms); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.GetGormDB().Create(&acronyms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create acronyms"})
		return
	}

	c.JSON(http.StatusCreated, acronyms)
}

func updateAcronym(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var acronym db.Acronym
	if err := db.GetGormDB().Where("uuid = ?", uuid.String()).First(&acronym).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acronym not found"})
		return
	}

	if err := c.ShouldBindJSON(&acronym); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.GetGormDB().Save(&acronym).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update acronym"})
		return
	}

	c.JSON(http.StatusOK, acronym)
}

func deleteAcronym(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	if err := db.GetGormDB().Where("uuid = ?", uuid.String()).Delete(&db.Acronym{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete acronym"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Acronym deleted successfully"})
}

func getAcronym(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var acronym db.Acronym
	if err := db.GetGormDB().Where("uuid = ?", uuid.String()).First(&acronym).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acronym not found"})
		return
	}

	c.JSON(http.StatusOK, acronym)
}

func searchAcronyms(c *gin.Context) {
	query := c.Query("query")
	var acronyms []db.Acronym

	if err := db.GetGormDB().Where("short_form ILIKE ? OR long_form ILIKE ?", 
		"%"+query+"%", "%"+query+"%").Find(&acronyms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search acronyms"})
		return
	}

	c.JSON(http.StatusOK, acronyms)
}




