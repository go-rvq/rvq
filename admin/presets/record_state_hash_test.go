package presets

import (
	"database/sql/driver"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// jsonColumn is a column the database stores whole, like seo.Setting: a struct
// that says for itself what it is worth.
type jsonColumn struct {
	Title string
}

func (c jsonColumn) Value() (driver.Value, error) { return c.Title, nil }

type stateCategory struct {
	ID    uint
	Label string
}

type stateItem struct {
	ID    uint
	Label string
	Pos   int
}

type stateModel struct {
	ID        uint
	Title     string
	Settings  jsonColumn
	At        time.Time
	Untouched string // a column, but not one this form edits

	CategoryID uint
	Category   *stateCategory // related: the form carries the choice
	Items      []*stateItem   // nested: edited through its own fields
}

// stateApp registers stateModel with an edit form over Title, Settings, At,
// Category and Items — Untouched is deliberately left out of the form.
func stateApp(t *testing.T) (*Builder, *ModelBuilder) {
	t.Helper()

	b := New(i18n.New())
	mb := b.Model(&stateModel{}).URIName("states")

	itemMB := NewModelBuilder(b, &stateItem{})
	itemEd := itemMB.Editing("Label", "Pos")

	ed := mb.Editing("Title", "Settings", "At", "Category", "Items")
	ed.Field("Items").Nested(NestedSlice(itemMB, &itemEd.FieldsBuilder))

	return b, mb
}

func baseStateModel() *stateModel {
	at := time.Date(2026, 7, 27, 10, 30, 0, 0, time.UTC)
	return &stateModel{
		ID:         1,
		Title:      "a",
		Settings:   jsonColumn{Title: "s"},
		At:         at,
		Untouched:  "before",
		CategoryID: 7,
		Category:   &stateCategory{ID: 7, Label: "c"},
		Items: []*stateItem{
			{ID: 1, Label: "A", Pos: 0},
			{ID: 2, Label: "B", Pos: 1},
		},
	}
}

// The same record hashes the same — otherwise the guard would refuse every save.
func TestRecordStateHashIsStable(t *testing.T) {
	_, mb := stateApp(t)
	ed := mb.Editing()

	if ed.RecordStateHash(baseStateModel()) != ed.RecordStateHash(baseStateModel()) {
		t.Error("dois registros iguais deram hashes diferentes")
	}
}

// What the form edits is what it locks.
func TestRecordStateHashNoticesWhatTheFormEdits(t *testing.T) {
	_, mb := stateApp(t)
	ed := mb.Editing()

	cases := map[string]func(m *stateModel){
		"a column":               func(m *stateModel) { m.Title = "b" },
		"a valuer column":        func(m *stateModel) { m.Settings = jsonColumn{Title: "other"} },
		"a time":                 func(m *stateModel) { m.At = m.At.Add(time.Second) },
		"the related choice":     func(m *stateModel) { m.Category = &stateCategory{ID: 9, Label: "c"} },
		"the related gone":       func(m *stateModel) { m.Category = nil },
		"a nested field":         func(m *stateModel) { m.Items[0].Label = "other" },
		"a nested row added":     func(m *stateModel) { m.Items = append(m.Items, &stateItem{ID: 3, Label: "C"}) },
		"a nested row removed":   func(m *stateModel) { m.Items = m.Items[:1] },
		"a nested row reordered": func(m *stateModel) { m.Items[0], m.Items[1] = m.Items[1], m.Items[0] },
	}

	before := ed.RecordStateHash(baseStateModel())
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			m := baseStateModel()
			change(m)
			if ed.RecordStateHash(m) == before {
				t.Error("a alteração não mudou o hash")
			}
		})
	}
}

// What the form does not carry is not part of what it locks.
func TestRecordStateHashIgnoresWhatTheFormDoesNotCarry(t *testing.T) {
	_, mb := stateApp(t)
	ed := mb.Editing()

	cases := map[string]func(m *stateModel){
		"a column outside the form": func(m *stateModel) { m.Untouched = "after" },
		"the related's own fields":  func(m *stateModel) { m.Category.Label = "renamed" },
	}

	before := ed.RecordStateHash(baseStateModel())
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			m := baseStateModel()
			change(m)
			if ed.RecordStateHash(m) != before {
				t.Error("algo que o form não carrega entrou no hash")
			}
		})
	}
}

// The guard reaches a model with no UpdatedAt: the form carries the hash of
// what it rendered, and an update over a record that moved since is refused.
func TestVerifyRecordStampByStateHash(t *testing.T) {
	b, mb := stateApp(t)

	stored := baseStateModel()
	mb.Editing().FetchFunc(func(obj any, id ID, ctx *web.EventContext) error {
		*obj.(*stateModel) = *stored
		return nil
	})

	verify := func(r *http.Request) error {
		return mb.Editing().VerifyRecordStamp(ID{}, &web.EventContext{R: r})
	}

	signed, err := b.SignForm(multipartRequest(t), stored)
	if err != nil {
		t.Fatal(err)
	}
	if err := verify(signed); err != nil {
		t.Errorf("o registro não mudou e o save foi recusado: %v", err)
	}

	// a column the form does not carry is not the form's business
	stored.Untouched = "after"
	if err := verify(signed); err != nil {
		t.Errorf("uma coluna fora do form derrubou o save: %v", err)
	}

	// somebody else saves over a field this form edits
	stored.Title = "after"
	err = verify(signed)
	if !errors.Is(err, ErrRecordChanged) {
		t.Fatalf("err = %v, want ErrRecordChanged", err)
	}
	// with no UpdatedAt there is no instant to name
	if got, want := err.Error(), string(Messages_en_US.ErrRecordChangedUnknownWhen); got != want {
		t.Errorf("mensagem = %q, want %q", got, want)
	}

	// and it is still required
	if err := verify(multipartRequest(t)); !errors.Is(err, ErrRecordStampMissing) {
		t.Errorf("sem o stamp: err = %v, want ErrRecordStampMissing", err)
	}
}

// A model that turned the fallback off is not guarded at all: no stamp is
// carried, and none is required.
func TestRecordStateStampCanBeTurnedOff(t *testing.T) {
	b, mb := stateApp(t)
	mb.SetRecordStateStamp(false)

	stored := baseStateModel()
	mb.Editing().FetchFunc(func(obj any, id ID, ctx *web.EventContext) error {
		*obj.(*stateModel) = *stored
		return nil
	})

	r := multipartRequest(t)
	signed, err := b.SignForm(r, stored)
	if err != nil {
		t.Fatal(err)
	}
	if signed != r {
		t.Error("um modelo sem trava não devia ganhar stamp")
	}

	stored.Title = "after"
	if err := mb.Editing().VerifyRecordStamp(ID{}, &web.EventContext{R: r}); err != nil {
		t.Errorf("o save foi recusado num modelo sem trava: %v", err)
	}
}
