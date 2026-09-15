package seo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/l10n"
	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	"golang.org/x/text/language"
)

const (
	I18nSeoKey         i18n.ModuleKey = "I18nSeoKey"
	SeoDetailFieldName                = "SEO"
)

var permVerifier *perm.Verifier

type myTd struct {
	td    *h.HTMLTagBuilder
	child h.MutableAttrHTMLComponent
}

func (mtd *myTd) SetAttr(k string, v interface{}) {
	mtd.td.SetAttr(k, v)
	mtd.child.SetAttr(k, v)
}

func (mtd *myTd) Write(ctx *h.Context) (err error) {
	mtd.td.Children(mtd.child)
	return mtd.td.Write(ctx)
}

func (b *Builder) Install(pb *presets.Builder) error {
	// The registration of FieldDefaults for writing Setting here
	// must be executed before `pb.Model(&RvqSEOSetting{})...`,
	pb.FieldDefaults(presets.WRITE).
		FieldType(Setting{}).
		ComponentFunc(b.EditingComponentFunc).
		SetterFunc(EditSetterFunc)

	// The help overlay event (the "?" button next to Save in any SEO editor).
	b.registerHelpEvent(pb)

	seoModel := pb.Model(&RvqSEOSetting{}).
		Label("SEO").
		RightDrawerWidth("1000").
		LayoutConfig(&presets.LayoutConfig{
			SearchBoxInvisible:          true,
			NotificationCenterInvisible: true,
		})

	// Configure Listing Page
	b.configListing(seoModel)
	// Configure Editing Page
	b.configEditing(seoModel)
	// b.ConfigDetailing(pb)

	pb.I18n().
		RegisterForModule(language.English, I18nSeoKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nSeoKey, Messages_zh_CN).
		RegisterForModule(language.BrazilianPortuguese, I18nSeoKey, Messages_pt_BR)

	permVerifier = perm.NewVerifier("seo", pb.GetPermission())
	return nil
}

func (b *Builder) configListing(seoModel *presets.ModelBuilder) {
	listing := seoModel.Listing("Name").Title("SEO")
	// disable new btn globally, no one can add new SEO record after the server start up.
	listing.NewButtonFunc(func(ctx *web.EventContext) h.HTMLComponent {
		return nil
	})
	listing.DisablePagination(true)

	// Remove the row menu from each row
	listing.RowMenu().Empty()

	// Configure the indentation for Name field to display hierarchy.
	listing.Field("Name").ComponentFunc(
		func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			seoSetting := field.Obj.(*RvqSEOSetting)
			icon := "mdi-folder"
			priority := b.GetSEOPriority(seoSetting.Name)
			return &myTd{
				td: h.Td(),
				child: h.Div(
					VIcon(icon).Size(SizeSmall).Class("mb-1"),
					h.Text(seoSetting.Name),
				).Style(fmt.Sprintf("padding-left: %dpx;", 32*(priority-1))),
			}
		},
	)

	listing.WrapSearchFunc(func(in presets.SearchFunc) presets.SearchFunc {
		return func(model interface{}, params *presets.SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error) {
			locale, _ := l10n.IsLocalizableFromContext(ctx.R.Context())
			var seoNames []string
			for name := range b.registeredSEO {
				if name, ok := name.(string); ok {
					seoNames = append(seoNames, name)
				}
			}
			cond := presets.SQLCondition{
				Query: "locale_code = ? and name in (?)",
				Args:  []interface{}{locale, seoNames},
			}

			params.SQLConditions = append(params.SQLConditions, &cond)
			r, totalCount, err = in(model, params, ctx)
			if totalCount == 0 {
				panic("The localization of SEO is not configured correctly. " +
					"Please check if you correctly configured the `WithLocales` option when initializing the SEO Builder.")
			}
			b.SortSEOs(r.([]*RvqSEOSetting))
			return
		}
	})
}

