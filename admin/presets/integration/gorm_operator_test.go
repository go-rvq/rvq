package integration_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/theplant/gofixtures"
)

type TestVariant struct {
	ProductCode string
	ColorCode   string
	Name        string
}

var emptyData = gofixtures.Data(gofixtures.Sql(``, []string{"test_variants"}))

func (tv *TestVariant) PrimarySlug() string {
	return fmt.Sprintf("%s_%s", tv.ProductCode, tv.ColorCode)
}

func (tv *TestVariant) PrimaryColumnValuesBySlug(slug string) map[string]string {
	segs := strings.Split(slug, "_")
	if len(segs) != 2 {
		panic("wrong slug")
	}

	return map[string]string{
		"product_code": segs[0],
		"color_code":   segs[1],
	}
}

func TestPrimarySlugger(t *testing.T) {
	db := TestDB
	db.AutoMigrate(&TestVariant{})
	rawDB, _ := db.DB()
	emptyData.TruncatePut(rawDB)
	op := gorm2op.DataOperator(db)
	ctx := new(web.EventContext)

	// the operator takes a parsed record id now, and fills the object it is
	// given instead of returning one
	schema, err := op.Schema(&TestVariant{})
	if err != nil {
		panic(err)
	}
	id, err := presets.ParseRecordID(schema, "P01_C01")
	if err != nil {
		panic(err)
	}

	if err = op.Save(&TestVariant{ProductCode: "P01", ColorCode: "C01", Name: "Product 1"}, presets.ID{}, ctx); err != nil {
		panic(err)
	}

	if err = op.Save(&TestVariant{ProductCode: "P01", ColorCode: "C01", Name: "Product 2"}, id, ctx); err != nil {
		panic(err)
	}

	tv := &TestVariant{}
	if err = op.Fetch(tv, id, ctx); err != nil {
		panic(err)
	}

	if tv.Name != "Product 2" {
		t.Error("didn't update product 2", tv)
	}

	if err = op.Delete(&TestVariant{}, id, false, ctx); err != nil {
		panic(err)
	}

	if err = op.Fetch(&TestVariant{}, id, ctx); err != presets.ErrRecordNotFound {
		t.Error("didn't return not found after delete", err)
	}
}
