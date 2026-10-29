// Package markdown is the MarkdownBuilder: fields of a model in Markdown
// (GitHub's: tables, footnotes) edited with rvq mdkit's editor
// (vx-markdown-editor: the table builder, the footnotes builder, files of
// Markdown included, pictures, preview, help), shown rendered (x/markup: the
// same rendering and sanitizing as where the text is shown).
//
//	markdown.New().Model(mb).Fields("Body").
//		Store(markdown.MediaLibraryStorer(mediaBuilder, db)).
//		Build(presets.FORM | presets.DETAIL)
//
// Its preview and its pictures are events of the model (Install): the text
// rendered by the server; a picture stored by its Storer — the media library
// (MediaLibraryStorer), or the application's own place (a record's files) —,
// its address put in the text at once. Without a Storer, no picture.
package markdown

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/fields/tiptap"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/markup"
	"github.com/go-rvq/rvq/x/perm"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// Storer stores a picture of a text (its format — png, jpeg, gif, webp —, its
// bytes, its name), giving its address: the TipTap fields' (tiptap.Storer).
// The context it is given holds the StoreContext.
type Storer = tiptap.Storer

// StoreFunc is a Storer of a function.
type StoreFunc = tiptap.StoreFunc

// MaxPicture is the most bytes of a picture sent.
const MaxPicture = 10 << 20

// StoreContext is where a picture is sent from: the event, the model, the
// record — zero, a new one — and the field.
type StoreContext struct {
	Event *web.EventContext
	Model *presets.ModelBuilder
	ID    presets.ID
	Field string
}

type storeContextKey struct{}

// StoreContextOf is the StoreContext of the context a Storer is given; nil,
// none.
func StoreContextOf(ctx context.Context) *StoreContext {
	sc, _ := ctx.Value(storeContextKey{}).(*StoreContext)
	return sc
}

// Builder is the MarkdownBuilder.
type Builder struct {
	mb         *presets.ModelBuilder
	fieldNames []string
	store      Storer
	render     func(ctx *web.EventContext, text string) string
	rows       int
	wrapComp   func(field *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent
	wrapEditor func(field *presets.FieldContext, comp *vx.VXMarkdownEditorBuilder)
	installed  bool
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Model(mb *presets.ModelBuilder) *Builder {
	b.mb = mb
	return b
}

func (b *Builder) Fields(fieldNames ...string) *Builder {
	b.fieldNames = fieldNames
	return b
}

// Store is where the pictures go; none, no picture.
func (b *Builder) Store(s Storer) *Builder {
	b.store = s
	return b
}

func (b *Builder) WrapStore(wrap func(s Storer) Storer) *Builder {
	b.store = wrap(b.store)
	return b
}

// Render is the HTML of a text (the preview, the detail); by default
// x/markup's Render of Markdown.
func (b *Builder) Render(f func(ctx *web.EventContext, text string) string) *Builder {
	b.render = f
	return b
}

// Rows is the height of the box, in lines.
func (b *Builder) Rows(n int) *Builder {
	b.rows = n
	return b
}

func (b *Builder) WrapComp(f func(field *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent) *Builder {
	b.wrapComp = f
	return b
}

func (b *Builder) WrapEditor(f func(field *presets.FieldContext, comp *vx.VXMarkdownEditorBuilder)) *Builder {
	if old := b.wrapEditor; old != nil {
		b.wrapEditor = func(field *presets.FieldContext, comp *vx.VXMarkdownEditorBuilder) {
			old(field, comp)
			f(field, comp)
		}
	} else {
		b.wrapEditor = f
	}
	return b
}

// HTML is text rendered (Render).
func (b *Builder) HTML(ctx *web.EventContext, text string) string {
	if b.render != nil {
		return b.render(ctx, text)
	}
	return markup.Render(markup.Markdown, text)
}

// eventID is the event of kind of the builder (its fields: two builders of
// a model, two Storers).
func (b *Builder) eventID(kind string) string {
	return "markdown_" + kind + "_" + strings.Join(b.fieldNames, "_")
}

// Install registers the events of the builder on its model (the preview,
// the pictures), once: Build does; a component of the builder alone
// (Editor), its caller.
func (b *Builder) Install() *Builder {
	if b.installed {
		return b
	}
	b.installed = true
	b.mb.RegisterEventFunc(b.eventID("preview"), b.previewEvent)
	b.mb.RegisterEventFunc(b.eventID("upload"), b.uploadEvent)
	return b
}

// previewEvent is the HTML of the text posted (body): {html}.
func (b *Builder) previewEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	r.Data = map[string]string{"html": b.HTML(ctx, ctx.R.FormValue("body"))}
	return
}

// pictureTypes are the formats of a picture accepted, by their media type.
var pictureTypes = map[string]string{"image/png": "png", "image/jpeg": "jpeg", "image/gif": "gif", "image/webp": "webp"}

// uploadEvent stores the picture posted (file) of the record (id; none, a
// new one) — who may edit it, or create one —: {url}, or {message}.
func (b *Builder) uploadEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	fail := func(msg string) (web.EventResponse, error) {
		r.Data = map[string]string{"message": msg}
		return r, nil
	}
	if b.store == nil {
		return fail("no pictures here")
	}
	var id presets.ID
	if id, err = b.mb.ParseRecordID(ctx.R.FormValue("id")); err != nil {
		return fail(err.Error())
	}
	var verifier *perm.Verifier
	if id.IsZero() {
		verifier = b.mb.Info().Permissioner().Creator(ctx.R, presets.ParentsModelID(ctx.R)...)
	} else {
		verifier = b.mb.Info().Permissioner().Updater(ctx.R, id, presets.ParentsModelID(ctx.R)...)
	}
	if verifier.Denied() {
		return fail(perm.PermissionDenied.Error())
	}
	file, header, err := ctx.R.FormFile("file")
	if err != nil {
		return fail(err.Error())
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, MaxPicture+1))
	if err != nil {
		return
	}
	if len(data) > MaxPicture {
		return fail(fmt.Sprintf("a picture of at most %d MB", MaxPicture>>20))
	}
	format, ok := pictureTypes[http.DetectContentType(data)]
	if !ok {
		return fail("not a picture")
	}
	alt := strings.TrimSuffix(header.Filename, path.Ext(header.Filename))
	sc := &StoreContext{Event: ctx, Model: b.mb, ID: id, Field: strings.Join(b.fieldNames, ",")}
	url, err := b.store.Store(context.WithValue(ctx.R.Context(), storeContextKey{}, sc), format, data, alt)
	if err != nil {
		return fail(err.Error())
	}
	r.Data = map[string]string{"url": url}
	return
}

