package db

import (
	"fmt"
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
	Name      string    `gorm:"unique" json:"name"`
	LegalName string    `gorm:"unique" json:"legal_name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Session   Session   `gorm:"foreignKey:UserID"`
	Settings  []UserSettings `gorm:"foreignKey:UserID"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type UserSettings struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Setting   string    `json:"setting"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Event     string    `gorm:"not null" json:"event"`
	Action    string    `gorm:"not null" json:"action"`
	Data      string    `gorm:"not null" json:"data"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Organization struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type OrganizationMember struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UUID           string    `gorm:"type:text" json:"uuid"`
	OrganizationID uint      `json:"organization_id"`
	UserID         uint      `json:"user_id"`
	IsAdmin        bool      `json:"is_admin"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Folder struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text;unique;index" json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uint      `json:"owner_id"`
	OwnerType   string    `gorm:"type:text" json:"owner_type"` // "user" or "organization"
	ParentID    *uint     `json:"parent_id"`  // Nullable for root folders
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type FolderGrant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:text" json:"uuid"`
	FolderID  uint      `json:"folder_id"`
	GranteeID uint      `json:"grantee_id"`
	GranteeType string  `json:"grantee_type"` // "user", "organization", or "group"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Document struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	OwnerID     uint      `json:"owner_id"`
	OwnerType   string    `gorm:"type:text" json:"owner_type"` // "user" or "organization"
	FolderID    uint      `json:"folder_id"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type DocumentGrant struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	DocumentID  uint      `json:"document_id"`
	GranteeID   uint      `json:"grantee_id"`
	GranteeType string    `json:"grantee_type"` // "user", "organization", or "group"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Acronym struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:uuid;unique;index" json:"uuid"`
	Acronym     string    `json:"acronym"`
	Meaning     string    `json:"meaning"`
	Description string    `json:"description"`
	OwnerID     uint      `json:"owner_id"`
	OwnerType   string    `gorm:"type:text" json:"owner_type"` // "user" or "organization"

	Related     []Acronym `gorm:"many2many:acronym_relations;"`
	Labels      []Label   `gorm:"many2many:acronym_labels;"`
	Comments    []AcronymComment `gorm:"foreignKey:AcronymID"`
	History     []AcronymHistory `gorm:"foreignKey:AcronymID"`
	Grants      []AcronymGrant   `gorm:"foreignKey:AcronymID"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UUID        string    `gorm:"type:text" json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Label struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:text" json:"uuid"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type AcronymCategory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UUID       string    `gorm:"type:text" json:"uuid"`
	AcronymID  uint      `json:"acronym_id"`
	CategoryID uint      `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type AcronymHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:text" json:"uuid"`
	AcronymID uint      `json:"acronym_id"`
	Action    string    `json:"action"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type AcronymComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:text" json:"uuid"`
	AcronymID uint      `json:"acronym_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type AcronymGrant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:text" json:"uuid"`
	AcronymID uint      `json:"acronym_id"`
	GranteeID uint      `json:"grantee_id"`
	GranteeType string    `json:"grantee_type"` // "user", "organization", or "group"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// UUIDable interface for models that need UUID generation
type UUIDable interface {
	SetUUID(uuid string)
}

// GenerateUUID generates a UUID for any model that implements UUIDable
func GenerateUUID(tx *gorm.DB, model UUIDable) error {
	
	model.SetUUID(uuid.New().String())

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// Implement UUIDable for all models that need UUID generation
func (u *User) SetUUID(uuid string) { u.UUID = uuid }
func (o *Organization) SetUUID(uuid string) { o.UUID = uuid }
func (om *OrganizationMember) SetUUID(uuid string) { om.UUID = uuid }
func (f *Folder) SetUUID(uuid string) { f.UUID = uuid }
func (fg *FolderGrant) SetUUID(uuid string) { fg.UUID = uuid }
func (d *Document) SetUUID(uuid string) { d.UUID = uuid }
func (dg *DocumentGrant) SetUUID(uuid string) { dg.UUID = uuid }
func (a *Acronym) SetUUID(uuid string) { a.UUID = uuid }
func (c *Category) SetUUID(uuid string) { c.UUID = uuid }
func (l *Label) SetUUID(uuid string) { l.UUID = uuid }
func (ac *AcronymCategory) SetUUID(uuid string) { ac.UUID = uuid }
func (ah *AcronymHistory) SetUUID(uuid string) { ah.UUID = uuid }
func (ac *AcronymComment) SetUUID(uuid string) { ac.UUID = uuid }
func (ag *AcronymGrant) SetUUID(uuid string) { ag.UUID = uuid }

// BeforeCreate hooks using the generic GenerateUUID function
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, u)
}

func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, o)
}

func (om *OrganizationMember) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, om)
}

func (f *Folder) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, f)
}

func (fg *FolderGrant) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, fg)
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, d)
}

func (dg *DocumentGrant) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, dg)
}

func (a *Acronym) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, a)
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, c)
}

func (l *Label) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, l)
}

func (ac *AcronymCategory) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, ac)
}

func (ah *AcronymHistory) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, ah)
}

func (ac *AcronymComment) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, ac)
}

func (ag *AcronymGrant) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, ag)
}

func (a *Acronym) AfterSave(tx *gorm.DB) error {
	// Create a new history record when an acronym is saved
	history := AcronymHistory{
		AcronymID: a.ID,
		Action:    "save",
		Data:      fmt.Sprintf("%+v", a),
	}
	
	return tx.Create(&history).Error
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
			&Folder{},
			&FolderGrant{},
			&Document{},
			&DocumentGrant{},
			&Acronym{},
			&Category{},
			&Label{},
			&AcronymCategory{},
			&AcronymHistory{},
			&AcronymComment{},
			&AcronymGrant{},
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
