package web

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/go-playground/form/v4"
	h "github.com/go-rvq/htmlgo"
	"github.com/mpvl/unique"
	"github.com/sunfmin/reflectutils"
)

type ComponentWrappers []func(comp h.HTMLComponent) h.HTMLComponent

func (w ComponentWrappers) Wrap(comp h.HTMLComponent) h.HTMLComponent {
	for _, f := range w {
		comp = f(comp)
	}
	return comp
}

type PageResponse struct {
	PageTitle   string
	Actions     h.HTMLComponents
	Menu        h.HTMLComponents
	Body        h.HTMLComponent
	Wrapers     ComponentWrappers
	RedirectURL string
}

func (pr *PageResponse) AddAction(action ...h.HTMLComponent) {
	pr.Actions = append(pr.Actions, action...)
}

func (pr *PageResponse) AddMenuItem(item ...h.HTMLComponent) {
	pr.Menu = append(pr.Menu, item...)
}

func (pr *PageResponse) Wrap(f ...func(comp h.HTMLComponent) h.HTMLComponent) {
	pr.Wrapers = append(pr.Wrapers, f...)
}

type PortalUpdateOptions struct {
	ResetScroll bool     `json:"resetScroll,omitempty"`
	Mounted     []string `json:"mounted,omitempty"`
	Unmounted   []string `json:"unmounted,omitempty"`
}

func (u *PortalUpdateOptions) OnMounted(script string) *PortalUpdateOptions {
	u.Mounted = append(u.Mounted, script)
	return u
}

func (u *PortalUpdateOptions) OnUnmounted(script string) *PortalUpdateOptions {
	u.Unmounted = append(u.Unmounted, script)
	return u
}

type PortalUpdate struct {
	Name    string               `json:"name,omitempty"`
	Body    h.HTMLComponent      `json:"body,omitempty"`
	Defer   bool                 `json:"defer,omitempty"`
	Options *PortalUpdateOptions `json:"options,omitempty"`
}

func (u *PortalUpdate) Wrap(f func(comp h.HTMLComponent) h.HTMLComponent) *PortalUpdate {
	u.Body = f(u.Body)
	return u
}

func (u *PortalUpdate) WithOptions(f func(opts *PortalUpdateOptions)) *PortalUpdate {
	f(u.Options)
	return u
}

// @snippet_begin(EventResponseDefinition)
type EventResponse struct {
	PageTitle     string           `json:"pageTitle,omitempty"`
	Body          h.HTMLComponent  `json:"body,omitempty"`
	Reload        bool             `json:"reload,omitempty"`
	PushState     *LocationBuilder `json:"pushState"`             // This we don't omitempty, So that {} can be kept when use url.Values{}
	RedirectURL   string           `json:"redirectURL,omitempty"` // change window url without push state
	ReloadPortals []string         `json:"reloadPortals,omitempty"`
	UpdatePortals []*PortalUpdate  `json:"updatePortals,omitempty"`
	Data          interface{}      `json:"data,omitempty"` // used for return collection data like TagsInput data source
	RunScript     string           `json:"runScript,omitempty"`
	// used with InitContextVars to set values for example vars.show to used by v-model

	deferedPortals map[string]bool
}

func (r *EventResponse) UpdatePortal(name string, body h.HTMLComponent, doOptions ...func(opts *PortalUpdateOptions)) *EventResponse {
	r.UpdatePortalR(name, body, doOptions...)
	return r
}

func (r *EventResponse) UpdatePortalR(name string, body h.HTMLComponent, doOptions ...func(opts *PortalUpdateOptions)) (pu *PortalUpdate) {
	for _, p := range r.UpdatePortals {
		if p.Name == name {
			panic("Duplicate Portal '" + name + "' Update")
		}
	}
	pu = &PortalUpdate{
		Name:    name,
		Body:    body,
		Defer:   r.deferedPortals[name],
		Options: &PortalUpdateOptions{},
	}
	for _, f := range doOptions {
		f(pu.Options)
	}
	r.UpdatePortals = append(r.UpdatePortals, pu)
	return pu
}

func (r *EventResponse) DeferedPortal(name string) *EventResponse {
	if r.deferedPortals == nil {
		r.deferedPortals = make(map[string]bool)
	}
	r.deferedPortals[name] = true
	return r
}

func (r *EventResponse) AppendRunScript(script string) {
	r.RunScript += "; " + script
}

// @snippet_end

// @snippet_begin(PageFuncAndEventFuncDefinition)
type (
	PageFunc        func(ctx *EventContext) (r PageResponse, err error)
	PageFuncWrapper func(old PageFunc) PageFunc
	EventHandler    interface {
		Handle(ctx *EventContext) (r EventResponse, err error)
	}
	EventFunc func(ctx *EventContext) (r EventResponse, err error)
)

