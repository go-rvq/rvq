package examples_presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

type brand struct{}

func PresetsBrandTitle(b *presets.Builder, db *gorm.DB) (
	mb *presets.ModelBuilder,
	cl *presets.ListingBuilder,
	ce *presets.EditingBuilder,
	dp *presets.DetailingBuilder,
) {
	// @snippet_begin(BrandTitleSample)
	b.BrandTitle("QOR5 Admin")
	// @snippet_end
	b.Model(&brand{}).Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.Body = vuetify.VContainer()
		return
	})
	return
}

func PresetsBrandFunc(b *presets.Builder, db *gorm.DB) (
	mb *presets.ModelBuilder,
	cl *presets.ListingBuilder,
	ce *presets.EditingBuilder,
	dp *presets.DetailingBuilder,
) {
	// @snippet_begin(BrandFuncSample)
	b.BrandFunc(func(ctx *web.EventContext) h.HTMLComponent {
		return vuetify.VCardText(
			h.H1("Admin").Style("color: red;"),
		).Class("pa-0")
	})
	// @snippet_end
	b.Model(&brand{}).Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.Body = vuetify.VContainer()
		return
	})
	return
}
