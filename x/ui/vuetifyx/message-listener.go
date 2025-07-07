package vuetifyx

import (
	"context"

	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type VXMessageListenerBuilder struct {
	v.VTagBuilder[*VXMessageListenerBuilder]
	listenFunc string
}

func VXMessageListener() (r *VXMessageListenerBuilder) {
	return v.VTag(&VXMessageListenerBuilder{}, "vx-messagelistener")
}

func (b *VXMessageListenerBuilder) ListenFunc(v string) (r *VXMessageListenerBuilder) {
	b.listenFunc = v
	return b
}

func (b *VXMessageListenerBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	if b.listenFunc != "" {
		b.Attr(":listen-func", b.listenFunc)
	}

	return b.GetHTMLTagBuilder().MarshalHTML(ctx)
}