func (b *Builder) configEditing(seoModel *presets.ModelBuilder) {
	editing := seoModel.Editing("Variables", "Setting")

	// The "?" help button next to Save, for the Global SEO context.
	b.AddEditingHelpButton(editing, HelpContextGlobal)

	// Customize the Saver to trigger the invocation of the `afterSave` hook function (if available)
	// when updating the global seo.
	editing.SaveFunc(func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
		seoSetting := obj.(*RvqSEOSetting)
		if err = b.db.Updates(obj).Error; err != nil {
			return err
		}
		if b.afterSave != nil {
			if err = b.afterSave(ctx.R.Context(), seoSetting.Name, seoSetting.LocaleCode); err != nil {
				return err
			}
		}
		return nil
	})

	// configure variables field
	// Variables editor: its component resolves the variable set per row's Name,
	// and its setter writes back the map. (Reused by application-mounted SEO
	// models via VariablesComponentFunc/VariablesSetterFunc.)
	editing.Field("Variables").
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			ss := field.Obj.(*RvqSEOSetting)
			return b.VariablesComponentFunc(ss.Name)(field, ctx)
		}).
		SetterFunc(VariablesSetterFunc)

	editing.Field("Setting").ComponentFunc(
		func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			seoSetting := field.Obj.(*RvqSEOSetting)
			return b.vseo("Setting", b.GetSEO(seoSetting.Name), &seoSetting.Setting, ctx.R)
		},
	)
}

// SettingComponentFunc returns the editor component for a Setting field of a
// RvqSEOSetting, resolving the SEO by name (seoName) rather than by the object's
// type — so an application can mount its own model builder for a specific SEO
// (e.g. a singleton "Global SEO" scoped by locale) and still get the full SEO
// editor. Pair it with EditSetterFunc on the same field.
func (b *Builder) SettingComponentFunc(seoName string) func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		ss, ok := field.Obj.(*RvqSEOSetting)
		if !ok {
			return h.Div()
		}
		return b.vseo(field.Name, b.GetSEO(seoName), &ss.Setting, ctx.R)
	}
}

// SettingDetailComponentFunc is the read-only counterpart of
// SettingComponentFunc: it renders the Setting for a detail view, resolving the
// SEO by name. Use it in an application-mounted model's Detailing.
func (b *Builder) SettingDetailComponentFunc(seoName string) func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		ss, ok := field.Obj.(*RvqSEOSetting)
		if !ok {
			return h.Div()
		}
		return b.vseoReadonly(field.Name, b.GetSEO(seoName), &ss.Setting, ctx.R)
	}
}

// SettingModelDetailComponentFunc renders a record's Setting field read-only for
// a detail view — the read-only counterpart of the WRITE FieldDefaults editor,
// resolving the SEO from the record's own type (via GetSEO) rather than requiring
// a *RvqSEOSetting like SettingDetailComponentFunc. Use it for an app-mounted
// per-record SEO detailing whose model is the record itself (a Page/Post), which
// has no DETAIL FieldDefaults counterpart to the WRITE editor. The record's type
// must be registered (RegisterModel) or GetSEO is nil and this renders empty.
func (b *Builder) SettingModelDetailComponentFunc() presets.FieldComponentFunc {
	return b.detailShowComponent
}

// formKeyForVariablesField is the form prefix of the SEO setting variables.
const formKeyForVariablesField = "Variables"

// VariablesComponentFunc returns the editor for a RvqSEOSetting's setting
// variables (e.g. SiteName), resolving the variable set by seoName. Pair it with
// VariablesSetterFunc. Exposed so an application can mount its own SEO model.
func (b *Builder) VariablesComponentFunc(seoName string) func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		ss, ok := field.Obj.(*RvqSEOSetting)
		if !ok {
			return h.Div()
		}
		msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nSeoKey, Messages_en_US).(*Messages)
		theSEO := b.GetSEO(seoName)
		if theSEO == nil {
			return h.Div()
		}
		settingVars := theSEO.settingVars
		var comps h.HTMLComponents
		if len(settingVars) > 0 {
			comps = append(comps, h.H3(msgr.Variable).Style("margin-top:15px;font-weight: 500"))
			for varName := range settingVars {
				comps = append(comps, VTextField().
					Attr(web.VField(fmt.Sprintf("%s.%s", formKeyForVariablesField, varName), ss.Variables[varName])...).
					Label(i18n.PT(ctx.Context(), I18nSeoKey, "SettingVar", varName)))
			}
		}
		return comps
	}
}

