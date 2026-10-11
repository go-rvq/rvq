package schemaform

import "testing"

// A mask: the whole value written by it; its placeholder, a blank where a
// character goes; read from [mask=…], refused where it means nothing.
func TestMask(t *testing.T) {
	re := MaskRegexp("(###) ###-####")
	for v, want := range map[string]bool{"(781) 555-0100": true, "(781) 555-010": false, "781 555 0100": false, "(78a) 555-0100": false} {
		if re.MatchString(v) != want {
			t.Errorf("%q: %v", v, !want)
		}
	}
	if !MaskRegexp("A#A #A#").MatchString("K1A 0B1") || MaskRegexp("A#A").MatchString("11A") {
		t.Error("letters")
	}
	if p := MaskPlaceholder("(###) ###-####"); p != "(___) ___-____" {
		t.Errorf("placeholder: %q", p)
	}
	s, err := Parse(`interface Form { [mask="#####-###"] zip str }`)
	if err != nil || s.Fields[0].Mask != "#####-###" {
		t.Fatalf("read: %v %+v", err, s)
	}
	for _, src := range []string{`interface Form { [mask="#"] n int }`, `interface Form { [mask="xyz"] s str }`} {
		if _, err := Parse(src); err == nil {
			t.Errorf("%s: no error", src)
		}
	}
}
