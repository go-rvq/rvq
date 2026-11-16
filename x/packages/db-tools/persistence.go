package db_tools

import (
	"fmt"
	"sort"
	"time"
)

type PersistenceOther uint8

const (
	PersistenceOtherDays PersistenceOther = iota + 1
	PersistenceOtherWeeks
	PersistenceOtherMonths
	PersistenceOtherYears
)

func (o PersistenceOther) String() string {
	switch o {
	case 0:
		return ""
	case PersistenceOtherDays:
		return "days"
	case PersistenceOtherWeeks:
		return "weeks"
	case PersistenceOtherMonths:
		return "months"
	case PersistenceOtherYears:
		return "years"
	default:
		return "unknown"
	}
}

func (o *PersistenceOther) Parse(v string) (err error) {
	switch v {
	case "":
		*o = 0
	case "days":
		*o = PersistenceOtherDays
	case "weeks":
		*o = PersistenceOtherWeeks
	case "months":
		*o = PersistenceOtherMonths
	case "years":
		*o = PersistenceOtherYears
	default:
		err = fmt.Errorf("unknown persistence option: %s", v)
	}
	return
}

type Persistence struct {
	Enabled bool             `json:"Enabled,omitempty"`
	Years   int              `json:"Years,omitempty"`
	Months  int              `json:"Months,omitempty"`
	Weeks   int              `json:"Weeks,omitempty"`
	Days    int              `json:"Days,omitempty"`
	Other   PersistenceOther `json:"Other,omitempty"`
}

func (p *Persistence) IsZero() bool {
	return p.Years == 0 && p.Months == 0 && p.Weeks == 0 && p.Days == 0 && p.Other == 0
}

// Group groups items by what p keeps of them, newest first: the Days most
// recent days that have backups (one each: the newest of the day), then the
// Weeks weeks before them, the Months months, the Years years — each a count
// of periods with backups, not of the calendar —; Other, the unit the older
// ones are kept by (one each, forever); other, those of no group (removed).
func (p *Persistence) Group(items []time.Time) (g *IntervalGrouper, other []time.Time) {
	g = new(IntervalGrouper)
	if len(items) == 0 {
		return
	}
	items = newestFirst(items)
	// tier takes the items of the n most recent periods of key (all of them,
	// forever) — the items of a period are together, the newest first
	tier := func(key func(time.Time) time.Time, n int, forever bool) Groups {
		var (
			gs   Groups
			last *Group
			i    int
		)
		for ; i < len(items); i++ {
			k := key(items[i])
			if last == nil || !last.Key.Equal(k) {
				if !forever && len(gs) == n {
					break
				}
				last = &Group{Key: k}
				gs = append(gs, last)
			}
			last.Add(items[i])
		}
		items = items[i:]
		return gs
	}
	if p.Other == PersistenceOtherDays {
		g.Days = tier(ToDay, 0, true)
		return
	}
	if p.Days > 0 {
		g.Days = tier(ToDay, p.Days, false)
	}
	if p.Other == PersistenceOtherWeeks {
		g.Weeks = tier(ToWeek, 0, true)
		return
	}
	if p.Weeks > 0 {
		g.Weeks = tier(ToWeek, p.Weeks, false)
	}
	if p.Other == PersistenceOtherMonths {
		g.Months = tier(ToMonth, 0, true)
		return
	}
	if p.Months > 0 {
		g.Months = tier(ToMonth, p.Months, false)
	}
	if p.Other == PersistenceOtherYears {
		g.Years = tier(ToYear, 0, true)
		return
	}
	if p.Years > 0 {
		g.Years = tier(ToYear, p.Years, false)
	}
	other = items
	return
}

// Removable are the items p does not keep — none when it is off (not
// Enabled) or empty: a retention not set removes nothing —; the newest is
// always kept.
func (p *Persistence) Removable(items []time.Time) (rm []time.Time) {
	if p == nil || !p.Enabled || p.IsZero() || len(items) == 0 {
		return nil
	}
	g, rm := p.Group(items)
	g.MustLast(func(t []time.Time) { rm = append(rm, t...) })
	newest := items[0]
	for _, t := range items {
		if t.After(newest) {
			newest = t
		}
	}
	out := rm[:0]
	for _, t := range rm {
		if !t.Equal(newest) {
			out = append(out, t)
		}
	}
	return out
}

// newestFirst is a copy of items, the newest first.
func newestFirst(items []time.Time) []time.Time {
	out := append([]time.Time(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].After(out[j]) })
	return out
}
