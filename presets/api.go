package presets

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/qor5/admin/v3/presets/data"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
)

type (
	ComponentFunc             func(ctx *web.EventContext) h.HTMLComponent
	ObjectComponentFunc       func(obj interface{}, ctx *web.EventContext) h.HTMLComponent
	TabComponentFunc          func(obj interface{}, ctx *web.EventContext) (tab h.HTMLComponent, content h.HTMLComponent)
	EditingTitleComponentFunc func(obj interface{}, defaultTitle string, ctx *web.EventContext) string
)

type FieldComponentFunc func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent

type (
	ActionComponentFunc func(id string, ctx *web.EventContext) h.HTMLComponent
	ActionUpdateFunc    func(id string, ctx *web.EventContext) (err error)
)

type (
	BulkActionComponentFunc                  func(selectedIds []string, ctx *web.EventContext) h.HTMLComponent
	BulkActionUpdateFunc                     func(selectedIds []string, ctx *web.EventContext) (err error)
	BulkActionSelectedIdsProcessorFunc       func(selectedIds []string, ctx *web.EventContext) (processedSelectedIds []string, err error)
	BulkActionSelectedIdsProcessorNoticeFunc func(selectedIds []string, processedSelectedIds []string, unactionableIds []string) string
)

type MessagesFunc func(r *http.Request) *Messages

// Data Layer
type (
	DataOperator = data.DataOperator
	SQLCondition = data.SQLCondition
	SearchParams = data.SearchParams
)

type (
	SetterFunc         func(obj interface{}, ctx *web.EventContext)
	FieldSetterFunc    func(obj interface{}, field *FieldContext, ctx *web.EventContext) (err error)
	ValidateFunc       func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors)
	OnChangeActionFunc func(id string, ctx *web.EventContext) (s string)
)

type (
	SearchFunc func(model interface{}, params *SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error)
	FetchFunc  func(obj interface{}, id string, ctx *web.EventContext) (err error)
	SaveFunc   func(obj interface{}, id string, ctx *web.EventContext) (err error)
	DeleteFunc func(obj interface{}, id string, ctx *web.EventContext) (err error)
)

type SlugDecoder interface {
	PrimaryColumnValuesBySlug(slug string) map[string]string
}

type SlugEncoder interface {
	PrimarySlug() string
}

type FilterDataFunc func(ctx *web.EventContext) vuetifyx.FilterData

type FilterTab struct {
	ID    string
	Label string
	// render AdvancedLabel if it is not nil
	AdvancedLabel h.HTMLComponent
	Query         url.Values
}

type FilterTabsFunc func(ctx *web.EventContext) []*FilterTab

type Plugin interface {
	Install(pb *Builder) (err error)
}

type ModelPlugin interface {
	ModelInstall(pb *Builder, mb *ModelBuilder) (err error)
}

type (
	ModelInstallFunc func(pb *Builder, mb *ModelBuilder) error
	InstallFunc      func(pb *Builder) error
)

type ID struct {
	Model *ModelBuilder
	Value interface{}
}

type IDOfField struct {
	ID    ID
	Field string
}

func IDString(s string) ID {
	return ID{Value: s}
}

type ModelConfigurer interface {
	ConfigureModel(mb *ModelBuilder)
}

type ModelConfigurators []ModelConfigurer

func (mc ModelConfigurators) ConfigureModel(mb *ModelBuilder) {
	for _, m := range mc {
		m.ConfigureModel(mb)
	}
}

func (mc *ModelConfigurators) Append(m ...ModelConfigurer) {
	*mc = append(*mc, m...)
}

func (mc *ModelConfigurators) Insert(i int, m ...ModelConfigurer) {
	*mc = slices.Insert(*mc, i, m...)
}

type ModelConfiguratorFunc func(mb *ModelBuilder)

func (f ModelConfiguratorFunc) ConfigureModel(mb *ModelBuilder) {
	f(mb)
}
