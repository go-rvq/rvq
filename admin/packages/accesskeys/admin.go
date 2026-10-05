package accesskeys

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"github.com/sunfmin/reflectutils"
	"gorm.io/gorm"
)

// The ids of the models.
const (
	// UserKeysModelID is the keys of a user, under the users.
	UserKeysModelID = "access_keys"
	// MyKeysModelID is the keys of the user of the request: "My keys".
	MyKeysModelID = "my_access_keys"
	// UsesModelID is the history of a key, under it.
	UsesModelID = "access_key_uses"
)

// Builder is the access keys in the admin: of each user (under the users),
// the user's own (My keys), each with its history; and their authenticator.
type Builder struct {
	db   *gorm.DB
	auth *Authenticator
	// PermissionsHelpURL is the documentation of the permissions, opened by
	// the ? of a key's permissions ("" none).
	PermissionsHelpURL func(ctx *web.EventContext) string
	// DefaultValidity is how long a key made without an expiration lasts;
	// MaxValidity the most a key may.
	DefaultValidity, MaxValidity time.Duration
	// Profile, when set, is the model of the profile of the user
	// (my_profile): My keys under it, as what is the user's own; in the menu
	// when not set.
	Profile *presets.ModelBuilder

	p *presets.Builder
}

// New is the access keys of db.
func New(db *gorm.DB) *Builder {
	b := &Builder{db: db, DefaultValidity: 90 * 24 * time.Hour, MaxValidity: 366 * 24 * time.Hour}
	b.auth = &Authenticator{DB: db}
	return b
}

// Authenticator is the authentication by the keys (login.KeyAuth).
func (b *Builder) Authenticator() *Authenticator { return b.auth }

// AutoMigrate makes the tables of the keys.
func (b *Builder) AutoMigrate() error {
	return b.db.AutoMigrate(&AccessKey{}, &AccessKeyUse{})
}

// currentUserID is the id of the user of r: its login's.
func currentUserID(r *http.Request) (uuid.UUID, bool) {
	u := login.CurrentUserFromRequest(r)
	if u == nil {
		return uuid.Nil, false
	}
	id, err := reflectutils.Get(u, "ID")
	if err != nil {
		return uuid.Nil, false
	}
	switch t := id.(type) {
	case uuid.UUID:
		return t, t != uuid.Nil
	case string:
		u, err := uuid.Parse(t)
		return u, err == nil
	}
	return uuid.Nil, false
}

// parentID is the id of the record the request is under (a user, a key).
func parentID(r *http.Request) (uuid.UUID, bool) {
	ids := presets.ParentsModelID(r)
	if len(ids) == 0 {
		return uuid.Nil, false
	}
	switch t := ids.Last().Value().(type) {
	case uuid.UUID:
		return t, true
	case string:
		u, err := uuid.Parse(t)
		return u, err == nil
	}
	return uuid.Nil, false
}

// Install registers the models: the keys under users (each user's), My keys
// (the user's own: any user logged in may), and the history under each.
// Requests by a key never reach them: a key does not make keys.
func (b *Builder) Install(p *presets.Builder, users *presets.ModelBuilder) error {
	b.p = p
	ConfigureMessages(p.I18n())
	if b.auth.FindUser == nil {
		return errors.New("accesskeys: the Authenticator's FindUser is not set")
	}

	// the id by the config: the config's Apply sets all its attributes
	config := func(id string) *presets.ModelBuilderConfig {
		return presets.ModelConfig().SetId(id).SetModuleKey(MessagesKey)
	}
	// the children made by NewModelBuilder: p.Model would serve them on
	// their own too
	userKeys := presets.NewModelBuilder(p, &AccessKey{}, config(UserKeysModelID))
	users.AddChild(userKeys)
	b.configureKeys(userKeys, parentID)

	var myKeys *presets.ModelBuilder
	if b.Profile != nil {
		myKeys = presets.NewModelBuilder(p, &AccessKey{}, config(MyKeysModelID))
	} else {
		myKeys = p.Model(&AccessKey{}, config(MyKeysModelID))
	}
	// its labels before it goes under the profile: the item of its menu
	myKeys.Label("MyAccessKey").SetPluralLabel("MyAccessKeys").MenuIcon("mdi-key-chain")
	if b.Profile != nil {
		b.Profile.AddChild(myKeys)
	}
	b.configureKeys(myKeys, currentUserID)

	for _, keys := range []*presets.ModelBuilder{userKeys, myKeys} {
		uses := presets.NewModelBuilder(p, &AccessKeyUse{}, config(UsesModelID+"_"+keys.MenuID()))
		keys.AddChild(uses)
		b.configureUses(uses)
	}

	// any user may keep their own keys; no key may reach any key
	p.GetPermission().CreatePolicies(
		perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On("admin:"+MyKeysModelID+":*"),
		perm.PolicyFor("key:*").WhoAre(perm.Denied).ToDo(perm.Anything).
			On("admin:*"+UserKeysModelID+"*", "admin:*"+MyKeysModelID+"*", "admin:*"+UsesModelID+"*"),
	)
	return nil
}

