package vuetifyx

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	v "github.com/go-rvq/rvq/x/ui/vuetify"

	"github.com/go-rvq/rvq/web"
	h "github.com/theplant/htmlgo"
)

type VXTablePaginationBuilder struct {
	total           int64
	currPage        int64
	perPage         int64
	customPerPages  []int64
	noPerPagePart   bool
	onSelectPerPage interface{}
	onSelectPage    interface{}
	onPrevPage      interface{}
	onNextPage      interface{}
	pageInfoText    string
	pageText        string
	ofPageText      string
	perPageText     string
}

func VXTablePagination() *VXTablePaginationBuilder {
	return &VXTablePaginationBuilder{}
}

func (tpb *VXTablePaginationBuilder) PageInfoText(v string) *VXTablePaginationBuilder {
	tpb.pageInfoText = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) PerPageText(v string) *VXTablePaginationBuilder {
	tpb.perPageText = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) PageText(v string) *VXTablePaginationBuilder {
	tpb.pageText = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) OfPageText(v string) *VXTablePaginationBuilder {
	tpb.ofPageText = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) Total(v int64) *VXTablePaginationBuilder {
	tpb.total = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) CurrPage(v int64) *VXTablePaginationBuilder {
	tpb.currPage = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) PerPage(v int64) *VXTablePaginationBuilder {
	tpb.perPage = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) CustomPerPages(v []int64) *VXTablePaginationBuilder {
	tpb.customPerPages = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) NoPerPagePart(v bool) *VXTablePaginationBuilder {
	tpb.noPerPagePart = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) OnSelectPerPage(v interface{}) *VXTablePaginationBuilder {
	tpb.onSelectPerPage = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) OnSelectPage(v interface{}) *VXTablePaginationBuilder {
	tpb.onSelectPage = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) OnPrevPage(v interface{}) *VXTablePaginationBuilder {
	tpb.onPrevPage = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) OnNextPage(v interface{}) *VXTablePaginationBuilder {
	tpb.onNextPage = v
	return tpb
}

func (tpb *VXTablePaginationBuilder) MarshalHTML(ctx context.Context) ([]byte, error) {
	if tpb.onSelectPerPage == nil {
		tpb.OnSelectPerPage(web.Plaid().
			PushState(true).
			Query("per_page", web.Var("[$event]")).
			MergeQuery(true).
			Go())
	}
	if tpb.onSelectPage == nil {
		tpb.OnSelectPage(web.Plaid().
			PushState(true).
			Query("page", web.Var("[$event.target.value]")).
			MergeQuery(true).
			Go())
	}
	if tpb.onPrevPage == nil {
		tpb.OnPrevPage(web.Plaid().
			PushState(true).
			Query("page", tpb.currPage-1).
			MergeQuery(true).
			Go())
	}
	if tpb.onNextPage == nil {
		tpb.OnNextPage(web.Plaid().
			PushState(true).
			Query("page", tpb.currPage+1).
			MergeQuery(true).
			Go())
	}

	var sItems []string
	{
		perPagesM := map[int64]struct{}{
			10:  {},
			20:  {},
			50:  {},
			100: {},
			500: {},
		}
		if tpb.perPage > 0 {
			perPagesM[tpb.perPage] = struct{}{}
		}
		for _, v := range tpb.customPerPages {
			if v <= 0 {
				continue
			}
			perPagesM[v] = struct{}{}
		}
		perPages := make([]int, 0, len(perPagesM))
		for k := range perPagesM {
			perPages = append(perPages, int(k))
		}
		sort.Ints(perPages)
		for _, v := range perPages {
			sItems = append(sItems, fmt.Sprint(v))
		}
	}

	currPageStart := (tpb.currPage-1)*tpb.perPage + 1
	currPageEnd := tpb.currPage * tpb.perPage
	if currPageEnd > tpb.total {
		currPageEnd = tpb.total
	}

	totalPages := int(math.Ceil(float64(tpb.total) / float64(tpb.perPage)))

	canNext := false
	canPrev := false
	if tpb.currPage*tpb.perPage < tpb.total {
		canNext = true
	}
	if tpb.currPage > 1 {
		canPrev = true
	}
	var nextIconStyle string
	var prevIconStyle string
	if canNext {
		nextIconStyle = "cursor: pointer;"
	}
	if canPrev {
		prevIconStyle = "cursor: pointer;"
	}

	rowsPerPageText := "Rows per page: "
	if tpb.perPageText != "" {
		rowsPerPageText = tpb.perPageText
	}

	pageInfoText := "{currPageStart}-{currPageEnd} of {total}"
	if tpb.pageInfoText != "" {
		pageInfoText = tpb.pageInfoText
	}

	pageText := "Page:"
	if tpb.pageText != "" {
		pageText = tpb.pageText
	}

	ofPageText := "of {totalPages}"
	if tpb.ofPageText != "" {
		ofPageText = tpb.ofPageText
	}
	ofPageText = strings.ReplaceAll(ofPageText, "{total}", fmt.Sprint(totalPages))

	pageInfoText = strings.NewReplacer(
		"{currPageStart}", fmt.Sprint(currPageStart),
		"{currPageEnd}", fmt.Sprint(currPageEnd),
		"{total}", fmt.Sprint(tpb.total),
	).Replace(pageInfoText)

	return h.Div(
		v.VRow().Justify("end").Align("center").Class("ma-0").
			Children(
				h.If(!tpb.noPerPagePart,
					h.Div(
						h.Text(rowsPerPageText),
					),
					h.Div(
						v.VSelect().Items(sItems).Variant("underlined").ModelValue(fmt.Sprint(tpb.perPage)).
							HideDetails(true).Density("compact").Attr("style", "margin-top: -8px").
							Attr("@update:model-value", tpb.onSelectPerPage),
					).Style("width: 64px;").Class("ml-3"),
				),
				h.Div(
					h.Text(pageInfoText),
				).Class("ml-3"),
				h.If(totalPages > 1,
					h.Div(
						h.Text(pageText),
					).Class("ml-6"),
					h.Div(
						v.VTextField().Variant("underlined").ModelValue(fmt.Sprint(tpb.currPage)).
							HideDetails(true).Density("compact").Attr("style", "margin-top: -8px").
							Attr("@keyup.enter", tpb.onSelectPage),
					).Style("width: 40px;").Class("ml-3"),
					h.Div(
						h.Text(ofPageText),
					).Class("ml-3"),
				),
				h.Div(
					h.Span("").Style(prevIconStyle).Children(
						v.VBtn("").Variant("text").Icon("mdi-chevron-left").Size(32).Disabled(!canPrev).
							Attr("@click", tpb.onPrevPage),
					),
					h.Span("").Style(nextIconStyle).Children(
						v.VBtn("").Variant("text").Icon("mdi-chevron-right").Size(32).Disabled(!canNext).
							Attr("@click", tpb.onNextPage),
					).Class("ml-3"),
				).Class("ml-6"),
			)).MarshalHTML(ctx)
}
