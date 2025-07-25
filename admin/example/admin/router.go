package admin

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-rvq/rvq/admin/example/models"
	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/sitemap"
	"gorm.io/gorm"
)

//go:embed assets/favicon.ico
var favicon []byte

const (
	logoutURL = "/auth/logout"

	exportOrdersURL = "/export-orders"
)

func TestHandler(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()
	c := NewConfig(db)
	u := &models.User{
		Model: gorm.Model{ID: 888},
		Roles: []role.Role{
			{
				Name: "admin",
			},
		},
	}
	m := login.MockCurrentUser(u)
	mux.Handle("/page_builder/", m(c.pageBuilder))
	mux.Handle("/", m(c.pb))
	return mux
}

func Router(db *gorm.DB) http.Handler {
	c := NewConfig(db)

	mux := http.NewServeMux()
	loginBuilder.Mount(mux)
	//	mux.Handle("/frontstyle.css", c.pb.GetWebBuilder().PacksHandler("text/css", web.ComponentsPack(`
	// :host {
	//	all: initial;
	//	display: block;
	// div {
	//	background-color:orange;
	// }
	// `)))
	// example of seo
	mux.Handle("/posts/first", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var post models.Post
		db.First(&post)
		seodata, _ := seoBuilder.Render(post, r).MarshalHTML(r.Context())
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><head>%s</head><body>%s</body></html>`, seodata, post.Body)
	}))

	mux.Handle("/page_builder/", c.pageBuilder)
	mux.Handle("/", c.pb)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(favicon)
		return
	})

	mux.Handle(exportOrdersURL, exportOrders(db))

	// example of sitemap and robot
	sitemap.SiteMap("product").RegisterRawString("https://dev.qor5.com/admin", "/product").MountTo(mux)
	robot := sitemap.Robots()
	robot.Agent(sitemap.AlexaAgent).Allow("/product1", "/product2").Disallow("/admin")
	robot.Agent(sitemap.GoogleAgent).Disallow("/admin")
	robot.MountTo(mux)

	cr := chi.NewRouter()
	cr.Use(
		loginBuilder.Middleware(),
		validateSessionToken(db),
		withRoles(db),
		withNoteContext(),
		securityMiddleware(),
	)
	cr.Mount("/", mux)
	return cr
}
