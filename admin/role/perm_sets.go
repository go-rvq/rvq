package role

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/l10n"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/packages/user/userref"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// The sets of permissions of the admin (presets.PermSetSource): made in
// Administrador › Conjuntos de permissões — a name, what it gives by kind of
// resource —, beside the sets of the code (presets.RegisterPermSets), whose
// names they may not take. A set is localized by the languages of the admin
// (l10n.NewAdmin), as pages and posts by the site's: a record in each
// language, its title and description in it; its name and definition are the
// set's, the same in all of them (saved in one, written in the others). A set
// saved or deleted, the stacks of the users are made again (Changed).

// PermSet is a set of permissions of the admin, in one language.
type PermSet struct {
	uuidkey.Model
	l10n.Locale

	// Name is how a line writes it, after "%": "readers".
	Name string `gorm:"size:64;index" admin:"required"`
	// Definition is what it gives by kind of resource (perm.SetBlocks).
	Definition string `gorm:"type:text" admin:"required"`
	// Title and Description: in the language of the record (LocaleCode).
	Title       string `gorm:"size:255"`
	Description string `gorm:"type:text"`
	// System: a set of the code (presets.RegisterPermSets), registered here on
	// boot — its name and definition the code's, read only; its title and
	// description translated in each language.
	System bool `gorm:"not null;default:false"`

	// who made it and who changed it last (CreatedAt, UpdatedAt: when): the
	// users, read by a LEFT JOIN (Fetch, Search); a user deleted, nil
	CreatedByID model.UUID   `gorm:"type:uuid"`
	CreatedBy   *userref.Ref `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	UpdatedByID model.UUID   `gorm:"type:uuid"`
	UpdatedBy   *userref.Ref `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName is the table of the sets.
func (PermSet) TableName() string { return "perm_sets" }

// String is the name.
func (s *PermSet) String() string { return s.Name }

// PrimarySlug is the id and the language: "<uuid>_pt-BR".
func (s *PermSet) PrimarySlug() string {
	if s.ID == uuid.Nil {
		return ""
	}
	return s.ID.String() + "_" + s.LocaleCode
}

// PrimaryColumnValuesBySlug reads PrimarySlug.
func (s *PermSet) PrimaryColumnValuesBySlug(slug string) map[string]string {
	id, locale, _ := strings.Cut(slug, "_")
	return map[string]string{"ID": id, "LocaleCode": locale}
}

// stamp writes in s who changes it, by r: made (isNew) and changed.
func (s *PermSet) stamp(r *http.Request, isNew bool) {
	id := userref.IDOf(r)
	if isNew {
		s.CreatedByID = id
	}
	s.UpdatedByID = id
	s.CreatedBy, s.UpdatedBy = nil, nil
}

// localizedOf is the record of rows in the language of ctx: of its tag, else
// of its base language ("pt" of "pt-BR"), else of def, else the first.
func localizedOf(ctx context.Context, rows []*PermSet, def string) *PermSet {
	if len(rows) == 0 {
		return nil
	}
	lang := ""
	if d := i18n.DynaFromContext(ctx); d != nil {
		lang = d.GetLanguage()
	}
	base, _, _ := strings.Cut(lang, "-")
	var byBase, byDef *PermSet
	for _, r := range rows {
		if r.LocaleCode == lang {
			return r
		}
		if rb, _, _ := strings.Cut(r.LocaleCode, "-"); byBase == nil && base != "" && strings.EqualFold(rb, base) {
			byBase = r
		}
		if r.LocaleCode == def {
			byDef = r
		}
	}
	if byBase != nil {
		return byBase
	}
	if byDef != nil {
		return byDef
	}
	return rows[0]
}

// presetsSet is the set of the records of one id (its languages) as a set of
// the presets (a definition that cannot be read: nil — it gives nothing).
func presetsSet(rows []*PermSet, def string) *presets.PermSet {
	blocks, err := perm.SetBlocks(rows[0].Definition)
	if err != nil {
		return nil
	}
	text := func(get func(*PermSet) string) func(ctx context.Context) string {
		return func(ctx context.Context) string {
			if r := localizedOf(ctx, rows, def); r != nil && get(r) != "" {
				return get(r)
			}
			for _, r := range rows {
				if get(r) != "" {
					return get(r)
				}
			}
			return ""
		}
	}
	return &presets.PermSet{Name: rows[0].Name, Blocks: blocks,
		Title:       text(func(s *PermSet) string { return s.Title }),
		Description: text(func(s *PermSet) string { return s.Description })}
}

// permSetCache holds the sets of the admin, read again at a generation of the
// roles not the one they were read at.
type permSetCache struct {
	mu   sync.Mutex
	gen  int64
	ok   bool
	sets []*presets.PermSet
	// sync registers the sets of the code, once, before they are first read
	sync   func() error
	synced bool
}

func (c *permSetCache) get(db *gorm.DB, def func() string) []*presets.PermSet {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.synced && c.sync != nil {
		if err := c.sync(); err != nil {
			log.Printf("role: the sets of permissions of the system: %v", err)
		}
		c.synced = true
	}
	if g := Generation(); c.ok && c.gen == g {
		return c.sets
	} else {
		c.gen = g
	}
	var rows []*PermSet
	if db.Session(&gorm.Session{NewDB: true}).Order("name, id, locale_code").Find(&rows).Error != nil {
		return c.sets
	}
	// by id: the languages of a set
	var ids []uuid.UUID
	byID := map[uuid.UUID][]*PermSet{}
	for _, r := range rows {
		if byID[r.ID] == nil {
			ids = append(ids, r.ID)
		}
		byID[r.ID] = append(byID[r.ID], r)
	}
	c.sets = c.sets[:0]
	for _, id := range ids {
		if s := presetsSet(byID[id], def()); s != nil {
			c.sets = append(c.sets, s)
		}
	}
	c.ok = true
	return c.sets
}

// permSetKinds are the kinds of resource a definition may name.
var permSetKinds = []string{string(presets.PermNodeModel), string(presets.PermNodeSingleton), string(presets.PermNodePage),
	string(presets.PermNodeCheck), string(presets.PermNodeRecord), string(presets.PermNodeSection),
	string(presets.PermNodeField), string(presets.PermNodeInline)}

// installPermSets is the model of the sets of the admin, beside the roles,
// localized by the languages of the admin, and the source of the sets of pb.
func (b *Builder) installPermSets(pb *presets.Builder) error {
	if err := migrate.Auto(b.db, func(db *gorm.DB) error { return db.AutoMigrate(&PermSet{}) }); err != nil {
		return err
	}
	lb := l10n.NewAdmin(b.db, pb.I18n()).Activity(b.ab)
	b.permSetL10n = lb
	cache := &permSetCache{sync: func() error { return b.syncSystemPermSets(pb, lb) }}
	pb.PermSetSource(func() []*presets.PermSet { return cache.get(b.db, lb.GetDefaultLocaleCode) })
	// the sets of the code registered on boot, all of them there by then —
	// as the application serves when the graph runs the migrations
	// (migrate.Manage), after them
	if migrate.Managed(b.db) {
		pb.Lifecycle().Add(lifecycle.Node{
			Key: PermSetsKey, After: []lifecycle.Key{migrate.Key}, Optional: []lifecycle.Key{migrate.Key},
			Start: func(context.Context) (io.Closer, error) {
				cache.get(b.db, lb.GetDefaultLocaleCode)
				return nil, nil
			},
		})
	} else {
		pb.AfterBuild(func() { cache.get(b.db, lb.GetDefaultLocaleCode) })
	}

	mb := pb.Model(&PermSet{}, presets.ModelConfig().SetModuleKey(MessagesKey), presets.ModelWithID("perm_sets")).
		MenuIcon("mdi-shield-star-outline")
	mb.Listing("Name", "Title", "Definition", "System").
		SearchColumns(`"perm_sets".name`, `"perm_sets".title`, `"perm_sets".definition`)
	// who made and changed it: the users by a LEFT JOIN
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		g, ok := do.(*gorm2op.DataOperatorBuilder)
		if !ok {
			return do
		}
		return g.WithModeSplitCallbacks(gorm2op.Fetch|gorm2op.Search, func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
			cb.Pre(func(state *gorm2op.CallbackState) error {
				state.DB.Joins("CreatedBy").Joins("UpdatedBy")
				return nil
			})
		})
	})
	userField := func(get func(*PermSet) *userref.Ref) presets.FieldComponentFunc {
		return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
				h.Div(h.Text(get(field.Obj.(*PermSet)).String())).Class("text-body-1")).Class("mb-3")
		}
	}
	mb.Detailing().Field("CreatedBy").ComponentFunc(userField(func(s *PermSet) *userref.Ref { return s.CreatedBy }))
	mb.Detailing().Field("UpdatedBy").ComponentFunc(userField(func(s *PermSet) *userref.Ref { return s.UpdatedBy }))
	ed := mb.Editing("Name", "Title", "Description", "Definition")
	mb.Detailing("Name", "Title", "Description", "Definition", "System", "CreatedAt", "CreatedBy", "UpdatedAt", "UpdatedBy")

	system := func(field *presets.FieldContext) bool { return field.Obj.(*PermSet).System }
	ed.Field("Name").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextField().Variant(v.FieldVariantUnderlined).Label(field.Label).
			Hint(GetMessages(ctx.Context()).PermSetName_Hint).PersistentHint(true).Prefix("%").
			Readonly(system(field)).
			Attr(web.VField(field.FormKey, field.StringValue())...).ErrorMessages(field.Errors...)
	})
	ed.Field("Description").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextarea().Variant(v.FieldVariantUnderlined).Label(field.Label).AutoGrow(true).Rows(1).
			Attr(web.VField(field.FormKey, field.StringValue())...).ErrorMessages(field.Errors...)
	})
	ed.Field("Definition").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextarea().Variant(v.FieldVariantOutlined).Label(field.Label).AutoGrow(true).Rows(3).
			Hint(GetMessages(ctx.Context()).PermSetDefinition_Hint).PersistentHint(true).Class("vx-monospace").
			Readonly(system(field)).
			Attr("style", "font-family:monospace").
			Attr(web.VField(field.FormKey, field.StringValue())...).ErrorMessages(field.Errors...)
	})
	mb.Detailing().Field("Definition").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
			h.Pre(field.StringValue()).Class("text-body-2").Attr("data-perm-set-definition", true)).Class("mb-3")
	})

	ed.Validators.AppendFunc(func(obj interface{}, _ presets.FieldModeStack, ctx *web.EventContext) (errs web.ValidationErrors) {
		s := obj.(*PermSet)
		msgs := GetMessages(ctx.Context())
		if s.System {
			// the code's: its name and definition as they are (SaveFunc)
			return
		}
		s.Name = strings.TrimPrefix(strings.TrimSpace(s.Name), perm.SetPrefix)
		if !perm.ValidSetName(s.Name) {
			errs.FieldError("Name", msgs.ErrPermSetName)
		} else if code := pb.PermSet(s.Name); code != nil && code.System {
			errs.FieldError("Name", msgs.ErrPermSetTaken)
		} else {
			// another set's — its languages are the same set
			var n int64
			q := b.db.Model(&PermSet{}).Where("name = ?", s.Name)
			if s.ID != uuid.Nil {
				q = q.Where("id <> ?", s.ID)
			}
			if q.Count(&n); n > 0 {
				errs.FieldError("Name", msgs.ErrPermSetTaken)
			}
		}
		blocks, err := perm.SetBlocks(s.Definition)
		if err != nil {
			errs.FieldError("Definition", err.Error())
			return
		}
		for kind := range blocks {
			if !slices.Contains(permSetKinds, kind) {
				errs.FieldError("Definition", fmt.Sprintf(msgs.ErrPermSetKind, kind))
			}
		}
		return
	})
	// saved: its name and definition in its other languages too; the stacks
	// made again
	ed.SaveFunc(func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
		s := obj.(*PermSet)
		return b.db.Transaction(func(tx *gorm.DB) error {
			// a set of the system: its name and definition the code's, whatever
			// the form said
			var old PermSet
			if tx.Where("id = ? AND locale_code = ?", s.ID, s.LocaleCode).First(&old).Error == nil {
				if old.System {
					s.Name, s.Definition, s.System = old.Name, old.Definition, true
				}
				// who made it is not the form's
				s.CreatedByID, s.CreatedAt = old.CreatedByID, old.CreatedAt
			}
			s.stamp(ctx.R, false)
			// the users only read: never written through a set
			if err := gorm2op.DataOperator(tx.Omit(clause.Associations)).Save(obj, id, ctx); err != nil {
				return err
			}
			if err := tx.Model(&PermSet{}).Where("id = ? AND locale_code <> ?", s.ID, s.LocaleCode).
				Updates(map[string]any{"name": s.Name, "definition": s.Definition}).Error; err != nil {
				return err
			}
			b.changed()
			return nil
		})
	})
	// created (a localization too: l10n copies by the creator)
	ed.WrapCreateFunc(func(in presets.CreateFunc) presets.CreateFunc {
		return func(obj interface{}, ctx *web.EventContext) (err error) {
			obj.(*PermSet).stamp(ctx.R, true)
			if err = in(obj, ctx); err == nil {
				b.changed()
			}
			return
		}
	})
	mb.Listing().WrapDeleteFunc(func(in presets.DeleteFunc) presets.DeleteFunc {
		return func(obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error) {
			var n int64
			if b.db.Model(&PermSet{}).Where("id = ? AND system", id.GetValue("ID")).Count(&n); n > 0 {
				return errors.New(GetMessages(ctx.Context()).ErrPermSetSystem)
			}
			if err = in(obj, id, cascade, ctx); err == nil {
				b.changed()
			}
			return
		}
	})
	// localized by the languages of the admin: the switch, Localize, the
	// localizations (l10n)
	mb.Use(lb)
	// its languages: the column of the plugin's (their flags; its filter, the
	// plugin's too)
	mb.Listing("Name", "Title", "Definition", "System", l10n.FieldLocationLabel)
	if b.ab != nil {
		// its history
		mb.Use(b.ab)
	}
	b.permSetMb = mb
	return nil
}

