package admin

import (
	"github.com/go-rvq/rvq/admin/example/models"
	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
)

func configNestedFieldDemo(b *presets.Builder, db *gorm.DB) {
	cust := b.Model(&models.Customer{}).RightDrawerWidth("50%").
		Label("NestedFieldDemos").URIName("nested-field-demos")

	addFb := b.NewFieldsBuilder(presets.WRITE).Model(&models.Address{}).Only("Street", "HomeImage", "Phones")

	addFb.Field("HomeImage").WithContextValue(media.MediaBoxConfig, &media_library.MediaBoxConfig{
		AllowType: "image",
		Sizes: map[string]*base.Size{
			"thumb": {
				Width:  400,
				Height: 300,
			},
			"main": {
				Width:  800,
				Height: 500,
			},
		},
	})

	// nested slices bind the field to a model builder of the item type
	phoneFb := b.NewFieldsBuilder(presets.WRITE).Model(&models.Phone{}).Only("Number")
	phoneMB := presets.NewModelBuilder(b, &models.Phone{})
	addFb.Field("Phones").Nested(presets.NestedSlice(phoneMB, phoneFb).SetDisplayFieldInSorter("Number"))

	ed := cust.Editing("Name", "Addresses", "MembershipCard")
	addressMB := presets.NewModelBuilder(b, &models.Address{})
	ed.Field("Addresses").Nested(presets.NestedSlice(addressMB, addFb).SetDisplayFieldInSorter("Street"))

	cardFb := b.NewFieldsBuilder(presets.WRITE).Model(&models.MembershipCard{}).Only("Number", "ValidBefore")
	cardMB := presets.NewModelBuilder(b, &models.MembershipCard{})
	ed.Field("MembershipCard").Nested(presets.NestedSlice(cardMB, cardFb))

	ed.FetchFunc(func(obj interface{}, id presets.ID, ctx *web.EventContext) (err error) {
		return gorm2op.DataOperator(db.Preload("Addresses.Phones").Preload("MembershipCard")).Fetch(obj, id, ctx)
	})

	ed.SaveFunc(func(obj interface{}, id presets.ID, ctx *web.EventContext) (err error) {
		c := obj.(*models.Customer)
		err = db.Delete(&models.Phone{}, "address_id IN (select id from addresses where customer_id = ?)", c.ID).Error
		if err != nil {
			panic(err)
		}
		err = db.Delete(&models.Address{}, "customer_id = ?", c.ID).Error
		if err != nil {
			panic(err)
		}
		return gorm2op.DataOperator(db.Session(&gorm.Session{FullSaveAssociations: true})).Save(obj, id, ctx)
	})
}
