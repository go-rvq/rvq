package admin

import (
	"context"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// revertActionName is the ParamAction of the confirmation action registered on
// the revisions child's Detailing. All revert buttons open it (a dialog showing
// what will change, grouped by field); confirming runs the revert.
const revertActionName = "revert"

// installRevertConfirm registers the revert confirmation as a Detailing action
// (the ActionForm pattern: a ComponentFunc renders the preview, an UpdateFunc
// applies on confirm). The revert parameters travel as query values and are
// preserved to the confirm submit by the form's MergeQuery.
func (mh *ModelHistory) installRevertConfirm(child *presets.ModelBuilder) {
	child.Detailing().
		Action(revertActionName).
		Icon("mdi-history").
		SetI18nLabel(func(c context.Context) string { return getMessages(c).Revert }).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			return mh.revertPreview(ctx), nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			return mh.applyRevert(ctx)
		})
}

// revertButtonClick builds the @click that opens the revert confirmation with
// the given scope. hash is the target revision; field/fields/hunks scope it (all
// empty = whole record). hunksExpr, when set, is a JS expression (locals.hunks).
func (mh *ModelHistory) revertButtonClick(recordKey, targetHash, field, fields, hunksExpr string) string {
	b := web.Plaid().
		EventFunc(actions.Action).
		Query(presets.ParamAction, revertActionName).
		Query(presets.ParamID, targetHash).
		Query(presets.ParamTargetPortal, actions.Dialog.PortalName()).
		Query(presets.ParamOverlay, actions.Dialog).
		Query("record", recordKey).
		Query("hash", targetHash).
		MergeQuery(true)
	if field != "" {
		b.Query("field", field)
	}
	if fields != "" {
		b.Query("fields", fields)
	}
	if hunksExpr != "" {
		b.Query("hunks", web.Var(hunksExpr))
	}
	return b.Go()
}

