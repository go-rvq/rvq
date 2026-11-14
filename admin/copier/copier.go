// Package copier makes copies of the records of a model (a ModelPlugin): the
// action Copy — an item of the menus of the listing and of the detail, for a
// record, and a bulk action, for several —, customizable by the model.
//
// A copy is the record read (Fetcher), copied whole into a new one, its
// fields of Reset and its relations (gorm's) zeroed, its title the one of a
// copy (CopyTitle), prepared by the model (Prepare) and made by the model's
// Creator: what a creation does, it does. The switches of its dialogs take
// more along (Expand), counted before; the records are copied in the order
// Expand gives, and the versions of a record in each language take one same
// new id (Session.NewID).
//
// It registers nothing of its own in the permissions: it asks
// MODEL:@create — once, where its items and its bulk action are decided and
// before it copies — and <ID>:@get of each distinct record it takes
// (ModelPermissioner.ObjectsDo), all before it copies anything.
package copier

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	// the events of the action, of the model
	EventOpen    = "copier_OpenEvent"
	EventDo      = "copier_DoEvent"
	EventPreview = "copier_PreviewEvent"

	// ParamIDs are the records a preview is of, joined by ",".
	ParamIDs = "copier_ids"

	defaultName   = "Copy"
	fieldTitle    = "CopyTitle"
	switchesField = "CopySwitches."
	choicesField  = "CopyChoices."
	previewPortal = "copier_Preview"
)

// PrepareFunc prepares the copy c, before it is made, in the session s of the
// records copied.
type PrepareFunc func(s *Session, c *Copy) error

// Copy is a record being copied.
type Copy struct {
	// From is the record read, To its copy, FromID the id of From.
	From, To any
	FromID   model.ID
	// Chosen: the record was chosen, not taken along by a switch.
	Chosen bool
}

// Switch is a switch of the dialogs of a copy: what Expand is told is on.
type Switch struct {
	Name        string
	Label, Hint func(ctx context.Context) string
	Default     bool
	// Count says the n records it takes along, in the dialog.
	Count func(ctx context.Context, n int) string
	// Choices, when not nil, are the values the user picks of it (chips,
	// shown while it is on; all of them picked at first): of the records
	// mids.
	Choices func(ctx *web.EventContext, mids model.IDSlice) []Choice
}

// Choice is a value of a switch to pick.
type Choice struct {
	Value, Label string
}

// State is what the dialog says of the switches: those on, and, of those of
// choices, the values picked.
type State struct {
	On     map[string]bool
	Picked map[string][]string
}

// Expansion is what a copy takes: the records, in the order they are copied
// (a parent before its children), and, by switch, how many it added.
type Expansion struct {
	IDs   model.IDSlice
	Added map[string]int
}

// ExpandFunc is what a copy of mids takes, with the switches in st.
type ExpandFunc func(ctx *web.EventContext, mids model.IDSlice, st *State) (*Expansion, error)

// Builder is the plugin of the copies of a model.
type Builder struct {
	db         *gorm.DB
	name       string
	title      func(ctx context.Context) string
	reset      []string
	titleField string
	copyTitle  func(ctx context.Context, title string) string
	prepare    []PrepareFunc
	switches   []*Switch
	expand     ExpandFunc
	done       func(ctx context.Context, n int) string
	configure  []func(*presets.BulkActionCofirmatorBuilder)

	mb *presets.ModelBuilder
}

// New is the plugin of the copies of a model of db: its action Copy, its
// fields ID, CreatedAt and UpdatedAt reset, its title the field Title.
func New(db *gorm.DB) *Builder {
	return &Builder{db: db, name: defaultName, reset: []string{"ID", "CreatedAt", "UpdatedAt"}, titleField: "Title"}
}

// Name names the action and titles it.
func (b *Builder) Name(name string, title func(ctx context.Context) string) *Builder {
	b.name, b.title = name, title
	return b
}

// Reset adds fields a copy has zeroed.
func (b *Builder) Reset(fields ...string) *Builder {
	b.reset = append(b.reset, fields...)
	return b
}

// TitleField is the field of the title (Title); "" a copy keeps its title.
func (b *Builder) TitleField(name string) *Builder {
	b.titleField = name
	return b
}

// CopyTitle is the title of the copy of a record titled title ("Copy of %s").
func (b *Builder) CopyTitle(f func(ctx context.Context, title string) string) *Builder {
	b.copyTitle = f
	return b
}