var PageFuncDefaultWrap PageFuncWrapper = func(old PageFunc) PageFunc {
	return old
}

func (f EventFunc) Handle(ctx *EventContext) (r EventResponse, err error) {
	return f(ctx)
}

// @snippet_end

type LayoutFunc func(in PageFunc) PageFunc

// @snippet_begin(EventHandlerHubDefinition)
type EventHandlerHub interface {
	RegisterEventHandler(eventFuncId string, ef EventHandler) (key string)
}

// @snippet_end

func AppendRunScripts(er *EventResponse, scripts ...string) {
	if er.RunScript != "" {
		scripts = append([]string{er.RunScript}, scripts...)
	}
	er.RunScript = strings.Join(scripts, "; ")
}

type EventFuncID struct {
	ID string `json:"id,omitempty"`
}

type ContextValuePointer struct {
	dot, child context.Context
	key        any
	value      reflect.Value
}

func (p *ContextValuePointer) Get() interface{} {
	return p.value.Interface()
}

func (p *ContextValuePointer) Set(value interface{}) {
	p.value.Set(reflect.ValueOf(value))
}

func (p *ContextValuePointer) With(value interface{}) func() {
	old := p.Get()
	p.Set(value)
	return func() {
		p.Set(old)
	}
}

func (p *ContextValuePointer) Parent() context.Context {
	parent := reflect.Indirect(reflect.ValueOf(p.dot).Elem()).FieldByName("Context")
	if parent.IsValid() {
		return parent.Interface().(context.Context)
	}
	return nil
}

func (p *ContextValuePointer) Top() (top *ContextValuePointer) {
	parent := p.Parent()
	top = p
	if parent == nil {
		return
	}

	p = getContextValuer(p.dot, parent, p.key)
	for p != nil && parent != nil {
		top = p
		parent = top.Parent()
		if parent != nil {
			p = getContextValuer(top.dot, parent, p.key)
		}
	}
	return
}

func (p *ContextValuePointer) Delete() context.Context {
	parent := p.Parent()
	if p.child == nil {
		return parent
	}
	parentField := reflect.Indirect(reflect.ValueOf(p.child).Elem()).FieldByName("Context")
	parentField.Set(reflect.ValueOf(parent))
	return p.child
}

var valueCtxType = reflect.TypeOf(context.WithValue(context.Background(), "a", nil)).Elem()

func getContextValuer(child, ctx context.Context, key any) *ContextValuePointer {
	contextValues := reflect.Indirect(reflect.ValueOf(ctx))
	contextKeys := reflect.TypeOf(ctx)
	for contextKeys.Kind() == reflect.Ptr {
		contextKeys = contextKeys.Elem()
	}

	if contextValues.Type() == valueCtxType {
		keyField := contextValues.FieldByName("key")
		keyValue := reflect.NewAt(keyField.Type(), unsafe.Pointer(keyField.UnsafeAddr())).Elem()
		if keyValue.Interface() == key {
			if valueField := contextValues.FieldByName("val"); valueField.IsValid() {
				value := reflect.NewAt(valueField.Type(), unsafe.Pointer(valueField.UnsafeAddr())).Elem()
				return &ContextValuePointer{
					dot:   ctx,
					child: child,
					key:   key,
					value: value,
				}
			}
		}
	}

	if contextField := contextValues.FieldByName("Context"); contextField.IsValid() {
		return getContextValuer(ctx, contextField.Interface().(context.Context), key)
	}
	return nil
}

func GetContextValuer(ctx context.Context, key any) *ContextValuePointer {
	return getContextValuer(nil, ctx, key)
}

func WithContextValue(ctx *EventContext, key any, value interface{}) (done func()) {
	if ptr := GetContextValuer(ctx.R.Context(), key); ptr != nil {
		return ptr.With(value)
	}
	ctx.WithContextValue(key, value)
	return func() {
		ctx.R = ctx.R.WithContext(GetContextValuer(ctx.R.Context(), key).Top().Delete())
	}
}

func GetContexValue(key any, ctx ...context.Context) (value any) {
	for _, ctx := range ctx {
		if ctx != nil {
			if value = ctx.Value(key); value != nil {
				return
			}
		}
	}
	return
}

type ContextValuer interface {
	WithContextValue(key any, value any)
	ContextValue(key any) any
	Context() context.Context
}

type RequestContext interface {
	ContextValuer
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Param(key string) (r string)
}

type EventContext struct {
	R         *http.Request
	W         ResponseWriter
	Resp      *EventResponse
	Injector  *PageInjector
	Flash     interface{} // pass value from actions to index
	i         int64
	dataStack []any
}

func (e *EventContext) UrlQuery() *UrlQuery {
	return UrlQueryFromRequest(e.R)
}

