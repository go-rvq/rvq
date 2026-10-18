package dashboard

import (
	"fmt"
	"html"
	"math"
	"strings"
)

// The charts are SVG in the theme's colors: a series' color is a CSS color
// (the Vuetify theme's, by its name: "primary", "success", "warning", "error",
// "info", "secondary"), so they follow the light and the dark theme.

// palette is the order of the series' colors.
var palette = []string{"primary", "info", "warning", "success", "secondary", "error"}

// cssColor is the CSS color of a theme color's name (or a CSS color as it is).
func cssColor(c string) string {
	switch c {
	case "primary", "secondary", "success", "warning", "error", "info", "on-surface", "surface":
		return "rgb(var(--v-theme-" + c + "))"
	}
	return c
}

func seriesColor(s Series, i int) string {
	if s.Color != "" {
		return cssColor(s.Color)
	}
	return cssColor(palette[i%len(palette)])
}

const (
	muted = "rgba(var(--v-theme-on-surface), 0.55)"
	grid  = "rgba(var(--v-theme-on-surface), 0.10)"
)

type svgBuf struct{ strings.Builder }

func (b *svgBuf) f(format string, args ...any) { fmt.Fprintf(&b.Builder, format, args...) }

func svgOpen(w, h int, label string) *svgBuf {
	b := &svgBuf{}
	b.f(`<svg viewBox="0 0 %d %d" width="100%%" preserveAspectRatio="none" role="img" aria-label="%s" style="display:block;overflow:visible">`, w, h, html.EscapeString(label))
	return b
}

func (b *svgBuf) text(x, y float64, s string, size int, anchor string) {
	b.f(`<text x="%.1f" y="%.1f" font-size="%d" text-anchor="%s" style="fill:%s">%s</text>`, x, y, size, anchor, muted, html.EscapeString(s))
}

func maxOf(series []Series, stacked bool) float64 {
	m := 0.0
	if stacked && len(series) > 0 {
		for i := range series[0].Values {
			t := 0.0
			for _, s := range series {
				if i < len(s.Values) {
					t += s.Values[i]
				}
			}
			m = math.Max(m, t)
		}
		return m
	}
	for _, s := range series {
		for _, v := range s.Values {
			m = math.Max(m, v)
		}
	}
	return m
}

// lineSVG draws series as lines (the first filled), labels under the axis.
func lineSVG(c *Chart) string {
	const w, h, pad = 600, 180, 18.0
	b := svgOpen(w, h+18, c.Title)
	for i := 1; i <= 3; i++ {
		y := h - float64(i)*h/4
		b.f(`<line x1="0" y1="%.1f" x2="%d" y2="%.1f" style="stroke:%s"/>`, y, w, y, grid)
	}
	b.f(`<line x1="0" y1="%d" x2="%d" y2="%d" style="stroke:%s"/>`, h, w, h, grid)
	mx := maxOf(c.Series, false)
	if mx == 0 {
		mx = 1
	}
	for k, s := range c.Series {
		n := len(s.Values)
		if n == 0 {
			continue
		}
		pts := make([]string, n)
		for i, v := range s.Values {
			x := 0.0
			if n > 1 {
				x = float64(i) * w / float64(n-1)
			}
			pts[i] = fmt.Sprintf("%.1f,%.1f", x, h-v/mx*(h-pad))
		}
		col := seriesColor(s, k)
		if k == 0 {
			b.f(`<polygon points="0,%d %s %d,%d" style="fill:%s;opacity:.15"/>`, h, strings.Join(pts, " "), w, h, col)
		}
		b.f(`<polyline points="%s" style="fill:none;stroke:%s;stroke-width:2;vector-effect:non-scaling-stroke"/>`, strings.Join(pts, " "), col)
	}
	for i, l := range c.Labels {
		if l == "" || len(c.Labels) < 2 {
			continue
		}
		b.text(float64(i)*w/float64(len(c.Labels)-1), h+14, l, 10, "middle")
	}
	b.f(`</svg>`)
	return b.String()
}

