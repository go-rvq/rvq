// Package sharedtest is what the tests of the sharing (shared) use: models
// of its own — a folder (Pasta, with an owner) and its documents
// (Documento, under it) — and users and groups in memory (Users); the mails
// go to the mail server of the tests (mail-sender/mailtest).
package sharedtest

import (
	"context"
	"regexp"
	"slices"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Pasta is a folder: its owner, its documents under it.
type Pasta struct {
	ID      string `gorm:"primaryKey;size:64"`
	Nome    string
	OwnerID string `gorm:"size:64"`
}

func (Pasta) TableName() string { return "sharedtest_pastas" }

// Documento is a document of a folder.
type Documento struct {
	ID      string `gorm:"primaryKey;size:64"`
	PastaID string `gorm:"size:64;index"`
	Titulo  string
}

func (Documento) TableName() string { return "sharedtest_documentos" }

// Migrate makes the tables of the models.
func Migrate(db *gorm.DB) error { return db.AutoMigrate(&Pasta{}, &Documento{}) }

// NewPasta is a folder of the owner.
func NewPasta(db *gorm.DB, nome, owner string) *Pasta {
	p := &Pasta{ID: uuid.NewString(), Nome: nome, OwnerID: owner}
	if err := db.Create(p).Error; err != nil {
		panic(err)
	}
	return p
}

// NewDocumento is a document of the folder.
func NewDocumento(db *gorm.DB, pasta *Pasta, titulo string) *Documento {
	d := &Documento{ID: uuid.NewString(), PastaID: pasta.ID, Titulo: titulo}
	if err := db.Create(d).Error; err != nil {
		panic(err)
	}
	return d
}

// User is a user of the tests.
type User struct {
	ID, Name, Email, Language string
	// Groups are the user's, in their order.
	Groups []string
}

// Users are the users and the groups of the tests (shared.Users).
type Users struct {
	mu     sync.Mutex
	users  map[string]*User
	groups map[string]string
}

// NewUsers are no users, no groups.
func NewUsers() *Users { return &Users{users: map[string]*User{}, groups: map[string]string{}} }

// Add adds the user.
func (u *Users) Add(x *User) *User {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.users[x.ID] = x
	return x
}

// AddGroup adds the group id of the name.
func (u *Users) AddGroup(id, name string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.groups[id] = name
}

// Join puts the user in the group, last of theirs; Leave takes them out.
func (u *Users) Join(user, group string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.users[user].Groups = append(u.users[user].Groups, group)
}

func (u *Users) Leave(user, group string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	x := u.users[user]
	x.Groups = slices.DeleteFunc(x.Groups, func(g string) bool { return g == group })
}

func (u *Users) get(id string) *User {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.users[id]
}

func (u *Users) Exists(_ context.Context, id string) bool { return u.get(id) != nil }

func (u *Users) Name(_ context.Context, id string) string {
	if x := u.get(id); x != nil {
		return x.Name
	}
	return id
}

func (u *Users) Email(_ context.Context, id string) string {
	if x := u.get(id); x != nil {
		return x.Email
	}
	return ""
}

func (u *Users) Language(_ context.Context, id string) string {
	if x := u.get(id); x != nil {
		return x.Language
	}
	return ""
}

func (u *Users) Groups(_ context.Context, id string) []string {
	if x := u.get(id); x != nil {
		return slices.Clone(x.Groups)
	}
	return nil
}

func (u *Users) GroupExists(_ context.Context, id string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	_, ok := u.groups[id]
	return ok
}

func (u *Users) GroupName(_ context.Context, id string) string {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.groups[id]
}

// ConfirmURL is the link of the confirmations of the tests.
func ConfirmURL(token string) string { return "https://app.test/shared/confirm?t=" + token }

var tokenRe = regexp.MustCompile(`shared/confirm\?t=([A-Za-z0-9_-]+)`)

// Token is the token of the link of a mail of body; "" when none.
func Token(body string) string {
	if t := tokenRe.FindStringSubmatch(body); t != nil {
		return t[1]
	}
	return ""
}
