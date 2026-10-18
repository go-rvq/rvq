package dashboard

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// The periods offered: the last hour to the last 90 days; Custom is a range
// of days (From, To).
const (
	Period1h     = "1h"
	Period12h    = "12h"
	Period24h    = "24h"
	Period7d     = "7d"
	Period30d    = "30d"
	Period90d    = "90d"
	PeriodCustom = "custom"
	// PeriodDefault is the period with none chosen.
	PeriodDefault = Period30d
)

// Periods are the periods offered, in order.
var Periods = []string{Period1h, Period12h, Period24h, Period7d, Period30d, Period90d, PeriodCustom}

var periodSpans = map[string]time.Duration{
	Period1h:  time.Hour,
	Period12h: 12 * time.Hour,
	Period24h: 24 * time.Hour,
	Period7d:  7 * 24 * time.Hour,
	Period30d: 30 * 24 * time.Hour,
	Period90d: 90 * 24 * time.Hour,
}

// Period is a span of time [From, To) and the step its series are counted
// in (Step: 5 minutes for the last hour, an hour up to two days, a day
// beyond).
type Period struct {
	Key      string
	From, To time.Time
	Step     time.Duration
}

// stepOf is the step of a span: 5 minutes for an hour, an hour up to two
// days, a day beyond.
func stepOf(span time.Duration) time.Duration {
	switch {
	case span <= time.Hour:
		return 5 * time.Minute
	case span <= 48*time.Hour:
		return time.Hour
	default:
		return 24 * time.Hour
	}
}

// NewPeriod is the period of key ending at now; Custom is from..to (days,
// "2006-01-02", to included). An unknown key is the default.
func NewPeriod(key string, now time.Time, from, to string) Period {
	if key == PeriodCustom {
		f, err1 := time.ParseInLocation("2006-01-02", from, now.Location())
		t, err2 := time.ParseInLocation("2006-01-02", to, now.Location())
		if err1 == nil && err2 == nil && !t.Before(f) {
			t = t.AddDate(0, 0, 1)
			return Period{Key: key, From: f, To: t, Step: stepOf(t.Sub(f))}
		}
		key = PeriodDefault
	}
	span, ok := periodSpans[key]
	if !ok {
		key, span = PeriodDefault, periodSpans[PeriodDefault]
	}
	step := stepOf(span)
	// the end at the step's next boundary: the last bucket is now's
	end := now.Truncate(step).Add(step)
	if step == 24*time.Hour {
		y, m, d := now.Date()
		end = time.Date(y, m, d+1, 0, 0, 0, 0, now.Location())
	}
	return Period{Key: key, From: end.Add(-span), To: end, Step: step}
}

// PeriodOf is the period of the request (?period=…&from=…&to=…), else def.
func PeriodOf(r *http.Request, now time.Time, def string) Period {
	key := r.FormValue("period")
	if key == "" {
		key = def
	}
	return NewPeriod(key, now, r.FormValue("from"), r.FormValue("to"))
}

// Previous is the period of the same length just before (for the deltas).
func (p Period) Previous() Period {
	span := p.To.Sub(p.From)
	return Period{Key: p.Key, From: p.From.Add(-span), To: p.From, Step: p.Step}
}

// Buckets are the starts of the steps of the period.
func (p Period) Buckets() []time.Time {
	var out []time.Time
	for t := p.From; t.Before(p.To); t = t.Add(p.Step) {
		out = append(out, t)
	}
	return out
}

// Labels are the buckets' labels: their time ("15:04") within a day, their
// day ("02/01") beyond; a few, the rest empty (an axis).
func (p Period) Labels(n int) []string {
	bs := p.Buckets()
	out := make([]string, len(bs))
	if len(bs) == 0 {
		return out
	}
	every := max(1, len(bs)/max(1, n))
	layout := "02/01"
	if p.Step < 24*time.Hour {
		layout = "15:04"
	}
	for i, t := range bs {
		if i%every == 0 || i == len(bs)-1 {
			out[i] = t.In(p.From.Location()).Format(layout)
		}
	}
	return out
}

// interval is the step as a Postgres interval.
func (p Period) interval() string {
	return fmt.Sprintf("%d seconds", int64(p.Step/time.Second))
}

// Series counts, by bucket of the period, the rows of table whose column col
// (a time) is in the period, where (with args) holds; sum, when given, is
// the expression summed (else the rows are counted). A bucket with none is 0.
func (p Period) Series(db *gorm.DB, table, col, sum, where string, args ...any) ([]float64, error) {
	agg := "count(*)"
	if sum != "" {
		agg = "coalesce(sum(" + sum + "), 0)"
	}
	cond := ""
	if strings.TrimSpace(where) != "" {
		cond = " AND (" + where + ")"
	}
	q := fmt.Sprintf(`SELECT date_bin(?::interval, %[1]s, ?::timestamptz) AS b, %[2]s AS n FROM %[3]s WHERE %[1]s >= ? AND %[1]s < ?%[4]s GROUP BY 1`,
		col, agg, table, cond)
	all := append([]any{p.interval(), p.From, p.From, p.To}, args...)
	var rows []struct {
		B time.Time
		N float64
	}
	if err := db.Session(&gorm.Session{}).Raw(q, all...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	bs := p.Buckets()
	out := make([]float64, len(bs))
	for _, r := range rows {
		i := int(r.B.Sub(p.From) / p.Step)
		if i >= 0 && i < len(out) {
			out[i] += r.N
		}
	}
	return out, nil
}

// Count is how many rows of table have col in the period where holds (sum:
// the expression summed instead).
func (p Period) Count(db *gorm.DB, table, col, sum, where string, args ...any) (float64, error) {
	agg := "count(*)"
	if sum != "" {
		agg = "coalesce(sum(" + sum + "), 0)"
	}
	cond := ""
	if strings.TrimSpace(where) != "" {
		cond = " AND (" + where + ")"
	}
	var n float64
	err := db.Session(&gorm.Session{}).Raw(fmt.Sprintf(`SELECT %s FROM %s WHERE %s >= ? AND %s < ?%s`, agg, table, col, col, cond),
		append([]any{p.From, p.To}, args...)...).Scan(&n).Error
	return n, err
}

// Total sums values.
func Total(values []float64) (t float64) {
	for _, v := range values {
		t += v
	}
	return
}

// AddSeries sums series of the same length, bucket by bucket.
func AddSeries(series ...[]float64) []float64 {
	if len(series) == 0 {
		return nil
	}
	out := make([]float64, len(series[0]))
	for _, s := range series {
		for i := range out {
			if i < len(s) {
				out[i] += s[i]
			}
		}
	}
	return out
}
