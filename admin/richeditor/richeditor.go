package richeditor

import (
	"context"
	"fmt"

	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/ui/redactor"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/samber/lo"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"
)

const PluginImageInsert = "imageinsert"

// how to add more plugins from https://imperavi.com/redactor/plugins/
// 1. add {{plugin}}.min.js to redactor dir
// 2. add plugin name in Plugins array

// how to add own plugins
// 1. load plugin jss,css to PluginsJS,PluginsCSS
// 2. add plugin names in Plugins array
var (
	Plugins    = []string{"alignment", "fontcolor", PluginImageInsert, "table", "video"}
	PluginsJS  [][]byte
	PluginsCSS [][]byte
)

type RichEditorBuilder struct {
	db            *gorm.DB
	name          string
	value         string
	label         string
	placeholder   string
	plugins       []string
	setPlugins    bool
	rawConfig     interface{}
	imageDisabled bool
}

func RichEditor(db *gorm.DB, name string) (r *RichEditorBuilder) {
	r = &RichEditorBuilder{db: db, name: name}
	return
}

func (b *RichEditorBuilder) Value(v string) (r *RichEditorBuilder) {
	b.value = v
	return b
}

func (b *RichEditorBuilder) Label(v string) (r *RichEditorBuilder) {
	b.label = v
	return b
}

func (b *RichEditorBuilder) Placeholder(v string) (r *RichEditorBuilder) {
	b.placeholder = v
	return b
}

func (b *RichEditorBuilder) Plugins(v []string) (r *RichEditorBuilder) {
	b.plugins = v
	b.setPlugins = true
	return b
}

// Note: RawConfig overwrites Plugins
func (b *RichEditorBuilder) RawConfig(v interface{}) (r *RichEditorBuilder) {
	b.rawConfig = v
	return b
}

func (b *RichEditorBuilder) ImageDisabled() bool {
	return b.imageDisabled
}

func (b *RichEditorBuilder) SetImageDisabled(imageDisabled bool) *RichEditorBuilder {
	b.imageDisabled = imageDisabled
	return b
}

func (b *RichEditorBuilder) MarshalHTML(ctx context.Context) ([]byte, error) {
	p := Plugins
	if b.setPlugins {
		p = b.plugins
	}

	if b.imageDisabled {
		p = lo.Filter(p, func(item string, index int) bool {
			return item != PluginImageInsert
		})
	}

	redactorB := redactor.New().Placeholder(b.placeholder).Attr(web.VField(b.name, b.value)...)
	if b.rawConfig != nil {
		redactorB.RawConfig(b.rawConfig)
	} else {
		redactorB.Config(redactor.Config{Plugins: p})
	}
	r := h.Components(
		v.VSheet(
			h.Label(b.label).Class("v-label theme--light"),
			redactorB,
			h.If(!b.imageDisabled,
				h.Div(
					media.QMediaBox(b.db).FieldName(fmt.Sprintf("%s_richeditor_medialibrary", b.name)).
						Value(&media_library.MediaBox{}).Config(&media_library.MediaBoxConfig{
						AllowType: "image",
					}),
				).Class("hidden-screen-only"),
			),
		).Class("pb-4").Rounded(true).Attr("data-type", "redactor").Attr("style", "position: relative; z-index:1;"),
	)
	return r.MarshalHTML(ctx)
}
