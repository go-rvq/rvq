package db_tools

import (
	"errors"
	"fmt"
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
		err = errors.New(fmt.Sprintf("Unknown persistence option: %s", v))
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

// Group groups sorted items
func (p *Persistence) Group(items []time.Time) (g *IntervalGrouper, other []time.Time) {
	g = new(IntervalGrouper)

	if len(items) == 0 {
		return
	}

	var start = ToDay(items[0])

	if p.Other == PersistenceOtherDays {
		g.Days, _, _ = GroupByWithFilter(items, time.Time{}, ToDay)
		return
	}

	if p.Days > 0 {
		g.Days, items, start = GroupByWithFilter(items, ToDay(start).Add(DAY_INTERVAL*time.Duration(p.Days)), ToDay)
	}

	if p.Other == PersistenceOtherWeeks {
		g.Weeks, _, _ = GroupByWithFilter(items, time.Time{}, ToWeek)
		return
	}

	if p.Weeks > 0 {
		g.Weeks, items, start = GroupByWithFilter(items, ToWeek(start).Add(WEEK_INTERVAL*time.Duration(p.Weeks)), ToWeek)
	}

	if p.Other == PersistenceOtherMonths {
		g.Months, _, _ = GroupByWithFilter(items, time.Time{}, ToMonth)
		return
	}

	if p.Months > 0 {
		g.Months, items, start = GroupByWithFilter(items, ToMonth(start).AddDate(0, p.Months, 0), ToMonth)
	}

	if p.Other == PersistenceOtherYears {
		g.Years, _, _ = GroupByWithFilter(items, time.Time{}, ToYear)
		return
	}

	if p.Years > 0 {
		g.Years, items, start = GroupByWithFilter(items, ToYear(start).AddDate(p.Years, 0, 0), ToYear)
	}

	other = items

	return
}