// keysOf scopes the reads of mb to the keys of the user owner gives — the
// user under whom the keys are, or the user of the request —, and stamps it
// on the keys made.
func (b *Builder) scope(mb *presets.ModelBuilder, column string, owner func(r *http.Request) (uuid.UUID, bool)) {
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).
			WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
				cb.Pre(func(state *gorm2op.CallbackState) error {
					id, ok := owner(state.Ctx.R)
					if !ok {
						// no owner, no record
						state.DB = state.DB.Where("1 = 0")
						return nil
					}
					state.DB = state.DB.Where(column+" = ?", id)
					return nil
				})
			})
	})
}

func (b *Builder) configureKeys(mb *presets.ModelBuilder, owner func(r *http.Request) (uuid.UUID, bool)) {
	b.scope(mb, "user_id", owner)

	mb.Listing("Name", "Prefix", "Enabled", "ExpiresAt", "LastUsedAt")
	mb.Detailing("Name", "Description", "Prefix", "Enabled", "ExpiresAt", "LastUsedAt", "LastUsedIP", "CreatedAt",
		"Permissions")
	ed := mb.Editing("Name", "Description", "ExpiresAt", "Enabled", "Permissions")

	policyModel := presets.NewModelBuilder(b.p, &perm.DefaultDBPolicy{}, presets.ModelConfig().SetModuleKey(MessagesKey))
	permFb := &policyModel.Editing("Actions", "Resources").FieldsBuilder
	ed.Field("Permissions").AutoNested(policyModel, permFb)
	// lists of text — after AutoNested, which would set them otherwise: the
	// form posts each item as its __value
	permFb.Field("Actions").AsSlice()
	permFb.Field("Resources").AsSlice()
	ed.Field("Permissions").WrapComponentFunc(func(old presets.FieldComponentFunc) presets.FieldComponentFunc {
		return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			comp := old(field, ctx)
			if b.PermissionsHelpURL == nil {
				return comp
			}
			href := b.PermissionsHelpURL(ctx)
			if href == "" {
				return comp
			}
			return h.Div(
				h.Div(v.VBtn("").Icon("mdi-help-circle-outline").Variant(v.VariantText).Size(v.SizeSmall).
					Attr("href", href).Attr("target", "_blank").Attr("rel", "noopener").
					Attr("title", GetMessages(ctx.Context()).PermissionsHelp).Attr("data-permissions-help", true),
				).Style("position: absolute; top: -8px; right: 0; z-index: 1"),
				comp,
			).Style("position: relative")
		}
	})

	// a new key starts enabled — the form opened, not posted back
	ed.Field("Enabled").WrapComponentFunc(func(old presets.FieldComponentFunc) presets.FieldComponentFunc {
		return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			if k := field.Obj.(*AccessKey); k.ID == uuid.Nil && !posted(ctx.R) {
				k.Enabled = true
			}
			return old(field, ctx)
		}
	})

	ed.FetchFunc(func(obj any, id model.ID, ctx *web.EventContext) error {
		return gorm2op.DataOperator(b.db.Preload("Permissions")).Fetch(obj, id, ctx)
	})
	mb.Detailing().FetchFunc(func(obj any, id model.ID, ctx *web.EventContext) error {
		return gorm2op.DataOperator(b.db.Preload("Permissions")).Fetch(obj, id, ctx)
	})
	// the permissions in the detail: each resource and its actions
	mb.Detailing().Field("Permissions").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		k := field.Obj.(*AccessKey)
		m := GetMessages(ctx.Context())
		var items []h.HTMLComponent
		for _, pol := range k.Permissions {
			for _, res := range pol.Resources {
				items = append(items, h.Li(h.Code(res), h.Text(" — "+strings.Join(pol.Actions, ", "))))
			}
		}
		body := h.HTMLComponent(h.P(h.Text(m.NoPermissions)).Class("text-medium-emphasis"))
		if len(items) > 0 {
			body = h.Ul(items...).Class("ps-4")
		}
		return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"), body).
			Class("mb-4").Attr("data-key-permissions", true)
	})

	ed.Validators.AppendFunc(func(obj any, _ presets.FieldModeStack, ctx *web.EventContext) (errs web.ValidationErrors) {
		k := obj.(*AccessKey)
		m := GetMessages(ctx.Context())
		if k.Name == "" {
			errs.FieldError("Name", m.ErrNameRequired)
		}
		now := time.Now()
		if k.ExpiresAt.IsZero() {
			k.ExpiresAt = now.Add(b.DefaultValidity)
		}
		if k.ExpiresAt.After(now.Add(b.MaxValidity)) {
			errs.FieldError("ExpiresAt", fmt.Sprintf(m.ErrTooLong, int(b.MaxValidity.Hours()/24)))
		}
		return
	})

	// made (CreateFunc) or changed (SaveFunc): one way
	save := func(obj any, ctx *web.EventContext, isNew bool) error {
		k := obj.(*AccessKey)
		owner, ok := owner(ctx.R)
		if !ok {
			return errors.New("accesskeys: no user to make the key of")
		}
		if isNew {
			if k.ID == uuid.Nil {
				k.ID = uuid.New()
			}
			k.UserID = owner
			k.Code, k.Prefix = NewCode()
			k.CodeHash = CodeHash(k.Code)
			if by, ok := currentUserID(ctx.R); ok {
				k.CreatedByID = &by
			}
		} else if k.UserID != owner {
			return errors.New("accesskeys: not a key of this user")
		}
		// its permissions: allows, of its subject
		for _, pol := range k.Permissions {
			pol.ID = uuid.Nil
			pol.ReferID = k.ID.String()
			pol.Subject = k.Subject()
			pol.Effect = perm.Allowed
		}
		err := b.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("refer_id = ?", k.ID.String()).Delete(&perm.DefaultDBPolicy{}).Error; err != nil {
				return err
			}
			if isNew {
				return tx.Session(&gorm.Session{FullSaveAssociations: true}).Create(k).Error
			}
			if err := tx.Model(k).Select("Name", "Description", "ExpiresAt", "Enabled").Updates(k).Error; err != nil {
				return err
			}
			for _, pol := range k.Permissions {
				if err := tx.Create(pol).Error; err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		b.reloadPolicies()
		if isNew {
			// the event's flash: what its response shows
			ctx.Flash = codeMessage(ctx.Context(), k)
		}
		return nil
	}
	ed.CreateFunc(func(obj any, ctx *web.EventContext) error { return save(obj, ctx, true) })
	ed.SaveFunc(func(obj any, _ model.ID, ctx *web.EventContext) error { return save(obj, ctx, false) })

	mb.Listing().DeleteFunc(func(obj any, id model.ID, cascade bool, ctx *web.EventContext) error {
		owner, ok := owner(ctx.R)
		if !ok {
			return errors.New("accesskeys: no user")
		}
		err := b.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("refer_id = ?", id.String()).Delete(&perm.DefaultDBPolicy{}).Error; err != nil {
				return err
			}
			return tx.Where("id = ? AND user_id = ?", id.String(), owner).Delete(&AccessKey{}).Error
		})
		if err == nil {
			b.reloadPolicies()
		}
		return err
	})
}

