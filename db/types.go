package db

import (
	"fmt"
	"time"
)

/* ################################ */
/* -------------------------------- */
/* -------- ACRONYM TYPES --------- */
/* -------------------------------- */
/* ################################ */

type Owner struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type JsonLabel struct {
	UUID  string `json:"uuid"`
	Label string `json:"label"`
}

// AcronymResponse represents the public view of an Acronym
type AcronymResponse struct {
	UUID      string      `json:"uuid"`
	Acronym   string      `json:"acronym"`
	Meaning   string      `json:"meaning"`
	Owner     Owner       `json:"owner"`
	Labels    []JsonLabel `json:"labels"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (a Acronym) ToJson() AcronymResponse {
	r := AcronymResponse{
		UUID:      a.UUID,
		Acronym:   a.Acronym,
		Meaning:   a.Meaning,
		Owner:     Owner{UUID: a.Owner.UUID, Name: a.Owner.Name},
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	r.Labels = make([]JsonLabel, len(a.Labels))
	for i, l := range a.Labels {
		r.Labels[i] = JsonLabel{UUID: l.UUID, Label: l.Label}
	}

	return r
}

func (a Acronym) Validate() error {
	// Validations
	// - Acronym cannot be empty
	// - Meaning cannot be empty
	// -

	if a.Acronym == "" {
		return fmt.Errorf("acronym cannot be empty")
	}
	if a.Meaning == "" {
		return fmt.Errorf("meaning cannot be empty")
	}
	if len(a.Labels) == 0 {
		return fmt.Errorf("labels cannot be empty")
	}
	return nil
}


/* ################################ */
/* -------------------------------- */
/* ---------- USER TYPES ---------- */
/* -------------------------------- */
/* ################################ */

type UserResponse struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	LegalName string    `json:"legal_name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Session   Session   `json:"session,omitempty"`
	Settings  []UserSettings `json:"settings,omitempty"`
}

func (u User) ToMinimalJson() UserResponse {
	return UserResponse{
		UUID:      u.UUID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

func (u User) ToJson(userIsMe bool) UserResponse {
	r := UserResponse{
		UUID:      u.UUID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}

	if userIsMe {
		r.LegalName = u.LegalName
		r.Session = u.Session
		r.Settings = u.Settings
	}

	return r
}

func (u User) ToPublicJson() UserResponse {
	return UserResponse{
		UUID:      u.UUID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
