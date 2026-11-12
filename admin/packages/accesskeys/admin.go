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
	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/account"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	if err := b.db.AutoMigrate(&AccessKey{}, &AccessKeyUse{}); err != nil {
		return err
	}
	if err := MigrateLines(b.db); err != nil {
		return err
	}
	return MigrateUngroup(b.db)
}

// MigrateUngroup cleans the lines of the keys (perm.CleanLines): what they
// deny of resources inside the groups of the menu by their unique names, no
// placeholder ("deny admin:", which decides nothing). Only the keys whose
// lines change; once done, nothing more.
func MigrateUngroup(db *gorm.DB) error {
	var all []*AccessKey
	if err := db.Select("id", "permissions").Find(&all).Error; err != nil {
		return err
	}
	for _, k := range all {
		lines, changed := perm.CleanLines(k.Permissions)
		if !changed {
			continue
		}
		if err := db.Model(k).UpdateColumn("permissions", pq.StringArray(lines)).Error; err != nil {
			return fmt.Errorf("accesskeys: the lines of %s out of the groups: %w", k.ID, err)
		}
	}
	return nil
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
	// the stacks of the keys: registered where the requests are decided
	b.auth.Perm = p.GetPermission()
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

	// the menu of the user (the top bar): their keys
	p.AddUserMenuItem(&presets.UserMenuItem{
		Name:    "MyAccessKeys",
		Group:   presets.UserMenuQuick,
		Order:   20,
		Enabled: true,
		Handler: func(ctx *web.EventContext) h.HTMLComponent {
			href := account.URL(p.GetURIPrefix(), myKeys.Info().ListingHref(), account.IndexOf(ctx.R))
			return v.VListItem(
				web.Slot(v.VIcon("mdi-key-chain")).Name(v.VSlotPrepend),
				v.VListItemTitle(h.Text(myKeys.TTitlePlural(ctx.Context()))),
			).Href(href)
		},
	})

	for _, keys := range []*presets.ModelBuilder{userKeys, myKeys} {
		uses := presets.NewModelBuilder(p, &AccessKeyUse{}, config(UsesModelID+"_"+keys.MenuID()))
		keys.AddChild(uses)
		b.configureUses(uses)
	}

	// any user may keep their own keys; no key may reach any key
	// (ProtectedLines)
	p.GetPermission().CreatePolicies(
		perm.PolicyFor(perm.Anybody).WhoAre(perm.Allowed).ToDo(perm.Anything).On("admin:" + MyKeysModelID + ":*"),
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

	// its lines, denies only, edited as one value — the tree of the
	// permissions of the admin it touches (userdocs.PermissionInput) —, kept
	// to what the key's user may, what the key allows its effect (all but its
	// denies and the protected)
	event := mb.Info().URIName() + "_" + PermissionInputEvent
	permOpts := func(ctx *web.EventContext) (userdocs.PermissionInputOptions, error) {
		allows, err := b.ownerAllows(ctx, owner)
		if err != nil {
			return userdocs.PermissionInputOptions{}, err
		}
		// the help of the permissions: a ? in the bar of the field
		help := ""
		if b.PermissionsHelpURL != nil {
			help = b.PermissionsHelpURL(ctx)
		}
		return userdocs.PermissionInputOptions{URL: ctx.R.URL.Path, Event: event,
			HelpURL: help, HelpTitle: GetMessages(ctx.Context()).PermissionsHelp,
			Prefixes: []string{userdocs.PermPrefixDeny}, StringView: true,
			Value: userdocs.PermValueOptions{Base: allows, Effect: func(lines []string) func(userdocs.PermCheck) bool {
				st := (&AccessKey{Permissions: lines}).Stack().Resolve()
				return func(c userdocs.PermCheck) bool { return st.Allows(c.Resource, c.Preferred, c.Action) }
			}}}, nil
	}
	ed.Field("Permissions").ComponentFunc(userdocs.PermissionFieldComponentFunc(b.p, permOpts)).
		SetterFunc(userdocs.PermissionFieldSetter)
	mb.RegisterEventFunc(event, userdocs.PermissionInputEventFunc(b.p, permOpts))

	// a new key starts enabled — the form opened, not posted back
	ed.Field("Enabled").WrapComponentFunc(func(old presets.FieldComponentFunc) presets.FieldComponentFunc {
		return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			if k := field.Obj.(*AccessKey); k.ID == uuid.Nil && !posted(ctx.R) {
				k.Enabled = true
			}
			return old(field, ctx)
		}
	})

	// the permissions in the detail: one a line, the denials marked
	mb.Detailing().Field("Permissions").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		k := field.Obj.(*AccessKey)
		m := GetMessages(ctx.Context())
		var items []h.HTMLComponent
		for _, p := range k.Permissions {
			li := h.Li(h.Code(p))
			if strings.HasPrefix(p, perm.PrefixDeny) {
				li.Class("text-error").Attr("data-deny", true)
			}
			items = append(items, li)
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
		// denies only: "deny " and a resource
		for _, l := range k.Permissions {
			if CheckLine(l) != nil {
				errs.FieldError("Permissions", fmt.Sprintf(m.ErrNotDeny, strings.TrimSpace(l)))
			}
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
			k.UserID = model.UUID(owner)
			k.Code, k.Prefix = NewCode()
			k.CodeHash = CodeHash(k.Code)
			if by, ok := currentUserID(ctx.R); ok {
				k.CreatedByID = model.UUID(by)
			}
		} else if k.UserID != model.UUID(owner) {
			return errors.New("accesskeys: not a key of this user")
		}
		// its permissions: its rules, as written, the blank ones out
		k.Permissions = cleanLines(k.Permissions)
		err := b.db.Transaction(func(tx *gorm.DB) error {
			if isNew {
				return tx.Omit("Policies").Create(k).Error
			}
			return tx.Model(k).Select("Name", "Description", "ExpiresAt", "Enabled", "Permissions").Updates(k).Error
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

// codeMessage is the message of a key just made: it stays until it is
// closed (Persistent), its code in it with a button to copy — shown this
// time only. The code is in the body alone, never in the text: the text may
// go in the cookie of a page's flash.
func codeMessage(ctx context.Context, k *AccessKey) *presets.FlashMessage {
	m := GetMessages(ctx)
	return &presets.FlashMessage{
		Text:       m.CodeCreated,
		Persistent: true,
		Body: h.Div(
			h.P(h.Text(m.CodeOnce)).Class("mb-2 text-body-2"),
			v.VTextField().Label(m.Code).ModelValue(k.Code).Readonly(true).Variant(v.FieldVariantSolo).
				Density(v.DensityCompact).HideDetails(true).Attr("data-access-key-code", true).
				Attr("append-inner-icon", "mdi-content-copy").
				Attr("@click:append-inner", "copyToClipboard("+h.JSONString(k.Code)+")").
				Style("min-width: 360px; font-family: monospace"),
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

// MigrateLines makes the permissions of the keys of before the lines —
// policies of the key's subject — its lines, once: a key of no lines and of
// policies gets theirs (perm.PolicyLines: "deny " before a deny's; the
// others of no prefix, what the key may reach — kept so, the key a key of
// before the denies), and its policies go.
func MigrateLines(db *gorm.DB) error {
	if !db.Migrator().HasTable(&perm.DefaultDBPolicy{}) {
		return nil
	}
	var all []*AccessKey
	if err := db.Find(&all).Error; err != nil {
		return err
	}
	for _, k := range all {
		if err := db.Where("refer_id = ?", k.ID.String()).Order("created_at, id").Find(&k.Policies).Error; err != nil {
			return err
		}
		if len(k.Permissions) > 0 || len(k.Policies) == 0 {
			continue
		}
		var lines pq.StringArray
		for _, p := range k.Policies {
			lines = append(lines, perm.PolicyLines(p.Effect, p.Actions, p.Resources)...)
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(k).UpdateColumn("permissions", lines).Error; err != nil {
				return err
			}
			return tx.Where("refer_id = ?", k.ID.String()).Delete(&perm.DefaultDBPolicy{}).Error
		})
		if err != nil {
			return fmt.Errorf("accesskeys: the lines of %s: %w", k.ID, err)
		}
	}
	return nil
}

// PermissionInputEvent is the event of the field of the permissions of the
// keys (userdocs.PermissionInputEventFunc), after the URI name of the model
// (the user's keys, My keys).
const PermissionInputEvent = "PermissionInput"

// ownerAllows is what the stack of the user of the keys of the request
// (owner) allows — what a key of it may deny —; nil when the permissions are
// not by stacks.
func (b *Builder) ownerAllows(ctx *web.EventContext, owner func(r *http.Request) (uuid.UUID, bool)) (func(userdocs.PermCheck) bool, error) {
	pb := b.p.GetPermission()
	if pb == nil || pb.GetStacks() == nil || b.auth.FindUser == nil {
		return nil, nil
	}
	id, ok := owner(ctx.R)
	if !ok {
		return nil, errors.New("accesskeys: no user of the keys")
	}
	u, err := b.auth.FindUser(ctx.R.Context(), id)
	if err != nil {
		return nil, err
	}
	// the stack of the owner: the request as if it were the owner's
	r := ctx.R.WithContext(context.WithValue(ctx.R.Context(), login.UserKey, u))
	_, _, build := pb.GetStacks().StackOf(r)
	if build == nil {
		return nil, nil
	}
	st, err := build()
	if err != nil {
		return nil, err
	}
	st.Resolve()
	return func(c userdocs.PermCheck) bool { return st.Allows(c.Resource, c.Preferred, c.Action) }, nil
}
