package role

import (
	"net/http"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/ory/ladon"
	"gorm.io/gorm"
)

type Builder struct {
	db        *gorm.DB
	actions   []*DefaultOptionItem
	resources []*DefaultOptionItem
	// editorSubject is the subject that has permission to edit roles
	// empty value means anyone can edit roles
	editorSubject    string
	roleMb           *presets.ModelBuilder
	AfterInstallFunc presets.ModelInstallFunc
}

func New(db *gorm.DB) *Builder {
	return &Builder{
		db: db,
		actions: []*DefaultOptionItem{
			{Text: "All", Value: "*"},
			{Text: "List", Value: presets.PermList},
			{Text: "Get", Value: presets.PermGet},
			{Text: "Create", Value: presets.PermCreate},
			{Text: "Update", Value: presets.PermUpdate},
			{Text: "Delete", Value: presets.PermDelete},
		},
	}
}

func (b *Builder) Actions(vs []*DefaultOptionItem) *Builder {
	b.actions = vs
	return b
}

func (b *Builder) AfterInstall(v presets.ModelInstallFunc) *Builder {
	b.AfterInstallFunc = v
	return b
}

func (b *Builder) Resources(vs []*DefaultOptionItem) *Builder {
	b.resources = vs
	return b
}

func (b *Builder) EditorSubject(v string) *Builder {
	b.editorSubject = v
	return b
}

func (b *Builder) Install(pb *presets.Builder) (err error) {
	ConfigureMessages(pb.I18n())

	if b.editorSubject != "" {
		permB := pb.GetPermission()
		if permB == nil {
			panic("pb does not have a permission builder")
		}
		ctxf := permB.GetContextFunc()
		ssf := permB.GetSubjectsFunc()
		permB.ContextFunc(func(r *http.Request, objs []interface{}) perm.Context {
			c := make(perm.Context)
			if ctxf != nil {
				c = ctxf(r, objs)
			}
			ss := ssf(r)
			hasRoleEditorSubject := false
			for _, s := range ss {
				if s == b.editorSubject {
					hasRoleEditorSubject = true
					break
				}
			}
			c["has_role_editor_subject"] = hasRoleEditorSubject
			return c
		})
		permB.CreatePolicies(
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(perm.Anything).On("*:roles:*").Given(perm.Conditions{
				"has_role_editor_subject": &ladon.BooleanCondition{
					BooleanValue: false,
				},
			}),
			perm.PolicyFor(b.editorSubject).WhoAre(perm.Allowed).ToDo(perm.Anything).On("*:roles:*"),
		)
	}

	b.roleMb = pb.Model(&Role{}, presets.ModelConfig().
		SetModuleKey(MessagesKey)).
		MenuIcon("mdi-account-key")

	ed := b.roleMb.Editing(
		"Name",
		"Permissions",
	)

	policeModel := presets.NewModelBuilder(pb, &perm.DefaultDBPolicy{})
	permFb := &policeModel.Editing("Effect", "Actions", "Resources").FieldsBuilder
	ed.Field("Permissions").AutoNested(policeModel, permFb)

	permFb.Field("Effect").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		p := field.Obj.(*perm.DefaultDBPolicy)
		msgr := GetMessages(ctx.Context())

		return vue.FormField(vx.VXSelectOne().
			Label(field.Label).
			ItemText("text").ItemValue("value").
			Chips(true).
			Items([]DefaultOptionItem{
				{Value: perm.Allowed, Text: msgr.Allowed},
				{Value: perm.Denied, Text: msgr.Denied},
			}).ErrorMessages(field.Errors...)).Value(field.FormKey, p.Effect).Bind()
	}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		p := obj.(*perm.DefaultDBPolicy)
		p.Effect = ctx.R.FormValue(field.FormKey)
		return
	})

	permFb.Field("Actions").AsSlice()
	permFb.Field("Resources").AsSlice()

	ed.FetchFunc(func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
		return gorm2op.DataOperator(b.db.Preload("Permissions")).Fetch(obj, id, ctx)
	})

	ed.Validators.AppendFunc(func(obj interface{}, _ presets.FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*Role)
		for _, p := range u.Permissions {
			p.Subject = u.Name
		}
		return
	})

	ed.SaveFunc(func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
		r := obj.(*Role)
		if r.ID != 0 {
			if err = b.db.Delete(&perm.DefaultDBPolicy{}, "refer_id = ?", r.ID).Error; err != nil {
				return
			}
		}
		if err = gorm2op.DataOperator(b.db.Session(&gorm.Session{FullSaveAssociations: true})).Save(obj, id, ctx); err != nil {
			return
		}
		startFrom := time.Now().Add(-1 * time.Second)
		pb.GetPermission().LoadDBPoliciesToMemory(b.db, &startFrom)
		return
	})

	b.roleMb.Listing().DeleteFunc(func(obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error) {
		err = b.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&perm.DefaultDBPolicy{}, "refer_id = ?", id).Error; err != nil {
				return err
			}
			if err := tx.Delete(&Role{}, "id = ?", id).Error; err != nil {
				return err
			}

			return nil
		})

		return
	})

	b.roleMb.Detailing()

	if b.AfterInstallFunc != nil {
		return b.AfterInstallFunc(b.roleMb)
	}

	return nil
}
