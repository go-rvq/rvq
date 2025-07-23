package presets

import (
	"log"
	"net/http"

	"github.com/theplant/osenv"
)

var routesDebug = osenv.GetBool("ADMIN_ROUTES_DEBUG", "Debug mounted routes", false)

func (mb *ModelBuilder) SetupRoutes(mux *http.ServeMux) {
	var (
		info            = mb.Info()
		routePath       = info.ListingHref()
		listingPageFunc = mb.listing.GetPageFunc()
		itemRoutePath   = routePath
	)

	if mb.singleton {
		mb.itemRoutePath = itemRoutePath

		if mb.layoutConfig == nil {
			mb.layoutConfig = &LayoutConfig{}
		}

		mb.layoutConfig.SearchBoxInvisible = true
		editPath := routePath
		edit := mb.p.layoutFunc(mb.BindPageFunc(mb.editing.defaultPageFunc), mb.layoutConfig)

		if mb.hasDetailing {
			mux.Handle(
				routePath,
				mb.p.WrapModel(mb, mb.p.detailLayoutFunc(mb.BindPageFunc(mb.detailing.defaultPageFunc), mb.layoutConfig)),
			)
			editPath += "/edit"

			if routesDebug {
				log.Printf("mounted url: %s\n", routePath)
			}

			// no wrap model
			mux.Handle(editPath, mb.p.Wrap(edit))
		} else {
			mux.Handle(editPath, mb.p.WrapModel(mb, edit))
		}

		if routesDebug {
			log.Printf("mounted url: %s\n", editPath)
		}

		mb.listing.pagesRegistrator.Build().SetupRoutes(mux, func(pattern string, ph *PageHandler) {
			if routesDebug {
				log.Printf("mounted url: %s\n", pattern)
			}
		})

		mb.detailing.pageHandlers.WithPathPrefix(routePath).SetupRoutes(mux, func(pattern string, ph *PageHandler) {
			if routesDebug {
				log.Printf("mounted url: %s\n", pattern)
			}
		})

		/*
			mb.detailing.pagesRegistrator.Build().SetupRoutes(mux, func(pattern string, ph *PageHandler) {
				if routesDebug {
					log.Printf("mounted url: %s\n", pattern)
				}
			})
		*/
	} else {
		mb.routePath = routePath

		mux.Handle(
			routePath,
			mb.p.WrapModel(mb, mb.p.layoutFunc(mb.BindPageFunc(listingPageFunc), mb.layoutConfig)),
		)

		if routesDebug {
			log.Printf("mounted url: %s\n", routePath)
		}

		mb.listing.pages.WithPathPrefix(itemRoutePath).SetupRoutes(mux, func(pattern string, ph *PageHandler) {
			if routesDebug {
				log.Printf("mounted url: %s\n", pattern)
			}
		})

		mb.listing.pagesRegistrator.Build().SetupRoutes(mux, func(pattern string, ph *PageHandler) {
			if routesDebug {
				log.Printf("mounted url: %s\n", pattern)
			}
		})

		itemRoutePath += "/{id}"
		mb.itemRoutePath = itemRoutePath

		if mb.hasDetailing {
			mux.Handle(
				itemRoutePath,
				mb.p.Wrap(mb.p.detailLayoutFunc(mb.BindPageFunc(mb.detailing.GetPageFunc()), mb.layoutConfig)),
			)
			if routesDebug {
				log.Printf("mounted url: %s\n", itemRoutePath)
			}
		}

		mb.detailing.pageHandlers.WithPathPrefix(itemRoutePath).SetupRoutes(mux, func(pattern string, ph *PageHandler) {
			if routesDebug {
				log.Printf("mounted url: %s\n", pattern)
			}
		})

		mb.detailing.pagesRegistrator.Build().SetupRoutes(mux, func(pattern string, ph *PageHandler) {
			if routesDebug {
				log.Printf("mounted url: %s\n", pattern)
			}
		})

		if !mb.editingDisabled {
			{
				routePath := itemRoutePath + "/edit"
				mux.Handle(
					routePath,
					mb.p.Wrap(mb.p.detailLayoutFunc(mb.BindPageFunc(mb.editing.GetPageFunc()), mb.layoutConfig)),
				)

				if routesDebug {
					log.Printf("mounted url: %s\n", routePath)
				}
			}
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

	for _, f := range mb.itemRouteSetuper {
		f(mux, itemRoutePath)
	}
}
