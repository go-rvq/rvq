package admin

import (
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Recorder accumulates per-revision access hits in memory and flushes them to
// the <table>_revisions rows in batches, so a public request never waits on a
// write. A hit targets the record's current published revision.
type Recorder struct {
	db       *gorm.DB
	interval time.Duration

	mu      sync.Mutex
	pending map[hitKey]int64 // (table, recordKey) -> count since last flush

	stop chan struct{}
}

type hitKey struct {
	table     string
	recordKey string
}

// defaultRecorder is the process-wide recorder set by Configure; the per-model
// Build wires public hits to it.
var defaultRecorder *Recorder

func NewRecorder(db *gorm.DB) *Recorder {
	return &Recorder{
		db:       db,
		interval: 10 * time.Second,
		pending:  map[hitKey]int64{},
		stop:     make(chan struct{}),
	}
}

// Hit records one access of the current published revision of (table, recordKey).
// It is cheap: an in-memory increment, flushed later.
func (r *Recorder) Hit(table, recordKey string) {
	if r == nil {
		return
	}
	r.mu.Lock()
	r.pending[hitKey{table, recordKey}]++
	r.mu.Unlock()
}

// Start launches the background flusher. Safe to call once.
func (r *Recorder) Start() {
	go func() {
		t := time.NewTicker(r.interval)
		defer t.Stop()
		for {
			select {
			case <-r.stop:
				r.Flush()
				return
			case <-t.C:
				r.Flush()
			}
		}
	}()
}

// Stop ends the flusher after a final flush.
func (r *Recorder) Stop() { close(r.stop) }

// Flush writes the pending counts onto the current published revision of each
// touched record (the latest row with published = true for that record_key).
func (r *Recorder) Flush() {
	r.mu.Lock()
	batch := r.pending
	r.pending = map[hitKey]int64{}
	r.mu.Unlock()

	now := time.Now()
	for k, n := range batch {
		if n == 0 {
			continue
		}
		// Target the record's current published revision; if none is published
		// yet, target its latest revision so the count is not lost.
		sub := r.db.Table(k.table).
			Select("hash").
			Where("record_key = ?", k.recordKey).
			Order(clause.OrderBy{Columns: []clause.OrderByColumn{
				{Column: clause.Column{Name: "published"}, Desc: true},
				{Column: clause.Column{Name: "created_at"}, Desc: true},
			}}).
			Limit(1)
		r.db.Table(k.table).
			Where("record_key = ? AND hash = (?)", k.recordKey, sub).
			UpdateColumns(map[string]any{
				"access_count": gorm.Expr("access_count + ?", n),
				"last_access":  now,
			})
	}
}
