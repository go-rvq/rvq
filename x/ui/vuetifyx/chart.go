package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// VXChartBuilder renders <vx-chart>: a chart drawn by Apache ECharts — a
// line, bars, stacked bars, horizontal bars (a ranking) or a ring —, in the
// theme's colors, its texts never stretched. Passing the mouse shows the
// full label of the point (Hints) and every series with its value; the
// legend turns a series off; a range dragged zooms in it (↺ back to all); a
// point with a drill (Drill), clicked, emits point with it ("@point"). Its
// code is loaded from Src the first time a chart is shown (rvq
// x/ui/echarts, a lazy extra asset of the admin).
type VXChartBuilder struct {
	v.VTagBuilder[*VXChartBuilder]
}

// VXChartSeries is a named list of values, a color of the theme ("primary",
// "success"…) or a CSS color (none: the palette's).
type VXChartSeries struct {
	Name   string    `json:"name"`
	Color  string    `json:"color,omitempty"`
	Values []float64 `json:"values"`
}

// Chart kinds.
const (
	VXChartLine    = "line"
	VXChartBars    = "bars"
	VXChartStacked = "stacked"
	VXChartHBars   = "hbars"
	VXChartDonut   = "donut"
)

func VXChart() *VXChartBuilder {
	return v.VTag(&VXChartBuilder{}, "vx-chart")
}

// Src is where its code is: p.ExtraAssetRef(echarts.Path), the admin p
// serving it (LazyExtraAsset).
func (b *VXChartBuilder) Src(v string) *VXChartBuilder { return b.Attr("src", v) }

// Kind is line, bars, stacked, hbars or donut.
func (b *VXChartBuilder) Kind(v string) *VXChartBuilder { return b.Attr("kind", v) }

// Labels are the axis' (by step: some empty, not written) or the parts'
// (hbars, donut).
func (b *VXChartBuilder) Labels(v []string) *VXChartBuilder {
	return b.Attr(":labels", h.JSONString(nonNil(v)))
}

func (b *VXChartBuilder) Series(v ...VXChartSeries) *VXChartBuilder {
	for i := range v {
		if v[i].Values == nil {
			v[i].Values = []float64{}
		}
	}
	return b.Attr(":series", h.JSONString(nonNil(v)))
}

// Hints are the full label of each point (the tooltip's title).
func (b *VXChartBuilder) Hints(v []string) *VXChartBuilder {
	return b.Attr(":hints", h.JSONString(nonNil(v)))
}

// Colors are the donut's parts' (default the palette).
func (b *VXChartBuilder) Colors(v []string) *VXChartBuilder {
	return b.Attr(":colors", h.JSONString(nonNil(v)))
}

// Drill is what opens each point (nil: not opened), emitted with point as
// {index, …it}.
func (b *VXChartBuilder) Drill(v []any) *VXChartBuilder {
	return b.Attr(":drill", h.JSONString(nonNil(v)))
}

// Height is its height in pixels (0: by its kind).
func (b *VXChartBuilder) Height(v int) *VXChartBuilder { return b.Attr(":height", v) }

// Title names it for the readers of the screen.
func (b *VXChartBuilder) Title(v string) *VXChartBuilder { return b.Attr("title", v) }

// ResetLabel is the text of the ↺ that zooms out.
func (b *VXChartBuilder) ResetLabel(v string) *VXChartBuilder { return b.Attr("reset-label", v) }

func nonNil[T any](v []T) []T {
	if v == nil {
		return []T{}
	}
	return v
}
