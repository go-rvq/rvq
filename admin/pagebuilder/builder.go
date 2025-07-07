package pagebuilder

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/qor5/admin/v3/activity"
	"github.com/qor5/admin/v3/l10n"
	"github.com/qor5/admin/v3/media"
	"github.com/qor5/admin/v3/model"
	"github.com/qor5/admin/v3/note"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/admin/v3/presets/gorm2op"
	"github.com/qor5/admin/v3/publish"
	"github.com/qor5/admin/v3/richeditor"
	"github.com/qor5/admin/v3/seo"
	"github.com/qor5/admin/v3/utils"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/perm"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type RenderInput struct {
	Page        *Page
	IsEditor    bool
	IsReadonly  bool
	Device      string
	ContainerId string
	DisplayName string
}

type RenderFunc func(obj interface{}, input *RenderInput, ctx *web.EventContext) h.HTMLComponent

type PageLayoutFunc func(body h.HTMLComponent, input *PageLayoutInput, ctx *web.EventContext) h.HTMLComponent

type SubPageTitleFunc func(ctx *web.EventContext) string

type PageLayoutInput struct {
	SeoTags           h.HTMLComponent
	CanonicalLink     h.HTMLComponent
	StructuredData    h.HTMLComponent
	FreeStyleCss      []string
	FreeStyleTopJs    []string
	FreeStyleBottomJs []string
	Hreflang          map[string]string
	Header            h.HTMLComponent
	Footer            h.HTMLComponent
	IsEditor          bool
	LocaleCode        string
	EditorCss         []h.HTMLComponent
	IsPreview         bool
}

type Builder struct {
	prefix            string
	wb                *web.Builder
	db                *gorm.DB
	containerBuilders []*ContainerBuilder
	ps                *presets.Builder
	models            []*ModelBuilder
	templateModel     *presets.ModelBuilder
	l10n              *l10n.Builder
	mediaBuilder      *media.Builder
	note              *note.Builder
	ab                *activity.Builder
	publisher         *publish.Builder
	seoBuilder        *seo.Builder
	pageStyle         h.HTMLComponent
	pageLayoutFunc    PageLayoutFunc
	subPageTitleFunc  SubPageTitleFunc
	images            http.Handler
	imagesPrefix      string
	defaultDevice     string
	publishBtnColor   string
	duplicateBtnColor string
	templateEnabled   bool
	expendContainers  bool
	templateInstall   presets.ModelInstallFunc
	pageInstall       presets.ModelInstallFunc
	categoryInstall   presets.ModelInstallFunc
	devices           []Device
}

const (
	openTemplateDialogEvent = "openTemplateDialogEvent"
	selectTemplateEvent     = "selectTemplateEvent"
	// clearTemplateEvent               = "clearTemplateEvent"
	republishRelatedOnlinePagesEvent = "republish_related_online_pages"

	paramOpenFromSharedContainer = "open_from_shared_container"

	PageBuilderPreviewCard = "PageBuilderPreviewCard"
)

func New(prefix string, db *gorm.DB, i18nB *i18n.Builder) *Builder {
	err := db.AutoMigrate(
		&Page{},
		&Template{},
		&Container{},
		&DemoContainer{},
		&Category{},
	)
	if err != nil {
		panic(err)
	}
	// https://github.com/go-gorm/sqlite/blob/64917553e84d5482e252c7a0c8f798fb672d7668/ddlmod.go#L16
	// fxxk: newline is not allowed
	err = db.Exec(`
create unique index if not exists uidx_page_builder_demo_containers_model_name_locale_code on page_builder_demo_containers (model_name, locale_code) where deleted_at is null;
`).Error
	if err != nil {
		panic(err)
	}

	r := &Builder{
		db:                db,
		wb:                web.New(),
		prefix:            prefix,
		defaultDevice:     DeviceComputer,
		publishBtnColor:   "primary",
		duplicateBtnColor: "primary",
		templateEnabled:   true,
		expendContainers:  true,
	}
	r.templateInstall = r.defaultTemplateInstall
	r.categoryInstall = r.defaultCategoryInstall
	r.pageInstall = r.defaultPageInstall

	r.ps = presets.New(i18nB).
		BrandTitle("Page Builder").
		DataOperator(gorm2op.DataOperator(db)).
		URIPrefix(prefix).
		DetailLayoutFunc(r.pageEditorLayout)

	return r
}

func (b *Builder) Prefix(v string) (r *Builder) {
	b.ps.URIPrefix(v)
	b.prefix = v
	return b
}

func (b *Builder) PageStyle(v h.HTMLComponent) (r *Builder) {
	b.pageStyle = v
	return b
}

func (b *Builder) WrapPageInstall(w func(presets.ModelInstallFunc) presets.ModelInstallFunc) (r *Builder) {
	b.pageInstall = w(b.pageInstall)
	return b
}

func (b *Builder) WrapTemplateInstall(w func(presets.ModelInstallFunc) presets.ModelInstallFunc) (r *Builder) {
	b.templateInstall = w(b.templateInstall)
	return b
}

func (b *Builder) WrapCategoryInstall(w func(presets.ModelInstallFunc) presets.ModelInstallFunc) (r *Builder) {
	b.categoryInstall = w(b.categoryInstall)
	return b
}

func (b *Builder) PageLayout(v PageLayoutFunc) (r *Builder) {
	b.pageLayoutFunc = v
	return b
}

func (b *Builder) SubPageTitle(v SubPageTitleFunc) (r *Builder) {
	b.subPageTitleFunc = v
	return b
}

func (b *Builder) L10n(v *l10n.Builder) (r *Builder) {
	b.l10n = v
	return b
}

func (b *Builder) Media(v *media.Builder) (r *Builder) {
	b.mediaBuilder = v
	return b
}

func (b *Builder) Activity(v *activity.Builder) (r *Builder) {
	b.ab = v
	return b
}

func (b *Builder) SEO(v *seo.Builder) (r *Builder) {
	b.seoBuilder = v
	return b
}

func (b *Builder) Note(v *note.Builder) (r *Builder) {
	b.note = v
	return b
}

func (b *Builder) Publisher(v *publish.Builder) (r *Builder) {
	b.publisher = v
	return b
}

func (b *Builder) GetPageTitle() SubPageTitleFunc {
	if b.subPageTitleFunc == nil {
		b.subPageTitleFunc = defaultSubPageTitle
	}
	return b.subPageTitleFunc
}

func (b *Builder) Images(v http.Handler, imagesPrefix string) (r *Builder) {
	b.images = v
	b.imagesPrefix = imagesPrefix
	return b
}

func (b *Builder) DefaultDevice(v string) (r *Builder) {
	b.defaultDevice = v
	return b
}

func (b *Builder) GetPresetsBuilder() (r *presets.Builder) {
	return b.ps
}

func (b *Builder) PublishBtnColor(v string) (r *Builder) {
	b.publishBtnColor = v
	return b
}

func (b *Builder) DuplicateBtnColor(v string) (r *Builder) {
	b.duplicateBtnColor = v
	return b
}

func (b *Builder) TemplateEnabled(v bool) (r *Builder) {
	b.templateEnabled = v
	return b
}

func (b *Builder) ExpendContainers(v bool) (r *Builder) {
	b.expendContainers = v
	return b
}

func (b *Builder) Model(mb *presets.ModelBuilder) (r *ModelBuilder) {
	r = &ModelBuilder{
		mb:      mb,
		editor:  b.ps.Model(mb.NewModel()).URIName(mb.Info().URI() + "/editors"),
		builder: b,
		db:      b.db,
	}
	b.models = append(b.models, r)
	r.setName()
	r.registerFuncs()
	return r
}

