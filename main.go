package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

	r.Use(auth.IsAuthenticated(unprotectedRoutes))
	r.Use(auth.AuthorizationMiddleware(auth.InitFgaClient()))

	// Serve static files from the public directory
	r.Static("/public", "./public")

	/* ########################################## */
	/* ################# Website ################# */
	/* ########################################## */

	r.SetFuncMap(template.FuncMap{
		"getCurrentYear": func() int {
			return time.Now().Year()
		},
		"getAppName": func() string {
			return appName
		},
	})

	// Create a Template and parse files in order:
	tmpl := template.New("").Funcs(r.FuncMap)
	// parse base and partials
	tmpl = template.Must(tmpl.ParseGlob("templates/layouts/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/components/*.html"))
	// parse all views (they only define blocks)
	tmpl = template.Must(tmpl.ParseGlob("templates/views/*.html"))

	// Tell Gin to use this template
	r.SetHTMLTemplate(tmpl)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index_unauthenticated", gin.H{})
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
		c.HTML(http.StatusOK, "acronyms.html", gin.H{
			"Acronyms": []db.Acronym{
				{ShortForm: "API", LongForm: "Application Programming Interface"},
				{ShortForm: "HTTP", LongForm: "Hypertext Transfer Protocol"},
			},
		})
	})

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
	v1.GET("/acronyms/:id", getAcronym)
	v1.GET("/acronyms/search", searchAcronyms)

	authGroup := v1.Group("/auth")

	authGroup.POST("/logout", auth.LogoutHandler(database))
	authGroup.GET("/user", auth.UserHandler(database))
	authGroup.POST("/refresh", auth.RefreshHandler(database))

	// Routes for API
	/*
GET    /api/acronyms              → list all
POST   /api/acronyms              → create new
GET    /api/acronyms/{id}         → detail fragment
PUT    /api/acronyms/{id}         → update
DELETE /api/acronyms/{id}         → delete
GET    /api/acronyms/{id}/share   → share‑modal fragment
GET    /api/acronyms?tag={tag}    → list filtered by tag 

	*/

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
