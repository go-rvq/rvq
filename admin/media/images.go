package media

import (
	"fmt"
	"sort"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"gorm.io/gorm"
)

// A list of images (several of a record — a product's, an item's): a list
// editor of cards (ImagesCards), each a record with a media box (its field
// ImageCardComponentFunc). A card is the image itself: a click opens it in
// the lightbox, going through the images of the list; its dimensions and
// size float at its bottom, the delete over its top right; under it, its name
// and its description (MediaBoxConfig.DescriptionField); after the cards, a
// place to drop files.
//
//	ed.Field("Media").WithContextValue(media.MediaBoxConfig, cfg).
//	    SetterFunc(media.MediaBoxSetterFunc(db)).
//	    ComponentFunc(media.ImageCardComponentFunc(db))
//	parent.Field("Images").Nested(presets.NestedSlice(m, &ed.FieldsBuilder).
//	    SetComponentBuilder(media.ImagesCards("Media", cfg, "")))

// ImagesCards are the cards of a list of images: field, the media box of an
// item; cfg, its config (the drop's type and sizes); placeholder, the words of
// the place to drop files ("", the media library's).
func ImagesCards(field string, cfg *media_library.MediaBoxConfig, placeholder string) *presets.ListEditorCardsBuilder {
	return presets.NewListEditorCardsBuilder().
		FloatingDelete(true).
		LayoutFunc(func(*presets.FieldsBuilder) presets.FieldsLayout { return presets.FieldsLayout{field} }).
		AfterFunc(func(ctx *presets.ListEditorContainerContext) h.HTMLComponent {
			if ctx.FieldContext.ReadOnly {
				return nil
			}
			words := placeholder
			if words == "" {
				words = GetMessages(ctx.FieldContext.EventContext.Context()).DropImages
			}
			return FileDrop(cfg, words, fmt.Sprintf("$listEditorAddItems(boxes.map(m => ({%s: m})))", field))
		})
}

// ImagesGroup is the group, in the lightbox, of the images of the list of the
// media box of formKey: the list's form key ("Images[2].Media" → "Images").
func ImagesGroup(formKey string) string {
	if i := strings.LastIndex(formKey, "["); i > 0 {
		return formKey[:i]
	}
	return formKey
}

// ImageCardComponentFunc is the media box of an image of a list (a card of
// ImagesCards): its config made a card's (MediaBoxConfig.Card) — the image
// chosen, opened in the lightbox among the list's; none yet, the button that
// chooses it.
func ImageCardComponentFunc(db *gorm.DB) presets.FieldComponentFunc {
	box := MediaBoxComponentFunc(db, false)
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		cfg, _ := field.ContextValue(MediaBoxConfig).(*media_library.MediaBoxConfig)
		c := media_library.MediaBoxConfig{}
		if cfg != nil {
			c = *cfg
		}
		c.Card = true
		if mbx, _ := field.Value().(media_library.MediaBox); mbx.Url == "" && field.ReadOnly {
			return nil
		}
		field.SetContextValue(MediaBoxConfig, &c)
		return box(field, ctx)
	}
}