// VariablesSetterFunc writes the submitted setting variables back to the
// RvqSEOSetting (its Variables is a map, so it needs an explicit setter).
func VariablesSetterFunc(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
	ss, ok := obj.(*RvqSEOSetting)
	if !ok {
		return nil
	}
	if ss.Variables == nil {
		ss.Variables = make(Variables)
	}
	for fieldName := range ctx.R.Form {
		if strings.HasPrefix(fieldName, formKeyForVariablesField+".") {
			ss.Variables[strings.TrimPrefix(fieldName, formKeyForVariablesField+".")] = ctx.R.Form[fieldName][0]
		}
	}
	return nil
}

func EditSetterFunc(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
	var setting Setting
	mediaBox := media_library.MediaBox{}
	for fieldWithPrefix := range ctx.R.Form {
		// make sure OpenGraphImageFromMediaLibrary.Description set after OpenGraphImageFromMediaLibrary.Values
		if fieldWithPrefix == fmt.Sprintf("%s.%s", field.Name, "OpenGraphImageFromMediaLibrary.Values") {
			err = mediaBox.Scan(ctx.R.FormValue(fieldWithPrefix))
			if err != nil {
				return
			}
			break
		}
	}
	for fieldWithPrefix := range ctx.R.Form {
		if strings.HasPrefix(fieldWithPrefix, fmt.Sprintf("%s.%s", field.Name, "OpenGraphImageFromMediaLibrary")) {
			if fieldWithPrefix == fmt.Sprintf("%s.%s", field.Name, "OpenGraphImageFromMediaLibrary.Description") {
				mediaBox.Description = ctx.R.Form.Get(fieldWithPrefix)
				setting.OpenGraphImageFromMediaLibrary = mediaBox
			}
			continue
		}
		if fieldWithPrefix == fmt.Sprintf("%s.%s", field.Name, "OpenGraphMetadataString") {
			metadata := GetOpenGraphMetadata(ctx.R.Form.Get(fieldWithPrefix))
			setting.OpenGraphMetadata = metadata
			continue
		}
		if strings.HasPrefix(fieldWithPrefix, fmt.Sprintf("%s.", field.Name)) {
			reflectutils.Set(&setting, strings.TrimPrefix(fieldWithPrefix, fmt.Sprintf("%s.", field.Name)), ctx.R.Form.Get(fieldWithPrefix))
		}
	}
	return reflectutils.Set(obj, field.Name, setting)
}

func (b *Builder) EditingComponentFunc(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	var (
		obj         = field.Obj
		msgr        = i18n.MustGetModuleMessages(ctx.Context(), I18nSeoKey, Messages_en_US).(*Messages)
		fieldPrefix string
		setting     Setting
		db          = b.db
		locale, _   = l10n.IsLocalizableFromContext(ctx.R.Context())
	)
	seo := b.GetSEO(obj)
	if seo == nil {
		return h.Div()
	}

	value := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < value.NumField(); i++ {
		if s, ok := value.Field(i).Interface().(Setting); ok {
			setting = s
			fieldPrefix = value.Type().Field(i).Name
		}
	}
	if !setting.EnabledCustomize && setting.IsEmpty() {
		modelSetting := &RvqSEOSetting{}
		db.Where("name = ? AND locale_code = ?", seo.name, locale).First(modelSetting)
		setting = modelSetting.Setting
	}

	return web.Scope(
		h.Div(
			h.Div(h.Text(msgr.Seo)).Class("text-h4 mb-10"),
			// The Customize switch stands alone; the SEO fields below render only
			// while it is on — a v-if template rather than an expansion panel.
			VSwitch().
				Label(msgr.Customize).Color("primary").
				Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "EnabledCustomize"), setting.EnabledCustomize)...).
				Attr("@update:modelValue", "locals.enabledCustomize = $event"),
			h.Template(
				b.vseo(fieldPrefix, seo, &setting, ctx.R),
			).Attr("v-if", "locals.enabledCustomize"),
		).Class("pb-4"),
	).LocalsInit(fmt.Sprintf(`{enabledCustomize: %t}`, setting.EnabledCustomize)).
		Slot("{ locals }")
}