// Prepare adds f to what prepares a copy, in that order.
func (b *Builder) Prepare(f ...PrepareFunc) *Builder {
	b.prepare = append(b.prepare, f...)
	return b
}

// Switches adds switches to the dialogs, whose state Expand is told.
func (b *Builder) Switches(sw ...*Switch) *Builder {
	b.switches = append(b.switches, sw...)
	return b
}

// Expand is what a copy takes: the records chosen alone when nil.
func (b *Builder) Expand(f ExpandFunc) *Builder {
	b.expand = f
	return b
}

// Done is the message of n records copied.
func (b *Builder) Done(f func(ctx context.Context, n int) string) *Builder {
	b.done = f
	return b
}

// Configure adds f to what configures the confirmation of the bulk action
// (its table of the records chosen, …).
func (b *Builder) Configure(f ...func(*presets.BulkActionCofirmatorBuilder)) *Builder {
	b.configure = append(b.configure, f...)
	return b
}

// ActionName is the name of the action.
func (b *Builder) ActionName() string { return b.name }

// ActionTitle is the title of the action.
func (b *Builder) ActionTitle(ctx context.Context) string {
	if b.title != nil {
		return b.title(ctx)
	}
	return MustGetMessages(ctx).Copy
}

// ModelInstall installs the action in mb: its items of the menus of the
// listing and of the detail, its bulk action and its events.
func (b *Builder) ModelInstall(pb *presets.Builder, mb *presets.ModelBuilder) error {
	b.mb = mb
	pb.I18n().
		RegisterForModule(language.English, I18nCopierKey, Messages_en_US).
		RegisterForModule(language.BrazilianPortuguese, I18nCopierKey, Messages_pt_BR)

	mb.RegisterEventFunc(EventOpen, b.open)
	mb.RegisterEventFunc(EventDo, b.do)
	mb.RegisterEventFunc(EventPreview, b.preview)

	for _, rmb := range []*presets.RowMenuBuilder{mb.Listing().RowMenu(), mb.Detailing().RowMenu()} {
		rmb.RowMenuItem(b.name).ComponentFuncOf(b.menuItem).Title(b.ActionTitle)
	}

	// the bulk action asks MODEL:@create — where it is offered and where it
	// is done —, nothing of its own
	ba := presets.BulkActionConfirmator(mb.Listing().BulkAction(b.name).
		SetI18nLabel(b.ActionTitle).
		SetPerm(perm.PermVerifier().Func(func(v *perm.Verifier) *perm.Verifier { return v.Do(presets.PermCreate) }))).
		SetPreComp(func(ctx *web.EventContext) (h.HTMLComponent, error) {
			return b.switchesComponent(ctx, strings.Split(ctx.R.FormValue(presets.ParamSelectedIds), ",")), nil
		}).
		SetUpdator(func(ctx *web.EventContext, r *web.EventResponse, mids model.IDSlice) error {
			n, err := b.Copy(ctx, mids, "")
			if err == nil {
				presets.ShowMessage(r, b.doneMessage(ctx.Context(), n), "")
			}
			return err
		})
	for _, f := range b.configure {
		f(ba)
	}
	ba.Build()
	return nil
}

// menuItem is the item of a menu: none without MODEL:@create — asked once,
// for all its rows —; the read of the record is asked when it is copied.
func (b *Builder) menuItem(ctx *web.EventContext) presets.RecordMenuItemFunc {
	if b.mb.Info().Permissioner().ReqCreator(ctx.R).Denied() {
		return nil
	}
	return func(rctx *presets.RecordMenuItemContext) h.HTMLComponent {
		return VListItem(
			web.Slot(VIcon("mdi-content-copy")).Name("prepend"),
			VListItemTitle(h.Text(b.ActionTitle(rctx.Ctx.Context()))),
		).Attr("@click", web.Plaid().URL(b.mb.Info().ListingHrefCtx(rctx.Ctx)).EventFunc(EventOpen).
			Query(presets.ParamID, rctx.ID).Go())
	}
}

