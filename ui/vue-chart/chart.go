package vue_chart

import (
	"github.com/qor5/web/v3/tag"
)

type (
	ChartBuilderGetter[T any] interface {
		tag.TagBuilderGetter[T]
		GetChartBuilder() *ChartBuilder[T]
	}

	ChartBuilder[T any] struct {
		tag.TagBuilder[T]
	}

	BarBuilder struct {
		ChartBuilder[*BarBuilder]
	}

	DoughnutBuilder struct {
		ChartBuilder[*DoughnutBuilder]
	}

	LineBuilder struct {
		ChartBuilder[*LineBuilder]
	}

	PieBuilder struct {
		ChartBuilder[*PieBuilder]
	}

	PolarAreaBuilder struct {
		ChartBuilder[*PolarAreaBuilder]
	}

	RadarBuilder struct {
		ChartBuilder[*RadarBuilder]
	}

	BubbleBuilder struct {
		ChartBuilder[*BubbleBuilder]
	}

	ScatterBuilder struct {
		ChartBuilder[*ScatterBuilder]
	}
)

func (c *ChartBuilder[T]) GetChartBuilder() *ChartBuilder[T] {
	return c
}

func (c *ChartBuilder[T]) Data(v any) T {
	return c.Attr("data", v)
}

func (c *ChartBuilder[T]) DataExpr(v any) T {
	return c.Attr(":data", v)
}

func (c *ChartBuilder[T]) Options(v any) T {
	return c.Attr("options", v)
}

func (c *ChartBuilder[T]) OptionsExpr(v any) T {
	return c.Attr(":options", v)
}

func Chart[T ChartBuilderGetter[T]](dot T, typ string) T {
	tb := dot.GetChartBuilder()
	tb.TagBuilder = *tag.NewTag(dot, "chart-"+typ).GetTagBuilder()
	return dot
}

func Bar() *BarBuilder {
	return Chart(&BarBuilder{}, "bar")
}

func Doughnut() *DoughnutBuilder {
	return Chart(&DoughnutBuilder{}, "doughnut")
}

func Line() *LineBuilder {
	return Chart(&LineBuilder{}, "line")
}

func Pie() *PieBuilder {
	return Chart(&PieBuilder{}, "pie")
}

func PolarArea() *PolarAreaBuilder {
	return Chart(&PolarAreaBuilder{}, "polar-area")
}

func Radar() *RadarBuilder {
	return Chart(&RadarBuilder{}, "radar")
}

func Bubble() *BubbleBuilder {
	return Chart(&BubbleBuilder{}, "bubble")
}

func Scatter() *ScatterBuilder {
	return Chart(&ScatterBuilder{}, "scatter")
}