func (b *Builder) vseo(fieldPrefix string, seo *SEO, setting *Setting, req *http.Request) h.HTMLComponent {
	var (
		msgr = i18n.MustGetModuleMessages(req.Context(), I18nSeoKey, Messages_en_US).(*Messages)
		db   = b.db
	)

	var varComps []h.HTMLComponent
	for varName := range seo.getAvailableVars() {
		varComps = append(varComps,
			VChip(
				VIcon("mdi-plus-box").Class("mr-2"),
				h.Text(i18n.PT(req.Context(), I18nSeoKey, "SettingVar", varName)),
			).Variant(VariantText).Attr("@click", fmt.Sprintf("$refs.seo.addTags('%s')", varName)).Label(true).Variant(VariantOutlined),
		)
	}
	var variablesEle []h.HTMLComponent
	variablesEle = append(variablesEle, VChipGroup(varComps...).Column(true).Class("ma-4"))

	image := &setting.OpenGraphImageFromMediaLibrary
	if image.ID.String() == "0" {
		image.ID = json.Number("")
	}
	refPrefix := strings.ReplaceAll(strings.ToLower(fieldPrefix), " ", "_")
	return VSeo(
		h.H4(msgr.Basic).Style("margin-top:15px;font-weight: 500"),
		VRow(
			variablesEle...,
		),
		VCard(
			VCardText(
				VTextField().Variant(FieldVariantUnderlined).Attr("counter", true).Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "Title"), setting.Title)...).Label(msgr.Title).Attr("@focus", fmt.Sprintf("$refs.seo.tagInputsFocus($refs.%s)", fmt.Sprintf("%s_title", refPrefix))).Attr("ref", fmt.Sprintf("%s_title", refPrefix)),
				VTextField().Variant(FieldVariantUnderlined).Attr("counter", true).Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "Description"), setting.Description)...).Label(msgr.Description).Attr("@focus", fmt.Sprintf("$refs.seo.tagInputsFocus($refs.%s)", fmt.Sprintf("%s_description", refPrefix))).Attr("ref", fmt.Sprintf("%s_description", refPrefix)),
				VTextarea().Variant(FieldVariantUnderlined).Attr("counter", true).Rows(2).AutoGrow(true).Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "Keywords"), setting.Keywords)...).Label(msgr.Keywords).Attr("@focus", fmt.Sprintf("$refs.seo.tagInputsFocus($refs.%s)", fmt.Sprintf("%s_keywords", refPrefix))).Attr("ref", fmt.Sprintf("%s_keywords", refPrefix)),
			),
		).Variant(VariantOutlined).Flat(true),

		h.H4(msgr.OpenGraphInformation).Style("margin-top:15px;margin-bottom:15px;font-weight: 500"),
		VCard(
			VCardText(
				VRow(
					VCol(VTextField().Variant(FieldVariantUnderlined).Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "OpenGraphTitle"), setting.OpenGraphTitle)...).Label(msgr.OpenGraphTitle).Attr("@focus", fmt.Sprintf("$refs.seo.tagInputsFocus($refs.%s)", fmt.Sprintf("%s_og_title", refPrefix))).Attr("ref", fmt.Sprintf("%s_og_title", refPrefix))).Cols(6),
					VCol(VTextField().Variant(FieldVariantUnderlined).Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "OpenGraphDescription"), setting.OpenGraphDescription)...).Label(msgr.OpenGraphDescription).Attr("@focus", fmt.Sprintf("$refs.seo.tagInputsFocus($refs.%s)", fmt.Sprintf("%s_og_description", refPrefix))).Attr("ref", fmt.Sprintf("%s_og_description", refPrefix))).Cols(6),
				),
				VRow(
					VCol(VTextField().Variant(FieldVariantUnderlined).Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "OpenGraphURL"), setting.OpenGraphURL)...).Label(msgr.OpenGraphURL).Attr("@focus", fmt.Sprintf("$refs.seo.tagInputsFocus($refs.%s)", fmt.Sprintf("%s_og_url", refPrefix))).Attr("ref", fmt.Sprintf("%s_og_url", refPrefix))).Cols(6),
					VCol(VTextField().Variant(FieldVariantUnderlined).Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "OpenGraphType"), setting.OpenGraphType)...).Label(msgr.OpenGraphType).Attr("@focus", fmt.Sprintf("$refs.seo.tagInputsFocus($refs.%s)", fmt.Sprintf("%s_og_type", refPrefix))).Attr("ref", fmt.Sprintf("%s_og_type", refPrefix))).Cols(6),
				),
				VRow(
					VCol(VTextField().Variant(FieldVariantUnderlined).Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "OpenGraphImageURL"), setting.OpenGraphImageURL)...).Label(msgr.OpenGraphImageURL).Attr("@focus", fmt.Sprintf("$refs.seo.tagInputsFocus($refs.%s)", fmt.Sprintf("%s_og_imageurl", refPrefix))).Attr("ref", fmt.Sprintf("%s_og_imageurl", refPrefix))).Cols(12),
				),
				VRow(
					VCol(media.QMediaBox(db).Label(msgr.OpenGraphImage).
						FieldName(fmt.Sprintf("%s.%s", fieldPrefix, "OpenGraphImageFromMediaLibrary")).
						Value(image).
						Config(&media_library.MediaBoxConfig{
							AllowType: "image",
							Sizes: map[string]*base.Size{
								"og": {
									Width:  1200,
									Height: 630,
								},
								"twitter-large": {
									Width:  1200,
									Height: 600,
								},
								"twitter-small": {
									Width:  630,
									Height: 630,
								},
							},
						})).Cols(12)),
				VRow(
					VCol(VTextarea().Variant(FieldVariantUnderlined).Attr(web.VField(fmt.Sprintf("%s.%s", fieldPrefix, "OpenGraphMetadataString"), GetOpenGraphMetadataString(setting.OpenGraphMetadata))...).Label(msgr.OpenGraphMetadata).Attr("@focus", fmt.Sprintf("$refs.seo.tagInputsFocus($refs.%s)", fmt.Sprintf("%s_og_metadata", refPrefix))).Attr("ref", fmt.Sprintf("%s_og_metadata", refPrefix))).Cols(12),
				),
			),
		).Variant(VariantOutlined).Flat(true),
	).Attr("ref", "seo")
}

