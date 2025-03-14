package pagebuilder

import (
	"context"
	"fmt"
	"net/http/httptest"
	"path"
	"path/filepath"

	"github.com/qor5/admin/v3/l10n"
	"github.com/qor5/admin/v3/media/storage"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/publish"
	"gorm.io/gorm"
)

type contextKeyType int

const contextKey contextKeyType = iota

func (b *ModelBuilder) ContextValueProvider(in context.Context) context.Context {
	return context.WithValue(in, contextKey, b)
}

func builderFromContext(c context.Context) (b *ModelBuilder, ok bool) {
	b, ok = c.Value(contextKey).(*ModelBuilder)
	return
}

func (p *Page) GetPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	var b *ModelBuilder
	var ok bool
	if b, ok = builderFromContext(ctx); !ok || b == nil {
		return
	}
	content, err := p.getPublishContent(b, ctx)
	if err != nil {
		return
	}

	var localePath string
	if b.builder.l10n != nil {
		localePath = l10n.LocalePathFromContext(p, ctx)
	}

	var category Category
	category, err = p.GetCategory(db)
	if err != nil {
		return
	}
	objs = append(objs, &publish.PublishAction{
		Url:      p.getPublishUrl(localePath, category.Path),
		Content:  content,
		IsDelete: false,
	})
	p.OnlineUrl = p.getPublishUrl(localePath, category.Path)

	var liveRecord Page
	{
		lrdb := db.Where("id = ? AND status = ?", p.ID, publish.StatusOnline)
		if b.builder.l10n != nil {
			lrdb = lrdb.Where("locale_code = ?", p.LocaleCode)
		}
		lrdb.First(&liveRecord)
	}
	if liveRecord.ID == 0 {
		return
	}

	if liveRecord.OnlineUrl != p.OnlineUrl {
		objs = append(objs, &publish.PublishAction{
			Url:      liveRecord.OnlineUrl,
			IsDelete: true,
		})
	}

	return
}

func (p *Page) GetUnPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	objs = append(objs, &publish.PublishAction{
		Url:      p.OnlineUrl,
		IsDelete: true,
	})
	return
}

func generatePublishUrl(localePath, categoryPath, slug string) string {
	return path.Join("/", localePath, categoryPath, slug, "/index.html")
}

func (p *Page) getPublishUrl(localePath, categoryPath string) string {
	return generatePublishUrl(localePath, categoryPath, p.Slug)
}

func (p *Page) getAccessUrl(publishUrl string) string {
	return filepath.Dir(publishUrl)
}

func (p *Page) getPublishContent(b *ModelBuilder, _ context.Context) (r string, err error) {
	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", fmt.Sprintf("/?id=%s", p.PrimarySlug()), nil)
	b.preview.ServeHTTP(w, req)
	r = w.Body.String()
	return
}
