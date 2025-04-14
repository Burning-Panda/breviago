package db

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	gormDB *gorm.DB
)

type Session struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Token     string    `gorm:"unique;index" json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:text;unique;index" json:"uuid"`
	Username  string    `gorm:"unique" json:"username"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Session   Session   `gorm:"foreignKey:UserID"`
}

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Event     string    `gorm:"not null" json:"event"`
	Action    string    `gorm:"not null" json:"action"`
	Data      string    `gorm:"not null" json:"data"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

type Organization struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrganizationMember struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UUID           string    `gorm:"type:text" json:"uuid"`
	OrganizationID uint      `json:"organization_id"`
	UserID         uint      `json:"user_id"`
	IsAdmin        bool      `json:"is_admin"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Group struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:text" json:"uuid"`
	GroupID   uint      `json:"group_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Folder struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uint      `json:"owner_id"`
	OwnerType   string    `json:"owner_type"` // "user" or "organization"
	ParentID    *uint     `json:"parent_id"`  // Nullable for root folders
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FolderGrant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:text" json:"uuid"`
	FolderID  uint      `json:"folder_id"`
	GranteeID uint      `json:"grantee_id"`
	GranteeType string  `json:"grantee_type"` // "user", "organization", or "group"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Document struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	OwnerID     uint      `json:"owner_id"`
	OwnerType   string    `json:"owner_type"` // "user" or "organization"
	FolderID    uint      `json:"folder_id"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DocumentGrant struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	DocumentID  uint      `json:"document_id"`
	GranteeID   uint      `json:"grantee_id"`
	GranteeType string    `json:"grantee_type"` // "user", "organization", or "group"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}


type Acronym struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:uuid" json:"uuid"`
	ShortForm   string    `json:"short_form"`
	LongForm    string    `json:"long_form"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Label struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:text" json:"uuid"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AcronymCategory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UUID       string    `gorm:"type:text" json:"uuid"`
	AcronymID  uint      `json:"acronym_id"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func GetGormDB() *gorm.DB {
	if gormDB == nil {
		var err error
		gormDB, err = gorm.Open(sqlite.Open("acronyms.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		// Auto migrate the schema
		err = gormDB.AutoMigrate(
			&User{},
			&Organization{},
			&OrganizationMember{},
			&Group{},
			&GroupMember{},
			&Folder{},
			&FolderGrant{},
			&Document{},
			&DocumentGrant{},
		)
		if err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
	}

	return gormDB
}

func CloseGormDB() error {
	if gormDB != nil {
		sqlDB, err := gormDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// BeforeCreate hooks for UUID generation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UUID = uuid.New().String()
	return nil
}

func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	o.UUID = uuid.New().String()
	return nil
}

func (om *OrganizationMember) BeforeCreate(tx *gorm.DB) error {
	om.UUID = uuid.New().String()
	return nil
}

func (g *Group) BeforeCreate(tx *gorm.DB) error {
	g.UUID = uuid.New().String()
	return nil
}

func (gm *GroupMember) BeforeCreate(tx *gorm.DB) error {
	gm.UUID = uuid.New().String()
	return nil
}

func (f *Folder) BeforeCreate(tx *gorm.DB) error {
	f.UUID = uuid.New().String()
	return nil
}

func (fg *FolderGrant) BeforeCreate(tx *gorm.DB) error {
	fg.UUID = uuid.New().String()
	return nil
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	d.UUID = uuid.New().String()
	return nil
}

func (dg *DocumentGrant) BeforeCreate(tx *gorm.DB) error {
	dg.UUID = uuid.New().String()
	return nil
} 