// syncSystemPermSets registers the sets of the code of pb (System): one not
// here, in each language of lb, of the title and description the code gives
// it there; one here, of the code's definition (its titles, as translated
// here).
func (b *Builder) syncSystemPermSets(pb *presets.Builder, lb *l10n.Builder) error {
	db := b.db.Session(&gorm.Session{NewDB: true})
	return db.Transaction(func(tx *gorm.DB) error {
		for _, s := range pb.CodePermSets() {
			def := s.Definition()
			var n int64
			if err := tx.Model(&PermSet{}).Where("name = ?", s.Name).Count(&n).Error; err != nil {
				return err
			}
			if n > 0 {
				// only when the code changed it: its updated_at is a change's
				if err := tx.Model(&PermSet{}).Where("name = ? AND (definition <> ? OR NOT system)", s.Name, def).
					UpdateColumns(map[string]any{"system": true, "definition": def}).Error; err != nil {
					return err
				}
				continue
			}
			id := uuidkey.New()
			for _, code := range lb.GetSupportLocaleCodes() {
				ctx := pb.I18n().LanguageContext(context.Background(), code)
				r := &PermSet{Name: s.Name, Definition: def, System: true, Description: s.DescriptionOf(ctx)}
				r.ID, r.LocaleCode = id, code
				if t := s.TitleOf(ctx); t != s.Name {
					r.Title = t
				}
				if err := tx.Create(r).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// LocalizePermSets gives the sets of the admin in from a localization in to
// — a language added to the admin —, of their title and description
// translated by translate (texts in from to texts in to; nil, as they are).
// A set already in to is left as it is.
func (b *Builder) LocalizePermSets(from, to string, translate func(texts []string) ([]string, error)) (n int, err error) {
	db := b.db.Session(&gorm.Session{})
	var rows []*PermSet
	if err = db.Where("locale_code = ? AND id NOT IN (?)", from,
		db.Model(&PermSet{}).Select("id").Where("locale_code = ?", to)).Find(&rows).Error; err != nil || len(rows) == 0 {
		return
	}
	texts := make([]string, 0, 2*len(rows))
	for _, r := range rows {
		texts = append(texts, r.Title, r.Description)
	}
	if translate != nil {
		var out []string
		if out, err = translate(texts); err != nil {
			return
		}
		if len(out) != len(texts) {
			return 0, fmt.Errorf("perm sets: %d texts translated into %d", len(texts), len(out))
		}
		texts = out
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		for i, r := range rows {
			c := *r
			c.LocaleCode, c.Title, c.Description = to, texts[2*i], texts[2*i+1]
			c.CreatedAt, c.UpdatedAt = time.Time{}, time.Time{}
			if err := tx.Create(&c).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err == nil {
		n = len(rows)
		b.changed()
	}
	return
}

// PermSetsKey is the node of the sets of permissions of the code written to
// the database (the system's), as the application serves.
const PermSetsKey lifecycle.Key = "rvq/role/perm-sets"
