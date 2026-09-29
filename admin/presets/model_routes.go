package presets

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-rvq/rvq/x/osenv"
)

var routesDebug = osenv.GetBool("RVQ_ADMIN_ROUTES_DEBUG", "Debug mounted routes", false)

func (mb *ModelBuilder) SetupRoutes(mux *http.ServeMux) {
	var (
		info            = mb.Info()
		routePath       = info.ListingHref()
		listingPageFunc = mb.listing.GetPageFunc()
		itemRoutePath   = routePath
	)

	// The listing route, read by the pages registered on the listing to know
	// where they hang. Set before anything is mounted.
	mb.routePath = routePath

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

		// Only the detailing: a singleton has one record and no listing, so the
		// menu it shows is the record's. A page registered on the listing has
		// no menu to appear in and nothing to be a child of, so it would be
		// mounted nowhere — a page that looks registered and answers 404. That
		// is a mistake in the setup, and the boot stops on it.
		if paths := mb.listing.registeredPagePaths(); len(paths) > 0 {
			panic(fmt.Sprintf("presets: %q is a singleton and has no listing, so the page(s) "+
				"registered on its Listing() would be mounted nowhere and appear in no menu: %s. "+
				"Register them on Detailing() instead.", mb.id, strings.Join(paths, ", ")))
		}

		mb.detailing.pageHandlers.WithPathPrefix(routePath).SetupRoutes(mux, func(pattern string, ph *PageHandler) {
			if routesDebug {
				log.Printf("mounted url: %s\n", pattern)
			}
		})

		mb.detailing.pagesRegistrator.Build().SetupRoutes(mux, func(pattern string, ph *PageHandler) {
			if routesDebug {
				log.Printf("mounted url: %s\n", pattern)
			}
		})
	} else {
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

		if !mb.creatingDisabled {
			// the create form as a page: this is the address a create overlay
			// shows while it is open (see FormHostBuilder.URL), so reloading or
			// sharing it has to render the same form. Mounted before `/{id}` —
			// a literal segment wins over the wildcard.
			newPath := routePath + "/new"
			mux.Handle(
				newPath,
				mb.p.WrapModel(mb, mb.p.detailLayoutFunc(mb.BindPageFunc(mb.editing.GetCreatingPageFunc()), mb.layoutConfig)),
			)

			if routesDebug {
				log.Printf("mounted url: %s\n", newPath)
			}
		}

		itemRoutePath += "/{id}"
		mb.itemRoutePath = itemRoutePath

		if mb.hasDetailing {
			mux.Handle(
				itemRoutePath,
				mb.p.WrapModel(mb, mb.p.detailLayoutFunc(mb.BindPageFunc(mb.detailing.GetPageFunc()), mb.layoutConfig)),
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
					mb.p.WrapModel(mb, mb.p.detailLayoutFunc(mb.BindPageFunc(mb.editing.GetPageFunc()), mb.layoutConfig)),
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
