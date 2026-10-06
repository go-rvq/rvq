package schemaform

import (
	"strings"
	"testing"
)

// A field of Range[T] is a range: its bounds of T — a union of numbers, a
// number —, of options; a bound left open ([open="to"|"from"]); anything
// else of open refused.
func TestRange(t *testing.T) {
	s, err := Parse(`interface Form {
		[open="to"] stay Range[time.CalendarDate]
		budget? Range[int|float]
		[options=[1, 2, 3, 4]] rooms Range[int]
		at? Range[time.CalendarTime]
	}`)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range s.Fields {
		enum := ""
		if f.Range != nil && f.Range.Enum != nil {
			enum = strings.Join(f.Range.Enum.Names, ",")
		}
		got = append(got, f.Name+":"+f.Type+":"+f.Range.Type+":"+f.RangeOpen+":"+enum)
	}
	want := "stay:range:calendarDate:to: budget:range:float:: rooms:range:int::1,2,3,4 at:range:calendarTime::"
	if strings.Join(got, " ") != want {
		t.Errorf("ranges:\n got %s\nwant %s", strings.Join(got, " "), want)
	}

	for src, want := range map[string]string{
		`interface Form { [open="end"] r Range[int] }`: `the bound left open is "from" or "to"`,
		`interface Form { [open="to"] r int }`:         "only a range",
	} {
		if _, err := Parse(src); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v", src, err)
		}
	}
}
