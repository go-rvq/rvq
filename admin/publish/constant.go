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
	ModelChangeNotifyKey      ContextKey = "model_change_notify"
)

// Change says what happened to a record's published state.
type Change string

const (
	ChangePublished   Change = "published"
	ChangeUnpublished Change = "unpublished"
)

// ChangeCallback is told that a record went online or offline, inside the same
// transaction that made it so — a failure here rolls the publication back, and
// nothing outside sees a change that did not happen.
//
// It is what an application hangs a side effect on: refreshing a cache it keeps
// in memory, telling the other processes serving the same database. The builder
// calls it only for models that asked, with WithChangeNotify.
type ChangeCallback func(tx *gorm.DB, ctx context.Context, mb *presets.ModelBuilder, record any, change Change) error

type ModelPublishCallback func(db *gorm.DB, ctx context.Context, obj interface{}) (done func(err error) error, err error)
type ModelUnpublishCallback func(db *gorm.DB, ctx context.Context, obj interface{}) (done func(err error) error, err error)

// WithPublishCallback registers a callback to run when a record is published.
// Callbacks accumulate — several packages (e.g. an app's cache refresh and the
// history plugin's tag) can each add one, and Publish calls them all in order.
func WithPublishCallback(b *presets.ModelBuilder, f ModelPublishCallback) {
	b.SetData(ModelPublishCallbackKey, append(PublishCallbacks(b), f))
}

// PublishCallbacks returns the registered publish callbacks, in registration
// order.
func PublishCallbacks(b *presets.ModelBuilder) []ModelPublishCallback {
	cbs, _ := b.GetData(ModelPublishCallbackKey).([]ModelPublishCallback)
	return cbs
}

// WithUnpublishCallback registers a callback to run when a record is
// unpublished. Callbacks accumulate, like WithPublishCallback.
func WithUnpublishCallback(b *presets.ModelBuilder, f ModelUnpublishCallback) {
	b.SetData(ModelUnpublishCallbackKey, append(UnpublishCallbacks(b), f))
}

// UnpublishCallbacks returns the registered unpublish callbacks, in registration
// order.
func UnpublishCallbacks(b *presets.ModelBuilder) []ModelUnpublishCallback {
	cbs, _ := b.GetData(ModelUnpublishCallbackKey).([]ModelUnpublishCallback)
	return cbs
}

// WithChangeNotify marks a model whose publications the builder's ChangeCallback
// should hear about. It is off by default: most models are read from the
// database on every request and have nothing to be told about.
func WithChangeNotify(b *presets.ModelBuilder) {
	b.SetData(ModelChangeNotifyKey, true)
}

// ChangeNotifyEnabled reports whether a model asked for the callback.
func ChangeNotifyEnabled(b *presets.ModelBuilder) bool {
	enabled, _ := b.GetData(ModelChangeNotifyKey).(bool)
	return enabled
}
