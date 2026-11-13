package role

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// The sets of permissions of the admin (presets.PermSetSource): made in
// Administrador › Conjuntos de permissões — a name, what it gives by kind of
// resource, a title and a description in each language of the admin —,
// beside the sets of the code (presets.RegisterPermSets), whose names they
// may not take. A set saved or deleted, the stacks of the users are made
// again (Changed).

// PermSet is a set of permissions of the admin.
type PermSet struct {
	uuidkey.Model

	// Name is how a line writes it, after "%": "readers".
	Name string `gorm:"size:64;index" admin:"required"`
	// Definition is what it gives by kind of resource (perm.SetBlocks).
	Definition string `gorm:"type:text" admin:"required"`
	// Titles and Descriptions: in each language of the admin ("pt-BR":
	// "Leitura").
	Titles       map[string]string `gorm:"serializer:json;type:text"`
	Descriptions map[string]string `gorm:"serializer:json;type:text"`
}

// TableName is the table of the sets.
func (PermSet) TableName() string { return "perm_sets" }

// String is the name.
func (s *PermSet) String() string { return s.Name }

// textOf is the text of the language of ctx in texts: of its tag, else of
// its base language ("pt" of "pt-BR"), else the first there is.
func textOf(ctx context.Context, texts map[string]string) string {
	if len(texts) == 0 {
		return ""
	}
	lang := ""
	if d := i18n.DynaFromContext(ctx); d != nil {
		lang = d.GetLanguage()
	}
	if t := texts[lang]; t != "" {
		return t
	}
	base, _, _ := strings.Cut(lang, "-")
	keys := make([]string, 0, len(texts))
	for k := range texts {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		if kb, _, _ := strings.Cut(k, "-"); base != "" && strings.EqualFold(kb, base) && texts[k] != "" {
			return texts[k]
		}
	}
	for _, k := range keys {
		if texts[k] != "" {
			return texts[k]
		}
	}
	return ""
}

// presetsSet is s as a set of the presets (a definition that cannot be read:
// nil — it gives nothing).
func (s *PermSet) presetsSet() *presets.PermSet {
	blocks, err := perm.SetBlocks(s.Definition)
	if err != nil {
		return nil
	}
	titles, descriptions := s.Titles, s.Descriptions
	return &presets.PermSet{Name: s.Name, Blocks: blocks,
		Title:       func(ctx context.Context) string { return textOf(ctx, titles) },
		Description: func(ctx context.Context) string { return textOf(ctx, descriptions) }}
}

// permSetCache holds the sets of the admin, read again at a generation of the
// roles not the one they were read at.
type permSetCache struct {
	mu   sync.Mutex
	gen  int64
	ok   bool
	sets []*presets.PermSet
}