func (b *Builder) ModelInstall(pb *presets.Builder, mb *presets.ModelBuilder) (err error) {
	defer b.ps.Build()
	if b.publisher != nil {
		mb.Use(b.publisher)
	}
	if mb.HasDetailing() {
		fb := mb.Detailing().GetField(PageBuilderPreviewCard)
		if fb != nil && fb.GetCompFunc() == nil {
			fb.ComponentFunc(overview(b, nil, mb))
		}
	}
	r := b.Model(mb)
	// register model editors page
	b.installAsset(pb)
	if err = b.configEditor(r); err != nil {
		return
	}
	b.configPublish(r)
	b.pageLayoutFunc = DefaultPageLayoutFunc
	return nil
}

func (b *Builder) installAsset(pb *presets.Builder) {
	pb.I18n().
		RegisterForModule(language.English, I18nPageBuilderKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nPageBuilderKey, Messages_zh_CN).
		RegisterForModule(language.Japanese, I18nPageBuilderKey, Messages_ja_JP)

	pb.ExtraAsset("/redactor.js", "text/javascript", richeditor.JSComponentsPack())
	pb.ExtraAsset("/redactor.css", "text/css", richeditor.CSSComponentsPack())
}

func (b *Builder) Install(pb *presets.Builder) (err error) {
	defer b.ps.Build()
	b.preparePlugins()

	b.installAsset(pb)
	r := b.Model(pb.Model(&Page{}))
	err = b.configEditor(r)
	if err != nil {
		return
	}
	b.configTemplateAndPage(pb, r)
	b.configSharedContainer(pb, r)
	b.configDemoContainer(pb)
	categoryM := pb.Model(&Category{}).URIName("page_categories").Label("Categories")
	err = b.categoryInstall(pb, categoryM)

	return
}

func (b *Builder) configEditor(m *ModelBuilder) (err error) {
	b.useAllPlugin(m.editor)
	md := m.editor.Detailing()
	md.PageFunc(b.Editor(m))
	return
}

func (b *Builder) configPublish(r *ModelBuilder) {
	publisher := b.publisher
	if publisher != nil {
		publisher.ContextValueFuncs(r.ContextValueProvider).Activity(b.ab).AfterInstall(func() {
			r.mb.Editing().SidePanelFunc(nil).ActionsFunc(nil).TabsPanels()
		})
	}
}

func (b *Builder) configTemplateAndPage(pb *presets.Builder, r *ModelBuilder) {
	templateM := presets.NewModelBuilder(pb, &Template{})
	if b.ab != nil {
		templateM.Use(b.ab)
	}
	if b.templateEnabled {
		templateM = pb.Model(&Template{}).URIName("page_templates").Label("Templates")
		err := b.templateInstall(pb, templateM)
		if err != nil {
			panic(err)
		}
		b.templateModel = templateM
	}
	pm := r.mb
	err := b.pageInstall(pb, pm)
	if err != nil {
		panic(err)
	}
	b.configPublish(r)
	b.useAllPlugin(pm)

	// dp.TabsPanels()
}

func (b *Builder) defaultPageInstall(pb *presets.Builder, pm *presets.ModelBuilder) (err error) {
	db := b.db
	lb := pm.Listing("ID", publish.ListingFieldLive, "Title", "Path")
	lb.Field("Path").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		page := field.Obj.(*Page)
		category, err := page.GetCategory(db)
		if err != nil {
			panic(err)
		}
		return h.Td(h.Text(page.getAccessUrl(page.getPublishUrl(b.l10n.GetLocalePath(page.LocaleCode), category.Path))))
	})
	dp := pm.Detailing(PageBuilderPreviewCard)
	// register modelBuilder

	// pm detailing overview
	dp.Field(PageBuilderPreviewCard).ComponentFunc(overview(b, b.templateModel, pm))

	// pm detailing page  detail-field
	detailPageEditor(dp, b.db)
	// pm detailing side panel
	// pm.Detailing().SidePanelFunc(detailingSidePanel(b, pb))

	b.configDetailLayoutFunc(pb, pm, b.templateModel, db)

	if b.templateEnabled {
		pm.RegisterEventFunc(openTemplateDialogEvent, openTemplateDialog(db, b.prefix))
		pm.RegisterEventFunc(selectTemplateEvent, selectTemplate(db))
		// pm.RegisterEventFunc(clearTemplateEvent, clearTemplate(db))
	}

	eb := pm.Editing("TemplateSelection", "Title", "CategoryID", "Slug")
	eb.Validators.AppendFunc(func(obj interface{}, _ presets.FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors) {
		c := obj.(*Page)
		err = pageValidator(ctx.R.Context(), c, db, b.l10n)
		return
	})

	eb.Field("Slug").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		var vErr web.ValidationErrors
		if ve, ok := ctx.Flash.(*web.ValidationErrors); ok {
			vErr = *ve
		}

		return VTextField().
			Variant(FieldVariantUnderlined).
			Attr(web.VField(field.Name, strings.TrimPrefix(field.Value().(string), "/"))...).
			Prefix("/").
			ErrorMessages(vErr.GetFieldErrors("Page.Slug")...)
	}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		m := obj.(*Page)
		m.Slug = path.Join("/", m.Slug)
		return nil
	})
	eb.Field("CategoryID").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		p := field.Obj.(*Page)
		categories := []*Category{}
		locale, _ := l10n.IsLocalizableFromContext(ctx.R.Context())
		if err := db.Model(&Category{}).Where("locale_code = ?", locale).Find(&categories).Error; err != nil {
			panic(err)
		}

		var vErr web.ValidationErrors
		if ve, ok := ctx.Flash.(*web.ValidationErrors); ok {
			vErr = *ve
		}

		msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nPageBuilderKey, Messages_en_US).(*Messages)

		return vx.VXAutocomplete().Label(msgr.Category).
			Attr(web.VField(field.Name, p.CategoryID)...).
			Multiple(false).Chips(false).
			Items(categories).ItemText("Path").ItemValue("ID").
			ErrorMessages(vErr.GetFieldErrors("Page.Category")...)
	})

	eb.Field("TemplateSelection").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		if !b.templateEnabled {
			return nil
		}

		p := field.Obj.(*Page)

		selectedID := ctx.R.FormValue(templateSelectedID)
		body, err := getTplPortalComp(ctx, db, selectedID)
		if err != nil {
			panic(err)
		}

		// Display template selection only when creating a new page
		if p.ID == 0 {
			return h.Div(
				web.Portal().Name(templateSelectPortal),
				web.Portal(
					body,
				).Name(selectedTemplatePortal),
			).Class("my-2").
				Attr(web.VAssign("vars", `{showTemplateDialog: false}`)...)
		}
		return nil
	})

	eb.SaveFunc(func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
		localeCode, _ := l10n.IsLocalizableFromContext(ctx.R.Context())
		p := obj.(*Page)
		if p.Slug != "" {
			p.Slug = path.Clean(p.Slug)
		}
		if p.LocaleCode == "" {
			p.LocaleCode = localeCode
		}
		funcName := ctx.R.FormValue(web.EventFuncIDName)
		if funcName == publish.EventDuplicateVersion {
			// id := ctx.Param(presets.ParamID)
			var fromPage Page
			eb.Fetcher(&fromPage, id, ctx)
			p.SEO = fromPage.SEO
		}

		err = db.Transaction(func(tx *gorm.DB) (inerr error) {
			if inerr = gorm2op.DataOperator(tx).Save(obj, id, ctx); inerr != nil {
				return
			}

			if strings.Contains(ctx.R.RequestURI, publish.EventDuplicateVersion) {
				if inerr = b.copyContainersToNewPageVersion(tx, int(p.ID), p.LocaleCode, p.ParentVersion, p.Version.Version); inerr != nil {
					return
				}
				return
			}

			if v := ctx.R.FormValue(templateSelectedID); v != "" {
				var tplID int
				tplID, inerr = strconv.Atoi(v)
				if inerr != nil {
					return
				}
				if b.l10n == nil {
					localeCode = ""
				}
				if inerr = b.copyContainersToAnotherPage(tx, tplID, templateVersion, localeCode, int(p.ID), p.Version.Version, localeCode); inerr != nil {
					panic(inerr)
				}
			}
			if b.l10n != nil && strings.Contains(ctx.R.RequestURI, l10n.DoLocalize) {
				fromID := ctx.R.Context().Value(l10n.FromID).(string)
				fromVersion := ctx.R.Context().Value(l10n.FromVersion).(string)
				fromLocale := ctx.R.Context().Value(l10n.FromLocale).(string)

				var fromIDInt int
				fromIDInt, err = strconv.Atoi(fromID)
				if err != nil {
					return
				}

				if inerr = b.localizeCategory(tx, p.CategoryID, fromLocale, p.LocaleCode); inerr != nil {
					panic(inerr)
				}

				if inerr = b.localizeContainersToAnotherPage(tx, fromIDInt, fromVersion, fromLocale, int(p.ID), p.Version.Version, p.LocaleCode); inerr != nil {
					panic(inerr)
				}
				return
			}
			return
		})

		return
	})
	return
}

