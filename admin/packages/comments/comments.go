// Package comments draws the comments of a record — a conversation: a tree,
// each answer under the one it answers, who wrote each, when, from where,
// its text (plain, or rich: HTML, sanitized), its attachments, its chips (a
// status, edited, how its mail went) — and the boxes that write them: a new
// one, an answer to one, an edit of one's own. What the comments are, who
// writes and what each change does are the app's (Builder.Messages, CanWrite,
// OnAdd, OnEdit, OnDelete, OnAction): the comments of a page, the answers to
// a notification. A change draws them again.
package comments

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/filecards"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/microcosm-cc/bluemonday"
)

// FileMaxSize is the largest attachment taken, in bytes (Attachments).
const FileMaxSize = 10 << 20

var (
	// ErrFileTooLarge is an attachment larger than FileMaxSize.
	ErrFileTooLarge = errors.New("comments: attachment too large")
	// ErrNotAllowed is a change the user may not do.
	ErrNotAllowed = errors.New("comments: not allowed")
)

// Attachment is a file sent with a comment.
type Attachment struct {
	Name        string
	ContentType string
	Data        []byte
}

// File is an attachment of a comment as it is shown: its name, its size,
// where it is (its link).
type File struct {
	Name string
	Size int64
	// Type is its media type: an image is seen in its card (filecards)
	Type string
	URL  string
}

// Action is a button of a comment: OnAction does it (approve, hide, send
// again).
type Action struct {
	Name, Label, Icon, Color string
}

// Comment is a comment of the tree.
type Comment struct {
	ID string
	// Author is who wrote it (empty: a removed user); AuthorHint, more of
	// them (their login, their email): its tooltip
	Author, AuthorHint string
	// Avatar, when given, is their picture (default: their initials)
	Avatar h.HTMLComponent
	Time   time.Time
	// From is where they were (a city, an address)
	From string
	// Body is its text: HTML when HTML (sanitized), plain otherwise
	Body string
	HTML bool
	// Title, when given, is over its text (a subject)
	Title string
	Files []File
	// Chips are after its time: its status, edited, how its mail went
	Chips []h.HTMLComponent
	// CanEdit, CanDelete: the user of the request changes it (theirs)
	CanEdit, CanDelete bool
	// Actions are its buttons of the app (OnAction)
	Actions []Action
	// Replies are the answers to it, the oldest first
	Replies []*Comment
}

// Write is a comment written: a new one (OnAdd: an answer to ParentID, when
// given) or an edit (OnEdit: of ID).
type Write struct {
	Key, ID, ParentID string
	// Body is its text (HTML with RichText)
	Body  string
	Files []Attachment
	// Options are the boxes of a new one (Builder.Options), by name
	Options map[string]bool
}

// Option is a box of the writing of a new comment ("Internal"): its value
// in Write.Options.
type Option struct {
	Name    string
	Label   func(ctx context.Context) string
	Default bool
}

// Builder is a kind of comments: their event, what they are, what changing
// them does.
type Builder struct {
	// EventName is the event of its changes (op: load, add, edit, delete,
	// an Action's name).
	EventName string
	// URL is where the events go (nil: the page's).
	URL func(ctx *web.EventContext) string
	// Comments are the tree of the record key.
	Comments func(ctx *web.EventContext, key string) ([]*Comment, error)
	// CanRead says whether the user sees them (nil: who reaches the event);
	// CanWrite, whether they write and answer (nil: no box).
	CanRead  func(ctx *web.EventContext, key string) bool
	CanWrite func(ctx *web.EventContext, key string) bool
	OnAdd    func(ctx *web.EventContext, w *Write) error
	OnEdit   func(ctx *web.EventContext, w *Write) error
	OnDelete func(ctx *web.EventContext, key, id string) error
	OnAction func(ctx *web.EventContext, key, id, action string) error
	// Options are the boxes of a new comment.
	Options []Option
	// Answers: each comment has its answer (the tree); without, a list.
	Answers bool
	// RichText: written in the rich text editor (HTML); Plugins, its.
	RichText bool
	Plugins  []string
	// Attachments says whether a comment of the record key takes files (nil:
	// none takes).
	Attachments func(ctx *web.EventContext, key string) bool
	// MaxLen is the longest text, in characters (0: no limit).
	MaxLen int
	// ValidateFile, when given, checks each file of a comment of the record
	// key (its size: MaxFileSize): an error refuses the comment, said.
	ValidateFile func(ctx *web.EventContext, key string, f *Attachment) error
	// ValidateMessage, when given, checks a comment written whole (the size
	// of its files together: MaxTotalSize): an error refuses it, said.
	ValidateMessage func(ctx *web.EventContext, w *Write) error
}