// Editor is the editor of the field: its preview and its pictures (when it
// has a Storer) by the events of the builder (Install).
func (b *Builder) Editor(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	field.LoadHint()
	var (
		url  = h.JSONString(b.mb.Info().ListingHrefCtx(ctx))
		comp = vx.VXMarkdownEditor().
			Label(field.Label).
			Hint(field.Hint).
			ErrorMessages(field.Errors).
			Preview(`(body) => plaid().url(` + url + `).eventFunc(` + h.JSONString(b.eventID("preview")) + `).formData({body}).go().then((r) => (r && r.data && r.data.html) || "")`)
	)
	if b.store != nil {
		id := ""
		if rid, err := b.mb.RecordID(field.Obj); err == nil && !rid.IsZero() {
			id = rid.String()
		}
		comp.Upload(`(fd) => plaid().url(` + url + `).eventFunc(` + h.JSONString(b.eventID("upload")) + `).formData({file: fd.get("file"), id: ` +
			h.JSONString(id) + `}).go().then((r) => (r && r.data) || {})`)
	}
	if b.rows > 0 {
		comp.Rows(b.rows)
	}
	if b.wrapEditor != nil {
		b.wrapEditor(field, comp)
	}
	var r h.HTMLComponent = vue.FormField(comp).Value(field.FormKey, field.Value()).Bind()
	if b.wrapComp != nil {
		r = b.wrapComp(field, r)
	}
	return r
}

// Detail is the text of the field rendered (HTML).
func (b *Builder) Detail(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	value := field.StringValue()
	if value == "" {
		return nil
	}
	var r h.HTMLComponent = h.Div(
		h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"),
		h.Div(h.RawHTML(b.HTML(ctx, value))).Class(web.StdContentClass),
	).Class("mb-4").Attr("data-markdown-detail", field.Name)
	if b.wrapComp != nil {
		r = b.wrapComp(field, r)
	}
	return r
}

// Cell is the text of the field in the listing: its beginning, as text (its
// marks gone: rendered, then its text).
func (b *Builder) Cell(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	s := []rune(strings.Join(strings.Fields(markup.Text(markup.Rich, b.HTML(ctx, field.StringValue()))), " "))
	if len(s) > 80 {
		s = append(s[:80], '…')
	}
	return h.Td(h.Text(string(s)))
}

// AutoComponentFunc is the field in a form: the editor; read only, the text
// rendered.
func (b *Builder) AutoComponentFunc(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	if field.ReadOnly || !field.Mode.IsWrite() {
		return b.Detail(field, ctx)
	}
	return b.Editor(field, ctx)
}

func (b *Builder) ComponentFunc(mode presets.FieldMode) presets.FieldComponentFunc {
	switch {
	case mode.Is(presets.EDIT, presets.NEW):
		return b.AutoComponentFunc
	case mode.Has(presets.DETAIL):
		return b.Detail
	case mode.Has(presets.LIST):
		return b.Cell
	}
	return nil
}

func (b *Builder) setFields(fb *presets.FieldsBuilder, mode presets.FieldMode) {
	for _, name := range b.fieldNames {
		fb.Field(name).ComponentFunc(b.ComponentFunc(mode))
	}
}

// Build puts the builder on the fields of the model in mode (EDIT, NEW,
// DETAIL, LIST) and registers its events.
func (b *Builder) Build(mode presets.FieldMode) *Builder {
	b.Install()
	ed := b.mb.Editing()
	if mode.Has(presets.EDIT) {
		b.setFields(&ed.FieldsBuilder, presets.EDIT)
	}
	if mode.Has(presets.NEW) && ed.HasCreatingBuilder() {
		b.setFields(&ed.Creating().FieldsBuilder, presets.NEW)
	}
	if mode.Has(presets.DETAIL) && b.mb.HasDetailing() {
		b.setFields(&b.mb.Detailing().FieldsBuilder, presets.DETAIL)
	}
	if mode.Has(presets.LIST) {
		b.setFields(&b.mb.Listing().FieldsBuilder, presets.LIST)
	}
	return b
}
