package gitedit

import (
	"net"
	"net/http"
	"strings"
	"sync"

	fs_tools "github.com/go-rvq/rvq/admin/packages/fs-tools"
	"github.com/go-rvq/rvq/x/perm"
)

// A draft is its owner's, reached by others it is shared with (DraftShare)
// and by whoever may see every draft (!drafts): the page of the user logged
// in is theirs (…/site-files), the page of another's under its key
// (…/site-files/u/<key>) — its editor, its IDE, its preview, its git —, each
// request asking whether it may.

// ActionShare is sharing one's draft with other users, and revoking it.
const ActionShare = "Share"

// ActionDrafts is seeing and reaching every draft — the IDE, git, the
// preview of each user's —, and revoking any sharing.
const ActionDrafts = "Drafts"

// userParam is the wildcard of the key of the owner of the draft in the URLs
// of another's (…/u/{user}).
const userParam = "user"

// User is a user of the editor, as the application knows them.
type User struct {
	// Key is their primary key: their draft's, their URL's (…/u/<key>).
	Key string
	// Login is their username.
	Login string
	// Name is their full name.
	Name string
	// Author is how their commits are signed (name and e-mail).
	Author Author
}

// Label is how the user is shown: the name and the login.
func (u *User) Label() string {
	if u.Name == "" || u.Name == u.Login {
		return u.Login
	}
	return u.Name + " (" + u.Login + ")"
}

// Users are the users the editor asks the application of, for a request (its
// host: the e-mail of a user with none, say): one by key, and searched (whom
// to share a draft with).
type Users interface {
	User(r *http.Request, key string) (*User, error)
	Search(r *http.Request, q string, limit int) ([]*User, error)
}

// ErrNoAccess refuses a draft the request may not reach.
var ErrNoAccess = perm.PermissionDenied

// ownerKey is the key of the owner of the draft r reaches: the user of its
// URL (…/u/{user}, or the key of a repository by git of each), or the one
// logged in; own when it is theirs.
func (b *Builder) ownerKey(r *http.Request) (key string, own bool) {
	actor, _ := b.Identity(r)
	k := r.PathValue(userParam)
	if k == "" {
		k = fs_tools.GitRepoKey(r)
	}
	if k == "" || k == actor {
		return actor, true
	}
	return k, false
}

// actor is the user of the request.
func (b *Builder) actor(r *http.Request) *User {
	key, author := b.Identity(r)
	u := b.user(r, key)
	u.Author = author // as the request signs (its credentials of git)
	return u
}

// user is the user of key, as the application knows them (Users); their key
// alone when it does not.
func (b *Builder) user(r *http.Request, key string) *User {
	if b.Users != nil {
		if u, err := b.Users.User(r, key); err == nil && u != nil {
			return u
		}
	}
	return &User{Key: key, Login: key, Name: key}
}

// mayReach says whether r may reach the draft of owner: their own; one
// shared with them; any, with !drafts — a user the application knows.
func (b *Builder) mayReach(r *http.Request, owner string) bool {
	actor, _ := b.Identity(r)
	if owner == actor {
		return true
	}
	if b.Users != nil {
		if u, err := b.Users.User(r, owner); err != nil || u == nil {
			return false
		}
	}
	if b.allowedAction(r, ActionDrafts) {
		return true
	}
	return b.sharedWith(r.Context(), owner, actor)
}

// isOwnerOrAdmin says whether r is the owner of the draft it reaches, or may
// see every draft: who resets it, revokes its sharing.
func (b *Builder) isOwnerOrAdmin(r *http.Request) bool {
	_, own := b.ownerKey(r)
	return own || b.allowedAction(r, ActionDrafts)
}

// lockDraft holds the draft of key while git changes it — a commit, an
// update, a publishing, a push, a reset —: two users on one draft, one at a
// time. Its unlock is returned.
func (b *Builder) lockDraft(key string) func() {
	v, _ := b.locks.LoadOrStore(key, &sync.Mutex{})
	m := v.(*sync.Mutex)
	m.Lock()
	return m.Unlock
}

// host is the site's, of r: what the commits say they were made on.
func (b *Builder) host(r *http.Request) string {
	if b.Host != nil {
		return b.Host(r)
	}
	h := r.Host
	if hh, _, err := net.SplitHostPort(h); err == nil {
		h = hh
	}
	return h
}

// scopeURL is the page of the draft r reaches: the own page, or the
// owner's (…/u/<key>).
func (b *Builder) scopeURL(r *http.Request) string {
	if key, own := b.ownerKey(r); !own {
		return b.UserURL(key)
	}
	return b.page.Page().FullPath()
}

// UserURL is the page of the draft of the user of key.
func (b *Builder) UserURL(key string) string {
	return strings.Replace(b.userPage.Page().FullPath(), "{"+userParam+"}", key, 1)
}

// previewURL is the site as the draft r reaches makes it.
func (b *Builder) previewURL(r *http.Request) string {
	if b.PreviewURL == "" {
		return ""
	}
	if key, own := b.ownerKey(r); !own && b.previewPrefix != "" {
		return b.previewPrefix + "/u/" + key
	}
	return b.PreviewURL
}