func (e *EventContext) GetUrlQueryValue(key string) string {
	return e.UrlQuery().Get(key)
}

func (e *EventContext) Data() any {
	if len(e.dataStack) > 0 {
		return e.dataStack[len(e.dataStack)-1]
	}
	return nil
}

func (e *EventContext) WithData(v any) (reset func()) {
	l := len(e.dataStack)
	e.dataStack = append(e.dataStack, v)
	return func() {
		e.dataStack = e.dataStack[:l]
	}
}

func (e *EventContext) WithContextValue(key any, value any) {
	e.R = e.R.WithContext(context.WithValue(e.R.Context(), key, value))
}

func (e *EventContext) ContextValue(key any) any {
	return e.R.Context().Value(key)
}

func (e *EventContext) Context() context.Context {
	return e.R.Context()
}

func (e *EventContext) Request() *http.Request {
	return e.R
}

func (e *EventContext) ResponseWriter() http.ResponseWriter {
	return e.W
}

func (e *EventContext) Param(key string) (r string) {
	r = e.R.PathValue(key)
	if len(r) == 0 {
		r = e.R.FormValue(key)
	}
	return
}

func (e *EventContext) ParamAsInt(key string) (r int) {
	strVal := e.Param(key)
	if len(strVal) == 0 {
		return
	}
	val, _ := strconv.ParseInt(strVal, 10, 64)
	r = int(val)
	return
}

func (e *EventContext) Queries() (r url.Values) {
	r = e.R.URL.Query()
	delete(r, EventFuncIDName)
	return
}

func (ctx *EventContext) MustUnmarshalForm(v interface{}) {
	err := ctx.UnmarshalForm(v)
	if err != nil {
		panic(err)
	}
}

func (e *EventContext) UID() string {
	return "_" + fmt.Sprint(time.Now().UnixNano())
}

type CustoFormTypeDecoder struct {
	Decoder form.DecodeCustomTypeFunc
	Types   []any
}

var FormTypeDecoders []CustoFormTypeDecoder

func (ctx *EventContext) FormSliceValues(key string) (r []string) {
	for _, key := range ctx.FormSliceKeys(key) {
		for _, sufix := range []string{"", ".__value"} {
			if v, _ := ctx.R.MultipartForm.Value[key.Key+sufix]; len(v) > 0 {
				if s := v[0]; s != "" {
					r = append(r, s)
				}
			}
		}
	}
	return
}

func (ctx *EventContext) FormSliceKeys(key string) (r []struct {
	Key   string
	Index int
}) {
	mf := ctx.R.MultipartForm
	if mf == nil {
		return
	}

	var index []int

loop:
	for fkey := range mf.Value {
		if strings.HasPrefix(fkey, key+"[") {
			var (
				s   = fkey[len(key)+1:]
				i   int
				err error
			)

			for i := range s {
				if s[i] >= '0' && s[i] <= '9' {
					continue
				}
				if s[i] != ']' {
					continue loop
				}
				s = s[:i]
				break
			}

			if i, err = strconv.Atoi(s); err == nil {
				index = append(index, i)
			}
		}
	}

	if len(index) > 0 {
		unique.Sort(unique.IntSlice{&index})
	}

	for _, i := range index {
		r = append(r, struct {
			Key   string
			Index int
		}{Key: key + "[" + strconv.Itoa(i) + "]", Index: i})
	}
	return
}

func (ctx *EventContext) UnmarshalFormValues(values url.Values, v interface{}) (err error) {
	dec := form.NewDecoder()

	for _, decoder := range FormTypeDecoders {
		dec.RegisterCustomTypeFunc(decoder.Decoder, decoder.Types...)
	}

	if err = dec.Decode(v, values); err != nil {
		return
	}

	// Reorder list-editor slices by the per-item `__pos` metadata (the ListEditor
	// submits the intended order as `<slice>[i].__pos`), so the decoded order
	// follows the user's sort instead of the raw array index. The form values are
	// reindexed to the new order too, so later per-item reads (e.g. `__deleted`)
	// stay aligned. Items without `__pos` are pushed to the end (stable order).
	ReorderSlicesByPos(v, values)
	return
}

// posFieldSuffix is the item metadata key carrying its position in a list editor.
const (
	posFieldSuffix     = ".__pos"
	deletedFieldSuffix = ".__deleted"
	newFieldSuffix     = ".__new"
)

var reItemPos = regexp.MustCompile(`^(.*)\[(\d+)\]\.__pos$`)