func (b *Builder) useAllPlugin(pm *presets.ModelBuilder) {
	if b.publisher != nil {
		pm.Use(b.publisher)
	}

	if b.ab != nil {
		pm.Use(b.ab)
	}

	if b.seoBuilder != nil {
		pm.Use(b.seoBuilder)
	}

	if b.l10n != nil {
		pm.Use(b.l10n)
	}

	if b.note != nil {
		pm.Use(b.note)
	}
}

func (b *Builder) preparePlugins() {
	l10nB := b.l10n
	// activityB := b.ab
	publisher := b.publisher
	if l10nB != nil {
		l10nB.Activity(b.ab)
	}
	seoBuilder := b.seoBuilder
	if seoBuilder != nil {
		seoBuilder.RegisterSEO("Page", &Page{}).RegisterContextVariable(
			"Title",
			func(object interface{}, _ *seo.Setting, _ *http.Request) string {
				if p, ok := object.(*Page); ok {
					return p.Title
				}
				return ""
			},
		)
	}

	if b.mediaBuilder == nil {
		b.mediaBuilder = media.New(b.db)
	}
	b.ps.Use(b.mediaBuilder, publisher, seoBuilder)
}

func (b *Builder) configDetailLayoutFunc(
	pb *presets.Builder,
	pm *presets.ModelBuilder,
	templateM *presets.ModelBuilder,
	db *gorm.DB,
) {
	oldDetailLayout := pb.GetDetailLayoutFunc()

	// change old detail layout
	pb.DetailLayoutFunc(func(in web.PageFunc, cfg *presets.LayoutConfig) (out web.PageFunc) {
		return func(ctx *web.EventContext) (pr web.PageResponse, err error) {
			if !strings.Contains(ctx.R.RequestURI, "/"+pm.Info().URI()+"/") && templateM != nil && !strings.Contains(ctx.R.RequestURI, "/"+templateM.Info().URIName()+"/") {
				pr, err = oldDetailLayout(in, cfg)(ctx)
				return
			}

			pb.InjectAssets(ctx)
			// call createMenus before in(ctx) to fill the menuGroupName for modelBuilders first
			menu := pb.CreateMenus(ctx)
			var profile h.HTMLComponent
			if pb.GetProfileFunc() != nil {
				profile = pb.GetProfileFunc()(ctx)
			}
			utilsMsgr := i18n.MustGetModuleMessages(ctx.Context(), utils.I18nUtilsKey, utils.Messages_en_US).(*utils.Messages)
			pvMsgr := i18n.MustGetModuleMessages(ctx.Context(), publish.I18nPublishKey, publish.Messages_en_US).(*publish.Messages)
			id := ctx.Param(presets.ParamID)

			if id == "" {
				return pb.DefaultNotFoundPageFunc(ctx)
			}
			var isTemplate bool
			isPage := strings.Contains(ctx.R.RequestURI, "/"+pm.Info().URI()+"/")
			if templateM != nil {
				isTemplate = strings.Contains(ctx.R.RequestURI, "/"+templateM.Info().URI()+"/")
			}

			if isTemplate {
				ctx.R.Form.Set(paramsTpl, "1")
			}
			var obj interface{}
			var dmb *presets.ModelBuilder
			if isPage {
				dmb = pm
			} else {
				dmb = templateM
			}
			obj = dmb.NewModel()
			if sd, ok := obj.(presets.SlugDecoder); ok {
				vs, err := presets.RecoverPrimaryColumnValuesBySlug(sd, id)
				if err != nil {
					return pb.DefaultNotFoundPageFunc(ctx)
				}
				if _, err := strconv.Atoi(vs["id"]); err != nil {
					return pb.DefaultNotFoundPageFunc(ctx)
				}
			} else {
				if _, err := strconv.Atoi(id); err != nil {
					return pb.DefaultNotFoundPageFunc(ctx)
				}
			}
			err = dmb.Detailing().GetFetchFunc()(obj, id, ctx)
			if err != nil {
				if errors.Is(err, presets.ErrRecordNotFound) {
					return pb.DefaultNotFoundPageFunc(ctx)
				}
				return
			}

			if l, ok := obj.(l10n.LocaleInterface); ok {
				locale := ctx.R.FormValue("locale")
				if ctx.R.FormValue(web.EventFuncIDName) == "__reload__" && locale != "" && locale != l.EmbedLocale().LocaleCode {
					// redirect to list page when change locale
					http.Redirect(ctx.W, ctx.R, dmb.Info().ListingHref(presets.ParentsModelID(ctx.R)...), http.StatusSeeOther)
					return
				}
			}
			var (
				tabContent   web.PageResponse
				versionBadge *VChipBuilder
			)

			if v, ok := obj.(PrimarySlugInterface); ok {
				ps := v.PrimaryColumnValuesBySlug(id)
				versionBadge = VChip(h.Text(fmt.Sprintf("%d versions", versionCount(b.db, pm.NewModel(), ps["id"], ps["localCode"])))).Color(ColorPrimary).Size(SizeSmall).Class("px-1 mx-1").Attr("style", "height:20px")
			}
			tabContent, err = in(ctx)
			if errors.Is(err, perm.PermissionDenied) {
				pr.Body = h.Text(MustGetMessages(ctx.Context()).ErrPermissionDenied.Error())
				return pr, nil
			}
			if err != nil {
				panic(err)
			}

			primaryID := ""
			if v, ok := obj.(presets.SlugEncoder); ok {
				primaryID = v.PrimarySlug()
			}
			var pageAppbarContent []h.HTMLComponent
			pageAppbarContent = append(pageAppbarContent,
				VProgressLinear().
					Attr(":active", "isFetching").
					Class("ml-4").
					Attr("style", "position: fixed; z-index: 99;").
					Indeterminate(true).
					Height(2).
					Color(pb.GetProgressBarColor()),
			)
			pageAppbarContent = h.Components(
				VAppBarNavIcon().
					Density(DensityCompact).
					Class("mr-2").
					Attr("v-if", "!vars.navDrawer").
					On("click.stop", "vars.navDrawer = !vars.navDrawer"),
				h.Div(
					VToolbarTitle(
						b.GetPageTitle()(ctx),
					),
				).Class("mr-auto"),
				VSpacer(),
				publish.DefaultVersionBar(db)(obj, ctx),
			)

			toolbar := VContainer(
				VRow(
					VCol(pb.RunBrandFunc(ctx)).Cols(8),
					VCol(
						pb.RunSwitchLanguageFunc(ctx),
					).Cols(2),

					VCol(
						VAppBarNavIcon().Icon("mdi-menu").
							Density(DensityCompact).
							Attr("@click", "vars.navDrawer = !vars.navDrawer"),
					).Cols(2),
				).Attr("align", "center").Attr("justify", "center"),
			)
			pr.Body = web.Scope(
				VApp(
					VNavigationDrawer(
						VLayout(
							VMain(
								toolbar,
								VCard(
									menu,
								).Class("ma-4").Variant(VariantText),
							),
							VAppBar(
								profile,
							).Location("bottom").Class("border-t-sm border-b-0").Elevation(0),
						).Class("ma-2 border-sm rounded-lg elevation-1").Attr("style", "height: calc(100% - 16px);"),
					).Width(320).
						ModelValue(true).
						Attr("v-model", "vars.navDrawer").
						Permanent(true).
						Floating(true).
						Elevation(0),

					VAppBar(
						h.Div(
							pageAppbarContent...,
						).Class("d-flex align-center  justify-space-between   border-b w-100").Style("height: 48px"),
					).
						Elevation(0).
						Density("compact").Class("px-6"),
					web.Portal().Name(actions.RightDrawer.PortalName()),
					web.Portal().Name(actions.Dialog.PortalName()),
					web.Portal().Name(presets.DeleteConfirmPortalName),
					web.Portal().Name(presets.DefaultConfirmDialogPortalName),
					web.Portal().Name(presets.ListingDialogPortalName),
					web.Portal().Name(dialogPortalName),
					utils.ConfirmDialog(pvMsgr.Areyousure, web.Plaid().EventFunc(web.Var("locals.action")).
						Query(presets.ParamID, primaryID).Go(),
						utilsMsgr),
					VProgressLinear().
						Attr(":active", "isFetching").
						Attr("style", "position: fixed; z-index: 99").
						Indeterminate(true).
						Height(2).
						Color(pb.GetProgressBarColor()),
					h.Template(
						VSnackbar(
							h.Text("{{vars.presetsMessage.message}}"),
							web.Slot(
								VBtn("").
									Icon("mdi-close").
									Attr("@click", `vars.presetsMessage.show = false`),
							).Name("actions"),
						).
							Attr("v-model", "vars.presetsMessage.show").
							Attr(":color", "vars.presetsMessage.color").
							Timeout(15000).
							ZIndex(1000000).
							Location(LocationTop),
					).Attr("v-if", "vars.presetsMessage"),
					VMain(
						VLayout(
							VMain(
								VContainer(
									VBtn("").Size(SizeXSmall).Icon("mdi-arrow-left").Tile(true).Variant(VariantOutlined).Attr("@click",
										web.GET().URL(pm.Info().ListingHref(presets.ParentsModelID(ctx.R)...)).PushState(true).Go(),
									),
									h.H1("{{vars.pageTitle}}").Class("ml-4"),
									versionBadge.Class("mt-2 ml-2"),
								).Class("d-inline-flex align-center"),
								tabContent.Body.(h.HTMLComponent),
							),
						).Class("pa-6"),
					),
				).Attr("id", "vt-app").Attr(web.VAssign("vars", `{presetsRightDrawer: false, presetsDialog: false, dialogPortalName: false, presetsListingDialog: false, presetsMessage: {show: false, color: "success", message: ""}}`)...),
			).Slot(" { locals } ").LocalsInit(fmt.Sprintf(`{action: "", commonConfirmDialog: false }`))
			return
		}
	})
	return
}

