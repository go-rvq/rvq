package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/login"
	"gorm.io/gorm"
)

const (
	RoleAdmin   = "Admin"
	RoleManager = "Manager"
	RoleEditor  = "Editor"
	RoleViewer  = "Viewer"

	OAuthProviderGoogle          = "google"
	OAuthProviderMicrosoftOnline = "microsoftonline"
	OAuthProviderGithub          = "github"
)

var DefaultRoles = []string{
	RoleAdmin,
	RoleManager,
	RoleEditor,
	RoleViewer,
}

var OAuthProviders = []string{
	OAuthProviderGoogle,
	OAuthProviderMicrosoftOnline,
	OAuthProviderGithub,
}

type User struct {
	// Users use a UUID primary key (not the shared uint gorm.Model), so user
	// references across the system (activity, notes, sessions) are UUIDs,
	// matching the rvq user identity.
	ID        uuid.UUID      `admin:"-" gorm:"type:uuid;primaryKey"`
	DeletedAt gorm.DeletedAt `sql:"index"`

	Name             string
	Company          string
	Roles            []role.Role `gorm:"many2many:user_role_join;"`
	Status           string
	UpdatedAt        time.Time
	CreatedAt        time.Time
	FavorPostID      uint
	RegistrationDate time.Time `gorm:"type:date"`

	// Username is email
	login.UserPass
	login.OAuthInfo
	login.SessionSecure
}

func (u User) GetName() string {
	return u.Name
}

// BeforeCreate assigns a random UUID when none was provided.
func (u *User) BeforeCreate(*gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u User) GetID() uuid.UUID {
	return u.ID
}

func (u User) GetRoles() (rs []string) {
	for _, r := range u.Roles {
		rs = append(rs, r.Name)
	}
	if len(rs) == 0 {
		rs = []string{RoleViewer}
	}
	return
}

func (u User) IsOAuthUser() bool {
	return u.OAuthProvider != "" && u.OAuthIdentifier != ""
}
