package main

import (
	"log"
	"net/http"
	"strconv"

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
	rows, err := database.Query("SELECT id, short_form, long_form, created_at, updated_at FROM acronyms ORDER BY short_form DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch acronyms"})
		return
	}
	defer rows.Close()

	var acronyms []gin.H
	for rows.Next() {
		var id int
		var shortForm, longForm, createdAt, updatedAt string

		err := rows.Scan(&id, &shortForm, &longForm, &createdAt, &updatedAt);
		if  err != nil {
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

	var acronym = db.Acronym{
		ShortForm: "",
		LongForm: "",
		Description: "",
	}

	if err := c.ShouldBindJSON(&acronym); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	db := db.GetDB()

	res, err := db.Exec(
		"INSERT INTO acronyms (short_form, long_form, description) VALUES (?, ?, ?)",
		acronym.ShortForm,
		acronym.LongForm,
		acronym.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create acronym"})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get last insert id"})
		return
	}

	db.QueryRow(
		"SELECT id, short_form, long_form, description FROM acronyms WHERE id = ?", id,
	).Scan(
		&acronym.ID,
		&acronym.ShortForm,
		&acronym.LongForm,
		&acronym.Description,
	)

	c.JSON(http.StatusOK, acronym)
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
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func searchAcronyms(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}






