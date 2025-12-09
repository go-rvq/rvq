package web

import (
	"context"
	"net/http"
	"net/url"
)

type contextKey int

func (k contextKey) String() string {
	return contextKeyNames[k]
}

const (
	EventContextKey contextKey = iota
	VueLocaleSetterContextKey
	UrlQueryKey
)

var contextKeyNames = [...]string{
	EventContextKey:           "EventContext",
	VueLocaleSetterContextKey: "VueLocaleSetter",
	UrlQueryKey:               "UrlQuery",
}

type UrlQuery struct {
	Values url.Values
}

func (u *UrlQuery) Get(key string) string {
	return u.Values.Get(key)
}

func (u *UrlQuery) Has(key string) bool {
	return u.Values.Has(key)
}

func (u *UrlQuery) Set(key, value string) {
	u.Values.Set(key, value)
}

func (u *UrlQuery) Del(key string) {
	u.Values.Del(key)
}

func UrlQueryFromRequest(r *http.Request) (q *UrlQuery) {
	q, _ = r.Context().Value(UrlQueryKey).(*UrlQuery)
	return
}

func MustGetFromContext[T any](ctx context.Context, key contextKey) (r T) {
	var ok bool
	r, ok = ctx.Value(EventContextKey).(T)
	if !ok {
		panic(key.String() + " required")
	}
	return
}

func ContextWithEventContext(parent context.Context, ctx *EventContext) (r context.Context) {
	r = context.WithValue(parent, EventContextKey, ctx)
	return
}

func EventContextFromContext(parent context.Context) (r *EventContext) {
	r, _ = parent.Value(EventContextKey).(*EventContext)
	return
}

func MustGetEventContext(c context.Context) (r *EventContext) {
	r, _ = c.Value(EventContextKey).(*EventContext)
	if r == nil {
		panic("EventContext required")
	}
	return
}

func InjectorFromContext(c context.Context) (r *PageInjector) {
	return MustGetFromContext[*EventContext](c, EventContextKey).Injector
}

func ContextWithVueLocaleSetter(ctx ContextValuer, v *VueLayoutLocaleSetter) {
	ctx.WithContextValue(VueLocaleSetterContextKey, v)
}

func VueLocaleSetterFromContext(ctx ContextValuer) (v *VueLayoutLocaleSetter) {
	v, _ = ctx.ContextValue(VueLocaleSetterContextKey).(*VueLayoutLocaleSetter)
	return
}
