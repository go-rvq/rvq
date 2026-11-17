package backup

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	fs_tools "github.com/go-rvq/rvq/admin/packages/fs-tools"
	"github.com/go-rvq/rvq/admin/packages/fs-tools/explorer"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/packages/backup"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LinkTTL is how long a download link works.
const LinkTTL = 24 * time.Hour

// BackupLink is a link that downloads a state whole as tar.gz, with no
// login, for LinkTTL: made by the action of the detailing (archive_link),
// shown once — its token is kept only as its sum —, used as many times as
// wanted (a tar.gz made as it is sent cannot be resumed: a download that
// fell starts again by it), each use counted; revoked, expired, or its
// state removed, it is gone.
type BackupLink struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt     time.Time
	StateID       string `gorm:"size:64;index;not null"`
	TokenSum      string `gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt     time.Time
	RevokedAt     *time.Time
	CreatedByID   model.UUID `gorm:"type:uuid"`
	CreatedByName string
	Uses          int
	LastUsedAt    *time.Time
	LastIP        string
	LastUserAgent string `gorm:"type:text"`
}

// TableName is backup_links.
func (BackupLink) TableName() string { return "backup_links" }

// Active says the link works at now.
func (l *BackupLink) Active(now time.Time) bool { return l.RevokedAt == nil && now.Before(l.ExpiresAt) }

func tokenSum(tok string) string {
	s := sha256.Sum256([]byte(tok))
	return hex.EncodeToString(s[:])
}

// NewLink makes a download link of the state id (ready), of the user of r:
// its token (shown once) and the link.
func (b *Builder) NewLink(r *http.Request, id string) (string, *BackupLink, error) {
	var s BackupState
	if err := b.db.First(&s, "id = ?", id).Error; err != nil {
		return "", nil, err
	}
	if s.Status != StatusReady {
		return "", nil, fmt.Errorf("backup: the state %s is %s", id, s.Status)
	}
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, err
	}
	tok := base64.RawURLEncoding.EncodeToString(raw)
	l := &BackupLink{ID: uuid.New(), StateID: id, TokenSum: tokenSum(tok), ExpiresAt: time.Now().Add(LinkTTL)}
	if uid, name := userOf(login.CurrentUserFromRequest(r)); name != "" {
		l.CreatedByID, l.CreatedByName = model.UUID(uid), name
	}
	if err := b.db.Create(l).Error; err != nil {
		return "", nil, err
	}
	return tok, l, nil
}

// LinkPath is the path of the download of the state id by the link tok.
func (b *Builder) LinkPath(id, tok string) string {
	return b.ClientPath() + "/" + id + "/archive/" + tok + ".tar.gz"
}

// serveArchive serves the state of a link whole as tar.gz: the link
// active, its state ready; else not found (the link does not say whether it
// was).
func (b *Builder) serveArchive(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tok, ok := strings.CutSuffix(r.PathValue("file"), ".tar.gz")
	now := time.Now()
	var l BackupLink
	if !ok || b.db.First(&l, "token_sum = ? AND state_id = ?", tokenSum(tok), id).Error != nil || !l.Active(now) {
		http.NotFound(w, r)
		return
	}
	var s BackupState
	if b.db.First(&s, "id = ?", id).Error != nil || s.Status != StatusReady {
		http.NotFound(w, r)
		return
	}
	_ = b.db.Model(&BackupLink{}).Where("id = ?", l.ID).Updates(map[string]any{
		"uses": gorm.Expr("uses + 1"), "last_used_at": now,
		"last_ip": hostOf(login.ClientAddr(r)), "last_user_agent": r.UserAgent()}).Error
	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", "backup-"+id+".tar.gz"))
	w.Header().Set("Cache-Control", "no-store")
	if _, err := b.core.WriteArchive(w, id, backup.ArchiveOptions{}); err != nil {
		b.log.Warn("archive of a link", "state", id, "err", err)
	}
}

// fsHandler is the explorer of the files of a state (the page Files),
// read only: whoever may see the state reads it all — checked once, of the
// state, not of each file.
func (b *Builder) fsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		mid, err := b.mb.ParseRecordID(id)
		if err != nil || b.mb.Info().Permissioner().Reader(r, mid).Denied() {
			http.NotFound(w, r)
			return
		}
		var s BackupState
		if b.db.First(&s, "id = ?", id).Error != nil || s.Status != StatusReady {
			explorer.Fail(w, http.StatusConflict, fmt.Errorf("the state is not ready"))
			return
		}
		sfs, err := b.stateFSOf(id)
		if err != nil {
			explorer.Fail(w, http.StatusNotFound, err)
			return
		}
		x := explorer.New(sfs, explorer.Options{
			Allow: func(_ *http.Request, _ string, write bool) bool { return !write },
			Archive: func(w io.Writer, _ *http.Request, base string, paths []string) error {
				_, err := b.core.WriteArchive(w, id, archiveOptions(base, paths))
				return err
			},
		})
		http.StripPrefix(b.ClientPath()+"/"+id+"/fs", x).ServeHTTP(w, r)
	})
}

// archiveOptions is what the explorer asks of an archive (its paths of the
// file system, "/public/a.txt"; the folder they are in) as the core's.
func archiveOptions(base string, paths []string) backup.ArchiveOptions {
	return backup.ArchiveOptions{Base: base, Paths: paths}
}

// FilesPath is the page of the files of a state: under its detailing.
const FilesPath = "/files"

// filesPage is the page Files of a state: its files in the explorer, read
// only (download a file, a folder, several).
func (b *Builder) filesPage(ctx *web.EventContext, obj any, _ model.ID, r *web.PageResponse) error {
	s := obj.(*BackupState)
	m := GetMessages(ctx.Context())
	if s.Status != StatusReady {
		r.Body = v.VAlert(h.Text(fill(m.NotReady, "status", statusText(m, s.Status)))).Type("warning").Variant(v.VariantTonal)
		return nil
	}
	r.Body = v.VContainer(
		h.Tag("vx-files").Attr("base", b.ClientPath()+"/"+s.ID+"/fs/").
			Attr("title", fill(m.FilesPageOf, "id", s.ID)).
			Attr(":fetch", "false").Attr(":extract", "false").
			Attr("start", ctx.R.URL.Query().Get("dir")),
	).Fluid(true)
	return nil
}

// configureLinks: the action of a link (the detailing of a ready state),
// the links of a state (its child Links: revoke), the page Files.
func (b *Builder) configureLinks(mb *presets.ModelBuilder) {
	fs_tools.UseAssets(b.p)

	// a new link: an action of the menu of a ready state, its dialog the link
	// made as it opens — nothing to run (no button), shown once: only its
	// sum is kept. No permission of its own (NoPerm): it asks the record's,
	// and it is there for who may download (PermDownload)
	mb.Detailing().Action("archive_link").Icon("mdi-link-variant").NoPerm().
		SetEnabledObj(func(obj any, id string, ctx *web.EventContext) (bool, error) {
			return obj.(*BackupState).Status == StatusReady && b.mayDownload(ctx.R), nil
		}).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			if !b.mayDownload(ctx.R) {
				return nil, perm.PermissionDenied
			}
			tok, l, err := b.NewLink(ctx.R, id)
			if err != nil {
				return nil, err
			}
			return newLinkMade(GetMessages(ctx.Context()), clientBase(ctx.R)+b.LinkPath(id, tok), l.ExpiresAt), nil
		})

	// the page of the files: under the detailing; its link, in the detailing
	mb.Detailing().AddPageFunc(nil, FilesPath, b.filesPage)

	// the links of a state: who made them, until when, their uses; revoked
	links := presets.NewModelBuilder(b.p, &BackupLink{}, presets.ModelConfig().SetId("backup_links").SetModuleKey(MessagesKey))
	links.Label("BackupLink").SetPluralLabel("BackupLinks").MenuIcon("mdi-link-variant")
	mb.AddChild(links)
	b.scopeToState(links)
	links.Listing("CreatedAt", "CreatedByName", "ExpiresAt", "RevokedAt", "Uses", "LastUsedAt", "LastIP").OrderBy("created_at DESC")
	links.Listing().Field("ExpiresAt").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		l := field.Obj.(*BackupLink)
		return h.Td(h.Text(l.ExpiresAt.Local().Format("2006-01-02 15:04")+" "), linkChip(GetMessages(ctx.Context()), l))
	})
	links.Detailing("CreatedAt", "CreatedByName", "ExpiresAt", "RevokedAt", "Uses", "LastUsedAt", "LastIP", "LastUserAgent")
	links.SetCreatingDisabled(true).SetEditingDisabled(true).SetDeletingDisabled(true)
	revoke := links.Detailing().Action("revoke").Icon("mdi-link-variant-off").
		SetEnabledObj(func(obj any, id string, ctx *web.EventContext) (bool, error) {
			return obj.(*BackupLink).Active(time.Now()), nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			now := time.Now()
			if err := b.db.Model(&BackupLink{}).Where("id = ? AND revoked_at IS NULL", id).Update("revoked_at", now).Error; err != nil {
				return err
			}
			ctx.Flash = GetMessages(ctx.Context()).LinkRevoked
			return nil
		})
	links.Listing().AppendItemAction(revoke)
}

// linkMade is a link just made, shown in its dialog: its address, copied or
// opened.
type linkMade struct{ h.HTMLComponent }

func newLinkMade(m *Messages, url string, expires time.Time) *linkMade {
	return &linkMade{h.Div(
		h.P(h.Text(fill(m.LinkText, "expires", expires.Local().Format("2006-01-02 15:04")))).Class("mb-3"),
		v.VTextField().Attr("model-value", url).Readonly(true).Variant(v.FieldVariantOutlined).
			Density(v.DensityCompact).HideDetails(true).Attr("data-backup-link", true),
		h.Div(
			v.VBtn(m.LinkCopy).PrependIcon("mdi-content-copy").Variant(v.VariantTonal).
				Attr("@click", fmt.Sprintf("navigator.clipboard.writeText(%s)", h.JSONString(url))),
			v.VBtn(m.LinkOpen).PrependIcon("mdi-open-in-new").Variant(v.VariantTonal).
				Href(url).Attr("target", "_blank").Attr("rel", "noopener noreferrer"),
		).Class("d-flex ga-2 mt-2"),
	)}
}

func linkChip(m *Messages, l *BackupLink) h.HTMLComponent {
	text, color := m.LinkStatusActive, "success"
	switch {
	case l.RevokedAt != nil:
		text, color = m.LinkStatusRevoked, "grey"
	case !l.Active(time.Now()):
		text, color = m.LinkStatusExpired, "grey"
	}
	return v.VChip(h.Text(text)).Size(v.SizeXSmall).Color(color).Variant(v.VariantTonal)
}
