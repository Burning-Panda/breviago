package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Burning-Panda/acronyms-vault/db"
	"github.com/google/uuid"

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
	r.Use(func(ctx *gin.Context) {
		// TODO: Add Authentication middleware here
		ctx.Next()
	})

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
	v1.POST("/acronyms/batch", createAcronyms)

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
	rows, err := database.Query("SELECT uuid, short_form, long_form, created_at, updated_at FROM acronyms ORDER BY short_form DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch acronyms"})
		return
	}
	defer rows.Close()

	var acronyms []gin.H
	for rows.Next() {
		var uuid uuid.UUID
		var shortForm, longForm, createdAt, updatedAt string

		err := rows.Scan(&uuid, &shortForm, &longForm, &createdAt, &updatedAt);
		if  err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan acronym"})
			return
		}
		acronyms = append(acronyms, gin.H{
			"id": uuid,
			"short_form": shortForm,
			"long_form": longForm,
			"created_at": createdAt,
			"updated_at": updatedAt,
		})
	}

	c.JSON(http.StatusOK, acronyms)
}

func createAcronym(c *gin.Context) {
	var acronym db.Acronym
	if err := c.ShouldBindJSON(&acronym); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.InsertAcronym(&acronym); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create acronym"})
		return
	}

	c.JSON(http.StatusCreated, acronym)
}

func createAcronyms(c *gin.Context) {
	var acronyms []db.Acronym
	if err := c.ShouldBindJSON(&acronyms); err != nil {
		fmt.Println("Error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	createdAcronyms, err := db.InsertAcronyms(acronyms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create acronyms: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, createdAcronyms)
}

func updateAcronym(c *gin.Context) {
	id := c.Param("id")

	i, err := uuid.Parse(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var acronym = db.Acronym{
		UUID: i,
	}

	if err := c.ShouldBindJSON(&acronym); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	database := db.GetDB()

	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM acronyms WHERE uuid = ?)", i).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if acronym exists"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acronym not found"})
		return
	}

	var returnedAcronym db.Acronym
	updateErr := database.QueryRow(
		"UPDATE acronyms SET short_form = ?, long_form = ?, description = ? WHERE uuid = ?",
		acronym.ShortForm,
		acronym.LongForm,
		acronym.Description,
		i,
	).Scan(
		&returnedAcronym.UUID,
		&returnedAcronym.ShortForm,
		&returnedAcronym.LongForm,
		&returnedAcronym.Description,
		&returnedAcronym.CreatedAt,
		&returnedAcronym.UpdatedAt,
	)

	if updateErr != nil {
		fmt.Println("Error: ", updateErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update acronym"})
		return
	}

	c.JSON(http.StatusOK, returnedAcronym)
}

func deleteAcronym(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func getAcronym(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var acronym db.Acronym
	err = db.GetDB().QueryRow(
		"SELECT uuid, short_form, long_form, description, created_at, updated_at FROM acronyms WHERE uuid = ?",
		uuid,
	).Scan(
		&acronym.UUID,
		&acronym.ShortForm,
		&acronym.LongForm,
		&acronym.Description,
		&acronym.CreatedAt,
		&acronym.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Acronym not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch acronym"})
		return
	}

	c.JSON(http.StatusOK, acronym)
}

func searchAcronyms(c *gin.Context) {
	searchTerm := c.Query("query")

	acronyms, err := db.SearchAcronyms(searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search acronyms"})
		return
	}
	c.JSON(http.StatusOK, acronyms)
}






