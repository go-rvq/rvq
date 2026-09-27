package media_library

import (
	"encoding/json"
	"testing"
)

func TestMediaIDJSON(t *testing.T) {
	for in, want := range map[string]MediaID{
		`{"ID":0}`:    "",
		`{"ID":12}`:   "12",
		`{"ID":""}`:   "",
		`{"ID":null}`: "",
		`{"ID":"019a8f3c-1111-7222-8333-444455556666"}`: "019a8f3c-1111-7222-8333-444455556666",
		`{}`: "",
	} {
		var b MediaBox
		if err := json.Unmarshal([]byte(in), &b); err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		if b.ID != want {
			t.Errorf("%s: ID = %q, want %q", in, b.ID, want)
		}
	}

	out, err := json.Marshal(MediaBox{ID: "019a8f3c-1111-7222-8333-444455556666"})
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	_ = json.Unmarshal(out, &m)
	if m["ID"] != "019a8f3c-1111-7222-8333-444455556666" {
		t.Errorf("marshal: %s", out)
	}
}
