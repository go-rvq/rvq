package admin

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/example/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/seo"
	"gorm.io/gorm"
)

// @snippet_begin(SeoExample)
var seoBuilder *seo.Builder

func configureSeo(pb *presets.Builder, db *gorm.DB, locales ...string) {
	seoBuilder = seo.New(db, seo.WithLocales(locales...))
	seoBuilder.RegisterSEO("Post", &models.Post{}).RegisterContextVariable(
		"Title",
		func(object interface{}, _ *seo.Setting, _ *http.Request) string {
			if article, ok := object.(models.Post); ok {
				return article.Title
			}
			return ""
		},
	).RegisterSettingVariables("Test")
	seoBuilder.RegisterSEO("Product")
	seoBuilder.RegisterSEO("Announcement")
	pb.Use(seoBuilder)
}

// @snippet_end
