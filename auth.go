package main

import (
	"github.com/Burning-Panda/breviago/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Temporary Authentication Middleware
// This middleware should add the user to the context
func AuthenticationMiddleware(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		/*
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		*/

		// TODO: Validate token
		// TODO: Add user to context

		var user db.User
		database.Where("uuid = ?", "00000000-0000-0000-0000-000000000000").First(&user)
		

		c.Set("user", gin.H{"user": user})
		c.Next()
	}
}
