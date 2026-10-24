package schemaform

import (
	"math"
	"slices"
	"strconv"
	"time"
	"unicode/utf8"
)

// The codes of what is wrong with a value: left empty, not of its type, off
// its limits (CheckLimits), a range whose from is after its to.
const (
	CodeRequired = "required"
	CodeInvalid  = "invalid"
	CodeMin      = "min"
	CodeMax      = "max"
	CodeStep     = "step"
	CodeTooShort = "too_short"
	CodeTooLong  = "too_long"
	CodePattern  = "pattern"
	CodeMask     = "mask"
	CodeRange    = "range"
)

// CheckLimits is what is wrong with s, the value of f — of its type
// already —, by its limits: [minlength, maxlength], [pattern], [mask], [min,
// max, step]; a code, "" when nothing. The same the browser's HTML5 checks.
func CheckLimits(f *Field, s string) string {
	if f.MinLength > 0 && utf8.RuneCountInString(s) < f.MinLength {
		return CodeTooShort
	}
	if f.MaxLength > 0 && utf8.RuneCountInString(s) > f.MaxLength {
		return CodeTooLong
	}
	if f.Pattern != "" {
		if re, err := PatternRegexp(f.Pattern); err == nil && !re.MatchString(s) {
			return CodePattern
		}
	}
	if f.Mask != "" && !MaskRegexp(f.Mask).MatchString(s) {
		return CodeMask
	}
	if f.Min == "" && f.Max == "" && f.Step == "" {
		return ""
	}
	v, ok := LimitValue(f.Type, s)
	if !ok {
		return ""
	}
	if min, ok := LimitValue(f.Type, f.Min); ok && v < min {
		return CodeMin
	}
	if max, ok := LimitValue(f.Type, f.Max); ok && v > max {
		return CodeMax
	}
	if step, err := strconv.ParseFloat(f.Step, 64); err == nil && step > 0 {
		// from the min, as HTML5 counts it (else from zero: a time from
		// midnight, a date from 1970-01-01)
		base, _ := LimitValue(f.Type, f.Min)
		n := (v - base) / step
		if math.Abs(n-math.Round(n)) > 1e-9 {
			return CodeStep
		}
	}
	return ""
}

// LimitValue is s, a value or a limit of the type t, as a number to compare:
// a number its value; a date its days, a time of day its seconds, a date and
// time its seconds (the units of HTML5's step); false when it is none.
func LimitValue(t, s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	switch t {
	case "int", "uint", "float", "decimal":
		n, err := strconv.ParseFloat(s, 64)
		return n, err == nil
	case "date", "cdate":
		d, err := time.Parse("2006-01-02", s)
		return float64(d.Unix() / 86400), err == nil
	case "time", "ctime":
		for _, layout := range []string{"15:04", "15:04:05"} {
			if d, err := time.Parse(layout, s); err == nil {
				return float64(d.Hour()*3600 + d.Minute()*60 + d.Second()), true
			}
		}
	case "cdtime":
		for _, layout := range []string{"2006-01-02T15:04", "2006-01-02T15:04:05"} {
			if d, err := time.Parse(layout, s); err == nil {
				return float64(d.Unix()), true
			}
		}
	}
	return 0, false
}

// After says a, the from of a range, is after b, its to: numbers whichever
// their kind (a Range[int|float]), texts — dates, times — as written,
// members of an enum by their order.
func After(enum *Enum, a, b any) bool {
	if enum != nil {
		return slices.Index(enum.Names, toText(a)) > slices.Index(enum.Names, toText(b))
	}
	if x, ok := asFloat(a); ok {
		if y, ok := asFloat(b); ok {
			return x > y
		}
	}
	if x, ok := a.(string); ok {
		if y, ok := b.(string); ok {
			return x > y
		}
	}
	return false
}

func asFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case int64:
		return float64(x), true
	case uint64:
		return float64(x), true
	case float64:
		return x, true
	case int:
		return float64(x), true
	}
	return 0, false
}

func toText(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
