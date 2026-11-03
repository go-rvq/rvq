package utils

import (
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

// ColumnCapture holds the values of the columns Names of the first row a
// query read: columns the model has no field for, which gorm's scan would
// discard (a total of the whole filter by a window function, say).
type ColumnCapture struct {
	Names  []string
	Values map[string]any
	done   bool
}

const captureKey = "rvq:capture_columns"

// CaptureColumns is db, its query keeping the values of the columns names of
// the first row in the ColumnCapture returned (RegisterColumnCapture first).
func CaptureColumns(db *gorm.DB, names ...string) (*gorm.DB, *ColumnCapture) {
	c := &ColumnCapture{Names: names, Values: map[string]any{}}
	return db.InstanceSet(captureKey, c), c
}

// RegisterColumnCapture puts the query of the capture in place of gorm's
// (gorm:query): a query that asks for no capture runs gorm's own.
func RegisterColumnCapture(db *gorm.DB) error {
	return db.Callback().Query().Replace("gorm:query", queryCapturing)
}

func queryCapturing(db *gorm.DB) {
	v, ok := db.InstanceGet(captureKey)
	c, _ := v.(*ColumnCapture)
	if !ok || c == nil {
		callbacks.Query(db)
		return
	}
	if db.Error != nil {
		return
	}
	callbacks.BuildQuerySQL(db)
	if db.DryRun || db.Error != nil {
		return
	}
	rows, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
	if err != nil {
		db.AddError(err)
		return
	}
	defer func() {
		db.AddError(rows.Close())
	}()
	gorm.Scan(&capturingRows{Rows: rows, c: c}, db, 0)
	if db.Statement.Result != nil {
		db.Statement.Result.RowsAffected = db.RowsAffected
	}
}

// capturingRows reads, after gorm's scan of a row, the values of the
// columns of c from the destinations gorm gave them.
type capturingRows struct {
	gorm.Rows
	c *ColumnCapture
}

func (r *capturingRows) Scan(dest ...any) error {
	if err := r.Rows.Scan(dest...); err != nil || r.c.done {
		return err
	}
	r.c.done = true
	cols, err := r.Rows.Columns()
	if err != nil {
		return err
	}
	for i, col := range cols {
		for _, name := range r.c.Names {
			if col == name && i < len(dest) {
				r.c.Values[name] = deref(dest[i])
			}
		}
	}
	return nil
}

// deref is the value a destination of a scan points to; bytes (a JSON or a
// text of the driver, whose buffer it may reuse) as a string, copied.
func deref(d any) any {
	switch v := d.(type) {
	case *any:
		if b, ok := (*v).([]byte); ok {
			return string(b)
		}
		return *v
	case *[]byte:
		return string(*v)
	}
	return d
}
