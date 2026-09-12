package admin

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rvq/rvq/web"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/sunfmin/reflectutils"
)

// Revert applies a past revision to a loaded record and saves it, recording a
// NEW revision (git-revert style: the history chain is never rewritten). The
// obj is the record as loaded by the admin (Detailing), so no re-fetch is
// needed. Four granularities:
//
//   RevertRecord       — every versioned field, to `hash`.
//   RevertFields       — only the named fields, to `hash`.
//   RevertField        — one field's whole content, to `hash` (always allowed).
//   RevertFieldPartial — only selected hunks of one field (AcceptsPartial only).

// RevertRecord restores every versioned field of obj to revision `hash`.
func (h *ModelHistory) RevertRecord(obj any, hash []byte, ctx *web.EventContext) error {
	return h.RevertFields(obj, hash, h.resolved, ctx)
}

// RevertFields restores only the named fields of obj to revision `hash`.
func (h *ModelHistory) RevertFields(obj any, hash []byte, fields []string, ctx *web.EventContext) error {
	recordKey := h.mb.MustRecordID(obj).String()
	rev, err := h.Revision(recordKey, hash)
	if err != nil {
		return err
	}
	m, err := fieldMap(rev)
	if err != nil {
		return err
	}
	for _, f := range fields {
		raw, ok := m[f]
		if !ok {
			continue
		}
		if err := setFieldFromJSON(obj, f, raw); err != nil {
			return fmt.Errorf("history: revert field %q: %w", f, err)
		}
	}
	return h.saveAndCapture(obj, ctx)
}

// RevertField restores one field's whole content to revision `hash`. Always
// available, including for whole-only fields.
func (h *ModelHistory) RevertField(obj any, hash []byte, field string, ctx *web.EventContext) error {
	return h.RevertFields(obj, hash, []string{field}, ctx)
}

// RevertFieldPartial restores only part of a field's content: it applies the
// given diffmatchpatch patch text (the hunks the user selected) to the field's
// current value. Only for fields that accept partial revert.
func (h *ModelHistory) RevertFieldPartial(obj any, field, patchText string, ctx *web.EventContext) error {
	if !h.AcceptsPartial(field) {
		return fmt.Errorf("history: field %q does not accept partial revert", field)
	}
	cur, err := reflectutils.Get(obj, field)
	if err != nil {
		return err
	}
	curStr, _ := cur.(string)
	dmp := diffmatchpatch.New()
	patches, err := dmp.PatchFromText(patchText)
	if err != nil {
		return fmt.Errorf("history: bad patch: %w", err)
	}
	res, _ := dmp.PatchApply(patches, curStr)
	if err := reflectutils.Set(obj, field, res); err != nil {
		return err
	}
	return h.saveAndCapture(obj, ctx)
}

// setFieldFromJSON sets obj.field from the field's JSON snapshot value,
// unmarshaling into the field's own Go type (taken from the struct, so a struct
// or foreign-key field is rebuilt as itself — not left as the generic
// map[string]any that a plain unmarshal into `any` would give, which is not
// assignable to e.g. models.PageOptions).
func setFieldFromJSON(obj any, field string, raw json.RawMessage) error {
	fv, ok := settableFieldValue(obj, field)
	if !ok {
		// Not a resolvable struct field (virtual/dotted): best-effort via reflectutils.
		ft, typed := structFieldType(obj, field)
		if !typed {
			var v any
			if err := json.Unmarshal(raw, &v); err != nil {
				return err
			}
			return reflectutils.Set(obj, field, v)
		}
		ptr := reflect.New(ft)
		if err := json.Unmarshal(raw, ptr.Interface()); err != nil {
			return err
		}
		return reflectutils.Set(obj, field, ptr.Elem().Interface())
	}
	// Unmarshal into the field's own type and assign through the addressable,
	// settable reflect.Value — this handles a promoted pointer field (e.g.
	// Page.PageOptions *PageOptions via an embedded struct), which reflectutils
	// cannot set ("reflect.Value.Set on zero Value").
	ptr := reflect.New(fv.Type())
	if err := json.Unmarshal(raw, ptr.Interface()); err != nil {
		return err
	}
	fv.Set(ptr.Elem())
	return nil
}

// settableFieldValue resolves the addressable, settable reflect.Value of obj's
// field, walking a dotted path, following pointers (allocating a nil one so a
// nested field can be set) and promoted embedded fields. Returns ok=false when
// the path can't be resolved to a settable field.
func settableFieldValue(obj any, field string) (reflect.Value, bool) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return reflect.Value{}, false
	}
	v = v.Elem()
	parts := strings.Split(field, ".")
	for i, name := range parts {
		for v.Kind() == reflect.Ptr {
			if v.IsNil() {
				if !v.CanSet() {
					return reflect.Value{}, false
				}
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return reflect.Value{}, false
		}
		fv := v.FieldByName(name)
		if !fv.IsValid() || !fv.CanSet() {
			return reflect.Value{}, false
		}
		if i == len(parts)-1 {
			return fv, true
		}
		v = fv
	}
	return reflect.Value{}, false
}

// structFieldType returns the reflect.Type of obj's field (following pointers
// and promoted embedded fields), even when the field's current value is nil. It
// walks a dotted path so a nested struct field resolves too (e.g.
// "PageOptions.Layout").
func structFieldType(obj any, field string) (reflect.Type, bool) {
	t := reflect.TypeOf(obj)
	for _, name := range strings.Split(field, ".") {
		for t != nil && t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t == nil || t.Kind() != reflect.Struct {
			return nil, false
		}
		sf, ok := t.FieldByName(name)
		if !ok {
			return nil, false
		}
		t = sf.Type
	}
	return t, true
}

// saveAndCapture persists the reverted record through the model's editing save
// pipeline, so the change is recorded in BOTH the activity log and the history
// (activity's and history's save wraps fire, exactly as a normal edit) — the new
// revision is created by history's wrap, not here.
func (h *ModelHistory) saveAndCapture(obj any, ctx *web.EventContext) error {
	return h.mb.Editing().Saver(obj, h.mb.MustRecordID(obj), ctx)
}
