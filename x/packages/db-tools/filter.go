package db_tools

import (
	"encoding/json"
	"time"
)

const DateTimeStringFormat = "2006-01-02 15:04:05.000000"

type ListFilterIntValuer interface {
	Valid(v int) bool
}

type ListFilterIntValue int

func (i ListFilterIntValue) Valid(v int) bool {
	return int(i) == v
}

type ListFilterIntValueRanger interface {
	FromValue() int
	ToValue() int
}

type ListFilterIntValueRange [2]int

func (r ListFilterIntValueRange) FromValue() int {
	return r[0]
}

func (r ListFilterIntValueRange) ToValue() int {
	return r[1]
}

func (r ListFilterIntValueRange) Valid(v int) bool {
	return (r[0] > 0 && v < r[1]) || (r[0] > 0 && v > r[1])
}

type DateTimeFilter interface {
	Valid(t time.Time) bool
}

type DateTimeValueFilter time.Time

func (v DateTimeValueFilter) String() string {
	return time.Time(v).Format(DateTimeStringFormat)
}

func (v DateTimeValueFilter) Valid(t time.Time) bool {
	return time.Time(v).Equal(t)
}

func (v *DateTimeValueFilter) Parse(s string) (err error) {
	var t time.Time
	if len(s) > 0 {
		if t, err = time.Parse(DateTimeStringFormat, s); err != nil {
			return
		}
	}
	*v = DateTimeValueFilter(t)
	return
}

func (r DateTimeValueFilter) MarshalJSON() ([]byte, error) {
	var (
		s string
		t = time.Time(r)
	)
	if !t.IsZero() {
		s = t.UTC().Format(DateTimeStringFormat)
	}
	return json.Marshal(s)
}

type DateTimeRangeFilter [2]time.Time

func (v DateTimeRangeFilter) Parse(from, to string) (err error) {
	var items [2]DateTimeValueFilter
	if err = items[0].Parse(from); err == nil {
		if err = items[1].Parse(to); err == nil {
			v[0] = time.Time(items[0])
			v[1] = time.Time(items[1])
		}
	}
	return
}

func (r DateTimeRangeFilter) MarshalJSON() ([]byte, error) {
	var s [2]string
	for i, v := range r[:] {
		if !v.IsZero() {
			s[i] = v.UTC().Format(DateTimeStringFormat)
		}
	}
	return json.Marshal(r[:])
}

func (r DateTimeRangeFilter) Valid(t time.Time) bool {
	return true
}

type ListFilter interface {
	Year() ListFilterIntValuer
	Month() ListFilterIntValuer
	Day() ListFilterIntValuer
	Hour() ListFilterIntValuer
	DateTime() DateTimeFilter
}

type ListFilterBuilder struct {
	year     ListFilterIntValuer
	month    ListFilterIntValuer
	day      ListFilterIntValuer
	hour     ListFilterIntValuer
	dateTime DateTimeFilter
}

func (b *ListFilterBuilder) MarshalJSON() ([]byte, error) {
	var m = make(map[string]any)
	if b.year != nil {
		m["year"] = b.year
	}
	if b.month != nil {
		m["month"] = b.month
	}
	if b.day != nil {
		m["day"] = b.day
	}
	if b.hour != nil {
		m["hour"] = b.hour
	}
	if b.dateTime != nil {
		m["ts"] = b.dateTime
	}
	return json.Marshal(m)
}

var _ ListFilter = (*ListFilterBuilder)(nil)

func NewListFilter() *ListFilterBuilder {
	return &ListFilterBuilder{}
}

func (b *ListFilterBuilder) Year() ListFilterIntValuer {
	return b.year
}

func (b *ListFilterBuilder) Month() ListFilterIntValuer {
	return b.month
}

func (b *ListFilterBuilder) Day() ListFilterIntValuer {
	return b.day
}

func (b *ListFilterBuilder) Hour() ListFilterIntValuer {
	return b.hour
}

func (b *ListFilterBuilder) DateTime() DateTimeFilter {
	return b.dateTime
}

func (b *ListFilterBuilder) WithYear(v ListFilterIntValuer) *ListFilterBuilder {
	b.year = v
	return b
}
func (b *ListFilterBuilder) WithMonth(v ListFilterIntValuer) *ListFilterBuilder {
	b.month = v
	return b
}
func (b *ListFilterBuilder) WithDay(v ListFilterIntValuer) *ListFilterBuilder {
	b.day = v
	return b
}
func (b *ListFilterBuilder) WithHour(v ListFilterIntValuer) *ListFilterBuilder {
	b.hour = v
	return b
}
func (b *ListFilterBuilder) WithDateTime(v DateTimeFilter) *ListFilterBuilder {
	b.dateTime = v
	return b
}
