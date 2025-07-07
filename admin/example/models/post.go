package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/media/storage"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/publish"
	"github.com/go-rvq/rvq/admin/seo"
	"github.com/go-rvq/rvq/admin/slug"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model

	Title         string
	TitleWithSlug slug.Slug
	Seo           seo.Setting
	Body          string
	HeroImage     media_library.MediaBox `sql:"type:text;"`
	BodyImage     media_library.MediaBox `sql:"type:text;"`
	UpdatedAt     time.Time
	CreatedAt     time.Time

	publish.Status
	publish.Schedule
	publish.Version
}

func (p *Post) PrimarySlug() string {
	return fmt.Sprintf("%v_%v", p.ID, p.Version.Version)
}

func (p *Post) PrimaryColumnValuesBySlug(slug string) map[string]string {
	segs := strings.Split(slug, "_")
	if len(segs) != 2 {
		panic("wrong slug")
	}

	return map[string]string{
		"id":      segs[0],
		"version": segs[1],
	}
}

func (p *Post) GetPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	return
}

func (p *Post) GetUnPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	return
}

func (p *Post) PermissionRN() []string {
	return []string{"posts", strconv.Itoa(int(p.ID)), p.Version.Version}
}
