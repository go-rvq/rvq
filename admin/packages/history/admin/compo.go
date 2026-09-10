package admin

import (
	"encoding/hex"
	"encoding/json"

	h "github.com/go-rvq/htmlgo"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
)

func shortHash(hash histmodels.Hash) string {
	s := hex.EncodeToString(hash)
	if len(s) > 14 {
		return s[:14]
	}
	return s
}

// diffSection renders the versioned fields like the detail form: an unchanged
// field is just its label, struck through; a changed HTML field (or plain text)
// gets an inline diff; anything else shows Before/After values.
func (mh *ModelHistory) diffSection(recordKey string, aHash, bHash histmodels.Hash, msgr *Messages) (h.HTMLComponent, error) {
	am, bm, err := mh.twoFieldMaps(recordKey, aHash, bHash)
	if err != nil {
		return nil, err
	}

	var rows h.HTMLComponents
	for _, f := range mh.resolved {
		ov, nv := fieldValue(am, f), fieldValue(bm, f)
		if ov == nv {
			rows = append(rows, h.Div(
				h.Span(f).Style("text-decoration: line-through"),
				h.Span(" "+msgr.Unchanged).Class("text-medium-emphasis text-caption ms-1"),
			).Class("mb-1"))
			continue
		}

		var val h.HTMLComponent
		if mh.IsHTML(f) || (isStr(am, f) && isStr(bm, f)) {
			val = h.RawHTML(HTMLDiff(ov, nv))
		} else {
			val = h.Div(
				h.Div(h.Strong(msgr.Old+": "), h.Text(ov)),
				h.Div(h.Strong(msgr.New+": "), h.Text(nv)),
			)
		}
		rows = append(rows, h.Div(
			h.Div(h.Strong(f)).Class("mb-1"),
			h.Div(val).Class("ps-2"),
		).Class("mb-3"))
	}
	return h.Div(rows...).Class("mt-2"), nil
}

// isStr reports whether a field's stored JSON value is a JSON string (so a text
// or HTML diff makes sense, as opposed to a structured/foreign-key value).
func isStr(m map[string]json.RawMessage, f string) bool {
	raw, ok := m[f]
	return ok && len(raw) > 0 && raw[0] == '"'
}
