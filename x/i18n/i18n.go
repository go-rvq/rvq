package i18n

import (
	"context"
	"fmt"
	"log"
	"maps"
	"net/http"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/text/language"
)

type ModuleKey string

type ErrorString string

func (e ErrorString) Error() string {
	return string(e)
}

type Builder struct {
	getSupportLanguagesFromRequestFunc func(R *http.Request) []language.Tag
	cookieName                         string
	queryName                          string

	// mu orders the writers; the readers load state, which is never changed
	// once stored — a writer stores a copy (copy-on-write), so an update at
	// run time (SetLanguage, SetSupportLanguages) races no request.
	mu    sync.Mutex
	state atomic.Pointer[builderState]
}

// builderState is what the Builder serves: the languages, and per language the
// messages registered by the code and those that replace some of them.
type builderState struct {
	supportLanguages []language.Tag
	matcher          language.Matcher
	// code is the context of the messages the code registered, per language.
	code map[language.Tag]context.Context
	// codeModules is code's messages of a ModuleKey, per language.
	codeModules map[language.Tag]map[ModuleKey]Messages
	// overrides is what SetLanguage gave, per language.
	overrides map[language.Tag]map[ModuleKey]Messages
	// effective is code with overrides on it: what a request gets.
	effective map[language.Tag]context.Context
}

func (s *builderState) clone() *builderState {
	c := *s
	c.supportLanguages = slices.Clone(s.supportLanguages)
	c.code = maps.Clone(s.code)
	c.codeModules = maps.Clone(s.codeModules)
	c.overrides = maps.Clone(s.overrides)
	c.effective = maps.Clone(s.effective)
	return &c
}

// build makes the effective messages of lang.
func (s *builderState) build(lang language.Tag) {
	c := s.code[lang]
	if c == nil {
		c = context.TODO()
	}
	for module, msg := range s.overrides[lang] {
		c = context.WithValue(c, module, msg)
	}
	s.effective[lang] = c
}

func (s *builderState) setSupportLanguages(vs []language.Tag) {
	s.supportLanguages = vs
	for _, l := range vs {
		if s.effective[l] == nil {
			s.build(l)
		}
	}
	s.matcher = language.NewMatcher(vs)
}

type Messages interface{}

func New() *Builder {
	b := &Builder{
		cookieName: "lang",
		queryName:  "lang",
	}
	st := &builderState{
		code:        map[language.Tag]context.Context{},
		codeModules: map[language.Tag]map[ModuleKey]Messages{},
		overrides:   map[language.Tag]map[ModuleKey]Messages{},
		effective:   map[language.Tag]context.Context{},
	}
	st.setSupportLanguages([]language.Tag{language.English})
	b.state.Store(st)

	b.RegisterForModules(language.English, DefaultKey, Default_en).
		RegisterForModules(language.BrazilianPortuguese, DefaultKey, Default_pt_BR)

	return b
}

// update stores a copy of the state changed by do.
func (b *Builder) update(do func(s *builderState)) *Builder {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.state.Load().clone()
	do(s)
	b.state.Store(s)
	return b
}

func (b *Builder) GetCookieName() string {
	return b.cookieName
}

func (b *Builder) GetQueryName() string {
	return b.queryName
}

// SupportLanguages sets the languages, the default first. It may be called at
// run time (SetSupportLanguages is it).
func (b *Builder) SupportLanguages(vs ...language.Tag) (r *Builder) {
	if len(vs) == 0 {
		panic("have to support at least one language")
	}
	vs = slices.Clone(vs)
	return b.update(func(s *builderState) { s.setSupportLanguages(vs) })
}

// SetSupportLanguages sets the languages at run time, the default first: the
// languages enabled in a database, say.
func (b *Builder) SetSupportLanguages(vs ...language.Tag) (r *Builder) {
	return b.SupportLanguages(vs...)
}

func (b *Builder) SupportLanguage(vs ...language.Tag) (r *Builder) {
	return b.update(func(s *builderState) {
		langs := s.supportLanguages
		for _, v := range vs {
			if !slices.Contains(langs, v) {
				langs = append(langs, v)
			}
		}
		s.setSupportLanguages(langs)
	})
}

func (b *Builder) GetSupportLanguages() []language.Tag {
	return b.state.Load().supportLanguages
}

func (b *Builder) GetSupportLanguagesFromRequest(R *http.Request) []language.Tag {
	if b.getSupportLanguagesFromRequestFunc != nil {
		return b.getSupportLanguagesFromRequestFunc(R)
	}
	return b.GetSupportLanguages()
}

func (b *Builder) GetSupportLanguagesFromRequestFunc(v func(R *http.Request) []language.Tag) (r *Builder) {
	b.getSupportLanguagesFromRequestFunc = v
	return b
}

