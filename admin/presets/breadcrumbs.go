package presets

import (
	"context"
	"net/http"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

type breadcrumbsKey string

// BreadcrumbsKey is the context key of the request's breadcrumbs.
const BreadcrumbsKey breadcrumbsKey = "BreadcrumbsKey"

func GetOrInitBreadcrumbs(r *http.Request) (bc *BreadcrumbsBuilder) {
	if bc = GetBreadcrumbs(r); bc == nil {
		bc = &BreadcrumbsBuilder{}
		r2 := r.WithContext(context.WithValue(r.Context(), BreadcrumbsKey, bc))
		*r = *r2
	}
	return bc
}

func GetBreadcrumbs(r *http.Request) *BreadcrumbsBuilder {
	if v := r.Context().Value(BreadcrumbsKey); v != nil {
		return v.(*BreadcrumbsBuilder)
	}
	return nil
}

type Breadcrumb struct {
	URI   string
	Label string
}

type BreadcrumbsBuilder struct {
	items             []*Breadcrumb
	youAreHereDisable bool
}

func (b *BreadcrumbsBuilder) YouAreHere(v bool) *BreadcrumbsBuilder {
	b.youAreHereDisable = !v
	return b
}

func (b *BreadcrumbsBuilder) Items() []*Breadcrumb {
	return b.items
}

func (b *BreadcrumbsBuilder) Labels() (s []string) {
	s = make([]string, len(b.items))
	for i, item := range b.items {
		s[i] = item.Label
	}
	return
}

func (b *BreadcrumbsBuilder) Append(item ...*Breadcrumb) {
	b.items = append(b.items, item...)
}

func (b *BreadcrumbsBuilder) Prepend(item ...*Breadcrumb) {
	b.items = append(item, b.items...)
}

func (b *BreadcrumbsBuilder) Empty() bool {
	return len(b.items) == 0
}

// Component is the trail: a compact line under the title — the home an icon,
// the steps quiet links between small chevrons, the page itself last, in
// evidence —; youAreHere names it to who reads the page aloud (aria-label).
func (b *BreadcrumbsBuilder) Component(youAreHere string) h.HTMLComponent {
	var children h.HTMLComponents
	last := len(b.items) - 1
	for i, item := range b.items {
		var comp h.HTMLComponent = h.Text(item.Label)
		if i == 0 && item.URI != "" && last > 0 {
			// the home: its icon, its name in its title
			comp = VIcon("mdi-home-outline").Size(SizeSmall)
		}
		if i < last {
			if item.URI != "" {
				comp = h.A(comp).Href(item.URI).Class("rvq-breadcrumbs__link").Attr("title", item.Label)
			} else {
				comp = h.Span("").Children(comp).Class("rvq-breadcrumbs__group")
			}
			children = append(children,
				VBreadcrumbsItem(comp).Class("rvq-breadcrumbs__item"),
				VBreadcrumbsDivider(VIcon("mdi-chevron-right").Size(SizeXSmall)).Class("rvq-breadcrumbs__divider"),
			)
			continue
		}
		children = append(children, VBreadcrumbsItem(comp).Class("rvq-breadcrumbs__item rvq-breadcrumbs__current").
			Active(true).Attr("aria-current", "page"))
	}
	crumbs := VBreadcrumbs(children...).Class("flex-wrap rvq-presets-breadcrumbs").Density(DensityCompact)
	if !b.youAreHereDisable && youAreHere != "" {
		crumbs.Attr("aria-label", youAreHere)
	}
	return h.Nav(crumbs).Class("rvq-presets-breadcrumbs-wraper")
}

func AddModelsTreeToBreadcrumb(rooted bool, ctx *web.EventContext, parents []*ModelBuilder, parentsID IDSlice, bc *BreadcrumbsBuilder) (records []any, err error) {
	parentsIDPtr := web.GetContextValuer(context.WithValue(ctx.R.Context(), ParentsModelIDKey, parentsID), ParentsModelIDKey)
	records = make([]any, len(parentsID))

	if rooted {
		if g := parents[0].MenuGroup(); g != nil {
			for _, g := range append(g.Ancestors(), g) {
				bc.Append(&Breadcrumb{Label: g.TTitle(ctx.Context())})
			}
		}
	}

	for i, id := range parentsID {
		if parents[i].singleton {
			bc.Append(&Breadcrumb{
				Label: parents[i].TTitle(ctx.Context()),
				URI:   parents[i].modelInfo.ListingHref(parentsID[:i]...),
			})
		} else {
			bc.Append(&Breadcrumb{
				Label: parents[i].TTitlePlural(ctx.Context()),
				URI:   parents[i].modelInfo.ListingHref(parentsID[:i]...),
			})

			var (
				r     = parents[i].NewModel()
				label string
			)

			id.SetTo(r)

			parentsIDPtr.Set(parentsID[:i])
			if label, err = parents[i].RecordTitleFetch(r, ctx); err != nil {
				return
			}

			bc.Append(&Breadcrumb{
				Label: label,
				URI:   parents[i].Info().DetailingHref(id.String(), parentsID[:i]...),
			})

			records[i] = r
		}
	}

	return
}
