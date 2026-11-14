package l10n

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/utils/db_utils"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

// What a localization takes: the records chosen and what its switches take
// along (ModelLocalizeOptions.Expand), localized in that order — by Do, the
// same from the dialog of a record and from a bulk action (LocalizeRecords).
// Its action is named, and titled, as the model says (SetAction).

const (
	// LocalizePreview answers the count of what the switches take along, in
	// the portal of the dialog (LocalizeSwitchesComponent).
	LocalizePreview = "l10n_LocalizePreviewEvent"
	// ParamLocalizeIDs are the records a preview is of, joined by ",".
	ParamLocalizeIDs = "l10n_ids"

	defaultActionName = "Localize"
	switchesField     = "LocalizeSwitches."
	previewPortal     = "l10n_LocalizePreview"
)

// OptionsOf are the localization options of mb (ModelWithOptions), nil of
// none.
func OptionsOf(mb *presets.ModelBuilder) *ModelLocalizeOptions {
	o, _ := mb.GetData(LocalizeOptions).(*ModelLocalizeOptions)
	return o
}

// ActionName is the name of the action that localizes the records of mb.
func ActionName(mb *presets.ModelBuilder) string {
	if o := OptionsOf(mb); o != nil && o.Name != "" {
		return o.Name
	}
	return defaultActionName
}

// ActionTitle is the title of the action that localizes the records of mb.
func ActionTitle(ctx context.Context, mb *presets.ModelBuilder) string {
	if o := OptionsOf(mb); o != nil && o.Title != nil {
		return o.Title(ctx)
	}
	return MustGetTranslation(ctx, "Localize")
}

// SetAction names the action that localizes the records of mb and titles
// it: its items of the menus of the listing and of the detail renamed.
func SetAction(mb *presets.ModelBuilder, name string, title func(ctx context.Context) string) {
	o := OptionsOf(mb)
	if o == nil {
		o = &ModelLocalizeOptions{}
		ModelWithOptions(mb, o)
	}
	old := ActionName(mb)
	o.Name, o.Title = name, title
	mb.Listing().RowMenu().Rename(old, name)
	if mb.HasDetailing() {
		mb.Detailing().RowMenu().Rename(old, name)
	}
}

// switchesOn are the switches of o on in the request: its fields, else their
// defaults.
func switchesOn(ctx *web.EventContext, o *ModelLocalizeOptions) map[string]bool {
	on := map[string]bool{}
	if o == nil {
		return on
	}
	for _, sw := range o.Switches {
		v := sw.Default
		if s, ok := ctx.R.Form[switchesField+sw.Name]; ok && len(s) > 0 {
			v, _ = strconv.ParseBool(s[0])
		}
		on[sw.Name] = v
	}
	return on
}

// expand is what a localization of mids to `to` takes (Expand), with the
// switches of the request.
func expand(ctx *web.EventContext, mb *presets.ModelBuilder, mids model.IDSlice, to []string) (*Expansion, error) {
	o := OptionsOf(mb)
	if o == nil || o.Expand == nil {
		return &Expansion{IDs: mids}, nil
	}
	return o.Expand(ctx, mids, to, switchesOn(ctx, o))
}

// LocalizeRecords localizes mids to `to` with what the switches of the
// request take along (Expand), in that order, by Do (DoLocalizeTo each, of
// none); every record taken along must be one the user may change, else
// nothing is localized and the error says which (a *web.ValidationErrors).
// How many records it localized.
func LocalizeRecords(db *gorm.DB, mb *presets.ModelBuilder, lb *Builder, ctx *web.EventContext, mids model.IDSlice, to []string) (n int, err error) {
	if len(to) == 0 {
		return
	}
	var exp *Expansion
	if exp, err = expand(ctx, mb, mids, to); err != nil {
		return
	}
	// what the switches take along: each distinct record asked once (the
	// versions of a record in each language are one)
	chosen := map[string]bool{}
	for _, mid := range mids {
		chosen[mid.String()] = true
	}
	var along model.IDSlice
	for _, mid := range exp.IDs {
		if !chosen[mid.String()] {
			along = append(along, mid)
		}
	}
	if denied := mb.Info().Permissioner().ObjectsDo(ctx.R, presets.PermUpdate, along); denied != nil {
		obj := mb.NewModel()
		if err = db_utils.ModelIdWhere(db.Session(&gorm.Session{}), mb.NewModel(), *denied).First(obj).Error; err != nil {
			return 0, fmt.Errorf("l10n: %s: %w", denied, err)
		}
		var vErr web.ValidationErrors
		vErr.GlobalError(fmt.Sprintf(MustGetMessages(ctx.Context()).ErrLocalizeDenied, mb.RecordTitle(obj, ctx)))
		return 0, &vErr
	}
	if o := OptionsOf(mb); o != nil && o.Do != nil {
		if err = o.Do(ctx, exp.IDs, to); err != nil {
			return
		}
		return len(exp.IDs), nil
	}
	for _, mid := range exp.IDs {
		cur, _ := mid.GetValue("LocaleCode").(string)
		var done []string
		if done, err = DoLocalizeTo(db, mb, lb, ctx, mid, slices.DeleteFunc(slices.Clone(to), func(l string) bool { return l == cur })); err != nil {
			return
		}
		if len(done) > 0 {
			n++
		}
	}
	return
}

