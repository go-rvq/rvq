package media

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/i18n"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

func fileChooser(p *presets.Builder, mb *Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nMediaLibraryKey, Messages_en_US).(*Messages)
		field := ctx.R.FormValue("field")
		cfg := stringToCfg(ctx.R.FormValue("cfg"))

		portalName := mainPortalName(field)

		body := VCard(
			VResponsive(
				VRow(
					VCol(
						VTextField().Class("ml-4").
							PrependIcon("mdi-magnify").
							Label(msgr.Search).
							Flat(true).
							Clearable(true).
							HideDetails(true).
							ModelValue("").
							Attr("@keyup.enter", web.Plaid().
								EventFunc(imageSearchEvent).
								Query("field", field).
								FieldValue("cfg", h.JSONString(cfg)).
								FieldValue(searchKeywordName(field), web.Var("$event")).
								Go()),
					),
				),
			),
			web.Portal().Name(deleteConfirmPortalName(field)),
			web.Portal(
				fileChooserDialogContent(mb, field, ctx, cfg),
			).Name(dialogContentPortalName(field)),
		)

		d := p.DialogPortal(portalName).
			SetScrollable(true)

		cb := &presets.ContentComponentBuilder{
			Title: msgr.ChooseAFile,
			Overlay: &presets.ContentComponentBuilderOverlay{
				Mode: actions.Dialog,
			},
			Body: body,
		}

		d.Respond(ctx, &r, cb.BuildOverlay())
		return
	}
}

