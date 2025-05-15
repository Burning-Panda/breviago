package db

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Session struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"OnDelete:CASCADE"`
	Token     string    `gorm:"unique;index" json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uint           `gorm:"primaryKey"`
	UUID      string         `gorm:"type:text;unique;index" json:"uuid"`
	Name      string         `gorm:"unique" json:"name"`
	LegalName string         `gorm:"unique" json:"legal_name"`
	Email     string         `gorm:"unique" json:"email"`
	Password  string         `gorm:"default:''"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Session   Session        `gorm:"foreignKey:UserID"`
	Settings  []UserSettings `gorm:"foreignKey:UserID"`
}

type UserSettings struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id" gorm:"OnDelete:CASCADE"`
	Setting   string    `json:"setting"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	Event     string         `gorm:"not null" json:"event"`
	Action    string         `gorm:"not null" json:"action"`
	Data      string         `gorm:"not null" json:"data"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Organization struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	UUID        string `gorm:"type:text" json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Members []OrganizationMember

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type OrganizationMember struct {
	ID      uint   `gorm:"primaryKey"`
	UUID    string `gorm:"type:text" json:"uuid"`
	IsAdmin bool   `json:"is_admin"`

	UserID         uint         `gorm:"uniqueIndex:idx_org_user" json:"user_id"`
	User           User         `gorm:"foreignKey:UserID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	OrganizationID uint         `gorm:"uniqueIndex:idx_org_user" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"organization"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// VisibilityType defines who can see an acronym
type VisibilityType string

const (
	VisibilityPrivate      VisibilityType = "private"      // Only visible to owner
	VisibilityPublic       VisibilityType = "public"       // Visible to everyone
	VisibilityOrganization VisibilityType = "organization" // Visible to specific organization
	VisibilityUser         VisibilityType = "user"         // Visible to specific user
)

type Acronym struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UUID        string         `gorm:"type:uuid;unique;index" json:"uuid"`
	Acronym     string         `json:"acronym"`
	Meaning     string         `json:"meaning"`
	Description *string        `json:"description"`
	OwnerID     uint           `json:"owner_id"`
	Visibility  VisibilityType `gorm:"type:text;default:'private'" json:"visibility"`

	Owner     User              `gorm:"foreignKey:OwnerID" json:"owner"`
	Related   []Acronym         `gorm:"many2many:acronym_relations;" json:"related"`
	Labels    []Label           `gorm:"many2many:acronym_labels;" json:"labels"`
	Notes     []Notes           `gorm:"foreignKey:AcronymID" json:"notes"`
	Grants    []AcronymGrant    `gorm:"foreignKey:AcronymID" json:"grants"`
	Revisions []AcronymRevision `gorm:"foreignKey:ForeignKeyID" json:"revisions"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Label struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UUID      string         `gorm:"type:text" json:"uuid"`
	Label     string         `json:"label"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Notes struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UUID      string         `gorm:"type:text" json:"uuid"`
	AcronymID uint           `json:"acronym_id"`
	UserID    uint           `json:"user_id"`
	Note      string         `json:"note"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
}

type AcronymGrant struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UUID        string         `gorm:"type:text" json:"uuid"`
	AcronymID   uint           `json:"acronym_id"`
	GranteeID   uint           `json:"grantee_id"`
	GranteeType string         `json:"grantee_type"` // "user", "organization", or "group"
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type AcronymRevision struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UUID         string         `gorm:"type:text" json:"uuid"`
	ForeignKeyID uint           `json:"foreign_key_id"`
	UserID       uint           `json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user"`
	Action       string         `json:"action"`
	OldValue     string         `json:"old_value"`
	NewValue     string         `json:"new_value"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
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
func (u *User) SetUUID(uuid string)                { u.UUID = uuid }
func (o *Organization) SetUUID(uuid string)        { o.UUID = uuid }
func (om *OrganizationMember) SetUUID(uuid string) { om.UUID = uuid }
func (a *Acronym) SetUUID(uuid string)             { a.UUID = uuid }
func (l *Label) SetUUID(uuid string)               { l.UUID = uuid }
func (ac *Notes) SetUUID(uuid string)              { ac.UUID = uuid }
func (ag *AcronymGrant) SetUUID(uuid string)       { ag.UUID = uuid }

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

func (a *Acronym) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, a)
}

func (l *Label) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, l)
}

func (ac *Notes) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, ac)
}

func (ag *AcronymGrant) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, ag)
}

// Revisionable interface for models that support revision history
// Must return the primary key, type name, and a summary of the change
// Optionally, you can add more methods as needed
//

type Revisionable interface {
	GetID() uint
	GetTypeName() string
}

// revisionCreateHelper creates a revision for any Revisionable model
func revisionCreateHelper(tx *gorm.DB, model Revisionable, action string, userID uint) error {
	id := model.GetID()

	// Use reflection to create a pointer to the model type
	oldValue := reflect.New(reflect.TypeOf(model).Elem()).Interface()
	if err := tx.Unscoped().First(oldValue, id).Error; err != nil {
		return err
	}

	// Save only the changed columns
	// Acronym
	// Meaning
	// Description
	// OwnerID

	oldJSON, _ := json.Marshal(oldValue)
	newJSON, _ := json.Marshal(model)

	revision := AcronymRevision{
		UUID:         uuid.New().String(),
		ForeignKeyID: id,
		UserID:       userID,
		Action:       action,
		OldValue:     string(oldJSON),
		NewValue:     string(newJSON),
	}
	return tx.Create(&revision).Error
}

// Example implementation for Acronym
func (a *Acronym) GetID() uint         { return a.ID }
func (a *Acronym) GetTypeName() string { return "Acronym" }
func (a *Acronym) BeforeUpdate(tx *gorm.DB) error {
	// Strip out the columns that is not necessary to save
	changed := Acronym{
		Acronym:     a.Acronym,
		Meaning:     a.Meaning,
		Description: a.Description,
		OwnerID:     a.OwnerID,
	}

	return revisionCreateHelper(tx, &changed, "update", 1) // Replace 0 with actual user ID if available
}

/* func (a *Acronym) AfterCreate(tx *gorm.DB) error {
	return revisionCreateHelper(tx, a, "create", 1) // Replace 0 with actual user ID if available
}
*/

func GetGormDB(database *gorm.DB) *gorm.DB {
	if database != nil {
		return database
	}

	gormDB, err := gorm.Open(sqlite.Open("acronyms.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil
	}

	// Auto migrate the schema
	err = gormDB.AutoMigrate(
		&User{},
		&Organization{},
		&OrganizationMember{},
		&Acronym{},
		&AcronymGrant{},
		&Label{},
		&AcronymRevision{},
		&Notes{},
		&Session{},
		&UserSettings{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return gormDB
}

func CloseGormDB(database *gorm.DB) error {
	if database != nil {
		sqlDB, err := database.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return nil
}