// LocalizedMessage is the message of n records localized of mb (Done).
func LocalizedMessage(ctx context.Context, mb *presets.ModelBuilder, n int) string {
	if o := OptionsOf(mb); o != nil && o.Done != nil {
		return o.Done(ctx, n)
	}
	return MustGetTranslation(ctx, "SuccessfullyLocalized")
}

// ErrorMessage is what err says to the user: its global errors, of a
// validation error.
func ErrorMessage(err error) string {
	var vErr *web.ValidationErrors
	if errors.As(err, &vErr) && vErr.HaveGlobalErrors() {
		return strings.Join(vErr.GetGlobalErrors(), " ")
	}
	return err.Error()
}

// PreviewOnChange is the script, of a field of the dialog of a localization
// of the records ids, that recounts what its switches take along once the
// field changed (LocalizePreview): "" when mb has no switches.
func PreviewOnChange(ctx *web.EventContext, mb *presets.ModelBuilder, ids []string) string {
	if o := OptionsOf(mb); o == nil || len(o.Switches) == 0 || o.Expand == nil {
		return ""
	}
	return "$nextTick(() => " + web.Plaid().URL(mb.Info().ListingHrefCtx(ctx)).EventFunc(LocalizePreview).
		Query(ParamLocalizeIDs, strings.Join(ids, ",")).Go() + ")"
}

// LocalizeSwitchesComponent are the switches of a localization of the
// records ids and the count of what they take along: nil when mb has none.
func LocalizeSwitchesComponent(ctx *web.EventContext, mb *presets.ModelBuilder, ids []string) h.HTMLComponent {
	o := OptionsOf(mb)
	if o == nil || len(o.Switches) == 0 {
		return nil
	}
	onChange := PreviewOnChange(ctx, mb, ids)
	var comps h.HTMLComponents
	for _, sw := range o.Switches {
		s := VSwitch().Attr(web.VField(switchesField+sw.Name, sw.Default)...).
			Density(DensityCompact).Color(ColorPrimary).
			Attr("data-localize-switch", sw.Name)
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
	}
	comps = append(comps, web.Portal(previewComponent(ctx, mb, parseIDs(mb, ids), nil)).Name(previewPortal))
	return comps
}

// previewComponent says what the switches take along, of a localization of
// mids to `to`: nothing when they take nothing.
func previewComponent(ctx *web.EventContext, mb *presets.ModelBuilder, mids model.IDSlice, to []string) h.HTMLComponent {
	o := OptionsOf(mb)
	if o == nil || o.Expand == nil || len(mids) == 0 {
		return nil
	}
	exp, err := expand(ctx, mb, mids, to)
	if err != nil {
		return VAlert(h.Text(ErrorMessage(err))).Type("error").Variant(VariantTonal).Density(DensityCompact)
	}
	var says []string
	for _, sw := range o.Switches {
		if n := exp.Added[sw.Name]; n > 0 && sw.Count != nil {
			says = append(says, sw.Count(ctx.Context(), n))
		}
	}
	if len(says) == 0 {
		return nil
	}
	return h.Div(VIcon("mdi-information-outline").Size(SizeSmall).Class("mr-1"), h.Text(strings.Join(says, " · "))).
		Class("text-caption text-medium-emphasis mt-2").Attr("data-localize-preview", true)
}

// parseIDs are the records of ids, of mb: those it cannot read left out.
func parseIDs(mb *presets.ModelBuilder, ids []string) (mids model.IDSlice) {
	for _, s := range ids {
		if s == "" {
			continue
		}
		if mid, err := mb.ParseRecordID(s); err == nil {
			mids = append(mids, mid)
		}
	}
	return
}

// localizePreview is LocalizePreview.
func localizePreview(mb *presets.ModelBuilder) web.EventFunc {
	return func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var form doLocalize
		ctx.UnmarshalForm(&form)
		mids := parseIDs(mb, strings.Split(ctx.R.FormValue(ParamLocalizeIDs), ","))
		if mb.Info().Permissioner().ReqLister(ctx.R).Denied() {
			return
		}
		r.UpdatePortal(previewPortal, previewComponent(ctx, mb, mids, form.LocalizeTo))
		return
	}
}