// barsSVG draws series as bars, side by side or stacked ((c.Kind == ChartStacked)).
func barsSVG(c *Chart) string {
	const w, h, pad = 600, 180, 12.0
	b := svgOpen(w, h+18, c.Title)
	b.f(`<line x1="0" y1="%d" x2="%d" y2="%d" style="stroke:%s"/>`, h, w, h, grid)
	mx := maxOf(c.Series, (c.Kind == ChartStacked))
	if mx == 0 {
		mx = 1
	}
	n := 0
	for _, s := range c.Series {
		n = max(n, len(s.Values))
	}
	if n == 0 {
		b.f(`</svg>`)
		return b.String()
	}
	slot := float64(w) / float64(n)
	bw := slot * 0.62
	for i := 0; i < n; i++ {
		x := float64(i)*slot + (slot-bw)/2
		base := float64(h)
		for k, s := range c.Series {
			if i >= len(s.Values) {
				continue
			}
			v := s.Values[i]
			bh := v / mx * (h - pad)
			if c.Kind == ChartStacked {
				base -= bh
				b.f(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="1.5" style="fill:%s"><title>%s: %s</title></rect>`,
					x, base, bw, bh, seriesColor(s, k), html.EscapeString(s.Name), formatNumber(v))
			} else {
				sw := bw / float64(len(c.Series))
				b.f(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="1.5" style="fill:%s"><title>%s: %s</title></rect>`,
					x+float64(k)*sw, h-bh, sw, bh, seriesColor(s, k), html.EscapeString(s.Name), formatNumber(v))
			}
		}
	}
	for i, l := range c.Labels {
		if l == "" {
			continue
		}
		b.text(float64(i)*slot+slot/2, h+14, l, 10, "middle")
	}
	b.f(`</svg>`)
	return b.String()
}

// sparkSVG is a small line of values (a Metric's).
func sparkSVG(values []float64) string {
	const w, h = 100, 28.0
	n := len(values)
	if n < 2 {
		return ""
	}
	mn, mx := values[0], values[0]
	for _, v := range values {
		mn, mx = math.Min(mn, v), math.Max(mx, v)
	}
	span := mx - mn
	if span == 0 {
		span = 1
	}
	pts := make([]string, n)
	for i, v := range values {
		pts[i] = fmt.Sprintf("%.1f,%.1f", float64(i)*w/float64(n-1), h-2-(v-mn)/span*(h-4))
	}
	return fmt.Sprintf(`<svg viewBox="0 0 %d %.0f" width="100" height="28" aria-hidden="true"><polyline points="%s" style="fill:none;stroke:%s;stroke-width:2"/></svg>`,
		w, h, strings.Join(pts, " "), cssColor("primary"))
}

// donutSVG is a ring of the parts of one series.
func donutSVG(c *Chart) string {
	const r = 40.0
	b := &svgBuf{}
	b.f(`<svg viewBox="0 0 100 100" width="120" height="120" role="img" aria-label="%s">`, html.EscapeString(c.Title))
	if len(c.Series) == 0 {
		b.f(`</svg>`)
		return b.String()
	}
	s := c.Series[0]
	total := Total(s.Values)
	if total == 0 {
		b.f(`<circle cx="50" cy="50" r="%.0f" style="fill:none;stroke:%s;stroke-width:16"/></svg>`, r, grid)
		return b.String()
	}
	a := -math.Pi / 2
	for i, v := range s.Values {
		if v <= 0 {
			continue
		}
		e := a + 2*math.Pi*v/total
		if v == total {
			e -= 0.0001
		}
		large := 0
		if e-a > math.Pi {
			large = 1
		}
		col := cssColor(palette[i%len(palette)])
		if i < len(c.Colors) && c.Colors[i] != "" {
			col = cssColor(c.Colors[i])
		}
		b.f(`<path d="M%.2f,%.2f A%.0f,%.0f 0 %d 1 %.2f,%.2f" style="fill:none;stroke:%s;stroke-width:16"><title>%s: %s</title></path>`,
			50+r*math.Cos(a), 50+r*math.Sin(a), r, r, large, 50+r*math.Cos(e), 50+r*math.Sin(e), col,
			html.EscapeString(labelAt(c.Labels, i)), formatNumber(v))
		a = e
	}
	b.f(`</svg>`)
	return b.String()
}

func labelAt(labels []string, i int) string {
	if i < len(labels) {
		return labels[i]
	}
	return ""
}

// formatNumber is n with thousands grouped, no decimals for integers.
func formatNumber(n float64) string {
	if n != math.Trunc(n) {
		return fmt.Sprintf("%.1f", n)
	}
	s := fmt.Sprintf("%d", int64(math.Abs(n)))
	var out []byte
	for i := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, '.')
		}
		out = append(out, s[i])
	}
	if n < 0 {
		return "-" + string(out)
	}
	return string(out)
}

// FormatBytes is n bytes as "1,4 GB".
func FormatBytes(n float64) string {
	units := []string{"B", "kB", "MB", "GB", "TB"}
	i := 0
	for n >= 1000 && i < len(units)-1 {
		n /= 1000
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%.0f %s", n, units[i])
	}
	return strings.Replace(fmt.Sprintf("%.1f %s", n, units[i]), ".", ",", 1)
}