func (b *Builder) vseoReadonly(fieldPrefix string, seo *SEO, setting *Setting, req *http.Request) h.HTMLComponent {
	var (
		msgr = i18n.MustGetModuleMessages(req.Context(), I18nSeoKey, Messages_en_US).(*Messages)
		db   = b.db
	)

	var varComps []h.HTMLComponent
	for varName := range seo.getAvailableVars() {
		varComps = append(varComps,
			VChip(
				VIcon("mdi-plus-box").Class("mr-2"),
				h.Text(i18n.PT(req.Context(), I18nSeoKey, "SettingVar", varName)),
			).Variant(VariantText).Attr("@click", fmt.Sprintf("$refs.seo.addTags('%s')", varName)).Label(true).Variant(VariantOutlined),
		)
	}

	image := &setting.OpenGraphImageFromMediaLibrary
	if image.ID.String() == "0" {
		image.ID = json.Number("")
	}
	// Two sections, each a VCard with a title (VCardTitle, so it never wraps like
	// the old fixed-width chip did): "Basic" and "Open Graph information".
	return h.Components(
		VCard(
			VCardTitle(h.Text(msgr.Basic)),
			VCardText(
				presets.FieldComponentContainer(msgr.Title, h.Text(setting.Title)),
				presets.FieldComponentContainer(msgr.Description, h.Text(setting.Description)),
				presets.FieldComponentContainer(msgr.Keywords, h.Text(setting.Keywords)),
			),
		).Variant(VariantOutlined).Class("mb-4"),

		VCard(
			VCardTitle(h.Text(msgr.OpenGraphInformation)),
			VCardText(
				presets.FieldComponentContainer(msgr.OpenGraphTitle, h.Text(setting.OpenGraphTitle)),
				presets.FieldComponentContainer(msgr.OpenGraphDescription, h.Text(setting.OpenGraphDescription)),
				presets.FieldComponentContainer(msgr.OpenGraphURL, h.Text(setting.OpenGraphURL)),
				presets.FieldComponentContainer(msgr.OpenGraphType, h.Text(setting.OpenGraphType)),
				presets.FieldComponentContainer(msgr.OpenGraphImageURL, h.Text(setting.OpenGraphImageURL)),
				presets.FieldComponentContainer(msgr.OpenGraphImage,
					VRow(
						VCol(media.QMediaBox(db).
							Readonly(true).
							FieldName(fmt.Sprintf("%s.%s", fieldPrefix, "OpenGraphImageFromMediaLibrary")).
							Value(image).
							Config(&media_library.MediaBoxConfig{
								AllowType: "image",
								Sizes: map[string]*base.Size{
									"og":            {Width: 1200, Height: 630},
									"twitter-large": {Width: 1200, Height: 600},
									"twitter-small": {Width: 630, Height: 630},
								},
							})).Cols(12))),
				presets.FieldComponentContainer(msgr.OpenGraphMetadata,
					h.Text(GetOpenGraphMetadataString(setting.OpenGraphMetadata))),
			),
		).Variant(VariantOutlined),
	)
}

