package l10n

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"slices"
	"time"

	"github.com/go-rvq/rvq/admin/activity"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/datafield"
	"github.com/sunfmin/reflectutils"
	. "github.com/theplant/htmlgo"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

var IncorrectLocaleErr = errors.New("incorrect locale")

type Builder struct {
	db *gorm.DB
	ab *activity.Builder
	// models                               []*presets.ModelBuilder
	locales                                          Locales
	getSupportLocaleCodesFromRequestFunc             func(R *http.Request) []string
	cookieName                                       string
	queryName                                        string
	defaultLocaleCode                                string
	disableDeletionForDefaultInternationalizedRecord bool
}

func (b *Builder) DisableDeletionForDefaultInternationalizedRecord() bool {
	return b.disableDeletionForDefaultInternationalizedRecord
}

func (b *Builder) SetDisableDeletionForDefaultInternationalizedRecord(disableDeletionForDefaultInternationalizedRecord bool) *Builder {
	b.disableDeletionForDefaultInternationalizedRecord = disableDeletionForDefaultInternationalizedRecord
	return b
}

type LocaleInfo struct {
	code  string
	path  string
	label string
	datafield.DataField[*LocaleInfo]
}

func (l *LocaleInfo) String() string {
	return l.code
}

func (l *LocaleInfo) Code() string {
	return l.code
}

func (l *LocaleInfo) Path() string {
	return l.path
}

func (l *LocaleInfo) Label() string {
	return l.label
}

func (l *LocaleInfo) Export() *PublicLocaleInfo {
	return &PublicLocaleInfo{
		Code:  l.code,
		Label: l.label,
		Path:  l.path,
	}
}

type PublicLocaleInfo struct {
	Code  string
	Path  string
	Label string
}

type Locales []*LocaleInfo

func (s Locales) Exported() (r []*PublicLocaleInfo) {
	r = make([]*PublicLocaleInfo, len(s))
	for i, info := range s {
		r[i] = info.Export()
	}
	return
}

func New(db *gorm.DB) *Builder {
	b := &Builder{
		db:         db,
		cookieName: "locale",
		queryName:  "locale",
	}
	return b
}

func (b *Builder) DefaultLocaleCode(localeCode string) *Builder {
	b.defaultLocaleCode = localeCode
	return b
}

func (b *Builder) GetDefaultLocaleCode() string {
	return b.defaultLocaleCode
}

func (b *Builder) GetDefaultLocale() *LocaleInfo {
	return b.GetLocale(b.defaultLocaleCode)
}

func (b *Builder) SupportedLocales() Locales {
	return b.locales
}

func (b *Builder) IsTurnedOn() bool {
	return len(b.GetSupportLocaleCodes()) > 0
}

func (b *Builder) GetCookieName() string {
	return b.cookieName
}

func (b *Builder) GetQueryName() string {
	return b.queryName
}

func (b *Builder) Activity(v *activity.Builder) (r *Builder) {
	b.ab = v
	return b
}

func (b *Builder) GetAtivity() *activity.Builder {
	return b.ab
}

func (b *Builder) RegisterLocale(localeCode, localePath, localeLabel string) (r *LocaleInfo) {
	if slices.ContainsFunc(b.locales, func(l *LocaleInfo) bool {
		return l.code == localeCode
	}) {
		return b.GetLocale(localeCode)
	}

	l := datafield.New(&LocaleInfo{
		code:  localeCode,
		path:  path.Join("/", localePath),
		label: localeLabel,
	})
	b.locales = append(b.locales, l)
	return l
}

func (b *Builder) UnRegisterLocales(localeCode ...string) (r *Builder) {
	var newLocales Locales

loop:
	for _, locale := range b.locales {
		for _, s := range localeCode {
			if locale.code == s {
				continue loop
			}
		}

		newLocales = append(newLocales, locale)
	}

	b.locales = newLocales
	return b
}

func (b *Builder) GetLocale(localeCode string) *LocaleInfo {
	for _, l := range b.locales {
		if l.code == localeCode {
			return l
		}
	}
	return nil
}

type contextKeyType int

const contextKey contextKeyType = iota

func (b *Builder) ContextValueProvider(in context.Context) context.Context {
	return context.WithValue(in, contextKey, b)
}

func builderFromContext(c context.Context) (b *Builder, ok bool) {
	b, ok = c.Value(contextKey).(*Builder)
	return
}