// reloadPolicies brings the policies of the keys to the memory at once.
func (b *Builder) reloadPolicies() {
	if pb := b.p.GetPermission(); pb != nil {
		startFrom := time.Now().Add(-time.Second)
		pb.LoadDBPoliciesToMemory(b.db, &startFrom)
	}
}

// codeMessage is the message of a key just made: its code, shown once, to
// copy.
func codeMessage(ctx context.Context, k *AccessKey) *presets.FlashMessage {
	m := GetMessages(ctx)
	return &presets.FlashMessage{
		Text: fmt.Sprintf(m.CodeCreated, k.Code),
		// long: the code is copied before it goes (the "created" message
		// goes first)
		Duration: 600,
		Detail: h.Div(
			h.P(h.Text(m.CodeOnce)).Class("mb-3"),
			v.VTextField().Label(m.Code).ModelValue(k.Code).Readonly(true).Variant(v.FieldVariantOutlined).
				Density(v.DensityCompact).HideDetails(true).Attr("data-access-key-code", true).
				Attr("append-inner-icon", "mdi-content-copy").
				Attr("@click:append-inner", "copyToClipboard("+h.JSONString(k.Code)+")"),
		),
	}
}

func (b *Builder) configureUses(mb *presets.ModelBuilder) {
	b.scope(mb, "key_id", parentID)
	mb.Listing("CreatedAt", "Kind", "Method", "Path", "Status", "IP", "Place", "UserAgent").
		OrderBy("created_at DESC")
	mb.Detailing("CreatedAt", "Kind", "Method", "Path", "Status", "IP", "Place", "UserAgent")
	// the history is kept, not made nor changed
	mb.SetCreatingDisabled(true).SetEditingDisabled(true).SetDeletingDisabled(true)
	mb.Editing().SaveFunc(func(any, model.ID, *web.EventContext) error {
		return errors.New("accesskeys: the history is read only")
	})
	mb.Listing().DeleteFunc(func(any, model.ID, bool, *web.EventContext) error {
		return errors.New("accesskeys: the history is read only")
	})
}

// posted says whether the form of r was posted back (its fields there): not
// the form just opened.
func posted(r *http.Request) bool {
	_ = r.ParseMultipartForm(32 << 20)
	_, ok := r.Form["Name"]
	return ok
}