// applyRevert reads the scope from the request and reverts the parent record to
// the target revision (whole record, named field(s), or selected hunks of one
// field), then navigates to the parent detail. Each revert records a new
// revision (git-revert).
func (mh *ModelHistory) applyRevert(ctx *web.EventContext) (err error) {
	recordKey := ctx.R.FormValue("record")
	if recordKey == "" {
		recordKey = parentRecordKey(ctx)
	}
	hash, err := decodeHash(ctx.R.FormValue("hash"))
	if err != nil {
		return
	}
	id, err := mh.mb.ParseRecordID(recordKey)
	if err != nil {
		return
	}
	obj := mh.mb.NewModel()
	id.SetTo(obj)
	if err = mh.db.First(obj).Error; err != nil {
		return
	}

	field := ctx.R.FormValue("field")
	switch {
	case field != "" && ctx.R.FormValue("hunks") != "":
		var target string
		if rev, rerr := mh.Revision(recordKey, hash); rerr == nil {
			if m, merr := fieldMap(rev); merr == nil {
				target = fieldValue(m, field)
			}
		}
		selected := map[int]bool{}
		_ = ctx.R.ParseForm()
		for _, s := range ctx.R.Form["hunks"] {
			for _, part := range strings.Split(s, ",") {
				if nnum, e := strconv.Atoi(strings.TrimSpace(part)); e == nil {
					selected[nnum] = true
				}
			}
		}
		value := applyHunks(fieldStringValue(obj, field), target, selected, mh.IsHTML(field))
		err = mh.RevertFieldContent(obj, field, value, ctx)
	case field != "":
		err = mh.RevertField(obj, hash, field, ctx)
	case ctx.R.FormValue("fields") != "":
		err = mh.RevertFields(obj, hash, splitCSV(ctx.R.FormValue("fields")), ctx)
	default:
		err = mh.RevertRecord(obj, hash, ctx)
	}
	if err != nil {
		return
	}

	if ctx.Resp != nil {
		presets.ShowMessage(ctx.Resp, getMessages(ctx.Context()).Reverted, "success")
		ctx.Resp.PushState = web.Location(nil).URL(mh.mb.Info().DetailingHref(recordKey))
	}
	return
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// changeRow is one line of the revert preview: an addition or a removal.
type changeRow struct {
	add  bool
	text string
}

// revertChanges computes, per affected field, the change rows between the current
// value and the value after the revert. Only fields that actually change appear.
func (mh *ModelHistory) revertChanges(ctx *web.EventContext) (recordKey string, byField [][2]any, err error) {
	recordKey = ctx.R.FormValue("record")
	if recordKey == "" {
		recordKey = parentRecordKey(ctx)
	}
	hash, err := decodeHash(ctx.R.FormValue("hash"))
	if err != nil {
		return
	}
	cur, err := mh.currentRecord(recordKey)
	if err != nil {
		return
	}
	curMap, _, err := mh.snapshot(cur)
	if err != nil {
		return
	}
	rev, err := mh.Revision(recordKey, hash)
	if err != nil {
		return
	}
	targetMap, err := fieldMap(rev)
	if err != nil {
		return
	}

	field := ctx.R.FormValue("field")
	var fields []string
	switch {
	case field != "":
		fields = []string{field}
	case ctx.R.FormValue("fields") != "":
		fields = splitCSV(ctx.R.FormValue("fields"))
	default:
		fields = mh.resolved
	}

	_ = ctx.R.ParseForm()
	hunksSel := map[int]bool{}
	partial := field != "" && ctx.R.FormValue("hunks") != ""
	if partial {
		for _, s := range ctx.R.Form["hunks"] {
			for _, part := range strings.Split(s, ",") {
				if nnum, e := strconv.Atoi(strings.TrimSpace(part)); e == nil {
					hunksSel[nnum] = true
				}
			}
		}
	}

	for _, f := range fields {
		curVal := fieldValue(curMap, f)
		var newVal string
		if partial && f == field {
			newVal = applyHunks(fieldStringValue(cur, f), fieldValue(targetMap, f), hunksSel, mh.IsHTML(f))
		} else {
			newVal = fieldValue(targetMap, f)
		}
		if curVal == newVal {
			continue
		}
		rows := mh.fieldChangeRows(curVal, newVal, mh.IsHTML(f))
		byField = append(byField, [2]any{f, rows})
	}
	return
}

// fieldChangeRows lists the additions and removals turning current into newVal.
func (mh *ModelHistory) fieldChangeRows(current, newVal string, htmlMode bool) []changeRow {
	var rows []changeRow
	for _, df := range diffOps(current, newVal, htmlMode) {
		switch df.Type {
		case diffmatchpatch.DiffDelete:
			rows = append(rows, changeRow{add: false, text: df.Text})
		case diffmatchpatch.DiffInsert:
			rows = append(rows, changeRow{add: true, text: df.Text})
		}
	}
	return rows
}

// revertPreview renders the confirmation body: the changes to be reverted,
// grouped by field, as a table identifying additions and removals.
func (mh *ModelHistory) revertPreview(ctx *web.EventContext) h.HTMLComponent {
	msgr := getMessages(ctx.Context())
	_, byField, err := mh.revertChanges(ctx)
	if err != nil {
		return v.VAlert(h.Text(err.Error())).Type(v.TypeError).Variant(v.VariantTonal)
	}
	if len(byField) == 0 {
		return v.VAlert(h.Text(msgr.NoChanges)).Type("info").Variant(v.VariantTonal)
	}

	var sections h.HTMLComponents
	for _, fe := range byField {
		field := fe[0].(string)
		rows := fe[1].([]changeRow)
		htmlMode := mh.IsHTML(field)

		var trs h.HTMLComponents
		for _, r := range rows {
			var content h.HTMLComponent
			if htmlMode {
				content = h.RawHTML(r.text)
			} else {
				content = h.Text(r.text)
			}
			sign, cls := "−", "text-error"
			if r.add {
				sign, cls = "+", "text-success"
			}
			trs = append(trs, h.Tr(
				h.Td(h.Strong(sign)).Class(cls).Style("width:1.5em;vertical-align:top"),
				h.Td(content).Class(cls),
			))
		}

		sections = append(sections, h.Div(
			h.Div(h.Strong(field)).Class("text-subtitle-2 mb-1"),
			v.VTable(h.Tbody(trs...)).Density(v.DensityCompact).Class("mb-3").
				Attr("style", "white-space:pre-wrap"),
		))
	}

	return h.Div(
		v.VAlert(h.Text(msgr.RevertConfirmHint)).Type("warning").Variant(v.VariantTonal).Density(v.DensityCompact).Class("mb-3"),
		h.Div(sections...),
	)
}