func (c *permSetCache) get(db *gorm.DB) []*presets.PermSet {
	c.mu.Lock()
	defer c.mu.Unlock()
	if g := Generation(); c.ok && c.gen == g {
		return c.sets
	} else {
		c.gen = g
	}
	var rows []*PermSet
	if db.Session(&gorm.Session{NewDB: true}).Order("name").Find(&rows).Error != nil {
		return c.sets
	}
	c.sets = c.sets[:0]
	for _, r := range rows {
		if s := r.presetsSet(); s != nil {
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
// and the source of the sets of pb.
func (b *Builder) installPermSets(pb *presets.Builder) error {
	if err := b.db.AutoMigrate(&PermSet{}); err != nil {
		return err
	}
	cache := &permSetCache{}
	pb.PermSetSource(func() []*presets.PermSet { return cache.get(b.db) })

	mb := pb.Model(&PermSet{}, presets.ModelConfig().SetModuleKey(MessagesKey), presets.ModelWithID("perm_sets")).
		MenuIcon("mdi-shield-star-outline")
	mb.Listing("Name", "Titles", "Definition")
	ed := mb.Editing("Name", "Titles", "Descriptions", "Definition")
	mb.Detailing("Name", "Titles", "Descriptions", "Definition")

	ed.Field("Name").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextField().Variant(v.FieldVariantUnderlined).Label(field.Label).
			Hint(GetMessages(ctx.Context()).PermSetName_Hint).PersistentHint(true).Prefix("%").
			Attr(web.VField(field.FormKey, field.StringValue())...).ErrorMessages(field.Errors...)
	})
	ed.Field("Definition").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextarea().Variant(v.FieldVariantOutlined).Label(field.Label).AutoGrow(true).Rows(3).
			Hint(GetMessages(ctx.Context()).PermSetDefinition_Hint).PersistentHint(true).Class("vx-monospace").
			Attr("style", "font-family:monospace").
			Attr(web.VField(field.FormKey, field.StringValue())...).ErrorMessages(field.Errors...)
	})
	texts := func(name string, get func(*PermSet) *map[string]string) {
		ed.Field(name).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			s := field.Obj.(*PermSet)
			values := *get(s)
			var rows []h.HTMLComponent
			for _, tag := range pb.I18n().GetSupportLanguages() {
				lang := tag.String()
				rows = append(rows, v.VTextField().Variant(v.FieldVariantUnderlined).Density(v.DensityCompact).
					Label(field.Label+" ("+lang+")").
					Attr(web.VField(field.FormKey+"."+lang, values[lang])...))
			}
			rows = append(rows, h.Input("").Type("hidden").Attr(web.VField(field.FormKey+".__sent", "1")...))
			return h.Div(rows...).Class("mb-2").Attr("data-perm-set-texts", name)
		}).SetterFunc(func(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
			if ctx.R.FormValue(field.FormKey+".__sent") == "" {
				return nil
			}
			out := map[string]string{}
			for _, tag := range pb.I18n().GetSupportLanguages() {
				if t := strings.TrimSpace(ctx.R.FormValue(field.FormKey + "." + tag.String())); t != "" {
					out[tag.String()] = t
				}
			}
			*get(obj.(*PermSet)) = out
			return nil
		})
		mb.Detailing().Field(name).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
				h.Div(h.Text(textOf(ctx.Context(), *get(field.Obj.(*PermSet))))).Class("text-body-1")).Class("mb-3")
		})
		mb.Listing().Field(name).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return h.Td(h.Text(textOf(ctx.Context(), *get(field.Obj.(*PermSet)))))
		})
	}
	texts("Titles", func(s *PermSet) *map[string]string { return &s.Titles })
	texts("Descriptions", func(s *PermSet) *map[string]string { return &s.Descriptions })
	mb.Detailing().Field("Definition").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
			h.Pre(field.StringValue()).Class("text-body-2").Attr("data-perm-set-definition", true)).Class("mb-3")
	})

	ed.Validators.AppendFunc(func(obj interface{}, _ presets.FieldModeStack, ctx *web.EventContext) (errs web.ValidationErrors) {
		s := obj.(*PermSet)
		msgs := GetMessages(ctx.Context())
		s.Name = strings.TrimPrefix(strings.TrimSpace(s.Name), perm.SetPrefix)
		if !perm.ValidSetName(s.Name) {
			errs.FieldError("Name", msgs.ErrPermSetName)
		} else if code := pb.PermSet(s.Name); code != nil && code.System {
			errs.FieldError("Name", msgs.ErrPermSetTaken)
		} else {
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
	// saved or deleted: the stacks made again, with it
	ed.SaveFunc(func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
		if err = gorm2op.DataOperator(b.db).Save(obj, id, ctx); err == nil {
			b.changed()
		}
		return
	})
	ed.WrapCreateFunc(func(in presets.CreateFunc) presets.CreateFunc {
		return func(obj interface{}, ctx *web.EventContext) (err error) {
			if err = in(obj, ctx); err == nil {
				b.changed()
			}
			return
		}
	})
	mb.Listing().DeleteFunc(func(obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error) {
		if err = b.db.Delete(&PermSet{}, "id = ?", id.String()).Error; err == nil {
			b.changed()
		}
		return
	})
	b.permSetMb = mb
	return nil
}
