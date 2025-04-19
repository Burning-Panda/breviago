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

// UUIDable interface for models that need UUID generation
type UUIDable interface {
	SetUUID(uuid string)
}

// GenerateUUID generates a UUID for any model that implements UUIDable
func GenerateUUID(tx *gorm.DB, model UUIDable) error {
	model.SetUUID(uuid.New().String())
	return nil
}

// Implement UUIDable for all models that need UUID generation
func (u *User) SetUUID(uuid string) { u.UUID = uuid }
func (o *Organization) SetUUID(uuid string) { o.UUID = uuid }
func (om *OrganizationMember) SetUUID(uuid string) { om.UUID = uuid }
func (g *Group) SetUUID(uuid string) { g.UUID = uuid }
func (gm *GroupMember) SetUUID(uuid string) { gm.UUID = uuid }
func (f *Folder) SetUUID(uuid string) { f.UUID = uuid }
func (fg *FolderGrant) SetUUID(uuid string) { fg.UUID = uuid }
func (d *Document) SetUUID(uuid string) { d.UUID = uuid }
func (dg *DocumentGrant) SetUUID(uuid string) { dg.UUID = uuid }
func (a *Acronym) SetUUID(uuid string) { a.UUID = uuid }
func (c *Category) SetUUID(uuid string) { c.UUID = uuid }

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

func (g *Group) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, g)
}

func (gm *GroupMember) BeforeCreate(tx *gorm.DB) error {
	return GenerateUUID(tx, gm)
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

// Timestampable interface for models that need timestamp updates
type Timestampable interface {
	SetUpdatedAt(time.Time)
}

// UpdateTimestamp updates the UpdatedAt field for any model that implements Timestampable
func UpdateTimestamp(tx *gorm.DB, model Timestampable) error {
	model.SetUpdatedAt(time.Now())
	return nil
}

// Implement Timestampable for all models that need timestamp updates
func (u *User) SetUpdatedAt(t time.Time) { u.UpdatedAt = t }
func (o *Organization) SetUpdatedAt(t time.Time) { o.UpdatedAt = t }
func (om *OrganizationMember) SetUpdatedAt(t time.Time) { om.UpdatedAt = t }
func (g *Group) SetUpdatedAt(t time.Time) { g.UpdatedAt = t }
func (gm *GroupMember) SetUpdatedAt(t time.Time) { gm.UpdatedAt = t }
func (f *Folder) SetUpdatedAt(t time.Time) { f.UpdatedAt = t }
func (fg *FolderGrant) SetUpdatedAt(t time.Time) { fg.UpdatedAt = t }
func (d *Document) SetUpdatedAt(t time.Time) { d.UpdatedAt = t }
func (dg *DocumentGrant) SetUpdatedAt(t time.Time) { dg.UpdatedAt = t }
func (a *Acronym) SetUpdatedAt(t time.Time) { a.UpdatedAt = t }
func (c *Category) SetUpdatedAt(t time.Time) { c.UpdatedAt = t }
func (l *Label) SetUpdatedAt(t time.Time) { l.UpdatedAt = t }
func (ac *AcronymCategory) SetUpdatedAt(t time.Time) { ac.UpdatedAt = t }

// BeforeUpdate hooks using the generic UpdateTimestamp function
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, u)
}

func (o *Organization) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, o)
}

func (om *OrganizationMember) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, om)
}

func (g *Group) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, g)
}

func (gm *GroupMember) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, gm)
}

func (f *Folder) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, f)
}

func (fg *FolderGrant) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, fg)
}

func (d *Document) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, d)
}

func (dg *DocumentGrant) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, dg)
}

func (a *Acronym) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, a)
}

func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, c)
}

func (l *Label) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, l)
}

func (ac *AcronymCategory) BeforeUpdate(tx *gorm.DB) error {
	return UpdateTimestamp(tx, ac)
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
