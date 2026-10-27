package web_test

import (
	"net/url"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/web"
	"github.com/google/uuid"
)

// A UUID of the form is read from its text: a key, a foreign key (a pointer),
// "" none — not taken for its 16 bytes, which failed the whole decoding.
func TestUnmarshalFormUUID(t *testing.T) {
	type child struct {
		ID     uuid.UUID
		TypeID *uuid.UUID
	}
	var v struct {
		Name  string
		Child *child
		None  uuid.UUID
	}
	id, typeID := uuid.New(), uuid.New()
	err := (&web.EventContext{}).UnmarshalFormValues(url.Values{
		"Name":         {"x"},
		"Child.ID":     {id.String()},
		"Child.TypeID": {typeID.String()},
		"None":         {""},
	}, &v)
	if err != nil {
		t.Fatal(err)
	}
	if v.Name != "x" || v.Child == nil || v.Child.ID != id || v.Child.TypeID == nil || *v.Child.TypeID != typeID || v.None != uuid.Nil {
		t.Errorf("decoded: %+v %+v", v, v.Child)
	}
	if err = (&web.EventContext{}).UnmarshalFormValues(url.Values{"None": {"not-a-uuid"}}, &v); err == nil {
		t.Error("a bad UUID: no error")
	}
}

// A model.UUID (a foreign key whose zero is none) is read from its text;
// "" is the zero — set before, it is none now.
func TestUnmarshalFormModelUUID(t *testing.T) {
	var v struct{ PageID, Other model.UUID }
	other := model.NewUUID()
	v.PageID = other
	if err := (&web.EventContext{}).UnmarshalFormValues(url.Values{"PageID": {""}, "Other": {other.String()}}, &v); err != nil {
		t.Fatal(err)
	}
	if !v.PageID.IsZero() || v.Other != other {
		t.Errorf("decoded: %v %v", v.PageID, v.Other)
	}
}
