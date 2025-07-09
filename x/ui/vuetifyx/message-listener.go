package vuetifyx

import (
	h "github.com/go-rvq/htmlgo"
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

func (b *VXMessageListenerBuilder) Write(ctx *h.Context) (err error) {
	if b.listenFunc != "" {
		b.Attr(":listen-func", b.listenFunc)
	}

	return b.VTagBuilder.Write(ctx)
}
