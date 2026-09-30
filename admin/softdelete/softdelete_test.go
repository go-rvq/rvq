package softdelete

import (
	"testing"

	"github.com/go-rvq/rvq/admin/origin"
	"github.com/google/uuid"
)

func TestColumns(t *testing.T) {
	by := uuid.New()
	lat, lng := 51.5, -0.1
	cols := Columns(&by, origin.Origin{IP: "81.2.69.160", City: "London", Latitude: &lat, Longitude: &lng})
	if cols["deleted_by_id"] != &by || cols["deleted_origin_ip"] != "81.2.69.160" ||
		cols["deleted_origin_city"] != "London" || cols["deleted_origin_latitude"] != &lat {
		t.Errorf("columns: %v", cols)
	}
	r := Restored()
	if v, ok := r["deleted_at"]; !ok || v != nil {
		t.Errorf("restored: deleted_at %v", v)
	}
	if r["deleted_by_id"] != (*uuid.UUID)(nil) || r["deleted_origin_ip"] != "" {
		t.Errorf("restored: %v", r)
	}
}