func versionCount(db *gorm.DB, obj interface{}, id string, localCode string) (count int64) {
	db.Model(obj).Where("id = ? and locale_code = ?", id, localCode).Count(&count)
	return
}

func scheduleCount(db *gorm.DB, p *Page) (count int64) {
	db.Model(&Page{}).Where("id = ? and version != ? and status = ? and (scheduled_start_at is not null or scheduled_end_at is not null)", p.ID, p.Version.Version, publish.StatusDraft).Count(&count)
	return
}

// cats should be ordered by path
func fillCategoryIndentLevels(cats []*Category) {
	for i, cat := range cats {
		if cat.Path == "/" {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			if strings.HasPrefix(cat.Path, cats[j].Path+"/") {
				cat.IndentLevel = cats[j].IndentLevel + 1
				break
			}
		}
	}
}

func (b *Builder) defaultCategoryInstall(pb *presets.Builder, pm *presets.ModelBuilder) (err error) {
	db := b.db

	lb := pm.Listing("Name", "Path", "Description")

	lb.WrapSearchFunc(func(in presets.SearchFunc) presets.SearchFunc {
		return func(model interface{}, params *presets.SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error) {
			r, totalCount, err = in(model, params, ctx)
			cats := r.([]*Category)
			sort.Slice(cats, func(i, j int) bool {
				return cats[i].Path < cats[j].Path
			})
			fillCategoryIndentLevels(cats)
			return
		}
	})

	lb.Field("Name").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		cat := field.Obj.(*Category)

		icon := "mdi-folder"
		if cat.IndentLevel != 0 {
			icon = "mdi-file"
		}

		return h.Td(
			h.Div(
				VIcon(icon).Size(SizeSmall).Class("mb-1"),
				h.Text(cat.Name),
			).Style(fmt.Sprintf("padding-left: %dpx;", cat.IndentLevel*32)),
		)
	})

	eb := pm.Editing("Name", "Path", "Description")

	eb.Field("Path").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		var vErr web.ValidationErrors
		if ve, ok := ctx.Flash.(*web.ValidationErrors); ok {
			vErr = *ve
		}

		return VTextField().Label("Path").
			Variant(FieldVariantUnderlined).
			Attr(web.VField("Path", strings.TrimPrefix(field.Value().(string), "/"))...).
			Prefix("/").
			ErrorMessages(vErr.GetFieldErrors("Category.Category")...)
	}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		m := obj.(*Category)
		m.Path = path.Join("/", m.Path)
		return nil
	})

	lb.DeleteFunc(func(obj interface{}, id string, ctx *web.EventContext) (err error) {
		cs := obj.(presets.SlugDecoder).PrimaryColumnValuesBySlug(id)
		ID := cs["id"]
		Locale := cs[l10n.SlugLocaleCode]

		var count int64
		if err = db.Model(&Page{}).Where("category_id = ? AND locale_code = ?", ID, Locale).Count(&count).Error; err != nil {
			return
		}
		if count > 0 {
			err = errors.New(unableDeleteCategoryMsg)
			return
		}
		if err = db.Model(&Category{}).Where("id = ? AND locale_code = ?", ID, Locale).Delete(&Category{}).Error; err != nil {
			return
		}
		return
	})

	eb.Validators.AppendFunc(func(obj interface{}, _ presets.FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors) {
		c := obj.(*Category)
		err = categoryValidator(c, db, b.l10n)
		return
	})

	eb.SaveFunc(func(obj interface{}, id string, ctx *web.EventContext) (err error) {
		c := obj.(*Category)
		c.Path = path.Clean(c.Path)
		err = db.Save(c).Error
		return
	})
	if b.ab != nil {
		pm.Use(b.ab)
	}
	if b.l10n != nil {
		pm.Use(b.l10n)
	}

	return
}