func LocalePathFromContext(m interface{}, ctx context.Context) (localePath string) {
	l10nBuilder, ok := builderFromContext(ctx)
	if !ok {
		return
	}

	if locale, ok := IsLocalizableFromContext(ctx); ok {
		localePath = l10nBuilder.GetLocale(locale).path
	}

	if localeCode, err := reflectutils.Get(m, "LocaleCode"); err == nil {
		localePath = l10nBuilder.GetLocale(localeCode.(string)).path
	}

	return
}

func (b *Builder) GetAllLocalePaths() (r []string) {
	for _, l := range b.locales {
		r = append(r, l.path)
	}
	return
}

func (b *Builder) GetLocaleLabel(localeCode string) string {
	for _, l := range b.locales {
		if l.code == localeCode {
			return l.label
		}
	}
	return "Unknown"
}

func (b *Builder) GetSupportLocaleCodes() (r []string) {
	for _, l := range b.locales {
		r = append(r, l.code)
	}
	return
}

func (b *Builder) GetSupportLocaleCodesFromRequest(R *http.Request) []string {
	if b.getSupportLocaleCodesFromRequestFunc != nil {
		return b.getSupportLocaleCodesFromRequestFunc(R)
	}
	return b.GetSupportLocaleCodes()
}

func (b *Builder) SupportLocalesFunc(v func(R *http.Request) []string) (r *Builder) {
	b.getSupportLocaleCodesFromRequestFunc = v
	return b
}

func (b *Builder) GetCurrentLocaleCodeFromCookie(r *http.Request) (localeCode string) {
	localeCookie, _ := r.Cookie(b.cookieName)
	if localeCookie != nil {
		localeCode = localeCookie.Value
	}
	return
}

func (b *Builder) GetCorrectLocaleCode(r *http.Request) string {
	localeCode := r.FormValue(b.queryName)
	if localeCode == "" {
		localeCode = b.GetCurrentLocaleCodeFromCookie(r)
	}

	supportLocaleCodes := b.GetSupportLocaleCodesFromRequest(r)
	for _, v := range supportLocaleCodes {
		if localeCode == v {
			return v
		}
	}

	return supportLocaleCodes[0]
}

type l10nContextKey int

const (
	LocaleCode l10nContextKey = iota
	SkipLocaleCode
	LocalizeOptions
)

func (b *Builder) EnsureLocale(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(b.GetSupportLocaleCodesFromRequest(r)) == 0 {
			in.ServeHTTP(w, r)
			return
		}

		localeCode := b.GetCorrectLocaleCode(r)

		maxAge := 365 * 24 * 60 * 60
		http.SetCookie(w, &http.Cookie{
			Name:    b.cookieName,
			Value:   localeCode,
			Path:    "/",
			MaxAge:  maxAge,
			Expires: time.Now().Add(time.Duration(maxAge) * time.Second),
		})
		ctx := context.WithValue(r.Context(), LocaleCode, localeCode)

		in.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (b *Builder) Install(pb *presets.Builder) error {
	pb.SetData(BuilderKey, b)

	db := b.db

	pb.FieldDefaults(presets.LIST).
		FieldType(Locale{}).
		ComponentFunc(localeListFunc(db, b))
	pb.FieldDefaults(presets.WRITE).
		FieldType(Locale{}).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			value := b.localeValue(field, ctx)
			return Input("").Type("hidden").Attr(web.VField("LocaleCode", value)...)
		}).
		SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
			value := EmbedLocale(obj).LocaleCode
			if !slices.Contains(b.GetSupportLocaleCodesFromRequest(ctx.R), value) {
				return IncorrectLocaleErr
			}

			return nil
		})

	pb.AddWrapHandler(WrapHandlerKey, b.EnsureLocale)
	pb.AddMenuTopItemFunc(MenuTopItemFunc, runSwitchLocaleFunc(b))
	pb.I18n().
		RegisterForModule(language.English, I18nLocalizeKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nLocalizeKey, Messages_zh_CN).
		RegisterForModule(language.Japanese, I18nLocalizeKey, Messages_ja_JP)

	return nil
}

func (b *Builder) localeValue(field *presets.FieldContext, ctx *web.EventContext) string {
	var (
		value string
		obj   = field.Obj
	)
	id, err := reflectutils.Get(obj, "ID")
	if err == nil && len(fmt.Sprint(id)) > 0 && fmt.Sprint(id) != "0" {
		value = EmbedLocale(obj).LocaleCode
	} else {
		value = b.GetCorrectLocaleCode(ctx.R)
	}
	return value
}
