package db_tools

import (
	"fmt"
	"testing"
	"time"
)

func GenTimes(start, end time.Time, interval time.Duration) (r []time.Time) {
	for ; start.Before(end); start = start.Add(interval) {
		r = append(r, start)
	}
	return
}

func init() {
	time.Local = time.UTC
}

func TestInterval_Group(t *testing.T) {
	var (
		start       = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		end         = start.AddDate(6, 0, 0)
		times       = GenTimes(start, end, time.Hour*12)
		printGroups = func(name, layout string, g []*Group) {
			if len(g) == 0 {
				return
			}
			fmt.Println("= " + name)
			for i, v := range g {
				fmt.Println(fmt.Sprintf("  => %03d. %v", i, v.Key.Format(layout)))
				for i, item := range v.Items {
					fmt.Println(fmt.Sprintf("     %03d. %v", i, item))
				}
			}
		}

		do = func(i *Persistence, times []time.Time) []time.Time {
			g, other := i.Group(times)
			g.MustLast()

			printGroups("days", "2006-01-02", g.Days)
			printGroups("wheeks", "2006-01-02 (Mon)", g.Weeks)
			printGroups("months", "2006-01", g.Months)
			printGroups("years", "2006", g.Years)
			return other
		}
	)

	times = do(&Persistence{
		Days:   3,
		Weeks:  3,
		Months: 2,
		Years:  4,
	}, times)

	if len(times) > 0 {
		fmt.Println("### other")
		do(&Persistence{
			Other: PersistenceOtherYears,
		}, times)
	}
	return
}
