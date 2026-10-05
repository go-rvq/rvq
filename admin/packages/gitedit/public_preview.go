package gitedit

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ActionPublicPreview is making the preview of the own draft seen by anyone
// who has its link, with no login: a public preview (PublicPreview).
const ActionPublicPreview = "PublicPreview"

// ActionPublicPreviewOff is ending a public preview: its owner's
// (!public_preview), any (!drafts).
const ActionPublicPreviewOff = "PublicPreviewOff"

// PublicPreview is the preview of a draft seen with no login, by its link —
// PublicPreviewPath, then Token —, until it expires or is ended. Ended, it is
// kept: who opened it, and when, is known; a link is never served again.
type PublicPreview struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	// Owner is the key of the owner of the draft.
	Owner string `gorm:"column:owner_key;index;size:128"`
	// Token is what its link carries: random, not guessed.
	Token string `gorm:"uniqueIndex;size:64"`
	// CreatedBy is the key of who opened it.
	CreatedBy string `gorm:"size:128"`
	// ExpiresAt, when set, is when it stops.
	ExpiresAt *time.Time
	// RevokedAt is when it was ended, and RevokedBy by whom.
	RevokedAt *time.Time
	RevokedBy string `gorm:"size:128"`
}

// TableName is the table of the public previews.
func (PublicPreview) TableName() string { return "gitedit_public_previews" }

// Active says whether it is served at now: not ended, not expired.
func (p *PublicPreview) Active(now time.Time) bool {
	return p.RevokedAt == nil && (p.ExpiresAt == nil || p.ExpiresAt.After(now))
}

// newToken is the token of a public preview: 32 random bytes.
func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// EnablePublicPreview opens the public preview of the draft of owner — by
// by —, until expires (nil: until ended). One each draft: the one open is
// ended, its link no more served.
func (b *Builder) EnablePublicPreview(ctx context.Context, owner, by string, expires *time.Time) (*PublicPreview, error) {
	if b.DB == nil {
		return nil, errors.New("gitedit: no database for the public previews")
	}
	token, err := newToken()
	if err != nil {
		return nil, err
	}
	p := &PublicPreview{ID: uuid.New(), Owner: owner, Token: token, CreatedBy: by, ExpiresAt: expires}
	err = b.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := endPublicPreview(tx, owner, by); err != nil {
			return err
		}
		return tx.Create(p).Error
	})
	return p, err
}

// DisablePublicPreview ends the public preview of the draft of owner, by by.
func (b *Builder) DisablePublicPreview(ctx context.Context, owner, by string) error {
	if b.DB == nil {
		return errors.New("gitedit: no database for the public previews")
	}
	return endPublicPreview(b.DB.WithContext(ctx), owner, by)
}

func endPublicPreview(db *gorm.DB, owner, by string) error {
	now := time.Now()
	return active(db.Model(&PublicPreview{}), now).Where("owner_key = ?", owner).
		Updates(map[string]any{"revoked_at": now, "revoked_by": by}).Error
}

// PublicPreviewOf is the public preview open of the draft of owner; nil when
// none.
func (b *Builder) PublicPreviewOf(ctx context.Context, owner string) *PublicPreview {
	if b.DB == nil {
		return nil
	}
	var p PublicPreview
	if active(b.DB.WithContext(ctx), time.Now()).Where("owner_key = ?", owner).
		Order("created_at DESC").First(&p).Error != nil {
		return nil
	}
	return &p
}

// publicPreview is the public preview open of token; nil when none.
func (b *Builder) publicPreview(ctx context.Context, token string) *PublicPreview {
	if b.DB == nil || token == "" {
		return nil
	}
	var p PublicPreview
	if active(b.DB.WithContext(ctx), time.Now()).Where("token = ?", token).First(&p).Error != nil {
		return nil
	}
	return &p
}

// publicPreviewURL is the link of p, of the site of r.
func (b *Builder) publicPreviewURL(r *http.Request, p *PublicPreview) string {
	return absURL(r, strings.TrimSuffix(b.PublicPreviewPath, "/")+"/"+p.Token)
}

// PublicPreviewHandler serves the public previews — no login: the token is
// the access —, mounted by the application at PublicPreviewPath, out of the
// admin: PublicPreviewPath/<token>/<the site's path>. A token not open is
// not found, as one never made.
func (b *Builder) PublicPreviewHandler() http.Handler {
	base := strings.TrimSuffix(b.PublicPreviewPath, "/")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rest, ok := strings.CutPrefix(r.URL.Path, base+"/")
		token, path, _ := strings.Cut(rest, "/")
		p := b.publicPreview(r.Context(), token)
		if !ok || p == nil || b.Preview == nil || b.Repo == nil {
			http.NotFound(w, r)
			return
		}
		d, err := b.Repo.Draft(r.Context(), p.Owner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// a draft, not the site: not indexed; its link not told to the
		// sites it links to
		w.Header().Set("X-Robots-Tag", "noindex, nofollow")
		w.Header().Set("Referrer-Policy", "no-referrer")
		prefix := base + "/" + token
		r2 := r.Clone(r.Context())
		r2.URL.Path, r2.URL.RawPath = "/"+path, ""
		b.Preview(w, r2, d, prefix)
	})
}

// PublicPreviewForm is the form of a public preview: until when ("" until
// ended; a date, "2026-12-31", to its end).
type PublicPreviewForm struct {
	ExpiresAt string
}

