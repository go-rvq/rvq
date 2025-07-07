package vuetify

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

type VPullToRefreshBuilder struct {
	VTagBuilder[*VPullToRefreshBuilder]
}

func VPullToRefresh(children ...h.HTMLComponent) *VPullToRefreshBuilder {
	return VTag(&VPullToRefreshBuilder{}, "v-pull-to-refresh", children...)
}

func (b *VPullToRefreshBuilder) PullDownThreshold(v int) (r *VPullToRefreshBuilder) {
	b.Attr(":pull-down-threshold", fmt.Sprint(v))
	return b
}

func (b *VPullToRefreshBuilder) On(name string, value string) (r *VPullToRefreshBuilder) {
	b.Attr(fmt.Sprintf("v-on:%s", name), value)
	return b
}

func (b *VPullToRefreshBuilder) Bind(name string, value string) (r *VPullToRefreshBuilder) {
	b.Attr(fmt.Sprintf("v-bind:%s", name), value)
	return b
}
