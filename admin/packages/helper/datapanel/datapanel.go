package datapanel

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

const (
	selectedEvent = "SelectedEvent"
)

type DataPanelBuilder struct {
	ModelBuilder *presets.ModelBuilder
	Load         func(ctx *web.EventContext, id string) (obj *Selected, err error)
}

func NewDataPanelBuilder(modelBuilder *presets.ModelBuilder) *DataPanelBuilder {
	return &DataPanelBuilder{ModelBuilder: modelBuilder}
}

func (b *DataPanelBuilder) Build(target *presets.ModelBuilder, fieldName, id string) *DataPanel {
	dp := &DataPanel{
		id:     id,
		b:      b,
		field:  fieldName,
		target: target,
	}
	target.EventsHub.RegisterEventFunc(dp.selectedEvent(), dp.chooseHandler)
	return dp
}

type DataPanel struct {
	id             string
	b              *DataPanelBuilder
	target         *presets.ModelBuilder
	field          string
	buildComponent func(ic *InputComponent) h.HTMLComponent
	getLabelFunc   func(obj any) string
	fieldID        bool
	parents        presets.ParentsModelIDResolver
}

func (dp *DataPanel) FieldID() bool {
	return dp.fieldID
}

func (dp *DataPanel) SetFieldID(fieldID bool) *DataPanel {
	dp.fieldID = fieldID
	return dp
}

func (dp *DataPanel) SetGetLabelFunc(getLabelFunc func(obj any) string) *DataPanel {
	dp.getLabelFunc = getLabelFunc
	return dp
}

func (dp *DataPanel) B() *DataPanelBuilder {
	return dp.b
}

func (dp *DataPanel) Field() string {
	return dp.field
}

func (dp *DataPanel) BuildComponent() func(ic *InputComponent) h.HTMLComponent {
	return dp.buildComponent
}

func (dp *DataPanel) SetBuildComponent(buildComponent func(ic *InputComponent) h.HTMLComponent) {
	dp.buildComponent = buildComponent
}

func (dp *DataPanel) Parents() presets.ParentsModelIDResolver {
	return dp.parents
}

func (dp *DataPanel) SetParents(parents presets.ParentsModelIDResolver) *DataPanel {
	dp.parents = parents
	return dp
}

func (dp *DataPanel) selectedEvent() string {
	return dp.id + dp.field + selectedEvent
}

func (dp *DataPanel) PortalName() string {
	return dp.field + "Portal"
}

type Selected struct {
	ID    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}