// setupPublicPreviewActions are the actions of the public preview of the
// own draft: opening it (!public_preview), until when; ending it (its owner,
// or !drafts).
func (b *Builder) setupPublicPreviewActions(p *presets.Builder, pg *presets.PageBuilder,
	action func(name, asks string) *presets.ActionBuilder, label func(f func(m *Messages) string) func(ctx context.Context) string) {
	open := action(ActionPublicPreview, ActionPublicPreview).Icon("mdi-earth").
		SetI18nLabel(label(func(m *Messages) string { return m.PublicPreviewAction }))
	form := presets.NewModelBuilder(p, &PublicPreviewForm{}, presets.ModelConfig().SetModuleKey(MessagesKey))
	ed := form.Editing("ExpiresAt")
	ed.Field("ExpiresAt").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		m := GetMessages(ctx.Context())
		return h.Div(
			h.P(h.Text(m.PublicPreviewWarn)).Class("text-body-2 mb-4"),
			v.VTextField().Type("date").Label(field.Label).Hint(m.PublicPreviewExpiresHint).PersistentHint(true).
				Variant(v.FieldVariantOutlined).Density(v.DensityCompact).
				Attr(web.VField(field.FormKey, "")...).ErrorMessages(field.Errors...),
		)
	})
	presets.ActionForm[*PublicPreviewForm](open, ed, func(c *presets.ActionFormContext[*PublicPreviewForm]) error {
		r := c.Context.R
		m := GetMessages(r.Context())
		owner, own := b.ownerKey(r)
		if !own || !b.allowedAction(r, ActionPreview) {
			return perm.PermissionDenied
		}
		var expires *time.Time
		if c.Form.ExpiresAt != "" {
			day, err := time.ParseInLocation("2006-01-02", c.Form.ExpiresAt, time.Local)
			if err != nil {
				return web.NewValidationErrors().FieldError("ExpiresAt", err.Error())
			}
			end := day.AddDate(0, 0, 1) // to the end of the day
			if !end.After(time.Now()) {
				return web.NewValidationErrors().FieldError("ExpiresAt", m.ErrShareExpired)
			}
			expires = &end
		}
		if _, err := b.EnablePublicPreview(r.Context(), owner, owner, expires); err != nil {
			return err
		}
		c.Context.Flash = m.PublicPreviewOpened
		return nil
	}).Build()

	// ending: the owner's (!public_preview), any (!drafts)
	pg.Action(ActionPublicPreviewOff).Icon("mdi-earth-off").
		SetVerifier(func(ctx *web.EventContext) *perm.Verifier {
			return b.page.Page().ActionVerifier(ctx.R, presets.PermGet)
		}).
		SetI18nLabel(label(func(m *Messages) string { return m.PublicPreviewOffAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			r := ctx.R
			owner, own := b.ownerKey(r)
			if !(own && b.allowedAction(r, ActionPublicPreview)) && !b.allowedAction(r, ActionDrafts) {
				return perm.PermissionDenied
			}
			actor, _ := b.Identity(r)
			if err := b.DisablePublicPreview(r.Context(), owner, actor); err != nil {
				return err
			}
			ctx.Flash = GetMessages(r.Context()).PublicPreviewEnded
			return nil
		})
}

// publicPreviewView is, on the tab of the draft of owner, its public
// preview: its link — to copy —, until when, and ending it; or opening it.
func (b *Builder) publicPreviewView(ctx *web.EventContext, m *Messages, portal, owner string, own bool) h.HTMLComponent {
	if b.PublicPreviewPath == "" || b.Preview == nil {
		return nil
	}
	r := ctx.R
	canOpen := own && b.allowedAction(r, ActionPublicPreview) && b.allowedAction(r, ActionPreview)
	canEnd := (own && b.allowedAction(r, ActionPublicPreview)) || b.allowedAction(r, ActionDrafts)
	p := b.PublicPreviewOf(r.Context(), owner)
	if p == nil {
		if !canOpen {
			return nil
		}
		return h.Div(
			h.H4(m.PublicPreview).Class("mt-4 mb-1"),
			h.P(h.Text(m.PublicPreviewNone)).Class("text-body-2 text-medium-emphasis mb-2"),
			v.VBtn(m.PublicPreviewAction).PrependIcon("mdi-earth").Variant(v.VariantTonal).Size(v.SizeSmall).
				Class("mb-2").Attr("data-public-preview-open", true).
				Attr("@click", b.actionOnClick(ctx, ActionPublicPreview, portal)),
		).Attr("data-public-preview", "off")
	}
	link := b.publicPreviewURL(r, p)
	until := m.ShareForever
	if p.ExpiresAt != nil {
		until = fmt.Sprintf(m.ShareUntil, p.ExpiresAt.Local().Format(dateFormat))
	}
	var end h.HTMLComponent
	if canEnd {
		end = v.VBtn(m.PublicPreviewOffAction).PrependIcon("mdi-earth-off").Variant(v.VariantText).Size(v.SizeSmall).
			Color("error").Attr("data-public-preview-off", true).
			Attr("@click", b.actionOnClick(ctx, ActionPublicPreviewOff, portal))
	}
	return h.Div(
		h.H4(m.PublicPreview).Class("mt-4 mb-1"),
		h.P(h.Text(m.PublicPreviewOn+" · "+until)).Class("text-body-2 mb-2"),
		v.VTextField().Label(m.PublicPreviewLink).ModelValue(link).Readonly(true).
			Variant(v.FieldVariantOutlined).Density(v.DensityCompact).HideDetails(true).Class("mb-2").
			Attr("append-inner-icon", "mdi-content-copy").Attr("data-public-preview-url", link).
			Attr("@click:append-inner", "copyToClipboard("+h.JSONString(link)+")"),
		v.VBtn(m.PublicPreviewVisit).PrependIcon("mdi-open-in-new").Variant(v.VariantTonal).Size(v.SizeSmall).
			Class("me-2").Href(link).Attr("target", "_blank"),
		end,
	).Attr("data-public-preview", "on")
}
