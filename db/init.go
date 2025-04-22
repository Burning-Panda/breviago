package db

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitDB(db *gorm.DB) {
	start := time.Now()
	fmt.Println("Initializing database...")

	var adminUser User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create admin user
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
			if err != nil {
				fmt.Printf("failed to hash admin password: %v", err)
				return
			}

			adminUser = User{
				Username: "admin",
				Password: string(hashedPassword),
				Email:    "admin@breviago.com",
			}

			if err := db.Create(&adminUser).Error; err != nil {
				fmt.Printf("failed to create admin user: %v", err)
				return
			}
			
			// Update the user's UUID to be a default UUID for easier testing
			if err := db.Model(&User{}).Where("email= ?", "admin@breviago.com").Update("uuid", "00000000-0000-0000-0000-000000000000").Error; err != nil {
				fmt.Printf("failed to update admin user UUID: %v", err)
				return
			}
			fmt.Println("Created default admin user")
		} else {
			fmt.Printf("failed to check for admin user: %v", err)
			return
		}
	}

	// Check if root organization exists
	var rootOrg Organization
	if err := db.Where("name = ?", "Root Organization").First(&rootOrg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create root organization
			rootOrg = Organization{
				Name:        "Root Organization",
				Description: "The root organization for Breviago",
			}

			if err := db.Create(&rootOrg).Error; err != nil {
				fmt.Printf("failed to create root organization: %v", err)
				return
			}
			fmt.Println("Created root organization")

			// Add admin as member of root organization
			orgMember := OrganizationMember{
				OrganizationID: rootOrg.ID,
				UserID:        adminUser.ID,
				IsAdmin:       true,
			}

			if err := db.Create(&orgMember).Error; err != nil {
				fmt.Printf("failed to add admin to root organization: %v", err)
				return
			}
			fmt.Println("Added admin to root organization")
		} else {
			fmt.Printf("failed to check for root organization: %v", err)
			return
		}
	}

	// Check if Breviago acronym exists
	var breviagoAcronym Acronym
	if err := db.Where("short_form = ?", "breviago").First(&breviagoAcronym).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create Breviago acronym
			breviagoAcronym = Acronym{
				UUID:        "00000000-0000-0000-0000-000000000000",
				ShortForm:   "breviago",
				LongForm:    "Is a application for remembering abbreviations",
				Description: "The main application for managing and remembering abbreviations",
			}

			if err := db.Create(&breviagoAcronym).Error; err != nil {
				fmt.Printf("failed to create Breviago acronym: %v", err)
				return
			}
			fmt.Println("Created Breviago acronym")
		} else {
			fmt.Printf("failed to check for Breviago acronym: %v", err)
			return
		}
	}

	// Check if root folder exists
	var rootFolder Folder
	if err := db.Where("name = ?", "Root Folder").First(&rootFolder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create root folder
			rootFolder = Folder{
				Name:        "Root Folder",
				Description: "The root folder for Breviago",
				OwnerID:     adminUser.ID,
				OwnerType:   "user",
			}

			if err := db.Create(&rootFolder).Error; err != nil {
				fmt.Printf("failed to create root folder: %v", err)
				return
			}
			fmt.Println("Created root folder")
		} else {
			fmt.Printf("failed to check for root folder: %v", err)
			return
		}
	}

	// Check if root document exists
	var rootDocument Document
	if err := db.Where("name = ?", "Root Document").First(&rootDocument).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create root document
			rootDocument = Document{
				Name:       "Root Document",
				Content:	"The root document for Breviago",
				OwnerID:    adminUser.ID,
				OwnerType:  "user",
			}

			if err := db.Create(&rootDocument).Error; err != nil {
				fmt.Printf("failed to create root document: %v", err)
				return
			}
			fmt.Println("Created root document")
		} else {
			fmt.Printf("failed to check for root document: %v", err)
			return
		}
	}

	fmt.Printf("Database initialized successfully in %s\n", time.Since(start))
}