func fileChooserDialogContent(mb *Builder, field string, ctx *web.EventContext,
	cfg *media_library.MediaBoxConfig,
) h.HTMLComponent {
	db := mb.db
	msgr := GetMessages(ctx.Context())

	keyword := ctx.R.FormValue(searchKeywordName(field))

	type selectItem struct {
		Text  string
		Value string
	}
	const (
		orderByKey           = "order_by"
		orderByCreatedAt     = "created_at"
		orderByCreatedAtDESC = "created_at_desc"

		typeKey   = "type"
		typeAll   = "all"
		typeImage = "image"
		typeVideo = "video"
		typeFile  = "file"
	)
	orderByVal := ctx.R.URL.Query().Get(orderByKey)
	typeVal := ctx.R.URL.Query().Get(typeKey)

	var files []*media_library.MediaLibrary
	wh := db.Model(&media_library.MediaLibrary{})

	switch orderByVal {
	case orderByCreatedAt:
		wh = wh.Order("created_at")
	default:
		orderByVal = orderByCreatedAtDESC
		wh = wh.Order("created_at DESC")
	}

	switch typeVal {
	case typeImage:
		wh = wh.Where("selected_type = ?", media_library.ALLOW_TYPE_IMAGE)
	case typeVideo:
		wh = wh.Where("selected_type = ?", media_library.ALLOW_TYPE_VIDEO)
	case typeFile:
		wh = wh.Where("selected_type = ?", media_library.ALLOW_TYPE_FILE)
	default:
		typeVal = typeAll
	}

	currentPageInt, _ := strconv.Atoi(ctx.R.FormValue(currentPageName(field)))
	if currentPageInt == 0 {
		currentPageInt = 1
	}

	if len(cfg.Sizes) > 0 {
		cfg.AllowType = media_library.ALLOW_TYPE_IMAGE
	}

	if len(cfg.AllowType) > 0 {
		wh = wh.Where("selected_type = ?", cfg.AllowType)
	}

	if len(keyword) > 0 {
		wh = wh.Where("file ILIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}

	var count int64
	err := wh.Count(&count).Error
	if err != nil {
		panic(err)
	}
	perPage := mb.mediaLibraryPerPage
	pagesCount := int(count/int64(perPage) + 1)
	if count%int64(perPage) == 0 {
		pagesCount--
	}

	wh = wh.Limit(perPage).Offset((currentPageInt - 1) * perPage)
	err = wh.Find(&files).Error
	if err != nil {
		panic(err)
	}

	fileAccept := "*/*"
	if cfg.AllowType == media_library.ALLOW_TYPE_IMAGE {
		fileAccept = "image/*"
	}

	row := VRow(
		h.If(mb.uploadIsAllowed(ctx.R) == nil,
			VCol(
				h.Label("").Children(
					VCard(
						VCardTitle(h.Text(msgr.UploadFiles)),
						VIcon("mdi-cloud-upload").Size(SizeLarge),
						h.Input("").
							Attr("accept", fileAccept).
							Type("file").
							Attr("multiple", true).
							Style("display:none").
							Attr("@change",
								"form.NewFiles = [...$event.target.files];"+
									web.Plaid().
										BeforeScript("locals.fileChooserUploadingFiles = $event.target.files").
										EventFunc(uploadFileEvent).
										Query("field", field).
										FieldValue("cfg", h.JSONString(cfg)).
										Go()),
					).
						Height(200).
						Class("d-flex align-center justify-center pa-6").
						Attr("role", "button").
						Attr("v-ripple", true),
				),
			).
				Cols(6).Sm(4).Md(3),

			VCol(
				VCard(
					VProgressCircular().
						Color("primary").
						Indeterminate(true),
				).
					Class("d-flex align-center justify-center").
					Height(200),
			).
				Attr("v-for", "f in locals.fileChooserUploadingFiles").
				Cols(6).Sm(4).Md(3),
		),
	)

	initCroppingVars := []string{fileCroppingVarName(0) + ": false"}

	for i, f := range files {
		_, needCrop := mergeNewSizes(f, cfg)
		croppingVar := fileCroppingVarName(f.ID)
		initCroppingVars = append(initCroppingVars, fmt.Sprintf("%s: false", croppingVar))
		imgClickVars := fmt.Sprintf("vars.mediaShow = '%s'; vars.mediaName = '%s'; vars.isImage = %s", f.File.URL(), f.File.FileName, strconv.FormatBool(base.IsImageFormat(f.File.FileName)))

		var (
			links h.HTMLComponent
		)

		if f.SelectedType == media_library.ALLOW_TYPE_IMAGE {
			var (
				items     []map[string]string
				sizeNames []string
			)

			if sizes := f.GetMediaOption().Sizes; len(sizes) > 0 {
				for name := range sizes {
					sizeNames = append(sizeNames, name)
				}
				sort.Slice(sizeNames, func(i, j int) bool {
					return sizes[sizeNames[i]].Width < sizes[sizeNames[j]].Width &&
						sizes[sizeNames[i]].Height < sizes[sizeNames[j]].Height
				})

				for _, name := range sizeNames {
					s := sizes[name]
					var url string
					if name == "default" {
						url = f.File.URL()
					} else {
						url = f.File.URL(name)
					}
					items = append(items, map[string]string{
						"label": fmt.Sprintf(`%s (%dx%d) %s`, name, s.Width, s.Height, base.ByteCountSI(f.File.FileSizes[name])),
						"uri":   url,
					})
				}
			}

			items = append(items,
				map[string]string{
					"label": fmt.Sprintf("original (%dx%d) %s", f.File.Width, f.File.Height, base.ByteCountSI(f.File.FileSizes["original"])),
					"uri":   f.File.URL("original"),
				},
			)

			links = vue.UserComponent(
				VMenu(
					web.Slot(
						VBtn("").
							Icon("mdi-link").
							Attr("v-bind", "props").
							AppendIcon("mdi-chevron-down"),
					).Name("activator").Scope("{ props }"),
					VList(
						VListItem(
							VListItemTitle(h.RawHTML(`{{ item.label }}`)),
						).
							PrependIcon("mdi-content-copy").
							Attr("v-for", "item in items").
							Attr("@click", `(e) => copy(e, item.uri)`).
							Attr("@click.middle", `(e) => e.view.window.open(item.uri, "_blank")`),
					).Density("compact"),
				),
			).Scope("items", items)
		} else {
			links = VBtn("").
				Title(msgr.CopyLink).
				Icon("mdi-link").
				Attr("@click", fmt.Sprintf(`(e) => copy(e, %q)`, f.File.URL())).
				Attr("@click.middle", fmt.Sprintf(`(e) => e.view.window.open(%q, "_blank")`, f.File.URL()))
		}

		row.AppendChild(
			VCol(
				VCard(
					h.Div(
						h.If(
							base.IsImageFormat(f.File.FileName),
							VImg(
								h.If(needCrop,
									h.Div(
										VProgressCircular().Indeterminate(true),
										h.Span(msgr.Cropping).Class("text-h6 pl-2"),
									).Class("d-flex align-center justify-center v-card--reveal white--text").
										Style("height: 100%; background: rgba(0, 0, 0, 0.5)").
										Attr("v-if", fmt.Sprintf("locals.%s", croppingVar)),
								),
							).Src(f.File.URL(media_library.QorPreviewSizeName)).Height(200),
							// .Contain(true),
						).Else(
							fileThumb(f.File.FileName),
						),
					).AttrIf("role", "button", field != mediaLibraryListField).
						AttrIf("@click", web.Plaid().
							BeforeScript(fmt.Sprintf("locals.%s = true", croppingVar)).
							EventFunc(chooseFileEvent).
							Query("field", field).
							Query("id", fmt.Sprint(f.ID)).
							FieldValue("cfg", h.JSONString(cfg)).
							Go(), field != mediaLibraryListField).
						AttrIf("@click", imgClickVars, field == mediaLibraryListField),
					VCardText(
						h.A().Text(f.File.FileName).
							Attr("@click", imgClickVars),
						h.Input("").
							Style("width: 100%;").
							Placeholder(msgr.DescriptionForAccessibility).
							Value(f.File.Description).
							Attr("@change", web.Plaid().
								EventFunc(updateDescriptionEvent).
								Query("field", field).
								Query("id", fmt.Sprint(f.ID)).
								FieldValue("cfg", h.JSONString(cfg)).
								FieldValue("CurrentDescription", web.Var("$event.target.value")).
								Go(),
							).Readonly(mb.updateDescIsAllowed(ctx.R, files[i]) != nil),
						h.If(base.IsImageFormat(f.File.FileName),
							fileChips(f),
						),
					),
					VCardActions(
						links,
						VSpacer(),
						h.If(mb.deleteIsAllowed(ctx.R, files[i]) == nil,
							VBtn(msgr.Delete).
								Variant(VariantText).
								Attr("@click",
									web.Plaid().
										EventFunc(deleteConfirmationEvent).
										Query("field", field).
										Query("id", fmt.Sprint(f.ID)).
										FieldValue("cfg", h.JSONString(cfg)).
										Go(),
								),
						),
					),
				),
			).Cols(6).Sm(4).Md(3),
		)
	}

	return vue.UserComponent(
		h.Div(
			VSnackbar(h.Text(msgr.DescriptionUpdated)).
				Attr("v-model", "vars.snackbarShow").
				Location("top").
				Color("primary").
				Timeout(5000),
			web.Scope(
				VContainer(
					h.If(field == mediaLibraryListField,
						VRow(
							VCol(
								VSelect().Items([]selectItem{
									{Text: msgr.All, Value: typeAll},
									{Text: msgr.Images, Value: typeImage},
									{Text: msgr.Videos, Value: typeVideo},
									{Text: msgr.Files, Value: typeFile},
								}).ItemTitle("Text").ItemValue("Value").
									Attr(web.VField(typeKey, typeVal)...).
									Attr("@change",
										web.GET().PushState(true).
											Query(typeKey, web.Var("$event")).
											MergeQuery(true).Go(),
									).
									Density(DensityCompact).Variant(FieldVariantSolo).Class("mb-n8"),
							).Cols(3),
							VCol(
								VSelect().Items([]selectItem{
									{Text: msgr.UploadedAtDESC, Value: orderByCreatedAtDESC},
									{Text: msgr.UploadedAt, Value: orderByCreatedAt},
								}).ItemTitle("Text").ItemValue("Value").
									Attr(web.VField(orderByKey, orderByVal)...).
									Attr("@change",
										web.GET().PushState(true).
											Query(orderByKey, web.Var("$event")).
											MergeQuery(true).Go(),
									).
									Density(DensityCompact).Variant(FieldVariantSolo).Class("mb-n8"),
							).Cols(3),
						).Justify("end"),
					),
					row,
					VRow(
						VCol().Cols(1),
						VCol(
							VPagination().
								Length(pagesCount).
								ModelValue(int(currentPageInt)).
								Attr("@update:modelValue", web.Plaid().
									FieldValue(currentPageName(field), web.Var("$event")).
									EventFunc(imageJumpPageEvent).
									Query("field", field).
									FieldValue("cfg", h.JSONString(cfg)).
									Go()),
						).Cols(10),
					),
					VCol().Cols(1),
				).Fluid(true),
			).LocalsInit(fmt.Sprintf(`{fileChooserUploadingFiles: [], %s}`, strings.Join(initCroppingVars, ", "))).Slot("{ locals }"),
			VOverlay(
				h.Img("").Attr(":src", "vars.isImage? vars.mediaShow: ''").
					Style("max-height: 80vh; max-width: 80vw; background: rgba(0, 0, 0, 0.5)"),
				h.Div(
					h.A(
						VIcon("info").Size(SizeSmall).Class("mb-1"),
						h.Text("{{vars.mediaName}}"),
					).Attr(":href", "vars.mediaShow? vars.mediaShow: ''").Target("_blank").
						Class("white--text").Style("text-decoration: none;"),
				).Class("d-flex align-center justify-center pt-2"),
			).Attr("v-if", "vars.mediaName").Attr("@click", "vars.mediaName = null").ZIndex(10),
		).Attr(web.VAssign("vars", `{snackbarShow: false, mediaShow: null, mediaName: null, isImage: false}`)...),
	).
		Scope("copy").
		Setup(`({ scope }) => {
scope.copy = (e, uri) => {
	e.view.navigator.clipboard.writeText(uri)
	vars.presetsMessage = { show: true, message: ` + strconv.Quote(msgr.LinkCopied) + `, color: "info"}
	e.view.setTimeout(() => { vars.presetsMessage.show = false }, 1000)
}
}`)
}

func fileChips(f *media_library.MediaLibrary) h.HTMLComponent {
	g := VChipGroup().Column(true)
	text := "original"
	if f.File.Width != 0 && f.File.Height != 0 {
		text = fmt.Sprintf("%s(%dx%d)", "original", f.File.Width, f.File.Height)
	}
	if f.File.FileSizes["original"] != 0 {
		text = fmt.Sprintf("%s %s", text, base.ByteCountSI(f.File.FileSizes["original"]))
	}
	g.AppendChild(
		VChip(h.Text(text)).Size(SizeSmall),
	)
	// if len(f.File.Sizes) == 0 {
	//	return g
	// }

	// for k, size := range f.File.GetSizes() {
	//	g.AppendChildren(
	//		VChip(thumbName(k, size)).XSize(SizeSmall),
	//	)
	// }
	return g
}

type uploadFiles struct {
	NewFiles []*multipart.FileHeader
}

func uploadFile(b *Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		field := ctx.R.FormValue("field")
		cfg := stringToCfg(ctx.R.FormValue("cfg"))

		if err = b.uploadIsAllowed(ctx.R); err != nil {
			return
		}

		var uf uploadFiles
		ctx.MustUnmarshalForm(&uf)
		for _, fh := range uf.NewFiles {
			m := media_library.MediaLibrary{}

			if base.IsImageFormat(fh.Filename) {
				m.SelectedType = media_library.ALLOW_TYPE_IMAGE
			} else if base.IsVideoFormat(fh.Filename) {
				m.SelectedType = media_library.ALLOW_TYPE_VIDEO
			} else {
				m.SelectedType = media_library.ALLOW_TYPE_FILE
			}
			err = m.File.Scan(fh)
			if err != nil {
				panic(err)
			}

			err = base.SaveUploadAndCropImage(&b.Config, b.db, &m)
			if err != nil {
				presets.ShowMessage(&r, err.Error(), "error")
				return r, nil
			}
		}

		renderFileChooserDialogContent(ctx, &r, field, b, cfg)
		return
	}
}

func mergeNewSizes(m *media_library.MediaLibrary, cfg *media_library.MediaBoxConfig) (sizes map[string]*base.Size, r bool) {
	sizes = make(map[string]*base.Size)
	for k, size := range cfg.Sizes {
		if m.File.Sizes[k] != nil {
			sizes[k] = m.File.Sizes[k]
			continue
		}
		sizes[k] = size
		r = true
	}
	return
}

func chooseFile(b *Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		db := b.db
		field := ctx.R.FormValue("field")
		id := ctx.ParamAsInt("id")
		cfg := stringToCfg(ctx.R.FormValue("cfg"))

		var m media_library.MediaLibrary
		err = db.Find(&m, id).Error
		if err != nil {
			return
		}
		sizes, needCrop := mergeNewSizes(&m, cfg)

		if needCrop {
			err = m.ScanMediaOptions(media_library.MediaOption{
				Sizes: sizes,
				Crop:  true,
			})
			if err != nil {
				return
			}
			err = db.Save(&m).Error
			if err != nil {
				return
			}

			err = base.SaveUploadAndCropImage(&b.Config, b.db, &m)
			if err != nil {
				presets.ShowMessage(&r, err.Error(), "error")
				return r, nil
			}
		}

		mediaBox := media_library.MediaBox{
			ID:          json.Number(fmt.Sprint(m.ID)),
			Url:         m.File.Url,
			VideoLink:   "",
			FileName:    m.File.FileName,
			Description: m.File.Description,
			FileSizes:   m.File.FileSizes,
			Width:       m.File.Width,
			Height:      m.File.Height,
		}

		r.UpdatePortal(
			mediaBoxThumbnailsPortalName(field),
			mediaBoxThumbnails(ctx, &mediaBox, field, cfg, false, false),
		)
		r.RunScript = `closer.show = false`
		return
	}
}

func searchFile(mb *Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		field := ctx.R.FormValue("field")
		cfg := stringToCfg(ctx.R.FormValue("cfg"))

		ctx.R.Form[currentPageName(field)] = []string{"1"}

		renderFileChooserDialogContent(ctx, &r, field, mb, cfg)
		return
	}
}

func jumpPage(mb *Builder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		field := ctx.R.FormValue("field")
		cfg := stringToCfg(ctx.R.FormValue("cfg"))
		renderFileChooserDialogContent(ctx, &r, field, mb, cfg)
		return
	}
}

func renderFileChooserDialogContent(ctx *web.EventContext, r *web.EventResponse, field string, mb *Builder, cfg *media_library.MediaBoxConfig) {
	r.UpdatePortal(
		dialogContentPortalName(field),
		fileChooserDialogContent(mb, field, ctx, cfg),
	)
}
