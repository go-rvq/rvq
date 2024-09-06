package presets

import (
	"net/url"

	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/ui/vuetifyx"
)

func (b *ListingBuilder) FilterDataFunc(v FilterDataFunc) *ListingBuilder {
	if v == nil {
		b.filterDataFunc = nil
		return b
	}

	b.filterDataFunc = func(ctx *web.EventContext) vuetifyx.FilterData {
		fd := v(ctx)
		for _, k := range fd {
			k.Key = "f_" + k.Key
		}
		return fd
	}
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