const (
	templateSelectPortal   = "templateSelectPortal"
	selectedTemplatePortal = "selectedTemplatePortal"

	templateSelectedID = "TemplateSelectedID"
	templateID         = "TemplateID"
	templateBlankVal   = "blank"
)

func selectTemplate(db *gorm.DB) web.EventFunc {
	return func(ctx *web.EventContext) (er web.EventResponse, err error) {
		defer func() {
			web.AppendRunScripts(&er, "vars.showTemplateDialog=false")
		}()

		msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nPageBuilderKey, Messages_en_US).(*Messages)

		id := ctx.R.FormValue(templateID)
		if id == templateBlankVal {
			er.UpdatePortal(selectedTemplatePortal,
				VRow(
					VCol(
						h.Input("").Type("hidden").Attr(web.VField(templateSelectedID, "")...),
						VTextField().Readonly(true).Label(msgr.SelectedTemplateLabel).ModelValue(msgr.Blank).Density(DensityCompact).Variant(VariantOutlined),
					).Cols(5),
					VCol(
						VBtn(msgr.ChangeTemplate).Color(ColorPrimary).
							Attr("@click", web.Plaid().Query(templateSelectedID, "").EventFunc(openTemplateDialogEvent).Go()),
					).Cols(5),
				),
			)
			return
		}

		var body h.HTMLComponent
		if body, err = getTplPortalComp(ctx, db, id); err != nil {
			return
		}

		er.UpdatePortal(selectedTemplatePortal, body)

		return
	}
}

func getTplPortalComp(ctx *web.EventContext, db *gorm.DB, selectedID string) (h.HTMLComponent, error) {
	msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nPageBuilderKey, Messages_en_US).(*Messages)
	locale, _ := l10n.IsLocalizableFromContext(ctx.R.Context())

	name := msgr.Blank
	if selectedID != "" {
		if err := db.Model(&Template{}).Where("id = ? AND locale_code = ?", selectedID, locale).Pluck("name", &name).Error; err != nil {
			return nil, err
		}
	}

	return VRow(
		VCol(
			h.Input("").Type("hidden").Attr(web.VField(templateSelectedID, selectedID)...),
			VTextField().Readonly(true).Label(msgr.SelectedTemplateLabel).ModelValue(name).Density(DensityCompact).Variant(VariantOutlined),
		).Cols(5),
		VCol(
			VBtn(msgr.ChangeTemplate).Color(ColorPrimary).
				Attr("@click", web.Plaid().Query(templateSelectedID, selectedID).EventFunc(openTemplateDialogEvent).Go()),
		).Cols(5),
	), nil
}

// Unused
func clearTemplate(_ *gorm.DB) web.EventFunc {
	return func(ctx *web.EventContext) (er web.EventResponse, err error) {
		msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nPageBuilderKey, Messages_en_US).(*Messages)
		er.UpdatePortal(selectedTemplatePortal,
			VRow(
				VCol(
					h.Input("").Type("hidden").Attr(web.VField(templateSelectedID, "")...),
					VTextField().Readonly(true).Label(msgr.SelectedTemplateLabel).ModelValue(msgr.Blank).Density(DensityCompact).Variant(VariantOutlined),
				).Cols(5),
				VCol(
					VBtn(msgr.ChangeTemplate).Color(ColorPrimary).
						Attr("@click", web.Plaid().Query(templateSelectedID, "").EventFunc(openTemplateDialogEvent).Go()),
				).Cols(5),
			),
		)
		return
	}
}

func openTemplateDialog(db *gorm.DB, prefix string) web.EventFunc {
	return func(ctx *web.EventContext) (er web.EventResponse, err error) {
		gmsgr := presets.MustGetMessages(ctx.Context())
		locale, _ := l10n.IsLocalizableFromContext(ctx.R.Context())
		selectedID := ctx.R.FormValue(templateSelectedID)
		if selectedID == "" {
			selectedID = templateBlankVal
		}

		tpls := []*Template{}
		if err := db.Model(&Template{}).Where("locale_code = ?", locale).Find(&tpls).Error; err != nil {
			panic(err)
		}
		msgrPb := i18n.MustGetModuleMessages(ctx.Context(), I18nPageBuilderKey, Messages_en_US).(*Messages)

		var tplHTMLComponents []h.HTMLComponent
		tplHTMLComponents = append(tplHTMLComponents,
			VCol(
				VCard(
					h.Div(
						h.Iframe().
							Attr("width", "100%", "height", "150", "frameborder", "no").
							Style("transform-origin: left top; transform: scale(1, 1); pointer-events: none;"),
					),
					VCardTitle(h.Text(msgrPb.Blank)),
					VCardSubtitle(h.Text("")),
					h.Div(
						h.Input(templateID).Type("radio").
							Value(templateBlankVal).
							Attr(web.VField(templateID, selectedID)...).
							AttrIf("checked", "checked", selectedID == "").
							Style("width: 18px; height: 18px"),
					).Class("mr-4 float-right"),
				).Height(280).Class("text-truncate").Variant(VariantOutlined),
			).Cols(3),
		)
		for _, tpl := range tpls {
			tplHTMLComponents = append(tplHTMLComponents,
				getTplColComponent(ctx, prefix, tpl, selectedID),
			)
		}
		if len(tpls) == 0 {
			tplHTMLComponents = append(tplHTMLComponents,
				h.Div(h.Text(gmsgr.ListingNoRecordToShow)).Class("pl-4 text-center grey--text text--darken-2"),
			)
		}

		er.UpdatePortal(templateSelectPortal,
			VDialog(
				VCard(
					VCardTitle(
						h.Text(msgrPb.CreateFromTemplate),
						VSpacer(),
						VBtn("").Icon("mdi-close").
							Size(SizeLarge).
							On("click", fmt.Sprintf("vars.showTemplateDialog=false")),
					),
					VCardActions(
						VRow(tplHTMLComponents...),
					),
					VCardActions(
						VSpacer(),
						VBtn(gmsgr.Cancel).Attr("@click", "vars.showTemplateDialog=false"),
						VBtn(gmsgr.OK).Color(ColorPrimary).
							Attr("@click", web.Plaid().EventFunc(selectTemplateEvent).
								Go(),
							),
					).Class("pb-4"),
				).Tile(true),
			).MaxWidth("80%").
				Attr(web.VAssign("vars", `{showTemplateDialog: false}`)...).
				Attr("v-model", fmt.Sprintf("vars.showTemplateDialog")),
		)
		er.RunScript = `setTimeout(function(){ vars.showTemplateDialog = true }, 100)`
		return
	}
}

func getTplColComponent(ctx *web.EventContext, prefix string, tpl *Template, selectedID string) h.HTMLComponent {
	msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nPageBuilderKey, Messages_en_US).(*Messages)

	// Avoid layout errors
	var name string
	var desc string
	if tpl.Name == "" {
		name = msgr.Unnamed
	} else {
		name = tpl.Name
	}
	if tpl.Description == "" {
		desc = msgr.NotDescribed
	} else {
		desc = tpl.Description
	}
	id := fmt.Sprintf("%d", tpl.ID)

	src := fmt.Sprintf("./%s/preview?id=%d&tpl=1&locale=%s", prefix, tpl.ID, tpl.LocaleCode)

	return VCol(
		VCard(
			h.Div(
				h.Iframe().Src(src).
					Attr("width", "100%", "height", "150", "frameborder", "no").
					Style("transform-origin: left top; transform: scale(1, 1); pointer-events: none;"),
			),
			VCardTitle(h.Text(name)),
			VCardSubtitle(h.Text(desc)),
			VBtn(msgr.Preview).Variant(VariantText).Size(SizeSmall).Class("ml-2 mb-4").
				Href(src).
				Attr("target", "_blank").Color(ColorPrimary),
			h.Div(
				h.Input(templateID).Type("radio").
					Value(id).
					Attr(web.VField(templateID, selectedID)...).
					AttrIf("checked", "checked", selectedID == id).
					Style("width: 18px; height: 18px"),
			).Class("mr-4 float-right"),
		).Height(280).Class("text-truncate").Variant(VariantOutlined),
	).Cols(3)
}

