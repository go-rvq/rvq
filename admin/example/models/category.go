package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/admin/media/storage"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/publish"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model

	Name     string
	Products pq.StringArray `gorm:"type:text[]"`
	publish.Status
	publish.Schedule
	publish.Version
}

func (c *Category) PrimarySlug() string {
	return fmt.Sprintf("%v_%v", c.ID, c.Version.Version)
}

func (c *Category) PrimaryColumnValuesBySlug(slug string) map[string]string {
	segs := strings.Split(slug, "_")
	if len(segs) != 2 {
		panic("wrong slug")
	}
	return map[string]string{
		"id":      segs[0],
		"version": segs[1],
	}
}

func (c *Category) GetPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	return
}

func (c *Category) GetUnPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, _ storage.Storage) (objs []*publish.PublishAction, err error) {
	return
}

func (c *Category) PermissionRN() []string {
	return []string{"categories", strconv.Itoa(int(c.ID)), c.Version.Version}
}
