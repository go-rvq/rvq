package helper

import (
	"net/http"
	"net/url"

	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	gormutils "github.com/go-rvq/rvq/thirdpart/gorm/utils"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func selEventContext() *web.EventContext {
	return &web.EventContext{R: httptestRequest()}
}

func httptestRequest() *http.Request {
	r, _ := http.NewRequest(http.MethodPost, "/", nil)
	return r
}

// postFormRequest builds a POST request whose urlencoded body carries the given
// form values, so EditingBuilder.RunSetterFunc → ParseForm populates PostForm.
func postFormRequest(vals url.Values) *http.Request {
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

// selType is a foreign model with a COMPOSITE primary key (ID + LocaleCode) and
// a slug decoder, standing in for a locale-scoped model whose record is named
// "<id>_<locale>" (e.g. "1_pt-BR") — the form the admin uses in its URLs.
type selType struct {
	ID         uint   `gorm:"primaryKey"`
	LocaleCode string `gorm:"primaryKey;type:varchar(8)"`
	Name       string
}

// PrimaryColumnValuesBySlug splits "1_pt-BR" into its key columns, so a model
// selector can resolve a record chosen by that slug (SlugDecoder).
func (selType) PrimaryColumnValuesBySlug(slug string) map[string]string {
	parts := strings.SplitN(slug, "_", 2)
	m := map[string]string{"ID": parts[0]}
	if len(parts) > 1 {
		m["LocaleCode"] = parts[1]
	}
	return m
}

// selOwner belongs to selType through a composite foreign key.
type selOwner struct {
	ID             uint     `gorm:"primaryKey"`
	Type           *selType `gorm:"foreignKey:TypeID,TypeLocaleCode;references:ID,LocaleCode"`
	TypeID         uint
	TypeLocaleCode string
}

// selManyType and selManyOwner model a many-to-many with NO "<Field>ID" column.
type selManyType struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type selManyOwner struct {
	ID    uint           `gorm:"primaryKey"`
	Types []*selManyType `gorm:"many2many:selmany_owner_types"`
}

// selSimpleType is a foreign model with a single integer primary key.
type selSimpleType struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

type selSimpleOwner struct {
	ID     uint `gorm:"primaryKey"`
	Type   *selSimpleType
	TypeID uint
}

func newSelDB(t *testing.T, models ...any) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(models...); err != nil {
		t.Fatal(err)
	}
	return db
}

// TestForeignKeyFields_Single resolves the single FK column of a belongs-to.
func TestForeignKeyFields_Single(t *testing.T) {
	db := newSelDB(t, &selSimpleType{}, &selSimpleOwner{})
	got := gormutils.ForeignKeyFields(db, &selSimpleOwner{}, "Type")
	if len(got) != 1 || got[0] != "TypeID" {
		t.Fatalf("ForeignKeyFields = %v, want [TypeID]", got)
	}
}

// TestForeignKeyFields_Composite resolves BOTH FK columns of a composite
// belongs-to, in related-primary-key order.
func TestForeignKeyFields_Composite(t *testing.T) {
	db := newSelDB(t, &selType{}, &selOwner{})
	got := gormutils.ForeignKeyFields(db, &selOwner{}, "Type")
	want := map[string]bool{"TypeID": true, "TypeLocaleCode": true}
	if len(got) != 2 || !want[got[0]] || !want[got[1]] {
		t.Fatalf("ForeignKeyFields = %v, want [TypeID TypeLocaleCode] (any order)", got)
	}
}

// TestModelSelector_SlugMultiFieldKey is the core behaviour the selector relies
// on: a record chosen by a multi-field slug ("1_pt-BR") is parsed into its key
// columns and mapped onto the owner's composite foreign key — exactly what the
// selector's setter does with the submitted value.
func TestModelSelector_SlugMultiFieldKey(t *testing.T) {
	db := newSelDB(t, &selType{}, &selOwner{})
	b := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(db))
	ownerMB := b.Model(&selOwner{})
	typeMB := b.Model(&selType{})

	// The value the form carries for the select is the record's slug.
	id, err := typeMB.ParseRecordID("1_pt-BR")
	if err != nil {
		t.Fatalf("ParseRecordID(%q): %v", "1_pt-BR", err)
	}
	if id.IsZero() {
		t.Fatal("ParseRecordID returned a zero id")
	}

	// Map the related record's composite key onto the owner's FK columns, in
	// related-primary-key order — the mapping the selector performs on save.
	fks := gormutils.ForeignKeyFields(db, &selOwner{}, "Type")
	owner := &selOwner{}
	id.Related(ownerMB.Schema(), fks...).SetTo(owner)

	if owner.TypeID != 1 {
		t.Errorf("owner.TypeID = %d, want 1", owner.TypeID)
	}
	if owner.TypeLocaleCode != "pt-BR" {
		t.Errorf("owner.TypeLocaleCode = %q, want %q", owner.TypeLocaleCode, "pt-BR")
	}
}

