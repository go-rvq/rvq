package presets

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

type columnsRecord struct {
	ID     uint
	Title  string
	Status string
	Secret string
}

// Columns are what the listing's table shows: the fields of its layout, in
// that order, labelled, each drawing the cell the table draws — so a view that
// is not the table (a tree) shows the same thing.
func TestListingColumns(t *testing.T) {
	b := New(i18n.New()).URIPrefix("/admin")
	mb := b.Model(&columnsRecord{})
	l := mb.Listing().Only("Title", "Status", "ID")
	l.Field("Status").Label("Estado").ComponentFunc(func(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
		return h.Td(h.Text("status:" + field.Obj.(*columnsRecord).Status))
	})

	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/admin/columns-records", nil)}
	cols := l.Columns(ctx)

	var names []string
	for _, c := range cols {
		names = append(names, c.Name)
	}
	if got := strings.Join(names, ","); got != "Title,Status,ID" {
		t.Fatalf("columns = %s, want the layout's order Title,Status,ID", got)
	}
	if cols[1].Label != "Estado" {
		t.Errorf("label = %q, want the field's own label", cols[1].Label)
	}

	rec := &columnsRecord{ID: 7, Title: "Hello", Status: "online"}
	cell, err := h.Marshal(cols[1].Cell(rec), context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := string(cell); !strings.Contains(got, "<td>") || !strings.Contains(got, "status:online") {
		t.Errorf("cell = %s, want the field's LIST component (a <td>)", got)
	}
	title, _ := h.Marshal(cols[0].Cell(rec), context.Background())
	if !strings.Contains(string(title), "Hello") {
		t.Errorf("default cell = %s, want the value", title)
	}
}

// A field the listing leaves out, or turns off for the request (SetEnabled), is
// not a column — the table would not show it either.
func TestListingColumnsSkipsDisabledFields(t *testing.T) {
	b := New(i18n.New()).URIPrefix("/admin")
	mb := b.Model(&columnsRecord{})
	l := mb.Listing().Only("Title", "Secret")
	l.Field("Secret").SetEnabled(func(*FieldContext) bool { return false })

	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/admin/columns-records", nil)}
	for _, c := range l.Columns(ctx) {
		if c.Name == "Secret" || c.Name == "Status" {
			t.Errorf("column %s should not be there", c.Name)
		}
	}
}