// MaxFileSize is a ValidateFile: a file larger than max of the record (0:
// no limit of its own) refused.
func MaxFileSize(max func(ctx *web.EventContext, key string) int64) func(ctx *web.EventContext, key string, f *Attachment) error {
	return func(ctx *web.EventContext, key string, f *Attachment) error {
		if n := max(ctx, key); n > 0 && int64(len(f.Data)) > n {
			return fmt.Errorf(GetMessages(ctx.Context()).FileTooLarge, f.Name, HumanSize(int64(len(f.Data))), HumanSize(n))
		}
		return nil
	}
}

// MaxTotalSize is a ValidateMessage: the files of a comment larger together
// than max of the record (0: no limit of its own) refused.
func MaxTotalSize(max func(ctx *web.EventContext, key string) int64) func(ctx *web.EventContext, w *Write) error {
	return func(ctx *web.EventContext, w *Write) error {
		var total int64
		for _, f := range w.Files {
			total += int64(len(f.Data))
		}
		if n := max(ctx, w.Key); n > 0 && total > n {
			return fmt.Errorf(GetMessages(ctx.Context()).MessageTooLarge, HumanSize(total), HumanSize(n))
		}
		return nil
	}
}

// Validate is the checks of the comment w: each of its files
// (ValidateFile), it whole (ValidateMessage).
func (b *Builder) Validate(ctx *web.EventContext, w *Write) error {
	if b.ValidateFile != nil {
		for i := range w.Files {
			if err := b.ValidateFile(ctx, w.Key, &w.Files[i]); err != nil {
				return err
			}
		}
	}
	if b.ValidateMessage != nil {
		return b.ValidateMessage(ctx, w)
	}
	return nil
}

// New is a kind of comments of p, of the event name (unique in the admin).
func New(p *presets.Builder, eventName string) *Builder {
	ConfigureMessages(p.I18n())
	return &Builder{EventName: eventName, Answers: true}
}

// Register makes the events of the comments handled by hub (a model's: its
// pages; the web builder's: every page).
func (b *Builder) Register(hub interface {
	RegisterEventFunc(id string, f web.EventFunc) string
}) *Builder {
	hub.RegisterEventFunc(b.EventName, b.event)
	return b
}

func (b *Builder) portal(key string) string {
	return b.EventName + "_" + strings.NewReplacer("-", "_", ".", "_", ":", "_").Replace(key)
}

// Count is how many comments the tree has, at every depth.
func Count(cs []*Comment) (n int) {
	for _, c := range cs {
		n += 1 + Count(c.Replies)
	}
	return
}

// CountChip is the chip of how many comments the record key has, drawn
// again with them (Panel).
func (b *Builder) CountChip(ctx *web.EventContext, key string) h.HTMLComponent {
	cs, _ := b.Comments(ctx, key)
	n := Count(cs)
	return v.VChip(h.Text(fmt.Sprint(n))).Size(v.SizeXSmall).Variant(v.VariantTonal).Class("ms-2").
		Attr("data-comment-count", n)
}

// Panel is the comments of the record key folded in a panel of label: how
// many in its title, loaded when it opens.
func (b *Builder) Panel(ctx *web.EventContext, key, label string) h.HTMLComponent {
	loader := web.POST().EventFunc(b.EventName).Query("key", key).Query("op", "load")
	if b.URL != nil {
		loader.URL(b.URL(ctx))
	}
	return v.VExpansionPanels(
		v.VExpansionPanel(
			v.VExpansionPanelTitle(
				v.VIcon("mdi-comment-text-multiple-outline").Size(v.SizeSmall).Class("me-2"),
				h.Text(label),
				web.Portal(b.CountChip(ctx, key)).Name(b.portal(key)+"_count").Inline(),
			),
			v.VExpansionPanelText(
				// empty: a portal with content of its own does not load
				web.Portal().Name(b.portal(key)).Loader(loader),
			),
		).Attr("data-comment-panel", key),
	).Class("mb-4")
}

