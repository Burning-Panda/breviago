package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db *sql.DB

func main() {
	initDB()

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
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
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

func initDB() {
	db, err := sql.Open("sqlite3", "./acronyms.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}