func (b *Builder) ModelInstall(pb *presets.Builder, mb *presets.ModelBuilder) error {
	b.configDetailing(mb.Detailing())
	return nil
}

func (b *Builder) configDetailing(pd *presets.DetailingBuilder) {
	pd.Section(SeoDetailFieldName).
		Editing("SEO").
		SaveFunc(b.detailSaver).
		ViewComponentFunc(b.detailShowComponent).
		EditComponentFunc(b.EditingComponentFunc)
}

func (b *Builder) detailShowComponent(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	var (
		obj         = field.Obj
		msgr        = i18n.MustGetModuleMessages(ctx.Context(), I18nSeoKey, Messages_en_US).(*Messages)
		fieldPrefix string
		setting     Setting
		db          = b.db
		locale, _   = l10n.IsLocalizableFromContext(ctx.Context())
	)
	seo := b.GetSEO(obj)
	if seo == nil {
		return h.Div()
	}

	value := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < value.NumField(); i++ {
		if s, ok := value.Field(i).Interface().(Setting); ok {
			setting = s
			fieldPrefix = value.Type().Field(i).Name
		}
	}
	if !setting.EnabledCustomize && setting.IsEmpty() {
		modelSetting := &RvqSEOSetting{}
		db.Where("name = ? AND locale_code = ?", seo.name, locale).First(modelSetting)
		setting = modelSetting.Setting
	}

	return h.Div(
		h.Div(h.Text(msgr.Seo)).Class("text-h4 mb-10"),
		b.vseoReadonly(fieldPrefix, seo, &setting, ctx.R),
	).Class("pb-4")
}

func (b *Builder) detailSaver(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
	if err = EditSetterFunc(
		obj,
		&presets.FieldContext{
			ToComponentOptions: &presets.ToComponentOptions{},
			Mode:               presets.FieldModeStack{presets.DETAIL},
			Obj:                obj,
			Name:               SeoDetailFieldName,
		},
		ctx); err != nil {
		return
	}
	if err = b.db.Updates(obj).Error; err != nil {
		return err
	}
	return
}