// Inline is the comments of the record key, drawn now (not folded).
func (b *Builder) Inline(ctx *web.EventContext, key string) h.HTMLComponent {
	return web.Portal(b.Component(ctx, key)).Name(b.portal(key))
}

// Call is the script of the event op on the comment id of the record key —
// with the text, the files and the options of the box (add, edit) —: for a
// button of the app (a chip that sends again).
func (b *Builder) Call(ctx *web.EventContext, key, op, id string) string {
	q := "plaid()"
	if b.URL != nil {
		q += ".url(" + h.JSONString(b.URL(ctx)) + ")"
	}
	q += fmt.Sprintf(`.eventFunc(%s).query("key", %s).query("op", %s).query("id", %s)`,
		h.JSONString(b.EventName), h.JSONString(key), h.JSONString(op), h.JSONString(id))
	switch op {
	case "add":
		opts := ""
		for _, o := range b.Options {
			opts += fmt.Sprintf(", %s: String(locals.opts[%s])", h.JSONString(o.Name), h.JSONString(o.Name))
		}
		q += `.form({body: locals.text, files: locals.files` + opts + `})`
	case "edit":
		q += `.form({body: locals.text, files: locals.files})`
	}
	return `locals.sending = true; ` + q + `.go().finally(() => { locals.sending = false })`
}

// policy is what the HTML of a comment keeps (bluemonday's of the content
// of users: no scripts, no styles).
var policy = bluemonday.UGCPolicy()

// SanitizeHTML is html as a comment shows it.
func SanitizeHTML(html string) string { return policy.Sanitize(html) }

var strict = bluemonday.StrictPolicy()

// PlainText is the text of html, its tags dropped, its spaces folded.
func PlainText(html string) string {
	return strings.Join(strings.Fields(strict.Sanitize(strings.NewReplacer("</p>", " </p>", "<br", " <br").Replace(html))), " ")
}