// ReorderSlicesByPos reorders, for every list-editor slice found in values (a key
// matching `<prefix>[i].__pos`), the decoded slice in v by the `__pos` value, and
// reindexes the matching form keys in values to the new order. Items without a
// `__pos` keep their relative order and are appended after the positioned ones.
// Prefixes are processed deepest-first so nested lists are ordered before their
// parents.
func ReorderSlicesByPos(v interface{}, values url.Values) {
	// collect, per slice prefix, the (originalIndex -> pos) pairs.
	type item struct{ orig, pos int }
	groups := map[string][]item{}
	for k, vs := range values {
		m := reItemPos.FindStringSubmatch(k)
		if m == nil || m[1] == "" || len(vs) == 0 {
			continue
		}
		orig, _ := strconv.Atoi(m[2])
		pos, e := strconv.Atoi(strings.TrimSpace(vs[0]))
		if e != nil {
			pos = 1 << 30 // no/invalid pos -> end
		}
		groups[m[1]] = append(groups[m[1]], item{orig: orig, pos: pos})
	}
	if len(groups) == 0 {
		return
	}

	prefixes := make([]string, 0, len(groups))
	for p := range groups {
		prefixes = append(prefixes, p)
	}
	// deepest first: more '[' means more nested.
	sort.SliceStable(prefixes, func(i, j int) bool {
		return strings.Count(prefixes[i], "[") > strings.Count(prefixes[j], "[")
	})

	for _, prefix := range prefixes {
		items := groups[prefix]

		slice, e := reflectutils.Get(v, prefix)
		if e != nil || slice == nil {
			continue
		}
		sv := reflect.ValueOf(slice)
		if sv.Kind() != reflect.Slice {
			continue
		}
		n := sv.Len()

		// order: positioned items by (pos, orig), then any uncovered indexes.
		sort.SliceStable(items, func(i, j int) bool {
			if items[i].pos != items[j].pos {
				return items[i].pos < items[j].pos
			}
			return items[i].orig < items[j].orig
		})
		order := make([]int, 0, n)
		covered := make(map[int]bool, len(items))
		for _, it := range items {
			if it.orig >= 0 && it.orig < n && !covered[it.orig] {
				order = append(order, it.orig)
				covered[it.orig] = true
			}
		}
		for i := 0; i < n; i++ {
			if !covered[i] {
				order = append(order, i)
			}
		}
		if sameOrder(order) {
			continue
		}

		// reorder the decoded slice.
		newSlice := reflect.MakeSlice(sv.Type(), 0, n)
		for _, oi := range order {
			newSlice = reflect.Append(newSlice, sv.Index(oi))
		}
		if err := reflectutils.Set(v, prefix, newSlice.Interface()); err != nil {
			continue
		}

		// reindex the form keys of this slice to the new order.
		reindexSliceKeys(values, prefix, order)
	}
}

func sameOrder(order []int) bool {
	for i, o := range order {
		if i != o {
			return false
		}
	}
	return true
}

// reindexSliceKeys rewrites, in values, every key `prefix[oldIndex]<rest>` to
// `prefix[newIndex]<rest>` where newIndex is the position of oldIndex in order.
func reindexSliceKeys(values url.Values, prefix string, order []int) {
	remap := make(map[int]int, len(order))
	for newIdx, oldIdx := range order {
		remap[oldIdx] = newIdx
	}
	replacement := make(map[string][]string)
	var toDelete []string
	for k, vs := range values {
		if !strings.HasPrefix(k, prefix+"[") {
			continue
		}
		rest := k[len(prefix)+1:]
		close := strings.IndexByte(rest, ']')
		if close < 0 {
			continue
		}
		oldIdx, e := strconv.Atoi(rest[:close])
		if e != nil {
			continue
		}
		newIdx, ok := remap[oldIdx]
		if !ok || newIdx == oldIdx {
			continue
		}
		newKey := fmt.Sprintf("%s[%d]%s", prefix, newIdx, rest[close+1:])
		replacement[newKey] = vs
		toDelete = append(toDelete, k)
	}
	for _, k := range toDelete {
		delete(values, k)
	}
	for k, vs := range replacement {
		values[k] = vs
	}
}

func (ctx *EventContext) UnmarshalForm(v interface{}) (err error) {
	if ctx.R.MultipartForm == nil {
		return
	}

	mf := ctx.R.MultipartForm

	if err = ctx.UnmarshalFormValues(mf.Value, v); err != nil {
		return
	}

	if len(mf.File) > 0 {
		for k, vs := range mf.File {
			// set slice
			if err2 := reflectutils.Set(v, k, vs); err2 != nil &&
				err2.Error() == "reflect.Set: value of type []*multipart.FileHeader is not "+
					"assignable to type multipart.FileHeader" {
				if len(vs) == 0 {
					// set to nil
					reflectutils.Set(v, k, nil)
				} else {
					// set first value
					reflectutils.Set(v, k, vs[0])
				}
			}
		}
	}
	return
}