// imageCard is the media box of field as a card — a card of a list of
// images (MediaBoxConfig.Card: the lightbox through the list's, the list's
// delete), or a field's own (single: alone in the lightbox; its actions over
// the image: edit, replace, delete; the crops of its sizes) —: the image, its
// dimensions and size floating at its bottom, its name and description; a
// file not an image, its icon. None, the card to drop one on.
func imageCard(ctx *web.EventContext, mbx *media_library.MediaBox, field string, cfg *media_library.MediaBoxConfig, disabled, readonly, single bool) h.HTMLComponent {
	msgr := GetMessages(ctx.Context())
	if cfg == nil {
		cfg = &media_library.MediaBoxConfig{}
	}
	value := ""
	if mbx != nil && !mbx.ID.IsZero() {
		value = h.JSONString(mbx)
	}
	values := h.Input("").Type("hidden").Attr(web.VField(field+".Values", value)...)
	if mbx == nil || mbx.Url == "" {
		if readonly || (single && disabled) {
			return values
		}
		height := "150px"
		if single {
			height = "96px"
		}
		return h.Components(
			// the card to drop the file on (or to click, to choose it)
			h.Div(mediaBoxDrop(ctx, field, cfg, height)).Attr("data-image-card-empty", true),
			web.Portal().Name(cropperPortalName(field)),
			values,
		)
	}
	isImage := base.IsImageFormat(mbx.FileName) || base.IsImageFormat(mbx.Url)
	group := ImagesGroup(field)
	if single {
		group = ""
	}
	var shown h.HTMLComponent
	if isImage {
		// its thumbnail when it has one (the lightbox, the image whole)
		src := mbx.URL()
		if mbx.FileSizes[media_library.ThumbSizeName] > 0 {
			src = mbx.URL(media_library.ThumbSizeName)
		}
		shown = v.VImg().Src(fmt.Sprintf("%s?%d", src, time.Now().UnixNano()%1e6)).Height(150).Cover(true).Class("rounded")
	} else {
		shown = h.A(fileThumb(mbx.FileName)).Href(mbx.Url).Target("_blank").Style("display: block; height: 150px")
	}
	image := h.Div(
		shown,
		MediaBoxInfoChip(mbx, "").Variant(v.VariantFlat).
			Style("position: absolute; bottom: 6px; left: 50%; transform: translateX(-50%); opacity: .9"),
	).Style("position: relative").
		Attr("data-image-card-image", true)
	if cfg.BackgroundColor != "" {
		image.Style(fmt.Sprintf("position: relative; background-color: %s", cfg.BackgroundColor))
	}
	if isImage {
		image.Attr(vx.Lightbox(group, mbx.URL(), mbx.FileName)...)
	}
	// a field's own: its actions over the image
	if single && !readonly && !disabled {
		act := func(icon, title, click, data string) *v.VBtnBuilder {
			return v.VBtn("").Icon(icon).Size(v.SizeSmall).Color("surface").Variant(v.VariantElevated).
				Attr("title", title).Attr("@click.stop", click).Attr(data, true)
		}
		actions := h.Div().Class("d-flex ga-1").
			Style("position: absolute; top: 6px; right: 6px; z-index: 1").
			Attr("data-image-card-actions", true)
		if isImage && !cfg.DisableCrop {
			actions.AppendChildren(act("mdi-image-edit-outline", msgr.EditImage,
				editImageScript(field, fmt.Sprint(mbx.ID), base.DefaultSizeKey, cfg), "data-media-edit"))
		}
		actions.AppendChildren(
			act("mdi-swap-horizontal", msgr.ChooseFile, web.Plaid().EventFunc(openFileChooserEvent).
				Query("field", field).
				FieldValue("cfg", h.JSONString(cfg)).
				Go(), "data-media-replace"),
			// light, the icon red: seen over any image
			act("mdi-delete", msgr.Delete, web.Plaid().EventFunc(deleteFileEvent).
				Query("field", field).
				FieldValue("cfg", h.JSONString(cfg)).
				Go(), "data-media-delete").Attr("style", "color: rgb(var(--v-theme-error)) !important"),
		)
		image.AppendChildren(actions)
	}
	card := h.Div(
		image,
		h.Div(h.Text(mbx.FileName)).Class("text-caption text-truncate mt-1").Attr("title", mbx.FileName),
	).Attr("data-image-card", true)
	if !cfg.DisableDescription || cfg.DescriptionField != "" {
		key := MediaBoxDescriptionKey(field, cfg)
		description := ctx.R.FormValue(key)
		if description == "" {
			description = mbx.Description
		}
		placeholder := msgr.ImageDescription
		if single && cfg.DescriptionField == "" {
			placeholder = msgr.DescriptionForAccessibility
		}
		if readonly || disabled {
			if description != "" {
				card.AppendChildren(h.Div(h.Text(description)).Class("text-body-2"))
			}
		} else {
			card.AppendChildren(v.VTextarea().Variant(v.FieldVariantUnderlined).Rows(1).AutoGrow(true).
				Density(v.DensityCompact).HideDetails(true).
				Placeholder(placeholder).
				Attr(web.VField(key, description)...))
		}
	}
	var out h.HTMLComponent = card
	if single {
		// a card of its own (a list's is the list's)
		out = v.VCard(v.VCardText(card).Class("pa-2")).Variant(v.VariantOutlined).
			Style("max-width: 300px; min-width: 140px; flex: 1 1 220px").Attr("data-media-card", field)
		// the crops of its sizes: a card each, beside it
		if isImage && len(cfg.Sizes) > 0 {
			out = h.Div(append([]h.HTMLComponent{out}, sizeCards(mbx, field, cfg, msgr, disabled || readonly)...)...).
				Class("d-flex flex-wrap ga-2 align-start")
		}
	}
	return h.Components(out, web.Portal().Name(cropperPortalName(field)), values)
}

// sizeCards are the crops of the sizes of a media box (cfg.Sizes), a card
// each in the order of their names: the crop, its measures and size floating
// at its bottom, the size's name; a click, the editor on it (its fixed crop:
// the size's), unless locked.
func sizeCards(mbx *media_library.MediaBox, field string, cfg *media_library.MediaBoxConfig, msgr *Messages, locked bool) []h.HTMLComponent {
	keys := make([]string, 0, len(cfg.Sizes))
	for k := range cfg.Sizes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]h.HTMLComponent, 0, len(keys))
	for _, k := range keys {
		sz := cfg.Sizes[k]
		if sz == nil {
			continue
		}
		info := fmt.Sprintf("%d × %d", sz.Width, sz.Height)
		if n := mbx.FileSizes[k]; n > 0 {
			info += " · " + base.ByteCountSI(n)
		}
		image := h.Div(
			v.VImg().Src(fmt.Sprintf("%s?%d", mbx.URL(k), time.Now().UnixNano()%1e6)).Height(110).Class("rounded").
				Attr("style", "background: repeating-conic-gradient(#ddd 0% 25%, #fff 0% 50%) 50% / 12px 12px"),
			infoChip(info).Variant(v.VariantFlat).
				Style("position: absolute; bottom: 4px; left: 50%; transform: translateX(-50%); opacity: .9"),
		).Style("position: relative")
		c := v.VCard(
			v.VCardText(
				image,
				h.Div(
					v.VIcon("mdi-crop").Size(v.SizeXSmall).Class("mr-1"),
					h.Text(k),
				).Class("text-caption text-truncate mt-1 d-flex align-center").Attr("title", k),
			).Class("pa-2"),
		).Variant(v.VariantOutlined).Style("width: 150px; flex: none").Attr("data-image-card-size", k)
		if !locked && !cfg.DisableCrop {
			c.Attr("@click", editImageScript(field, fmt.Sprint(mbx.ID), k, cfg)).
				Attr("title", msgr.EditImage).Attr("role", "button").
				Style("width: 150px; flex: none; cursor: pointer")
		}
		out = append(out, c)
	}
	return out
}