// Component is the tree of the comments of the record key, with the box that
// writes one.
func (b *Builder) Component(ctx *web.EventContext, key string) h.HTMLComponent {
	m := GetMessages(ctx.Context())
	cs, err := b.Comments(ctx, key)
	if err != nil {
		return v.VAlert(h.Text(err.Error())).Type("error").Density(v.DensityCompact)
	}
	canWrite := b.CanWrite != nil && b.CanWrite(ctx, key)
	loc := presets.TimeLocationOf(ctx)

	var item func(c *Comment) h.HTMLComponent
	item = func(c *Comment) h.HTMLComponent {
		id := c.ID
		author := c.Author
		if author == "" {
			author = m.RemovedUser
		}
		avatar := c.Avatar
		if avatar == nil {
			avatar = v.VAvatar(h.Text(Initials(author))).Size(32).Color("primary").Variant(v.VariantTonal).
				Class("text-caption font-weight-bold")
		}
		btn := func(label, icon string) *v.VBtnBuilder {
			return v.VBtn(label).Variant(v.VariantText).Size(v.SizeSmall).PrependIcon(icon).Class("text-none px-2")
		}
		var actions h.HTMLComponents
		if canWrite {
			if b.Answers && b.OnAdd != nil {
				actions = append(actions, btn(m.Reply, "mdi-reply-outline").
					Attr("@click", fmt.Sprintf(`locals.reply = %q; locals.edit = ""; locals.text = ""; locals.files = []`, id)).
					Attr("data-comment-reply", id))
			}
			if c.CanEdit && b.OnEdit != nil {
				actions = append(actions, btn(m.Edit, "mdi-pencil-outline").
					Attr("@click", fmt.Sprintf(`locals.edit = %q; locals.reply = ""; locals.files = []; locals.text = %s`, id, h.JSONString(c.Body))).
					Attr("data-comment-edit", id))
			}
			if c.CanDelete && b.OnDelete != nil {
				actions = append(actions,
					btn(m.Delete, "mdi-trash-can-outline").Color("error").
						Attr("v-if", fmt.Sprintf(`locals.del !== %q`, id)).
						Attr("@click", fmt.Sprintf(`locals.del = %q`, id)).
						Attr("data-comment-delete", id),
					h.Tag("span").Children(
						h.Text(m.DeleteConfirm),
						btn(m.DeleteYes, "mdi-check").Color("error").
							Attr(":loading", "locals.sending").
							Attr("@click", b.Call(ctx, key, "delete", id)).
							Attr("data-comment-delete-confirm", id),
						btn(m.Cancel, "mdi-close").Attr("@click", `locals.del = ""`),
					).Attr("v-else", true).Class("text-caption text-error d-inline-flex align-center"),
				)
			}
			for _, a := range c.Actions {
				actions = append(actions, btn(a.Label, a.Icon).Color(a.Color).
					Attr("@click", b.Call(ctx, key, a.Name, id)).Attr("data-comment-"+a.Name, id))
			}
		}
		var replies h.HTMLComponents
		for _, r := range c.Replies {
			replies = append(replies, item(r))
		}
		var children h.HTMLComponent
		if len(replies) > 0 {
			// the answers under it, a line at their left
			children = h.Div(replies...).Class("rvq-comment__replies ms-4 ps-4 mt-3 d-flex flex-column ga-3").
				Style("border-inline-start: 2px solid rgba(var(--v-border-color), var(--v-border-opacity))")
		}
		var body h.HTMLComponent
		if c.HTML {
			body = h.Div(h.RawHTML(SanitizeHTML(c.Body))).Class(web.StdContentClass + " text-body-2 mt-1")
		} else {
			// what a user wrote is never a template (v-pre)
			body = h.Div(h.Text(c.Body)).Attr("v-pre", true).Class("text-body-2 mt-1").
				Style("white-space:pre-line;overflow-wrap:anywhere")
		}
		var title h.HTMLComponent
		if c.Title != "" {
			title = h.Div(h.Text(c.Title)).Class("text-subtitle-2 mt-1")
		}
		authorSpan := h.Span(author).Class("font-weight-bold text-body-2")
		if c.AuthorHint != "" {
			authorSpan.Attr("title", c.AuthorHint)
		}
		var from h.HTMLComponent
		if c.From != "" {
			from = h.Span(c.From).Class("text-caption text-medium-emphasis d-none d-sm-inline")
		}
		head := h.HTMLComponents{
			authorSpan,
			h.Span(c.Time.In(loc).Format("02/01/2006 15:04")).Class("text-caption text-medium-emphasis").Attr("title", c.From),
			from,
		}
		for _, ch := range c.Chips {
			head = append(head, ch)
		}
		return h.Div(
			h.Div(
				avatar,
				h.Div(
					h.Div(head...).Class("d-flex flex-wrap align-center ga-2"),
					h.Div(title, body, b.files(c.Files)).Attr("v-show", fmt.Sprintf(`locals.edit !== %q`, id)),
					h.Div(b.box(ctx, m, key, "", b.Call(ctx, key, "edit", id), true, false)).
						Attr("v-if", fmt.Sprintf(`locals.edit === %q`, id)),
					h.Div(actions...).Class("d-flex flex-wrap align-center ms-n2 mt-1"),
					h.Div(b.box(ctx, m, key, m.ReplyPlaceholder, b.Call(ctx, key, "add", id), true, true)).
						Attr("v-if", fmt.Sprintf(`locals.reply === %q`, id)),
				).Class("flex-grow-1").Style("min-width:0"),
			).Class("d-flex ga-3"),
			children,
		).Class("rvq-comment").Attr("id", "comment-"+id).Attr("data-comment", id)
	}

	var items h.HTMLComponents
	for _, c := range cs {
		items = append(items, item(c))
	}
	var list h.HTMLComponent
	if len(items) > 0 {
		list = h.Div(items...).Class("d-flex flex-column ga-4")
	} else {
		list = h.Div(h.Text(m.None)).Class("text-body-2 text-medium-emphasis")
	}
	var top h.HTMLComponent
	if canWrite && b.OnAdd != nil {
		top = h.Div(b.box(ctx, m, key, m.Placeholder, b.Call(ctx, key, "add", ""), false, true)).
			Attr("v-if", `!locals.reply && !locals.edit`).Class("mb-4")
	}
	opts := map[string]bool{}
	for _, o := range b.Options {
		opts[o.Name] = o.Default
	}
	return web.Scope(
		h.Div(top, list).Class("rvq-comment-thread").Attr("data-comment-thread", key),
	).LocalsInit(fmt.Sprintf(`{reply: "", edit: "", del: "", text: "", files: [], sending: false, opts: %s}`, h.JSONString(opts))).
		Slot("{ locals }")
}

