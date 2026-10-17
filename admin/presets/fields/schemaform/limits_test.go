package schemaform

import (
	"strconv"
	"strings"
	"testing"
)

// The limits of a value — [min, max, step] of a number, a date, a time, a
// range's bounds; [minlength, maxlength, pattern] of a text; [placeholder]
// —: read, and refused where they mean nothing.
func TestLimits(t *testing.T) {
	s, err := Parse(`export class default {
	[min=1, max=10] rooms int
	[min=0.5, max=99.5, step=0.5] ratio float
	[min="2026-01-01", max="2026-12-31"] day date
	[min="08:00", max="18:00", step=900] at time
	[min=0, max=1000000, step=100] budget Range[int]
	[minlength=3, maxlength=80, placeholder="Your name"] name str
	[pattern="[0-9]{5}-?[0-9]{3}", placeholder="00000-000"] zip str
	[maxlength=500] about text
}`)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range s.Fields {
		l := f
		if f.Range != nil {
			l = f.Range
		}
		got = append(got, strings.Join([]string{f.Name, l.Min, l.Max, l.Step, itoa(l.MinLength), itoa(l.MaxLength), l.Pattern}, "|"))
	}
	want := []string{
		"rooms|1|10||0|0|", "ratio|0.5|99.5|0.5|0|0|", "day|2026-01-01|2026-12-31||0|0|",
		"at|08:00|18:00|900|0|0|", "budget|0|1000000|100|0|0|", "name||||3|80|",
		"zip||||0|0|[0-9]{5}-?[0-9]{3}", "about||||0|500|",
	}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("limits:\n got %v\nwant %v", got, want)
	}

	for src, want := range map[string]string{
		`export class default { [min=10, max=1] n int }`:             "is after",
		`export class default { [min=1.5] n int }`:                   "whole number",
		`export class default { [min="a"] n int }`:                   "want a number",
		`export class default { [min=1] s str }`:                     "only a number, a date or a time",
		`export class default { [min=2026] d date }`:                 "as text",
		`export class default { [step=0] n float }`:                  "more than 0",
		`export class default { [maxlength=5] n int }`:               "only a text",
		`export class default { [minlength=9, maxlength=5] s str }`:  "is more than",
		`export class default { [pattern="(?i)abc"] s str }`:         "(?…)",
		`export class default { [pattern="[0-9"] s str }`:            "no regular expression",
		`export class default { [pattern="[0-9]+"] s text }`:         "only a line of text",
		`export class default { [placeholder="x"] d date }`:          "only a text, a number or an address",
		`export class default { [options=["a", "b"], min=1] s str }`: "only a number, a date, a time or a text",
	} {
		if _, err := Parse(src); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v (want %q)", src, err, want)
		}
	}
	if _, err := Parse(`export class default { [pattern="(?:ab)+"] s str }`); err != nil {
		t.Errorf("(?:…) refused: %v", err)
	}
}

func itoa(n int) string { return strconv.Itoa(n) }
