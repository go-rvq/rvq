package presets

import (
	"net/url"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/ui/vuetifyx"
)

func (b *ListingBuilder) WrapFilterDataFunc(f func(old FilterDataFunc) FilterDataFunc) *ListingBuilder {
	b.filterDataFunc = f(b.filterDataFunc)
	return b
}

func (b *ListingBuilder) FilterDataFunc(v FilterDataFunc) *ListingBuilder {
	if v == nil {
		b.filterDataFunc = nil
		return b
	}

	b.filterDataFunc = func(ctx *web.EventContext) (r vuetifyx.FilterData) {
		fd := v(ctx)
		for _, k := range fd {
			if k == nil {
				continue
			}
			k.Key = "f_" + k.Key
			r = append(r, k)
		}
		return
	}
	return b
}

func (b *ListingBuilder) WrapFilterTabsFunc(f func(old FilterTabsFunc) FilterTabsFunc) *ListingBuilder {
	b.filterTabsFunc = f(b.filterTabsFunc)
	return b
}

func (b *ListingBuilder) FilterTabsFunc(v FilterTabsFunc) *ListingBuilder {
	if v == nil {
		b.filterTabsFunc = nil
		return b
	}

	b.filterTabsFunc = func(ctx *web.EventContext) []*FilterTab {
		fts := v(ctx)
		for _, ft := range fts {
			newQuery := make(url.Values)
			for k := range ft.Query {
				newQuery["f_"+k] = ft.Query[k]
			}
			ft.Query = newQuery
		}
		return fts
	}
	return b
}

// FilterTabsWrapper changes the listing's tabs — the FilterTabsFunc's, or
// none — when it renders: to add one of every listing of a kind, as the
// trash of the soft-deleted records, whatever tabs the model has, set before
// or after it.
type FilterTabsWrapper func(ctx *web.EventContext, tabs []*FilterTab) []*FilterTab

// AppendFilterTabsWrapper adds w to the wrappers of the tabs, applied in the
// order they were added, after the FilterTabsFunc.
func (b *ListingBuilder) AppendFilterTabsWrapper(w FilterTabsWrapper) *ListingBuilder {
	b.filterTabsWrappers = append(b.filterTabsWrappers, w)
	return b
}

// allFilterTabs are the listing's tabs: the FilterTabsFunc's, through the
// wrappers. nil when there is none.
func (b *ListingBuilder) allFilterTabs(ctx *web.EventContext) (tabs []*FilterTab) {
	if b.filterTabsFunc != nil {
		tabs = b.filterTabsFunc(ctx)
	}
	for _, w := range b.filterTabsWrappers {
		tabs = w(ctx, tabs)
	}
	return
}
