package publish

import (
	"context"
	"time"

	"github.com/qor5/admin/v3/media/storage"
	"github.com/qor5/admin/v3/model"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"gorm.io/gorm"
)

type PublishAction struct {
	Url      string
	Content  string
	IsDelete bool
}

// @snippet_begin(PublishList)
type List struct {
	PageNumber  int
	Position    int
	ListDeleted bool
	ListUpdated bool
}

// @snippet_end

// @snippet_begin(PublishSchedule)
type Schedule struct {
	ScheduledStartAt *time.Time `gorm:"index"`
	ScheduledEndAt   *time.Time `gorm:"index"`

	ActualStartAt *time.Time
	ActualEndAt   *time.Time
}

// @snippet_end

// @snippet_begin(PublishStatus)
const (
	StatusDraft   = "draft"
	StatusOnline  = "online"
	StatusOffline = "offline"
)

type Status struct {
	Status    string `gorm:"default:'draft'"`
	OnlineUrl string
}

// @snippet_end

// @snippet_begin(PublishVersion)
type Version struct {
	Version       string `gorm:"primary_key;size:128;not null;"`
	VersionName   string
	ParentVersion string
}

// @snippet_end

type PublishInterface interface {
	GetPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, Storage storage.Storage) (actions []*PublishAction, err error)
}

type UnPublishInterface interface {
	GetUnPublishActions(mb *presets.ModelBuilder, db *gorm.DB, ctx context.Context, Storage storage.Storage) (actions []*PublishAction, err error)
}

type AfterPublishInterface interface {
	AfterPublish(db *gorm.DB, Storage storage.Storage, ctx context.Context) error
}

type AfterUnPublishInterface interface {
	AfterUnPublish(db *gorm.DB, Storage storage.Storage, ctx context.Context) error
}

type StatusInterface interface {
	EmbedStatus() *Status
}

func (s *Status) EmbedStatus() *Status {
	return s
}

type PublicUrlInterface interface {
	GetPublicUrl(mb *presets.ModelBuilder, ctx *web.EventContext) string
}

func EmbedStatus(v any) *Status {
	return v.(StatusInterface).EmbedStatus()
}

type VersionInterface interface {
	EmbedVersion() *Version
	CreateVersion(db *gorm.DB, mid model.ID, obj interface{}) (string, error)
}

func (s *Version) EmbedVersion() *Version {
	return s
}

func EmbedVersion(v any) *Version {
	return v.(VersionInterface).EmbedVersion()
}

type ScheduleInterface interface {
	EmbedSchedule() *Schedule
}

func (s *Schedule) EmbedSchedule() *Schedule {
	return s
}

type ListInterface interface {
	EmbedList() *List
}

func (s *List) EmbedList() *List {
	return s
}
