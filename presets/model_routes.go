package presets

import (
	"log"
	"net/http"

	"github.com/theplant/osenv"
)

var routesDebug = osenv.GetBool("ADMIN_ROUTES_DEBUG", "Debug mounted routes", false)

func (mb *ModelBuilder) SetupRoutes(mux *http.ServeMux) {
	info := mb.Info()
	routePath := info.ListingHref()
	inPageFunc := mb.listing.GetPageFunc()

	if mb.singleton {
		inPageFunc = mb.BindPageFunc(mb.editing.defaultPageFunc)
		if mb.layoutConfig == nil {
			mb.layoutConfig = &LayoutConfig{}
		}
		mb.layoutConfig.SearchBoxInvisible = true
	}

	mux.Handle(
		routePath,
		mb.p.Wrap(mb, mb.p.layoutFunc(mb.BindPageFunc(inPageFunc), mb.layoutConfig)),
	)

	if routesDebug {
		log.Printf("mounted url: %s\n", routePath)
	}

	var itemRoutePath = routePath

	if !mb.singleton {
		itemRoutePath += "/{id}"
	}

	if mb.hasDetailing {
		mux.Handle(
			itemRoutePath,
			mb.p.Wrap(mb, mb.p.detailLayoutFunc(mb.BindPageFunc(mb.detailing.GetPageFunc()), mb.layoutConfig)),
		)
		if routesDebug {
			log.Printf("mounted url: %s\n", itemRoutePath)
		}
	}

	{
		routePath := itemRoutePath + "/edit"
		mux.Handle(
			routePath,
			mb.p.Wrap(mb, mb.p.detailLayoutFunc(mb.BindPageFunc(mb.editing.GetPageFunc()), mb.layoutConfig)),
		)

		if routesDebug {
			log.Printf("mounted url: %s\n", routePath)
		}
	}

	for _, child := range mb.children {
		child.SetupRoutes(mux)
	}

	if mb.subRoutesSetup != nil {
		mb.subRoutesSetup(mux, itemRoutePath)
	}

	if mb.routeSetuper != nil {
		mb.routeSetuper(mux, routePath)
	}

	if mb.itemRouteSetuper != nil {
		mb.itemRouteSetuper(mux, itemRoutePath)
	}
}
