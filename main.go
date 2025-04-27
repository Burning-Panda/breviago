package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Burning-Panda/acronyms-vault/auth"
	"github.com/Burning-Panda/acronyms-vault/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const appName = "Breviago"

var unprotectedRoutes = []string{"", "/", "/login", "/register", "/public", "/favicon.ico",
	// Temp
	"/home",
}

func initApplication(database *gorm.DB) error {
	// Check if admin user exists
	db.InitDB(database)

	return nil
}

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

	// Initialize application
	if err := initApplication(database); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	r := gin.Default()

	// Add debug middleware
	r.Use(func(c *gin.Context) {
		fmt.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Serve static files from the public directory
	r.Static("/public", "./public")

	// Add authentication middleware after static files
	// r.Use(auth.IsAuthenticated(unprotectedRoutes))
	// r.Use(auth.AuthorizationMiddleware(auth.InitFgaClient()))

	/* ########################################## */
	/* ################# Website ################# */
	/* ########################################## */

	// Create a Template and parse files in order:
	tmpl := template.New("")
	
	// Add functions to the template
	tmpl.Funcs(template.FuncMap{
		"getCurrentYear": func() int {
			return time.Now().Year()
		},
		"getAppName": func() string {
			return appName
		},
		"getRelatedAcronyms": func(acronyms []db.Acronym) []db.Acronym {
			mockRelatedAcronyms := []db.Acronym{
				{Acronym: "API", Meaning: "Application Programming Interface"},
				{Acronym: "HTTP", Meaning: "Hypertext Transfer Protocol"},
			}
			return mockRelatedAcronyms
		},
		"join": func(strs []string, sep string) string {
			return strings.Join(strs, sep)
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"getAllAcronyms": func() []db.Acronym {
			var acronyms []db.Acronym
			db.GetGormDB().Find(&acronyms)
			return acronyms
		},
	})

	// parse base and partials
	//tmpl = template.Must(tmpl.ParseGlob("templates/layouts/*.html"))
	//tmpl = template.Must(tmpl.ParseGlob("templates/components/*.html"))
	// parse all views (they only define blocks)
	tmpl = template.Must(tmpl.ParseGlob("templates/views/*.html"))

	// Tell Gin to use this template
	r.SetHTMLTemplate(tmpl)

	r.GET("/", func(c *gin.Context) {
		//c.HTML(http.StatusOK, "index_unauthenticated", gin.H{})
		c.HTML(http.StatusOK, "views/acronyms", gin.H{})
	})

	r.GET("/login", func(c *gin.Context) {
		fmt.Println("Rendering login page")
		c.HTML(http.StatusOK, "login", gin.H{
			"Request": c.Request,
		})
	})
	r.GET("/register", func(c *gin.Context) {
		fmt.Println("Rendering register page")
		c.HTML(http.StatusOK, "register", gin.H{
			"Request": c.Request,
		})
	})

	r.GET("/home", func(c *gin.Context) {
		fmt.Println("Rendering index page")
		c.HTML(http.StatusOK, "index", gin.H{})
	})

	r.POST("/login", auth.LoginHandler(database))
	r.POST("/register", auth.RegisterHandler(database))


	

	r.GET("/acronyms", func(c *gin.Context) {
		c.HTML(http.StatusOK, "testing", gin.H{
			"Acronyms": []db.Acronym{
				{Acronym: "API", Meaning: "Application Programming Interface"},
				{Acronym: "HTTP", Meaning: "Hypertext Transfer Protocol"},
			},
		})
	})

	r.GET("/acronyms/:id", getAcronym)
	/* ######################################### */
	/* ################## API ################## */
	/* ######################################### */

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
	v1.GET("/acronyms/:id", apiGetAcronym)
	v1.GET("/acronyms/search", searchAcronyms)

	authGroup := v1.Group("/auth")

	authGroup.POST("/logout", auth.LogoutHandler(database))
	authGroup.GET("/user", auth.UserHandler(database))
	authGroup.POST("/refresh", auth.RefreshHandler(database))

	r.Run(":8060")
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

func apiGetAcronym(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var acronym db.Acronym
	if err := db.GetGormDB().
		Preload("Related").
		Preload("Labels").
		Preload("Comments").
		Preload("History").
		Preload("Grants").
		Preload("Owner", "owner_type = ?", "user").
		Preload("Owner", "owner_type = ?", "organization").
		Where("uuid = ?", uuid.String()).
		First(&acronym).Error; err != nil {
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

/* ############################################## */
/* ###############      HTML      ############### */
/* ############################################## */

func getAcronym(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	
	var acronym db.Acronym
	if err := db.GetGormDB().
		Preload("Related").
		Preload("Labels").
		Preload("Comments").
		Preload("History").
		Preload("Grants").
		Preload("Owner", "owner_type = ?", "user").
		Preload("Owner", "owner_type = ?", "organization").
		Where("uuid = ?", uuid.String()).
		First(&acronym).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acronym not found"})
		return
	}

	c.HTML(http.StatusOK, "components/AcronymContent", gin.H{"Acronym": acronym})
}
