package admin

import (
	"net/http"

	. "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/docs/cmd/qor5/admin-template/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

func Initialize() *http.ServeMux {
	b := setupAdmin()
	mux := setupRouter(b)

	return mux
}

func setupAdmin() (b *presets.Builder) {
	db := ConnectDB()

	// Initialize the builder of QOR5
	b = presets.New()

	// Set up the project name, ORM and Homepage
	b.URIPrefix("/admin").
		BrandTitle("Admin").
		DataOperator(gorm2op.DataOperator(db)).
		HomePageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
			r.Body = VContainer(
				H1("Home"),
				P().Text("Change your home page here"))
			return
		})

	// Register Post into the builder
	// Use m to customize the model, Or config more models here.
	m := b.Model(&models.Post{})
	_ = m

	return
}