func (b *Builder) configSharedContainer(pb *presets.Builder, r *ModelBuilder) {
	db := b.db

	pm := pb.Model(&Container{}).URIName("shared_containers").Label("Shared Containers")

	pm.RegisterEventFunc(republishRelatedOnlinePagesEvent, republishRelatedOnlinePages(r.mb.Info().ListingHrefCtx))

	listing := pm.Listing("DisplayName").SearchColumns("display_name")
	listing.RowMenu("Rename").RowMenuItem("Rename").ComponentFunc(func(rctx *presets.RecordMenuItemContext) h.HTMLComponent {
		var (
			ctx = rctx.Ctx
			obj = rctx.Obj
		)
		c := obj.(*Container)
		msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nPageBuilderKey, Messages_en_US).(*Messages)
		return VListItem().PrependIcon("mdi-pencil-outline").Title(msgr.Rename).Attr("@click",
			web.Plaid().
				URL(b.ContainerByName(c.ModelName).mb.Info().ListingHref(presets.ParentsModelID(ctx.R)...)).
				EventFunc(RenameContainerDialogEvent).
				Query(paramContainerID, c.PrimarySlug()).
				Query(paramContainerName, c.DisplayName).
				Query("portal", "presets").
				Go(),
		)
	})
	if permB := pb.GetPermission(); permB != nil {
		permB.CreatePolicies(
			perm.PolicyFor(perm.Anybody).WhoAre(perm.Denied).ToDo(presets.PermCreate).On("*:shared_containers:*"),
		)
	}
	listing.Field("DisplayName").Label("Name")
	listing.SearchFunc(sharedContainerSearcher(db, r))
	listing.CellWrapperFunc(func(cell h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent {
		tdbind := cell
		c := obj.(*Container)

		tdbind.SetAttr("@click.self",
			web.Plaid().
				EventFunc(actions.Edit).
				URL(b.ContainerByName(c.ModelName).GetModelBuilder().Info().ListingHrefCtx(ctx)).
				Query(presets.ParamID, c.ModelID).
				Query(paramOpenFromSharedContainer, 1).
				Go()+fmt.Sprintf(`; vars.currEditingListItemID="%s-%d"`, dataTableID, c.ModelID))

		return tdbind
	})

	if b.ab != nil {
		pm.Use(b.ab)
	}
	if b.l10n != nil {
		pm.Use(b.l10n)
	}
	return
}

func (b *Builder) configDemoContainer(pb *presets.Builder) (pm *presets.ModelBuilder) {
	db := b.db

	pm = pb.Model(&DemoContainer{}).URIName("demo_containers").Label("Demo Containers")

	pm.RegisterEventFunc("addDemoContainer", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		modelID := ctx.ParamAsInt(presets.ParamOverlayUpdateID)
		modelName := ctx.R.FormValue("ModelName")
		locale, _ := l10n.IsLocalizableFromContext(ctx.R.Context())
		var existID uint
		{
			m := DemoContainer{}
			db.Where("model_name = ?", modelName).First(&m)
			existID = m.ID
		}
		db.Assign(DemoContainer{
			Model: gorm.Model{
				ID: existID,
			},
			ModelID: uint(modelID),
		}).FirstOrCreate(&DemoContainer{}, map[string]interface{}{
			"model_name":  modelName,
			"locale_code": locale,
		})
		r.Reload = true
		return
	})
	listing := pm.Listing("ModelName").SearchColumns("ModelName")
	listing.Field("ModelName").Label("Name")
	ed := pm.Editing("SelectContainer").ActionsFunc(func(obj interface{}, ctx *web.EventContext) h.HTMLComponent { return nil })
	ed.Field("ModelName")
	ed.Field("ModelID")
	ed.Field("SelectContainer").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		locale, localizable := l10n.IsLocalizableFromContext(ctx.R.Context())

		var demoContainers []DemoContainer
		db.Find(&demoContainers)

		var containers []h.HTMLComponent
		for _, builder := range b.containerBuilders {
			cover := builder.cover
			if cover == "" {
				cover = path.Join(b.prefix, b.imagesPrefix, strings.ReplaceAll(builder.name, " ", "")+".png")
			}

			c := VCol(
				VCard(
					VImg().Src(cover).Height(200),
					VCardActions(
						VCardTitle(h.Text(builder.name)),
						VSpacer(),
						VBtn("Select").
							Variant(VariantText).
							Color(ColorPrimary).Attr("@click",
							web.Plaid().
								EventFunc(actions.New).
								URL(builder.GetModelBuilder().Info().ListingHref(presets.ParentsModelID(ctx.R)...)).
								Query(presets.ParamPostChangeCallback, web.CallbackScript(
									web.POST().
										Query("ModelName", builder.name).
										EventFunc("addDemoContainer").
										AnyEvent().
										Go()).Encode(),
								).
								Go()),
					),
				),
			).Cols(6)

			var isExists bool
			var modelID uint
			for _, dc := range demoContainers {
				if dc.ModelName == builder.name {
					if localizable && dc.LocaleCode != locale {
						continue
					}
					isExists = true
					modelID = dc.ModelID
					break
				}
			}
			if isExists {
				c = VCol(
					VCard(
						VImg().Src(cover).Height(200),
						VCardActions(
							VCardTitle(h.Text(builder.name)),
							VSpacer(),
							VBtn("Edit").
								Variant(VariantText).
								Color(ColorPrimary).Attr("@click",
								web.Plaid().
									EventFunc(actions.Edit).
									URL(builder.GetModelBuilder().Info().ListingHref(presets.ParentsModelID(ctx.R)...)).
									Query(presets.ParamID, fmt.Sprint(modelID)).
									Go()),
						),
					),
				).Cols(6)
			}

			containers = append(containers, c)
		}
		return VSheet(
			VContainer(
				VRow(
					containers...,
				),
			),
		)
	})

	listing.CellWrapperFunc(func(cell h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent {
		tdbind := cell
		c := obj.(*DemoContainer)

		tdbind.SetAttr("@click.self",
			web.Plaid().
				EventFunc(actions.Edit).
				URL(b.ContainerByName(c.ModelName).GetModelBuilder().Info().ListingHrefCtx(ctx)).
				Query(presets.ParamID, c.ModelID).
				Go()+fmt.Sprintf(`; vars.currEditingListItemID="%s-%d"`, dataTableID, c.ModelID))

		return tdbind
	})

	ed.SaveFunc(func(obj interface{}, id string, ctx *web.EventContext) (err error) {
		this := obj.(*DemoContainer)
		err = db.Transaction(func(tx *gorm.DB) (inerr error) {
			if b.l10n != nil && strings.Contains(ctx.R.RequestURI, l10n.DoLocalize) {
				if inerr = b.createModelAfterLocalizeDemoContainer(tx, this); inerr != nil {
					panic(inerr)
				}
			}

			if inerr = gorm2op.DataOperator(tx).Save(this, id, ctx); inerr != nil {
				return
			}
			return
		})

		return
	})

	if b.ab != nil {
		pm.Use(b.ab)
	}
	if b.l10n != nil {
		pm.Use(b.l10n)
	}
	return
}