// files are the attachments of a comment, as cards (filecards).
func (b *Builder) files(fs []File) h.HTMLComponent {
	if len(fs) == 0 {
		return nil
	}
	cards := make([]filecards.File, len(fs))
	for i, f := range fs {
		cards[i] = filecards.File{Key: f.Name, Name: f.Name, Size: f.Size, Type: f.Type, URL: f.URL}
	}
	return h.Div(filecards.Cards(cards, "data-comment-file")).Class("mt-2")
}

// box writes a comment: its text, its files, the options of a new one (add),
// send — and cancel (cancel).
func (b *Builder) box(ctx *web.EventContext, m *Messages, key, placeholder, send string, cancel, add bool) h.HTMLComponent {
	withFiles := b.withFiles(ctx, key)
	var text h.HTMLComponent
	if b.RichText {
		ed := vx.VXTipTapEditor().Output("html").Attr("v-model", "locals.text").Attr("data-comment-text", true)
		if placeholder != "" {
			ed.Label(placeholder)
		}
		if len(b.Plugins) > 0 {
			ed.Plugins(b.Plugins...)
		}
		text = ed
	} else {
		ta := v.VTextarea().Attr("v-model", "locals.text").Placeholder(placeholder).Rows(2).AutoGrow(true).
			Variant(v.FieldVariantOutlined).Density(v.DensityCompact).HideDetails(true).Attr("data-comment-text", true)
		if b.MaxLen > 0 {
			ta.Counter(b.MaxLen).Attr(":maxlength", fmt.Sprint(b.MaxLen))
		}
		text = ta
	}
	var files h.HTMLComponent
	if withFiles {
		files = h.Div(filecards.Drop("locals.files", m.DropFiles, m.EditFile, m.RemoveFile)).Class("mt-2").Attr("data-comment-files", true)
	}
	var opts h.HTMLComponents
	if add {
		for _, o := range b.Options {
			opts = append(opts, v.VCheckbox().Attr("v-model", fmt.Sprintf("locals.opts[%s]", h.JSONString(o.Name))).Label(o.Label(ctx.Context())).
				Density(v.DensityCompact).HideDetails(true).Class("me-auto").Attr("data-comment-"+o.Name, true))
		}
	}
	var cancelBtn h.HTMLComponent
	if cancel {
		cancelBtn = v.VBtn(m.Cancel).Variant(v.VariantText).Size(v.SizeSmall).
			Attr("@click", `locals.reply = ""; locals.edit = ""; locals.text = ""; locals.files = []`)
	}
	empty := "!locals.text || !locals.text.trim()"
	if b.RichText {
		empty = "!locals.text || !locals.text.replace(/<[^>]*>/g, '').trim()"
	}
	if withFiles {
		empty = "(" + empty + ") && !(locals.files && locals.files.length)"
	}
	return h.Div(
		text,
		files,
		h.Div(
			append(opts,
				cancelBtn,
				v.VBtn(m.Send).Color(v.ColorPrimary).Variant(v.VariantFlat).Size(v.SizeSmall).
					PrependIcon("mdi-send").
					Attr(":loading", "locals.sending").
					Attr(":disabled", empty).
					Attr("@click", send).
					Attr("data-comment-send", true),
			)...,
		).Class("d-flex justify-end align-center ga-2 mt-2"),
	).Class("mt-2")
}

