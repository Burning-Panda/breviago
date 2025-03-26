package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create acronyms"})
		return
	}

	c.JSON(http.StatusCreated, createdAcronyms)
}

func updateAcronym(c *gin.Context) {
	id := c.Param("id")

	i, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var acronym = db.Acronym{
		ID: i,
	}

	if err := c.ShouldBindJSON(&acronym); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	database := db.GetDB()

	res, err := database.Exec(
		"UPDATE acronyms SET short_form = ?, long_form = ?, description = ? WHERE id = ?",
		acronym.ShortForm,
		acronym.LongForm,
		acronym.Description,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update acronym"})
		return
	}

	lastInsert, liErr := res.LastInsertId()

	if liErr != nil {
		c.JSON(http.StatusAccepted, gin.H{"error": "Acronym was updated, but it was not returned"})
	}

	var returnedAcronym = db.Acronym{
		ID: int(lastInsert),
	}

	err = database.QueryRow(
		"SELECT * FROM acronyms WHERE id = ?",
		lastInsert,
	).Scan(
		&returnedAcronym.ID,
		&returnedAcronym.ShortForm,
		&returnedAcronym.LongForm,
		&returnedAcronym.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated acronym"})
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
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}






