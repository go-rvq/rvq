package gorm2op

import (
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
)

type (
	CallbackMerger func(mode Mode, params *presets.SearchParams, ctx *web.EventContext) []CallbacksRegistrator[*DataOperatorBuilder]

	Callback func(state *CallbackState) (err error)

	CallbacksRegistrator[T any] struct {
		create,
		update,
		delete,
		fetch,
		search Callbacks[T]
	}

	CallbacksGetter[T any] interface {
		Callbacks() *CallbacksRegistrator[T]
	}
)

func NewCallbacks[T CallbacksGetter[T]](dot T) T {
	return dot.Callbacks().SetDot(dot)
}

func NewCallbacksRegistrator() *CallbacksRegistrator[*DataOperatorBuilder] {
	return &CallbacksRegistrator[*DataOperatorBuilder]{}
}

func (b *CallbacksRegistrator[T]) Callbacks() *CallbacksRegistrator[T] {
	return b
}

func (b *CallbacksRegistrator[T]) Dot() T {
	return b.create.dot
}

func (b *CallbacksRegistrator[T]) WithCallbacks(c *Callbacks[T], do func(cb *Callbacks[T])) T {
	do(c)
	return b.Dot()
}

func (b *CallbacksRegistrator[T]) EachModeSplitCallbacks(mode Mode, do func(cb *Callbacks[T])) *CallbacksRegistrator[T] {
	for _, m := range mode.Split() {
		do(b.ModeCallbacks(m))
	}
	return b
}

func (b *CallbacksRegistrator[T]) CreateCallbacks() *Callbacks[T] {
	return &b.create
}

func (b *CallbacksRegistrator[T]) UpdateCallbacks() *Callbacks[T] {
	return &b.update
}

func (b *CallbacksRegistrator[T]) DeleteCallbacks() *Callbacks[T] {
	return &b.delete
}

func (b *CallbacksRegistrator[T]) SearchCallbacks() *Callbacks[T] {
	return &b.search
}

func (b *CallbacksRegistrator[T]) FetchCallbacks() *Callbacks[T] {
	return &b.fetch
}

func (b *CallbacksRegistrator[T]) ModeCallbacks(mode Mode) (cb *Callbacks[T]) {
	switch mode {
	case Search:
		cb = b.SearchCallbacks()
	case Create:
		cb = b.CreateCallbacks()
	case Fetch:
		cb = b.FetchCallbacks()
	case FetchTitle:
		cb = b.FetchCallbacks()
	case Update:
		cb = b.UpdateCallbacks()
	case Delete:
		cb = b.DeleteCallbacks()
	}
	return
}

func (b *CallbacksRegistrator[T]) WithCreateCallbacks(do func(cb *Callbacks[T])) T {
	return b.WithCallbacks(&b.create, do)
}

func (b *CallbacksRegistrator[T]) WithUpdateCallbacks(do func(cb *Callbacks[T])) T {
	return b.WithCallbacks(&b.update, do)
}

func (b *CallbacksRegistrator[T]) WithDeleteCallbacks(do func(cb *Callbacks[T])) T {
	return b.WithCallbacks(&b.delete, do)
}

func (b *CallbacksRegistrator[T]) WithSearchCallbacks(do func(cb *Callbacks[T])) T {
	return b.WithCallbacks(&b.search, do)
}

func (b *CallbacksRegistrator[T]) WithFetchCallbacks(do func(cb *Callbacks[T])) T {
	return b.WithCallbacks(&b.fetch, do)
}

func (b *CallbacksRegistrator[T]) WithModeCallbacks(mode Mode, do func(cb *Callbacks[T])) T {
	return b.WithCallbacks(b.ModeCallbacks(mode), do)
}

func (b *CallbacksRegistrator[T]) WithModeSplitCallbacks(mode Mode, do func(cb *Callbacks[T])) T {
	b.EachModeSplitCallbacks(mode, do)
	return b.Dot()
}

func (b *CallbacksRegistrator[T]) SetDot(dot T) T {
	b.create.dot = dot
	b.update.dot = dot
	b.delete.dot = dot
	b.fetch.dot = dot
	b.search.dot = dot
	return dot
}

