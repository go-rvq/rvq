package slug

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/reflect_utils"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/i18n"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/gosimple/unidecode"
	"github.com/sunfmin/reflectutils"
	"golang.org/x/text/language"
)

type Slug string

const (
	syncEvent                  = "slug_sync"
	I18nSlugKey i18n.ModuleKey = "I18nSlugKey"
)

type Builder struct{}

func New() *Builder {
	return &Builder{}
}

func (sb *Builder) Install(b *presets.Builder) error {
	b.I18n().
		RegisterForModule(language.English, I18nSlugKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nSlugKey, Messages_zh_CN)
	b.GetWebBuilder().RegisterEventFunc(syncEvent, sync)
	return nil
}

func (sb *Builder) ModelInstall(b *presets.Builder, mb *presets.ModelBuilder) error {
	reflectType := reflect.Indirect(reflect.ValueOf(mb.NewModel())).Type()
	if reflectType.Kind() != reflect.Struct {
		panic("slug: model must be struct")
	}

	mb.RegisterEventFunc(syncEvent, sync)

	_, fields := reflect_utils.UniqueFieldsOfReflectType(reflectType)
	for _, field := range fields {
		if field.Type != reflect.TypeOf(Slug("")) {
			continue
		}

		fieldName := field.Name
		relatedFieldName := strings.TrimSuffix(fieldName, "WithSlug")
		if _, ok := reflectType.FieldByName(relatedFieldName); ok {
			eb := mb.Editing()
			eb.WrapPostSetterFunc(func(in presets.SetterFunc) presets.SetterFunc {
				return func(obj interface{}, ctx *web.EventContext) {
					if ctx.R.FormValue(fieldName+"SlugSync") == "true" {
						v := reflectutils.MustGet(obj, relatedFieldName).(string)
						reflectutils.Set(obj, fieldName, Slugify(v))
					}
					if in != nil {
						in(obj, ctx)
					}
				}
			})
			if f := eb.Field(relatedFieldName); f != nil {
				f.ComponentFunc(SlugEditingComponentFunc)
			}

			eb.Field(fieldName).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) (r h.HTMLComponent) { return })
			eb.Field(fieldName).SetterFunc(SlugEditingSetterFunc)
		}
	}
	return nil
}

func SlugEditingComponentFunc(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nSlugKey, Messages_en_US).(*Messages)
	slugFieldName := field.Name + "WithSlug"
	slugLabel := strings.TrimSpace(strings.TrimSuffix(field.Label, "*")) + " Slug"
	ckbName := checkBoxName(slugFieldName)

	sync := true

	if field.Mode.Dot().Is(presets.EDIT) {
		sync = false
	}

	return vue.UserComponent().
		Assign("form", ckbName, sync).
		Scope("checkboxName", vue.Var(strconv.Quote(checkBoxName(slugFieldName)))).
		Scope("sync").
		Setup(`({scope, debounce}) => {
			scope.sync = debounce(function(v) {
				if (form[scope.checkboxName]) {
					` + web.Plaid().EventFunc(syncEvent).Query("field_name", field.Name).Query("slug_label", slugLabel).String() + `.go()
				}
			}, 300)
		}`).
		AppendChild(VSheet(
			VTextField().
				Type("text").
				Attr(web.VField(field.Name, field.Value())...).
				Label(field.Label).
				Attr("@update:modelValue", `(e) => sync(e)`),

			VRow(
				VCol(
					web.Portal(
						VTextField().
							Type("text").
							Attr(web.VField(slugFieldName, reflectutils.MustGet(field.Obj, slugFieldName).(Slug))...).
							Label(slugLabel),
					).Name(portalName(slugFieldName)),
				).Cols(8),
				VCol(
					VCheckbox().
						Attr("v-model", fmt.Sprintf("form[%q]", ckbName)).
						Label(fmt.Sprintf(msgr.Sync, strings.ToLower(field.Label))),
				).Cols(4),
			),
		))
}

func SlugEditingSetterFunc(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
	v := ctx.R.FormValue(field.Name)
	err = reflectutils.Set(obj, field.Name, Slug(v))
	if err != nil {
		return
	}
	return
}

func sync(ctx *web.EventContext) (r web.EventResponse, err error) {
	fieldName := ctx.R.FormValue("field_name")
	if fieldName == "" {
		return
	}

	slugFieldName := fieldName + "WithSlug"
	if checked := ctx.R.FormValue(checkBoxName(slugFieldName)); checked != "true" {
		return
	}

	r.UpdatePortal(
		portalName(slugFieldName),
		VTextField().
			Type("text").
			Attr(web.VField(slugFieldName, Slugify(ctx.R.FormValue(fieldName)))...).
			Label(ctx.R.FormValue("slug_label")),
	)
	return
}

var (
	regexpNonAuthorizedChars = regexp.MustCompile("[^a-zA-Z0-9-_]")
	regexpMultipleDashes     = regexp.MustCompile("-+")
)

func Slugify(value string) string {
	value = strings.TrimSpace(value)
	value = unidecode.Unidecode(value)
	value = strings.ToLower(value)
	value = regexpNonAuthorizedChars.ReplaceAllString(value, "-")
	value = regexpMultipleDashes.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-_")
	return value
}

func portalName(field string) string {
	return fmt.Sprintf("%s_Portal", field)
}

func checkBoxName(field string) string {
	return fmt.Sprintf("%sSlugSync", field)
}
