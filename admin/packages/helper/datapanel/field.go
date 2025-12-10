package datapanel

import (
	"encoding/json"
	"fmt"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sunfmin/reflectutils"

	h "github.com/go-rvq/htmlgo"
)

type SeletorConfig struct {
	FieldName  string
	FieldLabel string
	FormKey    string
	Initial    Selected
	PortalName string
	RecordID   string
}

func (dp *DataPanel) FieldComponent(getLabel func(obj any) (string, error)) presets.FieldComponentFunc {
	return func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		obj := field.Obj
		cfg := &SeletorConfig{
			FieldName:  field.FormKey,
			FieldLabel: field.Label,
			FormKey:    field.FormKey,
			RecordID:   dp.b.ModelBuilder.MustRecordID(obj).String(),
			PortalName: dp.PortalName() + "/" + field.FormKey,
		}

		errs := field.Errors
		if len(errs) == 0 {
			var (
				err error
				id  string
			)

			if field.ValueOverride != nil {
				id = fmt.Sprint(field.ValueOverride)
			} else if dp.fieldID {
				id = fmt.Sprint(reflectutils.MustGet(obj, dp.field+"ID"))
			} else {
				id = fmt.Sprint(field.Value())
			}
			if id != "<nil>" && id != "" && id != "0" {
				cfg.Initial.ID = id
				if getLabel != nil {
					var label string
					if label, err = getLabel(obj); err == nil {
						cfg.Initial.Label = label
					}
				} else {
					var initial *Selected
					if initial, err = dp.load(ctx, id); err == nil {
						cfg.Initial = *initial
					}
				}
			}
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
		return web.Portal(dp.selector(cfg, cfg.Initial.ID, errs).Component(ctx)).Name(cfg.PortalName)
	}
}

func (dp *DataPanel) chooseHandler(ctx *web.EventContext) (r web.EventResponse, err error) {
	var cfg SeletorConfig
	if err = json.Unmarshal([]byte(ctx.R.FormValue(presets.SelectedEventConfigParamName)), &cfg); err != nil {
		return
	}

	r.UpdatePortal(cfg.PortalName, dp.selector(&cfg, ctx.R.FormValue(presets.ParamSelectedID), nil).Component(ctx))
	return
}

func (dp *DataPanel) selector(cfg *SeletorConfig, id string, errs []string) *InputComponent {
	return &InputComponent{DataPanel: dp, Config: cfg, Value: id, Errors: errs}
}

func (dp *DataPanel) load(ctx *web.EventContext, id string) (obj *Selected, err error) {
	if dp.parents != nil {
		var parents presets.IDSlice
		if parents, err = dp.parents(ctx.R); err != nil {
			return
		}

		defer web.WithContextValue(ctx, presets.ParentsModelIDKey, parents)()
	}

	if dp.b.Load != nil {
		return dp.b.Load(ctx, id)
	} else {
		var (
			p   = dp.b.ModelBuilder.NewModel()
			id_ presets.ID
		)
		if id_, err = dp.b.ModelBuilder.ParseRecordID(id); err != nil {
			return
		}
		if err = dp.b.ModelBuilder.CurrentDataOperator().Fetch(p, id_, ctx); err != nil {
			return
		}

		obj = &Selected{
			ID: id,
		}

		if dp.getLabelFunc != nil {
			obj.Label = dp.getLabelFunc(p)
		} else {
			obj.Label = fmt.Sprint(p)
		}
	}
	return
}

type InputComponent struct {
	DataPanel     *DataPanel
	Config        *SeletorConfig
	Value         string
	Errors        []string
	ComponentFunc func(ic *InputComponent) h.HTMLComponent
}

func (i *InputComponent) defaultComponent(ctx *web.EventContext) h.HTMLComponent {
	var (
		items = make([]*Selected, 0)
		comp  = VAutocomplete().
			Label(i.Config.FieldLabel).
			Name(i.Config.FieldName).
			Readonly(true).
			Clearable(true).
			ItemTitle("label").
			ItemValue("id")
	)

	if len(i.Errors) > 0 {
		comp.ErrorMessages(i.Errors...)
	}
	if i.Value != "" {
		var (
			item *Selected
			err  error
		)

		if item, err = i.DataPanel.load(ctx, i.Value); err != nil {
			comp.ErrorMessages(err.Error())
		} else {
			if i.Config.RecordID != "" && i.DataPanel.target.ModelType() == i.DataPanel.b.ModelBuilder.ModelType() && item.ID == i.Config.RecordID {
				comp.ErrorMessages("Auto referências não são permitidas.")
			}
			items = append(items, item)
		}
	}

	comp.
		Attr(web.VField(i.Config.FormKey, i.Value)...).
		Items(items)

	parents, err := presets.ResolveParentsModelID(i.DataPanel.parents, ctx.R)
	if err != nil {
		comp.ErrorMessages(err.Error())
	}

	return i.RegisterSelectDialog(h.Div(comp), parents)
}

func (i *InputComponent) Component(ctx *web.EventContext) h.HTMLComponent {
	if i.ComponentFunc != nil {
		return i.ComponentFunc(i)
	}

	return i.defaultComponent(ctx)
}

func (i *InputComponent) RegisterSelectDialog(tag *h.HTMLTagBuilder, parentsID presets.IDSlice) *h.HTMLTagBuilder {
	return tag.
		Attr("@click", web.Plaid().
			EventFunc(actions.OpenListingDialogForSelection).
			URL(i.DataPanel.b.ModelBuilder.Info().ListingHref(parentsID...)).
			Query(presets.SelectedEventParamName, i.DataPanel.selectedEvent()).
			Query(presets.SelectedEventConfigParamName, h.JSONString(i.Config)).
			Go())
}
