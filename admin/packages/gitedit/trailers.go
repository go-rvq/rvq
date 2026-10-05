package gitedit

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	fs_tools "github.com/go-rvq/rvq/admin/packages/fs-tools"
)

// Every commit says who made it on which site — the panel's and the ones
// pushed by git alike —: a trailer, as Co-authored-by is,
//
//	Site-User: alice <8b2f…@example.com>
//
// the user's login, and their primary key at the site's host. A commit of
// the panel also says so (Committed-Via) and names, as co-authors, the other
// users who changed its files (Co-authored-by). A push of commits that do not
// say who made them is refused: the hook commit-msg of the page adds it.

// The trailers of the commits.
const (
	TrailerSiteUser = "Site-User"
	TrailerVia      = "Committed-Via"
	TrailerCoAuthor = "Co-authored-by"
)

// Signature is the Site-User of the user u on host: "alice <8b2f…@example.com>".
func Signature(u *User, host string) string {
	return u.Login + " <" + u.Key + "@" + host + ">"
}

// signatureRe is a Site-User: the login, the key and the host.
var signatureRe = regexp.MustCompile(`^(.+?) <([^@<>\s]+)@([^<>\s]+)>$`)

// ParseSignature is the login, the key and the host of a Site-User.
func ParseSignature(v string) (login, key, host string, ok bool) {
	m := signatureRe.FindStringSubmatch(strings.TrimSpace(v))
	if m == nil {
		return "", "", "", false
	}
	return m[1], m[2], m[3], true
}

// panelMessage is the message of a commit of the panel, by the user of r in
// the draft d: message, then its trailers — the other users who changed its
// files (Co-authored-by), who made it (Site-User), where (Committed-Via).
func (b *Builder) panelMessage(r *http.Request, d *Draft, message string) string {
	actor := b.actor(r)
	host := b.host(r)
	var trailers []string
	for _, key := range d.Editors() {
		if key == actor.Key {
			continue
		}
		u := b.user(r, key)
		if u.Author.Email == "" {
			continue // not known: no one to name
		}
		trailers = append(trailers, TrailerCoAuthor+": "+u.Author.Name+" <"+u.Author.Email+">")
	}
	trailers = append(trailers,
		TrailerSiteUser+": "+Signature(actor, host),
		TrailerVia+": admin panel ("+host+")")
	return strings.TrimRight(message, "\n ") + "\n\n" + strings.Join(trailers, "\n") + "\n"
}

// checkSignatures refuses a push of a commit that does not say who made it on
// this site: each commit new to the repository needs a Site-User of its host,
// of a user the site knows, their login with their key.
func (b *Builder) checkSignatures(ctx context.Context, rc *fs_tools.GitReceive) error {
	host := b.host(rc.Request)
	for _, u := range rc.Updates {
		if u.New == fs_tools.ZeroID {
			continue
		}
		out, err := rc.Git(ctx, "rev-list", u.New, "--not", "--all")
		if err != nil {
			return err
		}
		for _, sha := range strings.Fields(out) {
			values, err := rc.Git(ctx, "log", "-1", "--format=%(trailers:key="+TrailerSiteUser+",valueonly)", sha)
			if err != nil {
				return err
			}
			if !b.signedHere(rc.Request, values, host) {
				return b.unsignedError(rc.Request, sha, host)
			}
		}
	}
	return nil
}

// signedHere says whether one of values (Site-User, a line each) is of host
// and of a user the site knows.
func (b *Builder) signedHere(r *http.Request, values, host string) bool {
	for _, v := range strings.Split(values, "\n") {
		login, key, h, ok := ParseSignature(v)
		if !ok || !strings.EqualFold(h, host) {
			continue
		}
		if b.Users == nil {
			return true
		}
		if u, err := b.Users.User(r, key); err == nil && u != nil && u.Login == login {
			return true
		}
	}
	return false
}

// unsignedError says what a commit pushed lacks, and how to add it.
func (b *Builder) unsignedError(r *http.Request, sha, host string) error {
	short := sha
	if len(short) > 10 {
		short = short[:10]
	}
	line := TrailerSiteUser + ": " + Signature(b.actor(r), host)
	return fmt.Errorf("commit %s does not say who made it on %s: end its message with\n\n    %s\n\n"+
		"the hook commit-msg adds it to each commit:\n\n    curl -fsSL -u <login> %s -o .git/hooks/commit-msg && chmod +x .git/hooks/commit-msg\n\n"+
		"then sign the commits already made: git rebase -x 'git commit --amend --no-edit' <the commit before them>",
		short, host, line, absURL(r, b.hookPath()))
}

// hookPath is the hook commit-msg of the user of the request.
func (b *Builder) hookPath() string { return b.page.Page().FullPath() + "/commit-msg" }

// serveHook is the hook commit-msg of the user of r: it ends each message
// with their Site-User, unless it has one.
func (b *Builder) serveHook(w http.ResponseWriter, r *http.Request) {
	if !b.allowedAction(r, ActionGit) {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}
	line := TrailerSiteUser + ": " + Signature(b.actor(r), b.host(r))
	w.Header().Set("Content-Type", "text/x-shellscript; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="commit-msg"`)
	fmt.Fprintf(w, "#!/bin/sh\n# commit-msg of %s: each commit says who made it on the site (%s), as a push to it asks.\n"+
		"git interpret-trailers --in-place --if-exists doNothing --trailer %s \"$1\"\n",
		b.host(r), TrailerSiteUser, shellQuote(line))
}

// shellQuote is s quoted for sh.
func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

// absURL is path on the host of r.
func absURL(r *http.Request, path string) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + r.Host + path
}
