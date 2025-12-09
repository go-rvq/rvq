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

func (b *BreadcrumbsBuilder) Component(youAreHere string) h.HTMLComponent {
	var (
		children h.HTMLComponents
		i        = 1
	)

	for _, item := range b.items[:len(b.items)-1] {
		comp := h.Text(item.Label)
		if len(item.URI) > 0 {
			comp = h.A(comp).Href(item.URI)
		}
		children = append(children, VBreadcrumbsItem(comp).Style("text-wrap-mode: nowrap"), VBreadcrumbsDivider())
		i += 2
	}
	children = append(children, VBreadcrumbsItem(h.Text(b.items[len(b.items)-1].Label)).Style("text-wrap-mode: nowrap").Active(true))
	crumbs := VBreadcrumbs(children...).Class("flex-wrap rvq-presets-breadcrumbs").Style("padding:0")

	if !b.youAreHereDisable && len(youAreHere) > 0 {
		return h.Div(
			h.Div(h.Text(youAreHere), h.Span(":")).Style("text-wrap-mode: nowrap").Class("font-italic v-breadcrumbs-item-youarehere"),
			crumbs,
		).Class("d-flex rvq-presets-breadcrumbs-wraper")
	}
	return crumbs
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
