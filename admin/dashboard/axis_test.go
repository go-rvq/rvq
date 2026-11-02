package dashboard

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

var svgText = regexp.MustCompile(`<text x="([0-9.]+)" y="[0-9.]+" font-size="10" text-anchor="(\w+)"[^>]*>([^<]*)</text>`)

// The labels under the axis: none over another, none out of the chart —
// the first and the last anchored inside.
func TestAxisLabelsFit(t *testing.T) {
	labels := make([]string, 30)
	for i := range labels {
		labels[i] = fmt.Sprintf("01/11 21:%02d", i)
	}
	for _, c := range []*Chart{
		{Kind: ChartLine, Labels: labels, Series: []Series{{Values: make([]float64, 30)}}},
		{Kind: ChartBars, Labels: labels, Series: []Series{{Values: make([]float64, 30)}}},
	} {
		svg := lineSVG(c)
		if c.Kind == ChartBars {
			svg = barsSVG(c)
		}
		m := svgText.FindAllStringSubmatch(svg, -1)
		if len(m) < 2 || len(m) > 9 {
			t.Fatalf("%s: %d labels drawn", c.Kind, len(m))
		}
		right := -1.0
		for i, l := range m {
			x, _ := strconv.ParseFloat(l[1], 64)
			lw := float64(len(l[3])) * 6
			left := x - lw/2
			switch l[2] {
			case "start":
				left = x
			case "end":
				left = x - lw
			}
			if left < 0 || left+lw > 600.01 {
				t.Errorf("%s: %q out of the chart (%.1f)", c.Kind, l[3], left)
			}
			if i > 0 && left < right {
				t.Errorf("%s: %q over the one before", c.Kind, l[3])
			}
			right = left + lw
		}
		if m[0][2] != "start" {
			t.Errorf("%s: the first not anchored at its start: %v", c.Kind, m[0][2])
		}
	}
}
