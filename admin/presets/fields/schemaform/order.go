package schemaform

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// OrderOption is a field a list may be ordered by: the value saved, and the
// label shown — the application's, as its forms name the field.
type OrderOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// OrderOptionsFunc are the fields a type of order offers, in the words of the
// request (c.Event).
type OrderOptionsFunc func(c *Context) []OrderOption

// Order registers name as a type of ORDER: the fields options gives, by
// priority, each ascending, descending or not ordering — never text to write.
// It is saved as the list of the ordering ones, in priority order:
// `[{"field": "CreatedAt", "dir": "DESC"}, …]`.
//
// The form shows every field, one per line: the ordering ones first, in their
// priority (↑/↓ move them), each with ASC / DESC / — (not ordered). The detail
// shows the ordering ones: "Criado em ↓, Título ↑".
func (b *Builder) Order(name string, options OrderOptionsFunc) *Builder {
	return b.Type(name, OrderComponentFunc(options)).
		Display(name, OrderDisplayFunc(options)).
		TypeDecoder(name, DecodeOrder)
}

// OrderComponentFunc draws an order: see Builder.Order.
func OrderComponentFunc(options OrderOptionsFunc) ComponentFunc {
	return func(c *Context) h.HTMLComponent {
		opts := h.JSONString(options(c))
		val := c.Value
		// the ordering lines, of the fields offered, in priority; then the
		// others, not ordering
		rows := fmt.Sprintf(`[...(%[1]s || []).filter(r => r.dir && (%[2]s).some(o => o.value === r.field)), `+
			`...(%[2]s).filter(o => !(%[1]s || []).some(r => r.dir && r.field === o.value)).map(o => ({field: o.value, dir: ""}))]`,
			val, opts)
		label := fmt.Sprintf(`((%s).find(o => o.value === r.field) || {label: r.field}).label`, opts)
		// ASC or DESC keeps the field in its place (a new one goes last); none
		// takes it out. The handlers are functions (`(d) => …`, `() => …`):
		// Vue's runtime compiler takes a handler that starts like one — an
		// arrow called at once, `((a) => …)(x)` — for the function itself, and
		// runs it while rendering.
		setDir := fmt.Sprintf(`(d) => { const f = r.field; const cur = (%[1]s || []).filter(x => x.dir); `+
			`const at = cur.findIndex(x => x.field === f); `+
			`if (!d) { %[1]s = cur.filter(x => x.field !== f) } `+
			`else if (at >= 0) { const n = [...cur]; n[at] = {field: f, dir: d}; %[1]s = n } `+
			`else { %[1]s = [...cur, {field: f, dir: d}] } }`, val)
		move := func(delta int) string {
			return fmt.Sprintf(`() => { const n = [...(%[1]s || []).filter(x => x.dir)]; const j = i + (%[2]d); `+
				`if (j < 0 || j >= n.length) return; [n[i], n[j]] = [n[j], n[i]]; %[1]s = n }`, val, delta)
		}
		count := fmt.Sprintf(`(%s || []).filter(r => r.dir).length`, val)

		line := h.Div(
			v.VBtn("").Icon("mdi-arrow-up").Size("x-small").Variant("text").
				Attr(":disabled", "!r.dir || i === 0").Attr("@click", move(-1)),
			v.VBtn("").Icon("mdi-arrow-down").Size("x-small").Variant("text").
				Attr(":disabled", fmt.Sprintf("!r.dir || i >= %s - 1", count)).Attr("@click", move(1)),
			h.Span("").Attr("v-text", `r.dir ? i + 1 : ""`).
				Class("text-caption text-medium-emphasis mx-2").Style("min-width: 1.2em"),
			h.Span("").Attr("v-text", label).Class("flex-grow-1").Attr(":class", `{"text-disabled": !r.dir}`),
			v.VBtnToggle(
				v.VBtn("ASC").Value("ASC").Attr("prepend-icon", "mdi-arrow-up-thin"),
				v.VBtn("DESC").Value("DESC").Attr("prepend-icon", "mdi-arrow-down-thin"),
				v.VBtn("—").Value(""),
			).Density("compact").Variant("outlined").Divided(true).
				Attr(":model-value", "r.dir").Attr("@update:model-value", setDir),
		).Class("d-flex align-center py-1").
			Attr("v-for", "(r, i) in "+rows).Attr(":key", "r.field")

		comps := []h.HTMLComponent{}
		if l := c.Label(); l != "" {
			comps = append(comps, h.Div(h.Text(l)).Class("text-subtitle-2 mb-1"))
		}
		comps = append(comps, h.Div(line).Class("border rounded px-2"))
		if hint := c.Hint(); hint != "" {
			comps = append(comps, h.Div(h.Text(hint)).Class("text-caption text-medium-emphasis mt-1"))
		}
		return h.Div(comps...).Class("mb-4")
	}
}

// OrderDisplayFunc shows an order: its ordering fields by their labels, ↑
// ascending and ↓ descending; "—" when there is none.
func OrderDisplayFunc(options OrderOptionsFunc) ComponentFunc {
	return func(c *Context) h.HTMLComponent {
		labels := map[string]string{}
		for _, o := range options(c) {
			labels[o.Value] = o.Label
		}
		var parts []string
		for _, it := range orderItems(c.Data) {
			l := labels[it["field"]]
			if l == "" {
				l = it["field"]
			}
			arrow := "↑"
			if it["dir"] == "DESC" {
				arrow = "↓"
			}
			parts = append(parts, l+" "+arrow)
		}
		text := strings.Join(parts, ", ")
		if text == "" {
			return emptyDisplay(c)
		}
		if c.Compact {
			return h.Span(text)
		}
		return h.Div(
			h.Div(h.Text(displayLabel(c))).Class("text-caption text-medium-emphasis"),
			h.Div(h.Text(text)),
		).Class("mb-2")
	}
}

// DecodeOrder reads an order the form posted under key — `key[0].field`,
// `key[0].dir`, … —: the ordering items, ASC or DESC, in their order.
func DecodeOrder(values url.Values, key string) any {
	out := []any{}
	for i := 0; ; i++ {
		at := key + "[" + strconv.Itoa(i) + "]"
		if !hasKeyUnder(values, at) {
			break
		}
		field, dir := values.Get(at+".field"), strings.ToUpper(values.Get(at+".dir"))
		if field == "" || (dir != "ASC" && dir != "DESC") {
			continue
		}
		out = append(out, map[string]any{"field": field, "dir": dir})
	}
	return out
}

// orderItems is an order as it was decoded or saved: its items of a field and
// a direction.
func orderItems(v any) (out []map[string]string) {
	list, _ := v.([]any)
	for _, it := range list {
		m, _ := it.(map[string]any)
		field, _ := m["field"].(string)
		dir, _ := m["dir"].(string)
		if field != "" && dir != "" {
			out = append(out, map[string]string{"field": field, "dir": strings.ToUpper(dir)})
		}
	}
	return
}