// open is the dialog of the copy of a record: the title of the copy and the
// switches.
func (b *Builder) open(ctx *web.EventContext) (r web.EventResponse, err error) {
	if b.mb.Info().Permissioner().ReqCreator(ctx.R).Denied() {
		return r, perm.PermissionDenied
	}
	var (
		paramID = ctx.R.FormValue(presets.ParamID)
		mid     model.ID
		msgr    = MustGetMessages(ctx.Context())
	)
	if mid, err = b.mb.ParseRecordID(paramID); err != nil {
		return
	}
	from := b.mb.NewModel()
	if err = b.mb.Fetcher(from, mid, ctx); err != nil {
		return
	}
	var title h.HTMLComponent
	if b.titleField != "" {
		t, _ := reflectutils.Get(from, b.titleField)
		title = VTextField().Label(msgr.CopyTitle).Variant(FieldVariantUnderlined).
			Attr(web.VField(fieldTitle, b.titleOf(ctx.Context(), fmt.Sprint(t)))...)
	}
	cb := &presets.ContentComponentBuilder{
		Context: ctx,
		Title:   b.ActionTitle(ctx.Context()),
		Body:    h.HTMLComponents{VContainer(title, b.switchesComponent(ctx, []string{paramID}))},
		BottomActions: h.HTMLComponents{
			VSpacer(),
			VBtn(presets.MustGetMessages(ctx.Context()).OK).Color(ColorPrimary).Variant(VariantFlat).
				Attr("@click", web.Plaid().URL(b.mb.Info().ListingHrefCtx(ctx)).EventFunc(EventDo).
					Query(presets.ParamID, paramID).Go()),
		},
		Overlay: &presets.ContentComponentBuilderOverlay{Mode: actions.Dialog},
	}
	b.mb.Builder().Dialog().SetScrollable(true).SetTargetPortal(actions.Dialog.PortalName()).
		Respond(ctx, &r, cb.BuildOverlay())
	return
}

// do copies the record of the dialog: MODEL:@create asked once, then the
// copy; refused, the dialog stays, saying why.
func (b *Builder) do(ctx *web.EventContext) (r web.EventResponse, err error) {
	if b.mb.Info().Permissioner().ReqCreator(ctx.R).Denied() {
		return r, perm.PermissionDenied
	}
	var mid model.ID
	if mid, err = b.mb.ParseRecordID(ctx.R.FormValue(presets.ParamID)); err != nil {
		return
	}
	n, err := b.Copy(ctx, model.IDSlice{mid}, ctx.R.FormValue(fieldTitle))
	if err != nil {
		presets.ShowMessage(&r, errorMessage(err), "error")
		return r, nil
	}
	// only said: the dialog closed, a listing on the page reloaded, nothing
	// opened
	presets.ShowMessage(&r, b.doneMessage(ctx.Context(), n), "")
	web.AppendRunScripts(&r, "closer.show = false; presetsListing?.loader?.go()")
	return
}

// Copy copies mids and what the switches of the request take along, in the
// order Expand gives — MODEL:@create asked by the caller —: <ID>:@get of each
// distinct record first, a denial copying nothing (a *web.ValidationErrors
// saying which). title, when not empty, is the title of the copy of the
// first record chosen (the dialog's). How many records it made.
func (b *Builder) Copy(ctx *web.EventContext, mids model.IDSlice, title string) (n int, err error) {
	var exp *Expansion
	if exp, err = b.expandOf(ctx, mids); err != nil {
		return
	}
	if denied := b.mb.Info().Permissioner().ObjectsDo(ctx.R, presets.PermGet, exp.IDs); denied != nil {
		obj := b.mb.NewModel()
		name := denied.String()
		if b.mb.Fetcher(obj, *denied, ctx) == nil {
			name = b.mb.RecordTitle(obj, ctx)
		}
		var vErr web.ValidationErrors
		vErr.GlobalError(fmt.Sprintf(MustGetMessages(ctx.Context()).ErrDenied, name))
		return 0, &vErr
	}

	chosen := map[string]bool{}
	for _, mid := range mids {
		chosen[mid.String()] = true
	}
	s := &Session{Ctx: ctx, DB: b.db, ids: map[string]any{}, maps: map[string]map[string]any{}}
	for _, mid := range exp.IDs {
		c := &Copy{From: b.mb.NewModel(), To: b.mb.NewModel(), FromID: mid, Chosen: chosen[mid.String()]}
		if err = b.mb.Fetcher(c.From, mid, ctx); err != nil {
			return n, fmt.Errorf("copier: %s: %w", mid, err)
		}
		reflect.ValueOf(c.To).Elem().Set(reflect.ValueOf(c.From).Elem())
		if err = b.zero(c.To); err != nil {
			return
		}
		oldID := mid.GetValue("ID")
		if id, ok := s.NewID(oldID); ok {
			if err = reflectutils.Set(c.To, "ID", id); err != nil {
				return
			}
		}
		if c.Chosen && b.titleField != "" {
			// the dialog's title is of its record; the others, a copy's
			t := title
			if t == "" || mid.String() != mids[0].String() {
				from, _ := reflectutils.Get(c.From, b.titleField)
				t = b.titleOf(ctx.Context(), fmt.Sprint(from))
			}
			if err = reflectutils.Set(c.To, b.titleField, t); err != nil {
				return
			}
		}
		for _, f := range b.prepare {
			if err = f(s, c); err != nil {
				return
			}
		}
		if err = b.mb.Creator(c.To, ctx); err != nil {
			return n, fmt.Errorf("copier: %s: %w", mid, err)
		}
		if id, gerr := reflectutils.Get(c.To, "ID"); gerr == nil && oldID != nil {
			s.ids[fmt.Sprint(oldID)] = id
		}
		n++
	}
	return
}