// event is a change of the comments of a record: op "load" (the tree), "add"
// (an answer to id, when given), "edit", "delete" (of id), an Action's;
// then the tree again, and its count.
func (b *Builder) event(ctx *web.EventContext) (r web.EventResponse, err error) {
	key := ctx.R.FormValue("key")
	op := ctx.R.FormValue("op")
	if key == "" || (b.CanRead != nil && !b.CanRead(ctx, key)) {
		return r, ErrNotAllowed
	}
	if op == "load" {
		r.Body = b.Component(ctx, key)
		return
	}
	if b.CanWrite == nil || !b.CanWrite(ctx, key) {
		return r, ErrNotAllowed
	}
	m := GetMessages(ctx.Context())
	id := ctx.R.FormValue("id")
	w := &Write{Key: key, Body: strings.TrimSpace(ctx.R.FormValue("body"))}
	switch op {
	case "add", "edit":
		if b.withFiles(ctx, key) {
			if w.Files, err = Attachments(ctx.R, "files"); err != nil {
				return
			}
		}
		text := w.Body
		if b.RichText {
			text = PlainText(w.Body)
		}
		if text == "" && len(w.Files) == 0 {
			return r, errors.New(m.Empty)
		}
		if b.MaxLen > 0 && !b.RichText && utf8.RuneCountInString(w.Body) > b.MaxLen {
			w.Body = string([]rune(w.Body)[:b.MaxLen])
		}
		if err = b.Validate(ctx, w); err != nil {
			return
		}
	}
	switch op {
	case "add":
		if b.OnAdd == nil {
			return r, ErrNotAllowed
		}
		w.ParentID = id
		w.Options = map[string]bool{}
		for _, o := range b.Options {
			w.Options[o.Name] = o.Default
			if s := ctx.R.FormValue(o.Name); s != "" {
				w.Options[o.Name] = s == "true"
			}
		}
		err = b.OnAdd(ctx, w)
	case "edit":
		if b.OnEdit == nil || id == "" {
			return r, ErrNotAllowed
		}
		w.ID = id
		err = b.OnEdit(ctx, w)
	case "delete":
		if b.OnDelete == nil || id == "" {
			return r, ErrNotAllowed
		}
		err = b.OnDelete(ctx, key, id)
	default:
		if b.OnAction == nil || id == "" {
			return r, fmt.Errorf("comments: op %q", op)
		}
		err = b.OnAction(ctx, key, id, op)
	}
	if err != nil {
		return
	}
	r.UpdatePortals = append(r.UpdatePortals,
		&web.PortalUpdate{Name: b.portal(key), Body: b.Component(ctx, key)},
		&web.PortalUpdate{Name: b.portal(key) + "_count", Body: b.CountChip(ctx, key)})
	return
}

// withFiles says whether a comment of the record key takes files.
func (b *Builder) withFiles(ctx *web.EventContext, key string) bool {
	return b.Attachments != nil && b.Attachments(ctx, key)
}

// Attachments are the files of the form (multipart) of key (key, key[0]…):
// each at most FileMaxSize.
func Attachments(r *http.Request, key string) (out []Attachment, err error) {
	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(32 << 20); err != nil && !errors.Is(err, http.ErrNotMultipart) {
			return nil, err
		}
	}
	if r.MultipartForm == nil {
		return nil, nil
	}
	for k, fhs := range r.MultipartForm.File {
		if k != key && !strings.HasPrefix(k, key+"[") {
			continue
		}
		for _, fh := range fhs {
			a, err := attachmentOf(fh)
			if err != nil {
				return nil, err
			}
			out = append(out, a)
		}
	}
	return
}

func attachmentOf(fh *multipart.FileHeader) (Attachment, error) {
	if fh.Size > FileMaxSize {
		return Attachment{}, fmt.Errorf("%w: %s", ErrFileTooLarge, fh.Filename)
	}
	f, err := fh.Open()
	if err != nil {
		return Attachment{}, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, FileMaxSize+1))
	if err != nil {
		return Attachment{}, err
	}
	if len(data) > FileMaxSize {
		return Attachment{}, fmt.Errorf("%w: %s", ErrFileTooLarge, fh.Filename)
	}
	return Attachment{Name: filepath.Base(fh.Filename), ContentType: fh.Header.Get("Content-Type"), Data: data}, nil
}

// HumanSize is n bytes as people read it.
func HumanSize(n int64) string { return filecards.FormatSize(n) }

// Initials are the first letters of the first and the last words of name.
func Initials(name string) string {
	f := strings.Fields(name)
	if len(f) == 0 {
		return "?"
	}
	first, _ := utf8.DecodeRuneInString(f[0])
	s := string(first)
	if len(f) > 1 {
		last, _ := utf8.DecodeRuneInString(f[len(f)-1])
		s += string(last)
	}
	return strings.ToUpper(s)
}
