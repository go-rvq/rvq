package vuetify

import (
	"encoding/json"

	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"
)

type (
	DataTableHeaderBasic struct {
		Title       string                 `json:"title,omitempty"`
		Value       any                    `json:":value,omitempty"`
		Key         string                 `json:"key,omitempty"`
		Align       string                 `json:"align,omitempty"`
		Width       string                 `json:"width,omitempty"`
		Children    any                    `json:"children,omitempty"`
		Index       int                    `json:"index,omitempty"`
		Sortable    bool                   `json:"sortable,omitempty"`
		HeaderProps map[string]interface{} `json:"headerProps,omitempty"`
		CellProps   map[string]interface{} `json:"cellProps,omitempty"`
		Hidden      bool                   `json:"hidden,omitempty"`
	}

	DataTableHeaderBasicSlice []DataTableHeaderBasic

	DataTableHeaderBuilder struct {
		DataTableHeaderBasic
		Advanced map[string]interface{}
	}

	DataTableHeaderSlice []*DataTableHeaderBuilder
)

func (b *DataTableHeaderBasic) ToJsonMap() (m map[string]any) {
	m = make(map[string]any)
	if b.Title != "" {
		m["title"] = b.Title
	}
	if b.Value != nil {
		m[":value"] = b.Value
	}
	if b.Key != "" {
		m["key"] = b.Key
	}
	if b.Align != "" {
		m["align"] = b.Align
	}
	if b.Width != "" {
		m["width"] = b.Width
	}
	if b.Children != nil {
		m["children"] = b.Children
	}
	if b.Index != 0 {
		m["index"] = b.Index
	}
	if b.Sortable {
		m["sortable"] = b.Sortable
	}
	if b.HeaderProps != nil {
		m["headerProps"] = b.HeaderProps
	}
	if b.CellProps != nil {
		m["cellProps"] = b.CellProps
	}
	if b.Hidden {
		m["hidden"] = b.Hidden
	}
	return
}

func (b *DataTableHeaderBasic) FromJsonMap(m map[string]any) {
	b.Title, _ = m["title"].(string)
	b.Value, _ = m[":value"]
	b.Key, _ = m["key"].(string)
	b.Align, _ = m["align"].(string)
	b.Width, _ = m["width"].(string)
	b.Children, _ = m["children"].([]any)
	b.Index, _ = m["index"].(int)
	b.Sortable, _ = m["sortable"].(bool)
	b.HeaderProps, _ = m["headerProps"].(map[string]any)
	b.CellProps, _ = m["cellProps"].(map[string]any)
	b.Hidden, _ = m["hidden"].(bool)
}

func (s DataTableHeaderBasicSlice) MarshalJSON() (b []byte, err error) {
	var s2 []any
	for _, b := range s {
		if !b.Hidden {
			s2 = append(s2, b)
		}
	}
	return json.Marshal(s2)
}

func (s DataTableHeaderSlice) MarshalJSON() (b []byte, err error) {
	var s2 []any
	for _, b := range s {
		if !b.Hidden {
			s2 = append(s2, b)
		}
	}
	return json.Marshal(s2)
}

func (d *DataTableHeaderBuilder) MarshalJSON() (b []byte, err error) {
	m := d.DataTableHeaderBasic.ToJsonMap()
	if d.Advanced != nil {
		for k, v := range d.Advanced {
			m[k] = v
		}
	}
	return json.Marshal(m)
}

func (d DataTableHeaderBuilder) Clone() *DataTableHeaderBuilder {
	adv := make(map[string]interface{}, len(d.Advanced))
	for k, v := range d.Advanced {
		adv[k] = v
	}
	d.Advanced = adv
	return &d
}

func DataTableHeader(b DataTableHeaderBasic) *DataTableHeaderBuilder {
	return &DataTableHeaderBuilder{DataTableHeaderBasic: b}
}

func DataTableHeaders(b ...DataTableHeaderBasic) (s DataTableHeaderSlice) {
	s = make([]*DataTableHeaderBuilder, len(b))
	for i, basic := range b {
		s[i] = DataTableHeader(basic)
	}
	return
}

func DataTableWithHeaderControl[T any](dataTable VTagBuilderGetter[T], comp h.HTMLComponent, config ...DataTableHeaderBasic) T {
	const key = "$headerControl$"
	var (
		tag  = dataTable.GetVTagBuilder()
		dst  = tag.GetHTMLTagBuilder()
		attr = tag.GetAttr(":headers")
		cfg  DataTableHeaderBasic
	)

	if attr == nil {
		for _, c := range tag.GetChildren() {
			if t, _ := c.(*VDataTableHeadersBuilder); t != nil {
				attr = t.GetAttr(":headers")
				dst = t.GetHTMLTagBuilder()
				break
			}
		}
	}

	v := attr.Value

	for _, cfg = range config {
	}

	cfg.Key = key

	switch t := v.(type) {
	case DataTableHeaderBasicSlice:
		cfg.Children = DataTableHeaders(t...)
	case DataTableHeaderSlice:
		cfg.Children = t
	}

	attr.Value = DataTableHeaders(cfg)
	dst.AppendChildren(web.Slot(comp).Name("header." + key).Scope("scope"))
	return tag.Dot()
}