// expandOf is what a copy of mids takes, with the switches of the request.
func (b *Builder) expandOf(ctx *web.EventContext, mids model.IDSlice) (*Expansion, error) {
	if b.expand == nil {
		return &Expansion{IDs: mids, Added: map[string]int{}}, nil
	}
	return b.expand(ctx, mids, b.state(ctx, mids))
}

// state is what the request says of the switches: their fields, else their
// defaults (all the choices of a switch of choices picked).
func (b *Builder) state(ctx *web.EventContext, mids model.IDSlice) *State {
	st := &State{On: map[string]bool{}, Picked: map[string][]string{}}
	for _, sw := range b.switches {
		on := sw.Default
		if v, ok := ctx.R.Form[switchesField+sw.Name]; ok && len(v) > 0 {
			on = v[0] == "true"
		}
		st.On[sw.Name] = on
		if sw.Choices == nil {
			continue
		}
		if v, ok := ctx.R.Form[choicesField+sw.Name]; ok {
			st.Picked[sw.Name] = v
		} else {
			for _, c := range sw.Choices(ctx, mids) {
				st.Picked[sw.Name] = append(st.Picked[sw.Name], c.Value)
			}
		}
	}
	return st
}

// relations are the fields of the relations of a model (gorm's), by its
// type.
var relations sync.Map

// zero zeroes the fields of Reset and of the relations of obj: a copy holds
// its own, the model's Prepare makes them.
func (b *Builder) zero(obj any) error {
	v := reflect.ValueOf(obj).Elem()
	names, ok := relations.Load(v.Type())
	if !ok {
		s, err := schema.Parse(obj, &sync.Map{}, b.db.NamingStrategy)
		if err != nil {
			return fmt.Errorf("copier: %s: %w", v.Type(), err)
		}
		var ns []string
		for name := range s.Relationships.Relations {
			ns = append(ns, name)
		}
		relations.Store(v.Type(), ns)
		names = ns
	}
	for _, name := range append(append([]string(nil), names.([]string)...), b.reset...) {
		if f := v.FieldByName(name); f.IsValid() && f.CanSet() {
			f.Set(reflect.Zero(f.Type()))
		}
	}
	return nil
}

// titleOf is the title of the copy of a record titled title.
func (b *Builder) titleOf(ctx context.Context, title string) string {
	if b.copyTitle != nil {
		return b.copyTitle(ctx, title)
	}
	return fmt.Sprintf(MustGetMessages(ctx).CopyOf, title)
}

// doneMessage is the message of n records copied.
func (b *Builder) doneMessage(ctx context.Context, n int) string {
	if b.done != nil {
		return b.done(ctx, n)
	}
	return fmt.Sprintf(MustGetMessages(ctx).Done, n)
}

// errorMessage is what err says to the user: its global errors, of a
// validation error.
func errorMessage(err error) string {
	var vErr *web.ValidationErrors
	if errors.As(err, &vErr) && vErr.HaveGlobalErrors() {
		return strings.Join(vErr.GetGlobalErrors(), " ")
	}
	return err.Error()
}

// Session is a copy of records under way: the new ids of the records copied
// and what the model maps of its own (Map).
type Session struct {
	Ctx  *web.EventContext
	DB   *gorm.DB
	ids  map[string]any
	maps map[string]map[string]any
}

// NewID is the id of the copy of the record of the id old, when it was
// copied in this session (another of its languages).
func (s *Session) NewID(old any) (any, bool) {
	v, ok := s.ids[fmt.Sprint(old)]
	return v, ok
}

// Mapped is what the model mapped of old, in kind (Map).
func (s *Session) Mapped(kind string, old any) (any, bool) {
	v, ok := s.maps[kind][fmt.Sprint(old)]
	return v, ok
}