func (b *CallbacksRegistrator[T]) Merge(other *CallbacksRegistrator[T]) *CallbacksRegistrator[T] {
	b.create.Merge(&other.create)
	b.update.Merge(&other.update)
	b.delete.Merge(&other.delete)
	b.fetch.Merge(&other.fetch)
	b.search.Merge(&other.search)
	return b
}

type Callbacks[T any] struct {
	pre,
	post []Callback
	dot T
}

func (b *Callbacks[T]) Merge(other *Callbacks[T]) *Callbacks[T] {
	b.pre = append(b.pre, other.pre...)
	b.post = append(b.post, other.post...)
	return b
}

func (b *Callbacks[T]) Pre(f ...Callback) *Callbacks[T] {
	b.pre = append(b.pre, f...)
	return b
}

func (b *Callbacks[T]) Post(f ...Callback) *Callbacks[T] {
	b.post = append(b.post, f...)
	return b
}

func (b *Callbacks[T]) Build(do ...Callback) CallbackSlice {
	return append(append(b.pre, do...), b.post...)
}

func (b *Callbacks[T]) Clone() *Callbacks[T] {
	n := *b
	b.pre = make([]Callback, len(n.pre))
	b.post = make([]Callback, len(n.post))
	copy(n.pre, b.pre)
	copy(b.post, b.post)
	return &n
}

func (b *Callbacks[T]) Dot() T {
	return b.dot
}

type CallbackSlice []Callback

func (b CallbackSlice) Execute(state *CallbackState) (err error) {
	defer func() {
		var err2 error
		for _, done := range state.dones {
			if err2 = done(); err2 != nil {
				if err == nil {
					err = err2
				}
				break
			}
		}
	}()

	for _, f := range b {
		err = f(state)
		if err != nil {
			return
		}
	}
	return
}

type CallbacksKeyType string

const (
	NamedCallbacksRegistratorKey CallbacksKeyType = "callbacksRegistrator"
	CallbacksKey                 CallbacksKeyType = "callbacks"
)

type NamedCallbacksRegistrator map[string]*CallbacksRegistrator[*DataOperatorBuilder]

func (n NamedCallbacksRegistrator) When(name string) (v *CallbacksRegistrator[*DataOperatorBuilder]) {
	v = n[name]
	if v == nil {
		v = &CallbacksRegistrator[*DataOperatorBuilder]{}
		n[name] = v
	}
	return
}

func NamedCallbacksRegistratorOf(mb *presets.ModelBuilder) (v NamedCallbacksRegistrator) {
	var ok bool

	if v, ok = mb.GetData(NamedCallbacksRegistratorKey).(NamedCallbacksRegistrator); !ok {
		v = make(NamedCallbacksRegistrator)
		mb.SetData(NamedCallbacksRegistratorKey, v)
	}

	return
}

func AddCallbacksToContext(ctx *web.EventContext, cb ...*CallbacksRegistrator[*DataOperatorBuilder]) {
	ctx.WithContextValue(CallbacksKey, append(GetContextCallbacks(ctx), cb...))
}

func GetContextCallbacks(ctx *web.EventContext) (v []*CallbacksRegistrator[*DataOperatorBuilder]) {
	v, _ = ctx.ContextValue(CallbacksKey).([]*CallbacksRegistrator[*DataOperatorBuilder])
	return
}

func (b *DataOperatorBuilder) GetCallbacks(mode Mode, ctx *web.EventContext) *Callbacks[*DataOperatorBuilder] {
	cbs := []*Callbacks[*DataOperatorBuilder]{
		b.ModeCallbacks(mode),
	}
	cbsCtx := GetContextCallbacks(ctx)

	for _, c := range cbsCtx {
		cbs = append(cbs, c.ModeCallbacks(mode))
	}

	if len(cbs) == 1 {
		return cbs[0]
	}

	cb := cbs[0].Clone()
	for _, c := range cbs[1:] {
		cb.Merge(c)
	}
	return cb
}
