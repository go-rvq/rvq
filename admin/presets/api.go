package presets

import (
	"encoding/json"
	"net/http"
	"net/url"
	"slices"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/data"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/ui/vuetifyx"
)

type ContextStringer interface {
	ContextString(ctx *web.EventContext) string
}

type (
	OnClick struct {
		Raw     string
		Builder *web.VueEventTagBuilder
	}
	ComponentFunc             func(ctx *web.EventContext) h.HTMLComponent
	ButtonComponentFunc       func(ctx *web.EventContext, onclick *OnClick) h.HTMLComponent
	ObjectComponentFunc       func(obj interface{}, ctx *web.EventContext) h.HTMLComponent
	ModeObjectComponentFunc   func(mode FieldModeStack, obj interface{}, ctx *web.EventContext) h.HTMLComponent
	TabComponentFunc          func(obj interface{}, ctx *web.EventContext) (tab h.HTMLComponent, content h.HTMLComponent)
	EditingTitleComponentFunc func(obj interface{}, defaultTitle string, ctx *web.EventContext) string

	DoPageBuilder func(p *web.PageBuilder)
)

func (c *OnClick) MarshalJSON() ([]byte, error) {
	if c.Raw != "" {
		return json.Marshal(c.Raw)
	}
	return json.Marshal(c.Builder.Go())
}

func (c *OnClick) String() string {
	if c.Raw != "" {
		return c.Raw
	}
	return c.Builder.Go()
}

type FieldComponentFunc func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent

type (
	ActionButtomComponentFunc func(ctx *web.EventContext, title func() string, onclick *web.VueEventTagBuilder) h.HTMLComponent
	ActionComponentFunc       func(id string, ctx *web.EventContext) (comp h.HTMLComponent, err error)
	ActionUpdateFunc          func(id string, ctx *web.EventContext) (err error)
	ActionEnabledFunc         func(id string, ctx *web.EventContext) (ok bool, err error)
	ActionEnabledObjFunc      func(obj any, id string, ctx *web.EventContext) (ok bool, err error)

	ActionBuildComponentFunc func(id string, ctx *web.EventContext, b *ContentComponentBuilder) (err error)
)

type (
	BulkActionComponentFunc                  func(selectedIds []string, ctx *web.EventContext) (comp h.HTMLComponent, err error)
	BulkActionUpdateFunc                     func(selectedIds []string, ctx *web.EventContext, r *web.EventResponse) (err error)
	BulkActionSelectedIdsProcessorFunc       func(selectedIds []string, ctx *web.EventContext) (processedSelectedIds []string, err error)
	BulkActionSelectedIdsProcessorNoticeFunc func(selectedIds []string, processedSelectedIds []string, unactionableIds []string) string
	BulkActionComponentHandler               func(cb *ContentComponentBuilder, ctx *web.EventContext)
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
	FetchFunc  func(obj interface{}, id ID, ctx *web.EventContext) (err error)
	SaveFunc   func(obj interface{}, id ID, ctx *web.EventContext) (err error)
	CreateFunc func(obj interface{}, ctx *web.EventContext) (err error)
	DeleteFunc func(obj interface{}, id ID, cascade bool, ctx *web.EventContext) (err error)

	SaveCallbackFunc func(ctx *web.EventContext, obj any) (done func(success bool) error, err error)
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
	Default       bool
}

type FilterTabsFunc func(ctx *web.EventContext) []*FilterTab

type Plugin interface {
	Install(pb *Builder) (err error)
}

type ModelPlugin interface {
	ModelInstall(pb *Builder, mb *ModelBuilder) (err error)
}

type (
	ModelInstallFunc func(mb *ModelBuilder) error
	InstallFunc      func(pb *Builder) error
)

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
