package utils

import (
	"context"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"golang.org/x/text/language"
)

const I18nUtilsKey i18n.ModuleKey = "I18nUtilsKey"

func Install(b *presets.Builder) {
	b.I18n().
		RegisterForModule(language.English, I18nUtilsKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nUtilsKey, Messages_zh_CN).
		RegisterForModule(language.Japanese, I18nUtilsKey, Messages_ja_JP)
}

func MustGetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nUtilsKey, Messages_en_US).(*Messages)
}

func ConfirmDialog(msg string, okAction string, msgr *Messages) h.HTMLComponent {
	return VDialog(
		VCard(
			VCardTitle(h.Text(msg)),
			VCardActions(
				VSpacer(),
				VBtn(msgr.Cancel).
					Variant(VariantFlat).
					Class("ml-2").
					On("click", "locals.commonConfirmDialog = false"),

				VBtn(msgr.OK).
					Color("primary").
					Variant(VariantFlat).
					Theme(ThemeDark).
					Attr("@click", okAction),
			),
		),
	).MaxWidth("600px").
		Attr("v-model", "locals.commonConfirmDialog")
}

func DeleteDialog(msg string, okAction string, msgr *Messages) h.HTMLComponent {
	return web.Scope(
		VDialog(
			VCard(
				VCardTitle(h.Text(msg)),
				VCardActions(
					VSpacer(),
					VBtn(msgr.Cancel).
						Variant(VariantFlat).
						Class("ml-2").
						On("click", "locals.deleteConfirmation = false"),

					VBtn(msgr.OK).
						Color("primary").
						Variant(VariantFlat).
						Theme(ThemeDark).
						Attr("@click", okAction),
				),
			),
		).MaxWidth("600px").
			Attr("v-model", "locals.deleteConfirmation"),
	).Slot(" { locals }").LocalsInit(`{deleteConfirmation: true}`)
}

const CloseCustomDialog = "locals.customConfirmationDialog = false"

func CustomDialog(title h.HTMLComponent, content h.HTMLComponent, okAction string, msgr *Messages) h.HTMLComponent {
	Vcard := VCard()
	if title != nil {
		Vcard.AppendChild(VCardTitle(title))
	}
	if content != nil {
		Vcard.AppendChild(VCardText(content))
	}
	Vcard.AppendChild(
		VCardActions(
			VSpacer(),
			VBtn(msgr.Cancel).
				Variant(VariantFlat).
				Class("ml-2").
				On("click", CloseCustomDialog),

			VBtn(msgr.OK).
				Color("primary").
				Variant(VariantFlat).
				Theme(ThemeDark).
				Attr("@click", okAction),
		),
	)
	return web.Scope(
		VDialog(
			Vcard,
		).MaxWidth("600px").
			Attr("v-model", "locals.customConfirmationDialog"),
	).Slot(" { locals }").LocalsInit(`{ customConfirmationDialog: true }`)
}
