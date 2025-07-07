package publish

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets"
	"gorm.io/gorm"
)

type Model struct {
	Record  interface{}
	Builder *presets.ModelBuilder
}

var (
	NonVersionPublishModels map[string]*Model
	VersionPublishModels    map[string]*Model
	ListPublishModels       map[string]*Model
)

func init() {
	NonVersionPublishModels = make(map[string]*Model)
	VersionPublishModels = make(map[string]*Model)
	ListPublishModels = make(map[string]*Model)
}

type ContextKey string

const (
	ModelPublishCallbackKey   ContextKey = "mode_publish_callback"
	ModelUnpublishCallbackKey ContextKey = "mode_unpublish_callback"
)

type ModelPublishCallback func(db *gorm.DB, ctx context.Context, obj interface{}) (done func(err error) error, err error)
type ModelUnpublishCallback func(db *gorm.DB, ctx context.Context, obj interface{}) (done func(err error) error, err error)

func WithPublishCallback(b *presets.ModelBuilder, f ModelPublishCallback) {
	b.SetData(ModelPublishCallbackKey, f)
}

func WithUnpublishCallback(b *presets.ModelBuilder, f ModelUnpublishCallback) {
	b.SetData(ModelUnpublishCallbackKey, f)
}