// TestModelSelector_ManyReadsAssociationKey is the regression for the
// many-to-many form key. A many selector has no "<Field>ID" column, so its
// values are submitted under the association field key itself ("Types") and the
// setter reads exactly that key. A hardcoded "ID" suffix in the component put the
// values under "TypesID" (which the setter never reads), so a selection was never
// applied and the association was cleared on save. The submitted "Types" values
// must land in owner.Types.
func TestModelSelector_ManyReadsAssociationKey(t *testing.T) {
	db := newSelDB(t, &selManyType{}, &selManyOwner{})
	b := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(db))
	mb := b.Model(&selManyOwner{})
	NewModelSelectorBuilder(mb, "Types").SetMany(true).
		SetForeignModel(b.Model(&selManyType{})).Build()

	ctx := &web.EventContext{R: postFormRequest(url.Values{"Types": {"1", "2"}})}

	owner := &selManyOwner{ID: 10}
	if vErr := mb.Editing().RunSetterFunc(&presets.FieldsSetterOptions{}, ctx, false, owner); vErr.HaveErrors() {
		t.Fatalf("setter: %s", vErr.Error())
	}

	if len(owner.Types) != 2 {
		t.Fatalf("owner.Types = %d entries, want 2 (values must be read under %q, not %q)",
			len(owner.Types), "Types", "TypesID")
	}
	got := map[uint]bool{}
	for _, ty := range owner.Types {
		got[ty.ID] = true
	}
	if !got[1] || !got[2] {
		t.Fatalf("owner.Types ids = %v, want {1,2}", got)
	}

	// The other half of the contract: the rendered editor must submit under the
	// same key the setter reads. web.VField binds v-model="form[\"<key>\"]"; a
	// many field must bind form["Types"], never form["TypesID"] — otherwise the
	// selection is submitted under a key the setter ignores.
	var buf strings.Builder
	comp := mb.Editing().ToComponent(&presets.ToComponentOptions{SkipPermVerify: true},
		&selManyOwner{ID: 10}, presets.FieldModeStack{presets.EDIT}, ctx)
	if err := h.Fprint(&buf, comp, ctx.Context()); err != nil {
		t.Fatalf("render edit form: %v", err)
	}
	html := buf.String()
	if !strings.Contains(html, `form["Types"]`) {
		t.Errorf("edit form does not bind form[%q]", "Types")
	}
	if strings.Contains(html, `form["TypesID"]`) {
		t.Errorf("edit form binds form[%q]; a many field has no such column and the setter never reads it", "TypesID")
	}
}

// TestModelSelector_OmitsAssociationOnSave proves the selector writes the FK from
// its column, not from a loaded (possibly stale) association: an owner carrying
// TypeID=2 but a Type struct pointing at id 1 must persist type 2.
func TestModelSelector_OmitsAssociationOnSave(t *testing.T) {
	db := newSelDB(t, &selSimpleType{}, &selSimpleOwner{})
	db.Create(&selSimpleType{ID: 1, Name: "one"})
	db.Create(&selSimpleType{ID: 2, Name: "two"})
	db.Create(&selSimpleOwner{ID: 10, TypeID: 1})

	b := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(db))
	mb := b.Model(&selSimpleOwner{})
	NewModelSelectorBuilder(mb, "Type").SetForeignModel(b.Model(&selSimpleType{})).Build()

	// The FK says type 2; the loaded association still points at type 1 (stale).
	owner := &selSimpleOwner{ID: 10, TypeID: 2, Type: &selSimpleType{ID: 1, Name: "one"}}
	if err := mb.Editing().Saver(owner, mb.MustRecordID(owner), selEventContext()); err != nil {
		t.Fatalf("save: %v", err)
	}

	var stored selSimpleOwner
	if err := db.First(&stored, 10).Error; err != nil {
		t.Fatal(err)
	}
	if stored.TypeID != 2 {
		t.Fatalf("stored TypeID = %d, want 2 (FK authoritative, stale association omitted)", stored.TypeID)
	}
}
