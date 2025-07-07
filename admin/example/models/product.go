package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/media/storage"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/publish"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Code  string
	Name  string
	Price int
	Image media_library.MediaBox `sql:"type:text;"`
	publish.Status
	publish.Schedule
	publish.Version
}

func (p *Product) PrimarySlug() string {
	return fmt.Sprintf("%v_%v", p.ID, p.Version.Version)
}

func (p *Product) PrimaryColumnValuesBySlug(slug string) map[string]string {
	segs := strings.Split(slug, "_")
	if len(segs) != 2 {
		panic("wrong slug")
	}

	return map[string]string{
		"id":      segs[0],
		"version": segs[1],
	}
}

func (p *Product) GetPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	return
}

func (p *Product) GetUnPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	return
}

func (p *Product) PermissionRN() []string {
	return []string{"products", strconv.Itoa(int(p.ID)), p.Code, p.Version.Version}
}
