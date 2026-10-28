package mail_sender

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// subjectRow is a row of the table of the parts of the subject: a family
// registered (SubjectFamilies), or the others.
type subjectRow struct {
	name  string
	label string
}

// subjectRows are the rows of the table: the families, then the others.
func subjectRows(ctx context.Context) []subjectRow {
	var rows []subjectRow
	for _, f := range SubjectFamilies() {
		label := f.Name
		if f.Label != nil {
			label = f.Label(ctx)
		}
		rows = append(rows, subjectRow{f.Name, label})
	}
	return append(rows, subjectRow{SubjectFamilyOther, GetMessages(ctx).SubjectFamilyOther})
}

// subjectRulesComponent is the table of the parts of the subject (SubjectRules):
// a row by type of message, by whom it goes to (the team, outside) the parts
// it puts — boxes to mark, editing; marks, seen. Its value is the rules of
// every row as they are now (each its own, else the others', else the
// default), as JSON.
func subjectRulesComponent(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	rules, _ := field.Value().(SubjectRules)
	m := GetMessages(ctx.Context())
	edit := !field.ReadOnly && field.Mode.Dot().HasAny(presets.EDIT, presets.NEW)
	rows := subjectRows(ctx.Context())
	now := SubjectRules{}
	for _, r := range rows {
		now[r.name] = rules.Of(r.name)
	}
	parts := []struct {
		bit   SubjectParts
		label string
	}{{SubjectSite, m.SubjectPartSite}, {SubjectTopic, m.SubjectPartTopic}, {SubjectRef, m.SubjectPartRef}}
	cell := func(row, side string, set SubjectParts) h.HTMLComponent {
		var boxes h.HTMLComponents
		for _, p := range parts {
			name := row + "." + side + "." + strconv.Itoa(int(p.bit))
			if !edit {
				icon := "mdi-checkbox-blank-outline"
				if set.Has(p.bit) {
					icon = "mdi-checkbox-marked"
				}
				boxes = append(boxes, h.Span("").Children(v.VIcon(icon).Size(v.SizeSmall).Class("me-1"), h.Text(p.label)).
					Class("me-3 text-no-wrap").Attr("data-subject-part", name).Attr("data-subject-set", set.Has(p.bit)))
				continue
			}
			key, bit := strconv.Quote(row), p.bit
			get := fmt.Sprintf("(((JSON.parse(fieldValue.value || '{}')[%s] || {}).%s || 0) & %d) === %d", key, side, bit, bit)
			put := fmt.Sprintf(`(on) => { const r = JSON.parse(fieldValue.value || '{}'); const x = r[%s] || {team: 0, external: 0};`+
				` x.%s = on ? ((x.%s || 0) | %d) : ((x.%s || 0) & ~%d); r[%s] = x; fieldValue.value = JSON.stringify(r) }`,
				key, side, side, bit, side, bit, key)
			boxes = append(boxes, v.VCheckbox().Label(p.label).Density(v.DensityCompact).HideDetails(true).
				Attr(":model-value", get).Attr("@update:model-value", put).Class("me-3 flex-grow-0").
				Attr("data-subject-part", row+"."+side+"."+strconv.Itoa(int(p.bit))))
		}
		return h.Td(h.Div(boxes...).Class("d-flex flex-wrap align-center"))
	}
	var body h.HTMLComponents
	for _, r := range rows {
		rule := now[r.name]
		body = append(body, h.Tr(h.Td(h.Text(r.label)), cell(r.name, "team", rule.Team), cell(r.name, "external", rule.External)).
			Attr("data-subject-family", r.name))
	}
	table := v.VTable(
		h.Thead(h.Tr(h.Th(m.SubjectRulesFamily), h.Th(m.SubjectRulesTeam), h.Th(m.SubjectRulesExternal))),
		h.Tbody(body...),
	).Density(v.DensityCompact)
	field.LoadHint()
	if !edit {
		return vx.VXReadonlyField(table).Label(field.Label)
	}
	raw, _ := json.Marshal(now)
	return h.Div(
		h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
		vue.FormField(table).Value(field.FormKey, string(raw)),
		h.If(field.Hint != "", h.Div(h.Text(field.Hint)).Class("text-caption text-medium-emphasis")),
	).Class("mb-4").Attr("data-subject-rules", field.FormKey)
}

// subjectRulesSetter takes the rules posted (their JSON): of the rows of the
// table only.
func subjectRulesSetter(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
	raw := ctx.R.FormValue(field.FormKey)
	if raw == "" {
		return nil
	}
	var in SubjectRules
	if err := json.Unmarshal([]byte(raw), &in); err != nil {
		return fmt.Errorf("subject parts: %w", err)
	}
	out := SubjectRules{}
	for _, r := range subjectRows(ctx.Context()) {
		if rule, ok := in[r.name]; ok {
			out[r.name] = SubjectRule{Team: rule.Team & SubjectAll, External: rule.External & SubjectAll}
		}
	}
	obj.(*MailSender).SubjectRules = out
	return nil
}
