package db_tools

import (
	"testing"
	"time"
)

// daily are n backups, one a day at 03:00, the last today; ascending, as
// the backups are listed.
func daily(n int) []time.Time {
	today := time.Date(2026, 11, 16, 3, 0, 0, 0, time.UTC)
	out := make([]time.Time, n)
	for i := range out {
		out[i] = today.AddDate(0, 0, i-n+1)
	}
	return out
}

func kept(items, rm []time.Time) (out []time.Time) {
	gone := map[time.Time]bool{}
	for _, t := range rm {
		gone[t] = true
	}
	for _, t := range items {
		if !gone[t] {
			out = append(out, t)
		}
	}
	return
}

// 30 days of 60 daily backups: the newest 30 kept — not the oldest.
func TestRemovableKeepsTheNewest(t *testing.T) {
	items := daily(60)
	rm := (&Persistence{Enabled: true, Days: 30}).Removable(items)
	k := kept(items, rm)
	if len(k) != 30 || !k[0].Equal(items[30]) || !k[29].Equal(items[59]) {
		t.Fatalf("kept %d: %v … %v", len(k), k[0], k[len(k)-1])
	}
}

// A retention not set, or off, removes nothing.
func TestRemovableNotSet(t *testing.T) {
	items := daily(10)
	for _, p := range []*Persistence{nil, {}, {Enabled: true}, {Days: 3}} {
		if rm := p.Removable(items); len(rm) != 0 {
			t.Errorf("%+v: removed %d", p, len(rm))
		}
	}
}

// Two backups a day: one a day kept, the newest of it; the newest of all,
// always.
func TestRemovableOneADay(t *testing.T) {
	var items []time.Time
	for _, d := range daily(5) {
		items = append(items, d, d.Add(6*time.Hour))
	}
	k := kept(items, (&Persistence{Enabled: true, Days: 5}).Removable(items))
	if len(k) != 5 {
		t.Fatalf("kept %d", len(k))
	}
	for _, x := range k {
		if x.Hour() != 9 {
			t.Errorf("kept the older of a day: %v", x)
		}
	}
	if k := kept(items, (&Persistence{Enabled: true, Days: 1}).Removable(items)); len(k) != 1 || !k[0].Equal(items[len(items)-1]) {
		t.Errorf("one day: %v", k)
	}
}

// Days with gaps: the days are of backups, not of the calendar — 10 days
// of backups made every other day kept, 10 backups (20 days of calendar).
func TestRemovableDaysOfBackups(t *testing.T) {
	var items []time.Time
	for i, d := range daily(60) {
		if i%2 == 1 {
			items = append(items, d)
		}
	}
	k := kept(items, (&Persistence{Enabled: true, Days: 10}).Removable(items))
	if len(k) != 10 || !k[9].Equal(items[len(items)-1]) || !k[0].Equal(items[len(items)-10]) {
		t.Fatalf("kept %d: %v", len(k), k)
	}
}

// Days then weeks then the others by months: the 7 newest days, then one a
// week for the 4 weeks before them, then one a month.
func TestRemovableTiers(t *testing.T) {
	items := daily(120)
	k := kept(items, (&Persistence{Enabled: true, Days: 7, Weeks: 4, Other: PersistenceOtherMonths}).Removable(items))
	// 7 days + 4 weeks + the months before (4 or 5)
	if len(k) < 7+4+3 || len(k) > 7+4+5 {
		t.Errorf("kept %d", len(k))
	}
	if !k[len(k)-1].Equal(items[len(items)-1]) {
		t.Errorf("the newest not kept")
	}
}

// TimesAt: those before at, and the others; the first already at it, none
// before (no panic).
func TestTimesAt(t *testing.T) {
	items := daily(5)
	if before, after := TimesAt(items, items[0]); len(before) != 0 || len(after) != 5 {
		t.Errorf("at the first: %d %d", len(before), len(after))
	}
	if before, after := TimesAt(items, items[3]); len(before) != 3 || len(after) != 2 {
		t.Errorf("at the fourth: %d %d", len(before), len(after))
	}
}