func (b *Builder) defaultTemplateInstall(pb *presets.Builder, pm *presets.ModelBuilder) (err error) {
	db := b.db

	pm.Listing("ID", "Name", "Description")

	dp := pm.Detailing("Overview")
	dp.Field("Overview").ComponentFunc(templateSettings(db, pm))

	eb := pm.Editing("Name", "Description")

	eb.SaveFunc(func(obj interface{}, id string, ctx *web.EventContext) (err error) {
		this := obj.(*Template)
		err = db.Transaction(func(tx *gorm.DB) (inerr error) {
			if inerr = gorm2op.DataOperator(tx).Save(obj, id, ctx); inerr != nil {
				return
			}

			if b.l10n != nil && strings.Contains(ctx.R.RequestURI, l10n.DoLocalize) {
				fromID := ctx.R.Context().Value(l10n.FromID).(string)
				fromLocale := ctx.R.Context().Value(l10n.FromLocale).(string)

				var fromIDInt int
				fromIDInt, err = strconv.Atoi(fromID)
				if err != nil {
					return
				}

				if inerr = b.localizeContainersToAnotherPage(tx, fromIDInt, "tpl", fromLocale, int(this.ID), "tpl", this.LocaleCode); inerr != nil {
					panic(inerr)
				}
				return
			}
			return
		})

		return
	})

	return
}

func sharedContainerSearcher(db *gorm.DB, b *ModelBuilder) presets.SearchFunc {
	return func(obj interface{}, params *presets.SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error) {
		ilike := "ILIKE"
		if db.Dialector.Name() == "sqlite" {
			ilike = "LIKE"
		}

		wh := db.Model(obj)
		if len(params.KeywordColumns) > 0 && len(params.Keyword) > 0 {
			var segs []string
			var args []interface{}
			for _, c := range params.KeywordColumns {
				segs = append(segs, fmt.Sprintf("%s %s ?", c, ilike))
				args = append(args, fmt.Sprintf("%%%s%%", params.Keyword))
			}
			wh = wh.Where(strings.Join(segs, " OR "), args...)
		}

		for _, cond := range params.SQLConditions {
			wh = wh.Where(strings.Replace(cond.Query, " ILIKE ", " "+ilike+" ", -1), cond.Args...)
		}

		locale, _ := l10n.IsLocalizableFromContext(ctx.R.Context())
		var c int64
		if err = wh.Select("count(display_name)").Where("shared = true AND locale_code = ? and page_model_name = ? ", locale, b.name).Group("display_name, model_name, model_id, locale_code").Count(&c).Error; err != nil {
			return
		}
		totalCount = int(c)

		if params.PerPage > 0 {
			wh = wh.Limit(int(params.PerPage))
			page := params.Page
			if page == 0 {
				page = 1
			}
			offset := (page - 1) * params.PerPage
			wh = wh.Offset(int(offset))
		}

		if err = wh.Select("MIN(id) AS id, display_name, model_name, model_id, locale_code").Find(obj).Error; err != nil {
			return
		}
		r = reflect.ValueOf(obj).Elem().Interface()
		return
	}
}

func (b *Builder) ContainerByName(name string) (r *ContainerBuilder) {
	for _, cb := range b.containerBuilders {
		if cb.name == name {
			return cb
		}
	}
	panic(fmt.Sprintf("No container: %s", name))
}

type ContainerBuilder struct {
	builder      *Builder
	name         string
	mb           *presets.ModelBuilder
	modelBuilder *presets.ModelBuilder
	model        interface{}
	modelType    reflect.Type
	renderFunc   RenderFunc
	cover        string
	group        string
}

func (b *Builder) RegisterContainer(name string) (r *ContainerBuilder) {
	r = &ContainerBuilder{
		name:    name,
		builder: b,
	}
	b.containerBuilders = append(b.containerBuilders, r)
	return
}

func (b *Builder) RegisterModelContainer(name string, mb *presets.ModelBuilder) (r *ContainerBuilder) {
	r = &ContainerBuilder{
		name:         name,
		builder:      b,
		modelBuilder: mb,
	}
	b.containerBuilders = append(b.containerBuilders, r)
	return
}

func (b *ContainerBuilder) Model(m interface{}) *ContainerBuilder {
	b.model = m
	b.mb = b.builder.ps.Model(m)

	val := reflect.ValueOf(m)
	if val.Kind() != reflect.Ptr {
		panic("model pointer type required")
	}

	b.modelType = val.Elem().Type()

	b.configureRelatedOnlinePagesTab()
	b.registerEventFuncs()
	b.uRIName(inflection.Plural(strcase.ToKebab(b.name)))
	return b
}

func (b *ContainerBuilder) uRIName(uri string) *ContainerBuilder {
	if b.mb == nil {
		return b
	}
	b.mb.URIName(uri)
	return b
}

func (b *ContainerBuilder) GetModelBuilder() *presets.ModelBuilder {
	return b.mb
}

func (b *ContainerBuilder) RenderFunc(v RenderFunc) *ContainerBuilder {
	b.renderFunc = v
	return b
}

func (b *ContainerBuilder) Cover(v string) *ContainerBuilder {
	b.cover = v
	return b
}

func (b *ContainerBuilder) Group(v string) *ContainerBuilder {
	b.group = v
	return b
}

func (b *ContainerBuilder) NewModel() interface{} {
	return reflect.New(b.modelType).Interface()
}

func (b *ContainerBuilder) ModelTypeName() string {
	return b.modelType.String()
}

func (b *ContainerBuilder) Editing(vs ...interface{}) *presets.EditingBuilder {
	return b.mb.Editing(vs...)
}

func (b *ContainerBuilder) configureRelatedOnlinePagesTab() {
	eb := b.mb.Editing()
	eb.OnChangeActionFunc(func(id string, ctx *web.EventContext) (s string) {
		return web.Plaid().URL(ctx.R.URL.Path).EventFunc(AutoSaveContainerEvent).Query(presets.ParamID, id).Go()
	})
	eb.AppendTabsPanelFunc(func(obj interface{}, ctx *web.EventContext) (tab h.HTMLComponent, content h.HTMLComponent) {
		if ctx.R.FormValue(paramOpenFromSharedContainer) != "1" {
			return nil, nil
		}

		pmsgr := i18n.MustGetModuleMessages(ctx.Context(), presets.CoreI18nModuleKey, Messages_en_US).(*presets.Messages)
		msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nPageBuilderKey, Messages_en_US).(*Messages)

		id, err := reflectutils.Get(obj, "id")
		if err != nil {
			panic(err)
		}

		pages := []*Page{}
		pageTable := (&Page{}).TableName()
		containerTable := (&Container{}).TableName()
		err = b.builder.db.Model(&Page{}).
			Joins(fmt.Sprintf(`inner join %s on 
        %s.id = %s.page_id
        and %s.version = %s.page_version
        and %s.locale_code = %s.locale_code`,
				containerTable,
				pageTable, containerTable,
				pageTable, containerTable,
				pageTable, containerTable,
			)).
			// FIXME: add container locale condition after container supports l10n
			Where(fmt.Sprintf(`%s.status = ? and %s.model_id = ? and %s.model_name = ?`,
				pageTable,
				containerTable,
				containerTable,
			), publish.StatusOnline, id, b.name).
			Group(fmt.Sprintf(`%s.id,%s.version,%s.locale_code`, pageTable, pageTable, pageTable)).
			Find(&pages).
			Error
		if err != nil {
			panic(err)
		}

		var pageIDs []string
		var pageListComps h.HTMLComponents
		for _, p := range pages {
			pid := p.PrimarySlug()
			pageIDs = append(pageIDs, pid)
			statusVar := fmt.Sprintf(`republish_status_%s`, strings.Replace(pid, "-", "_", -1))
			pageListComps = append(pageListComps, web.Scope(
				VListItem(
					h.Text(fmt.Sprintf("%s (%s)", p.Title, pid)),
					VSpacer(),
					VIcon(fmt.Sprintf(`{{itemLocals.%s}}`, statusVar)),
				).
					Density(DensityCompact),
			).Slot(" { locals : itemLocals }").LocalsInit(fmt.Sprintf(`{%s: ""}`, statusVar)),
			)

			tab = VTab(h.Text(msgr.RelatedOnlinePages))
			content = VWindowItem(
				h.If(len(pages) > 0,
					VList(pageListComps),
					h.Div(
						VSpacer(),
						VBtn(msgr.RepublishAllRelatedOnlinePages).
							Color(ColorPrimary).
							Attr("@click",
								web.Plaid().
									EventFunc(presets.EventOpenConfirmDialog).
									Query(presets.ConfirmDialogConfirmEvent,
										web.Plaid().
											EventFunc(republishRelatedOnlinePagesEvent).
											Query("ids", strings.Join(pageIDs, ",")).
											Go(),
									).
									Go(),
							),
					).Class("d-flex"),
				).Else(
					h.Div(h.Text(pmsgr.ListingNoRecordToShow)).Class("text-center grey--text text--darken-2 mt-8"),
				),
			)
		}
		return
	})
}

