package presets

import (
	"context"
	"net/http"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

const BreadcrumbsKey = "BreadcrumbsKey"

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
	items []*Breadcrumb
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

func (b *BreadcrumbsBuilder) Component(youAreHere string) h.HTMLComponent {
	var (
		children = make([]h.HTMLComponent, (len(b.items) * 2))
		i        = 1
	)
	children[0] = VBreadcrumbsItem(h.Text(youAreHere), h.Span(":")).Class("font-italic v-breadcrumbs-item-youarehere")
	for _, item := range b.items[:len(b.items)-1] {
		children[i] = h.A(h.Text(item.Label)).Href(item.URI)
		children[i+1] = VBreadcrumbsDivider()
		i += 2
	}
	children[i] = VBreadcrumbsItem(h.Text(b.items[len(b.items)-1].Label)).Active(true)
	return VBreadcrumbs(children...).Style("padding:0")
}

func AddModelsTreeToBreadcrumb(rooted bool, ctx *web.EventContext, parents []*ModelBuilder, parentsID IDSlice, bc *BreadcrumbsBuilder) (records []any, err error) {
	parentsIDPtr := web.GetContextValuer(context.WithValue(ctx.R.Context(), ParentsModelIDKey, parentsID), ParentsModelIDKey)
	records = make([]any, len(parentsID))

	if rooted {
		if root := parents[0]; root.menuGroupName != "" {
			bc.Append(&Breadcrumb{
				Label: root.p.menuGroups.MenuGroup(root.menuGroupName).TTitle(ctx.Context()),
			})
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
