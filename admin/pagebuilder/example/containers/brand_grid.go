package containers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	. "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/pagebuilder"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
)

type BrandGrid struct {
	ID             uint
	AddTopSpace    bool
	AddBottomSpace bool
	AnchorID       string
	Brands         Brands `sql:"type:text;"`
}

type Brand struct {
	Image media_library.MediaBox `sql:"type:text;"`
	Name  string
}

func (*BrandGrid) TableName() string {
	return "container_brand_grids"
}

type Brands []*Brand

func (this Brands) Value() (driver.Value, error) {
	return json.Marshal(this)
}

func (this *Brands) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), this)
	case []byte:
		return json.Unmarshal(v, this)
	default:
		return errors.New("not supported")
	}
}

func RegisterBrandGridContainer(pb *pagebuilder.Builder, db *gorm.DB) {
	vb := pb.RegisterContainer("BrandGrid").Group("Content").
		RenderFunc(func(obj interface{}, input *pagebuilder.RenderInput, ctx *web.EventContext) HTMLComponent {
			v := obj.(*BrandGrid)
			return BrandGridBody(v, input)
		})
	mb := vb.Model(&BrandGrid{})
	eb := mb.Editing("AddTopSpace", "AddBottomSpace", "AnchorID", "Brands")

	fb := pb.GetPresetsBuilder().NewFieldsBuilder(presets.WRITE).Model(&Brand{}).Only("Image", "Name")
	fb.Field("Image").WithContextValue(media.MediaBoxConfig, &media_library.MediaBoxConfig{
		AllowType: "image",
	})

	eb.Field("Brands").Nested(fb, &presets.DisplayFieldInSorter{Field: "Name"})
}

func BrandGridBody(data *BrandGrid, input *pagebuilder.RenderInput) (body HTMLComponent) {
	body = ContainerWrapper(data.AnchorID, "container-brand_grid",
		"", "", "",
		"", data.AddTopSpace, data.AddBottomSpace, "",
		Div(
			BrandsBody(data.Brands, input),
		).Class("container-wrapper"),
	)
	return
}

func BrandsBody(brands []*Brand, input *pagebuilder.RenderInput) HTMLComponent {
	brandsDiv := Div().Class("container-brand_grid-wrap")
	for _, b := range brands {
		img := LazyImageHtml(b.Image)
		if input.IsEditor {
			img = ImageHtml(b.Image)
		}
		brandsDiv.AppendChildren(
			Div(
				img,
			).Class("container-brand_grid-item"),
		)
	}
	return brandsDiv
}