func (b *ContainerBuilder) registerEventFuncs() {
	b.mb.RegisterEventFunc(AutoSaveContainerEvent, b.autoSaveContainer)
}

func (b *ContainerBuilder) getContainerDataID(id int) string {
	return fmt.Sprintf(inflection.Plural(strcase.ToKebab(b.name))+"_%v", id)
}

func republishRelatedOnlinePages(pageURL func(ctx *web.EventContext) string) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		ids := strings.Split(ctx.R.FormValue("ids"), ",")
		for _, id := range ids {
			statusVar := fmt.Sprintf(`republish_status_%s`, strings.Replace(id, "-", "_", -1))
			web.AppendRunScripts(&r,
				web.Plaid().
					URL(pageURL(ctx)).
					EventFunc(publish.EventRepublish).
					Query("id", id).
					Query(publish.ParamScriptAfterPublish, fmt.Sprintf(`vars.%s = "done"`, statusVar)).
					Query("status_var", statusVar).
					Go(),
				fmt.Sprintf(`vars.%s = "pending"`, statusVar),
			)
		}
		return r, nil
	}
}

func (b *Builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, mb := range b.models {
		if strings.Index(r.RequestURI, b.prefix+"/"+mb.name+"/preview") >= 0 {
			mb.preview.ServeHTTP(w, r)
			return
		}
	}
	if b.images != nil {
		if strings.Index(r.RequestURI, path.Join(b.prefix, b.imagesPrefix)) >= 0 {
			b.images.ServeHTTP(w, r)
			return
		}
	}
	b.ps.ServeHTTP(w, r)
}

func (b *Builder) generateEditorBarJsFunction(_ *web.EventContext) string {
	editAction := web.Plaid().
		PushState(true).
		MergeQuery(true).
		Query(paramContainerDataID, web.Var("container_data_id")).
		Query(paramContainerID, web.Var("container_id")).
		RunPushState() + ";" +
		web.POST().
			EventFunc(actions.Edit).
			URL(web.Var(fmt.Sprintf(`"%s/"+arr[0]`, b.prefix))).
			Query(presets.ParamID, web.Var("arr[1]")).
			Query(presets.ParamOverlay, actions.Content).
			Query(presets.ParamTargetPortal, pageBuilderRightContentPortal).
			Go()
	addAction := fmt.Sprintf(`vars.containerTab="%s";`, EditorTabElements) +
		web.Plaid().
			PushState(true).
			MergeQuery(true).
			Query(paramContainerID, web.Var("container_id")).
			Query(paramTab, EditorTabElements).
			RunPushState()
	deleteAction := web.POST().
		EventFunc(DeleteContainerConfirmationEvent).
		Query(paramContainerID, web.Var("container_id")).
		Query(paramContainerName, web.Var("display_name")).
		Go()
	moveAction := web.Plaid().
		EventFunc(MoveUpDownContainerEvent).
		MergeQuery(true).
		Query(paramContainerID, web.Var("container_id")).
		Query(paramMoveDirection, web.Var("msg_type")).
		Query(paramModelID, web.Var("model_id")).
		Go()
	return fmt.Sprintf(`
function(e){
	const { msg_type,container_data_id, container_id,display_name } = e.data
	if (!msg_type || !container_data_id.split) {
		return
	} 
	let arr = container_data_id.split("_");
	if (arr.length != 2) {
		console.log(arr);
		return
	}
    switch (msg_type) {
	  case '%s':
		%s;
		break
      case '%s':
        %s;
        break
	  case '%s':
	  case '%s':
		%s;
		break
	  case '%s':
        %s;
        break
    }
	
}`,
		EventEdit, editAction,
		EventDelete, deleteAction,
		EventUp, EventDown, moveAction, EventAdd, addAction,
	)
}

func defaultSubPageTitle(ctx *web.EventContext) string {
	return i18n.MustGetModuleMessages(ctx.Context(), I18nPageBuilderKey, Messages_en_US).(*Messages).PageOverView
}

func (b *ContainerBuilder) autoSaveContainer(ctx *web.EventContext) (r web.EventResponse, err error) {
	var (
		id = ctx.R.FormValue(presets.ParamID)
		mb = b.mb.Editing()
	)
	obj, vErr := mb.FetchAndUnmarshal(id, true, ctx)
	if vErr.HaveErrors() {
		err = errors.New(vErr.Error())
		return
	}

	if vErr = mb.Validators.Validate(obj, presets.FieldModeStack{presets.EDIT}, ctx); vErr.HaveErrors() {
		err = errors.New(vErr.Error())
		return
	}

	if err = mb.Saver(obj, id, ctx); err != nil {
		return
	}
	r.RunScript = web.Plaid().EventFunc(ReloadRenderPageOrTemplateEvent).Go()
	return
}

type (
	Device struct {
		Name  string
		Width string
		Icon  string
	}
)

func (b *Builder) PreviewDevices(devices ...Device) {
	b.devices = devices
}

func (b *Builder) getDevices() []Device {
	if len(b.devices) == 0 {
		b.setDefaultDevices()
	}
	return b.devices
}

func (b *Builder) setDefaultDevices() {
	b.devices = []Device{
		{Name: DeviceComputer, Width: "", Icon: "mdi-laptop"},
		{Name: DeviceTablet, Width: "768px", Icon: "mdi-tablet"},
		{Name: DevicePhone, Width: "414px", Icon: "mdi-cellphone"},
	}
}

func (b *Builder) deviceToggle(ctx *web.EventContext) h.HTMLComponent {
	var comps []h.HTMLComponent
	ctx.R.Form.Del(web.EventFuncIDName)
	for _, d := range b.getDevices() {
		comps = append(comps,
			VBtn("").Icon(d.Icon).Color(ColorPrimary).Variant(VariantText).Class("mr-4").
				Attr("@click", web.Plaid().EventFunc(ReloadRenderPageOrTemplateEvent).
					PushState(true).Queries(ctx.R.Form).Query(paramsDevice, d.Name).Go()).Value(d.Name),
		)
	}
	return web.Scope(
		VBtnToggle(
			comps...,
		).Class("pa-2 rounded-lg ").Attr("v-model", "toggleLocals.activeDevice").Density(DensityCompact),
	).Slot("{ locals : toggleLocals}").LocalsInit(fmt.Sprintf(`{activeDevice: "%s"}`, ctx.Param(paramsDevice)))
}
