package db_tools

import (
	"sort"
	"time"
)

const (
	DAY_INTERVAL  = 24 * time.Hour
	WEEK_INTERVAL = DAY_INTERVAL * 7
)

type Group struct {
	Key   time.Time
	Items []time.Time
}

func (g *Group) MustLast(rm func(t []time.Time)) {
	if len(g.Items) > 1 {
		rm(g.Items[1:])
		g.Items = g.Items[:1]
	}
}

type Groups []*Group

func (s Groups) MustLast(rm func(t []time.Time)) {
	for _, g := range s {
		g.MustLast(rm)
	}
}

type BackupsMapper struct {
	backups map[time.Time]Backuper
	keys    []time.Time
}

func (b *BackupsMapper) Map() map[time.Time]Backuper {
	return b.backups
}

func (b *BackupsMapper) Keys() []time.Time {
	return b.keys
}

func Mapper(backups []Backuper) *BackupsMapper {
	var (
		times = make([]time.Time, len(backups))
		m     = make(map[time.Time]Backuper, len(backups))
	)

	for i, b := range backups {
		t := b.GetCreatedAt()
		times[i] = t
		m[t] = b
	}

	sort.Slice(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})

	return &BackupsMapper{
		backups: m,
		keys:    times,
	}
}

func (g *Group) Add(item ...time.Time) {
	g.Items = append(g.Items, item...)
}

func GroupBy(group func(t time.Time) time.Time, items []time.Time) (r []*Group) {
	if len(items) == 0 {
		return nil
	}

	m := make(map[time.Time]*Group)

	for _, item := range items {
		w := group(item)
		g := m[w]

		if g == nil {
			g = &Group{
				Key: w,
			}
			m[w] = g
			r = append(r, g)
		}

		g.Add(item)
	}

	sort.Slice(r, func(i, j int) bool {
		return r[i].Key.Before(r[j].Key)
	})

	return
}

func ToWeek(t time.Time) time.Time {
	day := ToDay(t)
	w := day.Weekday()
	if w != time.Sunday {
		day = day.AddDate(0, 0, -int(w-time.Sunday))
	}
	return day
}

func ToMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func ToYear(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
}

func ToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

type IntervalGrouper struct {
	Days   Groups
	Weeks  Groups
	Months Groups
	Years  Groups
}

func (i *IntervalGrouper) MustLast(rm func(t []time.Time)) {
	i.Days.MustLast(rm)
	i.Weeks.MustLast(rm)
	i.Months.MustLast(rm)
	i.Years.MustLast(rm)
}

func GroupByWithFilter(items []time.Time, end time.Time, key func(t time.Time) time.Time) (r []*Group, other []time.Time, end_ time.Time) {
	end_ = end
	filtered := items
	if !end.IsZero() {
		filtered, other = TimesAt(items, end)
	}
	if len(filtered) > 0 {
		r = GroupBy(key, filtered)
	}
	return
}

func TimesAt(items []time.Time, at time.Time) (filtered, other []time.Time) {
	for i, item := range items {
		if !item.Before(at) {
			return items[:i-1], items[i:]
		}
	}
	return items, nil
}

func FilterTimes(items []time.Time, filter func(t time.Time) bool) (filtered, other []time.Time) {
	for _, item := range items {
		if filter(item) {
			filtered = append(filtered, item)
		} else {
			other = append(other, item)
		}
	}
	return
}