func (b *Builder) RegisterForModule(lang language.Tag, module ModuleKey, msg Messages) (r *Builder) {
	return b.RegisterForModules(lang, module, msg)
}

// RegisterForModules registers the code's messages of lang: pairs of a key —
// a ModuleKey, usually — and its messages.
func (b *Builder) RegisterForModules(lang language.Tag, args ...any) (r *Builder) {
	if len(args)%2 != 0 {
		panic("invalid number of arguments")
	}
	return b.update(func(s *builderState) {
		c := s.code[lang]
		if c == nil {
			c = context.TODO()
		}
		modules := maps.Clone(s.codeModules[lang])
		if modules == nil {
			modules = map[ModuleKey]Messages{}
		}
		for ; len(args) > 0; args = args[2:] {
			c = context.WithValue(c, args[0], args[1])
			if module, ok := args[0].(ModuleKey); ok {
				modules[module] = args[1]
			}
		}
		s.code[lang] = c
		s.codeModules[lang] = modules
		s.build(lang)
	})
}

// SetLanguage sets the messages of lang that replace the code's — those of a
// database, say —, at run time: every module of lang not in modules is the
// code's again. Nil modules leaves lang with the code's alone.
func (b *Builder) SetLanguage(lang language.Tag, modules map[ModuleKey]Messages) (r *Builder) {
	modules = maps.Clone(modules)
	return b.update(func(s *builderState) {
		if len(modules) == 0 {
			delete(s.overrides, lang)
		} else {
			s.overrides[lang] = modules
		}
		s.build(lang)
	})
}

// CodeModules is the messages the code registered, of each ModuleKey, per
// language. The maps are the Builder's: they must not be changed.
func (b *Builder) CodeModules() map[language.Tag]map[ModuleKey]Messages {
	return b.state.Load().codeModules
}

// GetCodeModuleMessages is the messages of module the code registered for
// lang: those of GetModuleMessages without the ones SetLanguage gave.
func (b *Builder) GetCodeModuleMessages(lang language.Tag, module ModuleKey) any {
	c := b.state.Load().code[lang]
	if c == nil {
		return nil
	}
	return c.Value(module)
}

// GetModuleMessages is the messages of module that lang has: the code's, or
// those SetLanguage gave in their place.
func (b *Builder) GetModuleMessages(lang language.Tag, module ModuleKey) any {
	c := b.state.Load().effective[lang]
	if c == nil {
		return nil
	}
	return c.Value(module)
}

func MustGetModuleMessages(ctx context.Context, module ModuleKey, defaultMessages Messages) Messages {
	v := ctx.Value(moduleMessagesKey)
	if v == nil {
		return defaultMessages
	}

	msg := v.(context.Context).Value(module)
	if msg == nil {
		msg = defaultMessages
	}
	return msg
}

type i18nContextKey int

const (
	moduleMessagesKey i18nContextKey = iota
	dynaBuilderKey
)

func (b *Builder) EnsureLanguage(in http.Handler) (out http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := ""
		lang = r.FormValue(b.queryName)
		if len(lang) > 0 {
			maxAge := 365 * 24 * 60 * 60
			http.SetCookie(w, &http.Cookie{
				Name:    b.cookieName,
				Value:   lang,
				Path:    "/",
				MaxAge:  maxAge,
				Expires: time.Now().Add(time.Duration(maxAge) * time.Second),
			})
		} else {
			lang = b.GetCurrentLangFromCookie(r)
		}

		accept := r.Header.Get("Accept-Language")

		st := b.state.Load()
		var availableLanguages []language.Tag
		var matcher language.Matcher
		if len(lang) > 0 {
			availableLanguages = st.supportLanguages
			matcher = st.matcher
		} else {
			availableLanguages = b.GetSupportLanguagesFromRequest(r)
			matcher = language.NewMatcher(availableLanguages)
		}
		_, i := language.MatchStrings(matcher, lang, accept)
		tag := availableLanguages[i]

		moduleMsgs := st.effective[tag]
		if moduleMsgs == nil {
			moduleMsgs = st.effective[st.supportLanguages[0]]
		}
		if moduleMsgs == nil {
			panic(fmt.Sprintf("language %s not supported", tag.String()))
		}
		dyna := DynaNew().Language(tag.String())
		ctx := context.WithValue(r.Context(), moduleMessagesKey, moduleMsgs)
		ctx = context.WithValue(ctx, dynaBuilderKey, dyna)
		in.ServeHTTP(w, r.WithContext(ctx))
		if dyna.HaveMissingKeys() {
			log.Println(dyna.PrettyMissingKeys())
		}
	})
}

func (b *Builder) GetCurrentLangFromCookie(r *http.Request) (lang string) {
	langCookie, _ := r.Cookie(b.cookieName)
	if langCookie != nil {
		lang = langCookie.Value
	}
	return
}
