package activity

import (
	"sort"
	"strings"
	"testing"
)

type (
	relCustomer struct {
		ID   uint
		Name string
	}
	relLine struct {
		ID      uint
		OrderID uint
		Product string
		Qty     int
	}
	relTag struct {
		ID   uint
		Name string
	}
	relOrder struct {
		ID         uint
		Note       string
		CustomerID uint
		Customer   *relCustomer `gorm:"foreignKey:CustomerID"`
		Lines      []*relLine   `gorm:"foreignKey:OrderID"`
		Tags       []*relTag    `gorm:"many2many:rel_order_tags"`
	}
)

func (l *relLine) String() string { return l.Product }
func (t *relTag) String() string  { return t.Name }

// The relations of a record in its diff: the belongs-to not diffed (its key
// is); the lists, loaded as they were before the save, compared by key —
// what came, what left, what changed in one kept —; none, nil or empty, the
// same; a list whose elements are audited on their own, not in the record's.
func TestDiffRelations(t *testing.T) {
	freshPresets()
	if err := db.AutoMigrate(&relCustomer{}, &relTag{}, &relOrder{}, &relLine{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Migrator().DropTable("rel_order_tags", &relOrder{}, &relTag{}, &relLine{}, &relCustomer{}) })
	builder := New(db, &TestActivityLog{})
	mb := builder.RegisterModel(pb.Model(&relOrder{}))

	cus := &relCustomer{Name: "Acme"}
	a, b := &relTag{Name: "a"}, &relTag{Name: "b"}
	db.Create(cus)
	db.Create(a)
	db.Create(b)
	o := &relOrder{Note: "x", CustomerID: cus.ID, Lines: []*relLine{{Product: "brick", Qty: 1}, {Product: "sand", Qty: 2}}, Tags: []*relTag{a}}
	if err := db.Create(o).Error; err != nil {
		t.Fatal(err)
	}
	diff := func(now *relOrder) string {
		t.Helper()
		old, ok := findOld(now, db)
		if !ok {
			t.Fatal("no old")
		}
		ds, err := NewDiffBuilder(mb).diffRecord(db, old, now)
		if err != nil {
			t.Fatal(err)
		}
		var s []string
		for _, d := range ds {
			s = append(s, d.Field+":"+d.Old+">"+d.Now)
		}
		sort.Strings(s)
		return strings.Join(s, " ")
	}

	// the form's record, as it was: nothing
	same := &relOrder{ID: o.ID, Note: "x", CustomerID: cus.ID, Customer: cus,
		Lines: []*relLine{{ID: o.Lines[0].ID, OrderID: o.ID, Product: "brick", Qty: 1}, {ID: o.Lines[1].ID, OrderID: o.ID, Product: "sand", Qty: 2}},
		Tags:  []*relTag{{ID: a.ID, Name: "a"}}}
	if got := diff(same); got != "" {
		t.Errorf("nothing changed: %q", got)
	}

	// a line changed, one removed, one new; the tag a for b
	ch := &relOrder{ID: o.ID, Note: "x", CustomerID: cus.ID, Customer: cus,
		Lines: []*relLine{{ID: o.Lines[0].ID, OrderID: o.ID, Product: "brick", Qty: 5}, {Product: "lime", Qty: 1}},
		Tags:  []*relTag{{ID: b.ID, Name: "b"}}}
	if got, want := diff(ch), "Lines.brick.Qty:1>5 Lines:>lime Lines:sand> Tags:>b Tags:a>"; got != want {
		t.Errorf("the changes:\n got %q\nwant %q", got, want)
	}

	// no tags: an empty list, as none
	db.Model(o).Association("Tags").Clear()
	if got := diff(&relOrder{ID: o.ID, Note: "x", CustomerID: cus.ID, Tags: []*relTag{}}); got != "" {
		t.Errorf("empty, none: %q", got)
	}

	// the lines audited on their own: not in the order's diff
	builder.RegisterModel(pb.Model(&relLine{}))
	if got := diff(ch); got != "Tags:>b" {
		t.Errorf("lines on their own: %q", got)
	}
}