// Map maps old to new, in kind: what the copies of a record share (a record
// of its own, copied once).
func (s *Session) Map(kind string, old, new any) {
	if s.maps[kind] == nil {
		s.maps[kind] = map[string]any{}
	}
	s.maps[kind][fmt.Sprint(old)] = new
}

// previewOnChange is the script that recounts what the switches take along
// once a field of the dialog changed: "" without switches.
func (b *Builder) previewOnChange(ctx *web.EventContext, ids []string) string {
	if len(b.switches) == 0 || b.expand == nil {
		return ""
	}
	return "$nextTick(() => " + web.Plaid().URL(b.mb.Info().ListingHrefCtx(ctx)).EventFunc(EventPreview).
		Query(ParamIDs, strings.Join(ids, ",")).Go() + ")"
}

// switchesComponent are the switches of a copy of ids — the chips of the
// choices of one shown while it is on — and the count of what they take
// along: nil without switches.
func (b *Builder) switchesComponent(ctx *web.EventContext, ids []string) h.HTMLComponent {
	if len(b.switches) == 0 {
		return nil
	}
	mids := b.parseIDs(ids)
	onChange := b.previewOnChange(ctx, ids)
	var comps h.HTMLComponents
	for _, sw := range b.switches {
		s := VSwitch().Attr(web.VField(switchesField+sw.Name, sw.Default)...).
			Density(DensityCompact).Color(ColorPrimary).Attr("data-copy-switch", sw.Name)
		if sw.Label != nil {
			s.Label(sw.Label(ctx.Context()))
		}
		if sw.Hint != nil {
			s.Hint(sw.Hint(ctx.Context())).PersistentHint(true)
		} else {
			s.HideDetails(true)
		}
		if onChange != "" {
			s.Attr("@update:model-value", onChange)
		}
		comps = append(comps, s)
		if sw.Choices == nil {
			continue
		}
		var (
			chips  []h.HTMLComponent
			values = []string{}
		)
		for _, c := range sw.Choices(ctx, mids) {
			chips = append(chips, VChip(h.Text(c.Label)).Value(c.Value).Filter(true).Variant(VariantTonal))
			values = append(values, c.Value)
		}
		if len(chips) == 0 {
			continue
		}
		group := VChipGroup(chips...).Multiple(true).Column(true).Color(ColorPrimary).
			Attr(web.VField(choicesField+sw.Name, values)...).Attr("data-copy-choices", sw.Name)
		if onChange != "" {
			group.Attr("@update:model-value", onChange)
		}
		comps = append(comps, h.Div(group).Class("ml-12 mb-2").
			Attr("v-show", fmt.Sprintf("form[%q]", switchesField+sw.Name)))
	}
	comps = append(comps, web.Portal(b.previewComponent(ctx, mids)).Name(previewPortal))
	return comps
}

// previewComponent says what the switches take along of a copy of mids:
// nothing when they take nothing.
func (b *Builder) previewComponent(ctx *web.EventContext, mids model.IDSlice) h.HTMLComponent {
	if b.expand == nil || len(mids) == 0 {
		return nil
	}
	exp, err := b.expandOf(ctx, mids)
	if err != nil {
		return VAlert(h.Text(errorMessage(err))).Type("error").Variant(VariantTonal).Density(DensityCompact)
	}
	var says []string
	for _, sw := range b.switches {
		if n := exp.Added[sw.Name]; n > 0 && sw.Count != nil {
			says = append(says, sw.Count(ctx.Context(), n))
		}
	}
	if len(says) == 0 {
		return nil
	}
	return h.Div(VIcon("mdi-information-outline").Size(SizeSmall).Class("mr-1"), h.Text(strings.Join(says, " · "))).
		Class("text-caption text-medium-emphasis mt-2").Attr("data-copy-preview", true)
}

// preview is EventPreview: the count, of the records of ParamIDs — it asks
// nothing: the copy asks, when it is done.
func (b *Builder) preview(ctx *web.EventContext) (r web.EventResponse, err error) {
	r.UpdatePortal(previewPortal, b.previewComponent(ctx, b.parseIDs(strings.Split(ctx.R.FormValue(ParamIDs), ","))))
	return
}

// parseIDs are the records of ids: those it cannot read left out.
func (b *Builder) parseIDs(ids []string) (mids model.IDSlice) {
	for _, s := range ids {
		if s == "" {
			continue
		}
		if mid, err := b.mb.ParseRecordID(s); err == nil {
			mids = append(mids, mid)
		}
	}
	return
}
