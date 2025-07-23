package presets

import (
	"log"
	"net/http"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/printer"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/go-rvq/rvq/x/ui/vuetifyx"
	"go.uber.org/zap"
)

func (b *Builder) MuxSetup(f func(prefix string, r *http.ServeMux)) *Builder {
	b.muxSetup = append(b.muxSetup, f)
	return b
}

func (b *Builder) SetupRoutes(mux *http.ServeMux) {
	if routesDebug {
		b.logger.Info("initializing mux for", zap.Reflect("models", modelNames(b.models)), zap.String("prefix", b.prefix))
	}

	ub := b.builder

	mainJSPath := b.prefix + "/assets/main.js"
	mux.Handle("GET "+mainJSPath,
		ub.PacksHandler("text/javascript",
			Vuetify(),
			JSComponentsPack(),
			vuetifyx.JSComponentsPack(),
			web.JSComponentsPack(),
		),
	)

	if routesDebug {
		log.Println("mounted url:", mainJSPath)
	}

	vueJSPath := b.prefix + "/assets/vue.js"
	mux.Handle("GET "+vueJSPath,
		ub.PacksHandler("text/javascript",
			web.JSVueComponentsPack(),
		),
	)

	mux.Handle("GET "+b.prefix+"/assets/main.css",
		ub.PacksHandler("text/css",
			vuetifyx.CSSComponentsPack(),
		),
	)

	HandleMaterialDesignIcons(b.prefix, mux)

	if routesDebug {
		log.Println("mounted url:", vueJSPath)
	}

	for _, ea := range b.extraAssets {
		fullPath := b.extraFullPath(ea)
		mux.Handle("GET "+fullPath, ub.PacksHandler(
			ea.contentType,
			ea.body,
		))

		if routesDebug {
			log.Println("mounted url:", fullPath)
		}
	}

	homeURL := b.prefix

	if homeURL == "" {
		homeURL = "/{$}"
	}

	mux.Handle(
		homeURL,
		b.Wrap(b.layoutFunc(b.getHomePageFunc(), b.homePageLayoutConfig)),
	)

	for _, page := range b.pages {
		b.pageHandlers.Add(page.Build(b.prefix))
	}

	b.pageHandlers.SetupRoutes(mux, func(pattern string, ph *PageHandler) {
		if routesDebug {
			log.Printf("mounted url: %s\n", pattern)
		}
	})

	for _, m := range b.models {
		m.SetupRoutes(mux)
	}

	for _, f := range b.muxSetup {
		f(b.prefix, mux)
	}

	// b.handler = mux
	// Handle 404
	b.handler = b.middleware(mux)
}

func (b *Builder) middleware(handler http.Handler) http.Handler {
	notFoundHandler := b.Wrap(b.layoutFunc(b.getNotFoundPageFunc(), b.notFoundPageLayoutConfig))
	return printer.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedResponse := &wrapedResponseWriter{web.WrapResponseWriter(w)}
		handler.ServeHTTP(capturedResponse, r)
		if capturedResponse.StatusCode() == http.StatusNotFound {
			if !b.skipNotFoundHandler(r) {
				// If no other handler wrote to the response, assume 404 and write our custom response.
				notFoundHandler.ServeHTTP(w, r)
			}
		}
		return
	}))
}

func (b *Builder) AddWrapHandler(key string, f func(in http.Handler) (out http.Handler)) {
	b.wrapHandlers[key] = f
}

func (b *Builder) WrapModel(m *ModelBuilder, pf web.PageFunc) http.Handler {
	return b.Wrap(pf, func(p *web.PageBuilder) {
		if m == nil {
			return
		}
		m.registerDefaultEventFuncs()
		p.MergeHub(m.EventsHub.Wrap(WrapEventHandler))
	})
}

func (b *Builder) Wrap(pf web.PageFunc, doPage ...func(p *web.PageBuilder)) http.Handler {
	p := b.builder.Page(pf)

	for _, f := range doPage {
		f(p)
	}

	handlers := b.I18n().EnsureLanguage(p)

	for _, wrapHandler := range b.wrapHandlers {
		handlers = wrapHandler(handlers)
	}

	return handlers
}

type EventHandlerWrapper struct {
	web.EventHandler
}

func (w *EventHandlerWrapper) Handle(ctx *web.EventContext) (r web.EventResponse, err error) {
	r, err = w.EventHandler.Handle(ctx)
	if ctx.W.Writed() {
		return
	}

	var msgs FlashMessages
	if err != nil {
		if err == perm.PermissionDenied {
			err = MustGetMessages(ctx.Context()).ErrPermissionDenied
		}
		msgs.Append(NewFlashMessage(err))
		err = nil
	} else {
		msgs = GetFlashMessages(ctx)
		if ctx.Flash != nil {
			msgs.Append(NewFlashMessages(ctx.Flash)...)
		}
	}

	if v, ok := ctx.ContextValue(CtxEventHandlerWrapperNoFlash).(bool); !ok || !v {
		RespondFlash(&r, msgs)
	} else {
		if len(msgs) > 0 {
			ctx.Flash = msgs
			web.RespondFlashCookie(web.FlashCookieName, "", []byte{}, ctx, &web.PageResponse{})
		}
	}
	return
}

func WrapEventHandler(wh web.EventHandler) web.EventHandler {
	if hw, ok := wh.(*EventHandlerWrapper); ok {
		return hw
	}
	return &EventHandlerWrapper{EventHandler: wh}
}
