package main

import (
	"log"
	"net/http"

	"github.com/Burning-Panda/acronyms-vault/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database := db.GetDB()
	if database == nil {
		log.Fatal("Failed to initialize database")
	}
	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	r := gin.Default()

	// Middleware
	
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", getDefault)

	// API v1
	v1 := r.Group("/api/v1")

	// Public the routes
	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	v1.GET("/acronyms", getAcronyms)

	v1.POST("/acronyms", createAcronym)

	v1.PUT("/acronyms/:id", updateAcronym)

	v1.DELETE("/acronyms/:id", deleteAcronym)

	v1.GET("/acronyms/:id", getAcronym)

	v1.GET("/acronyms/search", searchAcronyms)

	
	r.Run(":8080")
}

func getDefault(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func getAcronyms(c *gin.Context) {
	database := db.GetDB()
	rows, err := database.Query("SELECT id, short_form, long_form, created_at, updated_at FROM acronyms")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch acronyms"})
		return
	}
	defer rows.Close()

	var acronyms []gin.H
	for rows.Next() {
		var id int
		var shortForm, longForm, createdAt, updatedAt string
		if err := rows.Scan(&id, &shortForm, &longForm, &createdAt, &updatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan acronym"})
			return
		}
		acronyms = append(acronyms, gin.H{
			"id": id,
			"short_form": shortForm,
			"long_form": longForm,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	c.JSON(http.StatusOK, acronyms)
}

func createAcronym(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
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

func getAcronym(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func searchAcronyms(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}






