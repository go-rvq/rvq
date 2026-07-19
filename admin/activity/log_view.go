package activity

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// LogViewAction is the name of the detail action that opens the activity log
// dialog of a record.
const LogViewAction = "activityLog"

// RecordLogs returns the activity logs of the record identified by the preset
// record id, most recent first.
func (mb *ModelBuilder) RecordLogs(id string) (logs []ActivityLogInterface, err error) {
	keys := mb.modelKeysFromRecordID(id)

	slicePtr := mb.activity.NewLogModelSlice()
	if err = mb.activity.db.
		Where("model_name = ? AND model_keys = ?", mb.typ.Name(), keys).
		Order("created_at DESC").
		Find(slicePtr).Error; err != nil {
		return
	}

	sliceValue := reflect.ValueOf(slicePtr).Elem()
	for i := 0; i < sliceValue.Len(); i++ {
		if log, ok := sliceValue.Index(i).Interface().(ActivityLogInterface); ok {
			logs = append(logs, log)
		}
	}
	return
}

// InstallLogViewAction adds a detail action to the preset model builder that
// opens, in a dialog, the activity log of the current record: who, when, from
// which IP and browser, and — for edits — the per-field before/after diff.
//
// It is installed automatically for every model registered with
// RegisterModel (see installModelBuilder); call it manually only for models
// configured without the activity preset wiring.
func (mb *ModelBuilder) InstallLogViewAction(pmb *presets.ModelBuilder) {
	if pmb == nil || !pmb.HasDetailing() {
		return
	}

	pmb.Detailing().Action(LogViewAction).
		Icon("mdi-history").
		SetI18nLabel(func(ctx context.Context) string {
			return getMessages(ctx).LogAction
		}).
		ComponentFunc(func(id string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
			logs, err := mb.RecordLogs(id)
			if err != nil {
				return nil, err
			}
			return mb.logTimeline(ctx, logs), nil
		})
}

// modelKeysFromRecordID maps a preset record id to the stored ModelKeys value.
// A record id joins the primary-key values with "_" (model.ID.String); the
// stored ModelKeys (KeysValue) joins them with ":". For the common single-key
// case both equal the raw value, so no reconstruction is needed.
func (mb *ModelBuilder) modelKeysFromRecordID(id string) string {
	if len(mb.keys) <= 1 {
		return id
	}
	return strings.ReplaceAll(id, "_", ":")
}

func getMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nActivityKey, Messages_en_US).(*Messages)
}

// logTimeline renders logs as a vuetify timeline, newest first.
func (mb *ModelBuilder) logTimeline(ctx *web.EventContext, logs []ActivityLogInterface) h.HTMLComponent {
	msgr := getMessages(ctx.Context())

	if len(logs) == 0 {
		return v.VAlert(h.Text(msgr.LogEmpty)).
			Type("info").Variant(v.VariantTonal).Density(v.DensityComfortable)
	}

	timeline := v.VTimeline().
		Attr("side", "end").
		Density(v.DensityComfortable).
		Class("pa-4")

	for _, log := range logs {
		timeline.AppendChild(mb.logItem(msgr, log))
	}
	return timeline
}

func (mb *ModelBuilder) logItem(msgr *Messages, log ActivityLogInterface) h.HTMLComponent {
	head := h.Div(
		h.Strong(actionLabel(msgr, log.GetAction())),
		h.Span(" · "+log.GetCreator()).Class("text-medium-emphasis"),
		h.Span(" · "+log.GetCreatedAt().Format("2006-01-02 15:04:05")).Class("text-medium-emphasis"),
	)

	meta := h.HTMLComponents{}
	if s, ok := log.(interface{ GetIP() string }); ok && s.GetIP() != "" {
		meta = append(meta, chip(msgr.ModelIP+": "+s.GetIP()))
	}
	if s, ok := log.(interface{ GetUserAgent() string }); ok && s.GetUserAgent() != "" {
		meta = append(meta, chip(msgr.ModelUserAgent+": "+s.GetUserAgent()))
	}

	body := h.HTMLComponents{head}
	if len(meta) > 0 {
		body = append(body, h.Div(meta...).Class("d-flex flex-wrap ga-1 mt-1"))
	}
	if diff := diffTable(msgr, log.GetModelDiffs()); diff != nil {
		body = append(body, diff)
	}

	return v.VTimelineItem(
		h.Div(body...),
	).DotColor(actionColor(log.GetAction())).Size(v.SizeSmall)
}

func chip(text string) h.HTMLComponent {
	return v.VChip(h.Text(text)).Size(v.SizeXSmall).Variant(v.VariantOutlined)
}

func actionLabel(msgr *Messages, action string) string {
	switch action {
	case ActivityCreate:
		return msgr.ActionCreate
	case ActivityEdit:
		return msgr.ActionEdit
	case ActivityDelete:
		return msgr.ActionDelete
	case ActivityView:
		return msgr.ActionView
	}
	return action
}

func actionColor(action string) string {
	switch action {
	case ActivityCreate:
		return "success"
	case ActivityEdit:
		return "warning"
	case ActivityDelete:
		return "error"
	}
	return "info"
}

// diffTable renders the JSON diff (from ModelDiffs) as a field/old/new table.
func diffTable(msgr *Messages, raw string) h.HTMLComponent {
	if raw == "" {
		return nil
	}
	var diffs []Diff
	if err := json.Unmarshal([]byte(raw), &diffs); err != nil || len(diffs) == 0 {
		return nil
	}
	sort.SliceStable(diffs, func(i, j int) bool { return diffs[i].Field < diffs[j].Field })

	rows := h.HTMLComponents{
		h.Tr(
			h.Th(msgr.DiffField),
			h.Th(msgr.DiffOld),
			h.Th(msgr.DiffNow),
		),
	}
	for _, d := range diffs {
		rows = append(rows, h.Tr(
			h.Td(h.Text(d.Field)),
			h.Td(h.Text(d.Old)).Class("text-medium-emphasis"),
			h.Td(h.Text(d.Now)),
		))
	}
	return v.VTable(h.Tbody(rows...)).Density(v.DensityCompact).Class("mt-2")
}
