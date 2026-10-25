package presets

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/utils/bcp47"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/account"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"golang.org/x/text/language"
)

// UserMenuGroup is a section of the menu of the user (the button of the user
// in the top bar): the quick access to what is the user's, the preferences.
type UserMenuGroup int

const (
	// UserMenuQuick: the user's things — the profile, the keys, the sessions,
	// the notifications, the activity.
	UserMenuQuick UserMenuGroup = iota
	// UserMenuPrefs: the preferences — the language, the theme.
	UserMenuPrefs
)

// Built-in user-menu items.
const (
	UserMenuItemLanguage = "Language"
	UserMenuItemTheme    = "Theme"

	// UserMenuThemeStorageKey is where the browser keeps the theme of the
	// admin (light, dark).
	UserMenuThemeStorageKey = "rvq.admin.theme"
)

// UserMenuItem is an item of the menu of the user: Handler renders it (a
// VListItem), nil to hide it; by Order, then as added.
type UserMenuItem struct {
	Name    string
	Group   UserMenuGroup
	Order   int
	Handler ComponentFunc
	Enabled bool
}

// AddUserMenuItem adds items to the menu of the user.
func (b *Builder) AddUserMenuItem(items ...*UserMenuItem) *Builder {
	b.userMenuItems = append(b.userMenuItems, items...)
	return b
}

// GetUserMenuItem is the user-menu item named name (nil if absent).
func (b *Builder) GetUserMenuItem(name string) *UserMenuItem {
	for _, it := range b.userMenuItems {
		if it.Name == name {
			return it
		}
	}
	return nil
}

// UserMenuItems renders the enabled items of group, by Order.
func (b *Builder) UserMenuItems(ctx *web.EventContext, group UserMenuGroup) (comps h.HTMLComponents) {
	var items []*UserMenuItem
	for _, it := range b.userMenuItems {
		if it.Group == group && it.Enabled && it.Handler != nil {
			items = append(items, it)
		}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Order < items[j].Order })
	for _, it := range items {
		if c := it.Handler(ctx); c != nil {
			comps = append(comps, c)
		}
	}
	return
}

// AddTopBarItem adds what the top bar shows before the user (the help's
// button, …): f renders it, nil to hide it.
func (b *Builder) AddTopBarItem(f ...ComponentFunc) *Builder {
	b.topBarItems = append(b.topBarItems, f...)
	return b
}

func (b *Builder) topBarItemsOf(ctx *web.EventContext) (comps h.HTMLComponents) {
	for _, f := range b.topBarItems {
		if c := f(ctx); c != nil {
			comps = append(comps, c)
		}
	}
	return
}

// SwitchAccountQuery marks the link of the menu of the user to another
// account of the browser (the page itself in it): when that account may not
// see it, its home (NoAccessQuery) instead.
const SwitchAccountQuery = "rvq_switch"

// NoAccessQuery: the home an account was sent to, the page it was switched to
// being one it may not see — said by a notice.
const NoAccessQuery = "rvq_no_access"

// URIPrefixCtx is the prefix of the admin for the account of the request:
// /admin, or /admin/u/N for the account N of a browser holding several.
func (b *Builder) URIPrefixCtx(ctx context.Context) string {
	return account.PrefixCtx(ctx, b.prefix)
}

// the built-in items of the user menu
func (b *Builder) addDefaultUserMenuItems() {
	b.AddUserMenuItem(
		&UserMenuItem{Name: UserMenuItemLanguage, Group: UserMenuPrefs, Order: 10, Enabled: true, Handler: b.userMenuLanguage},
		&UserMenuItem{Name: UserMenuItemTheme, Group: UserMenuPrefs, Order: 20, Enabled: true, Handler: b.userMenuTheme},
	)
}

// userMenuLanguage: the language of the admin, a submenu of the languages
// (each one the same page with ?lang=…, as the side menu's switch did).
func (b *Builder) userMenuLanguage(ctx *web.EventContext) h.HTMLComponent {
	supported := b.I18n().GetSupportLanguagesFromRequest(ctx.R)
	if len(b.I18n().GetSupportLanguages()) <= 1 || len(supported) <= 1 {
		return nil
	}
	msgr := MustGetMessages(ctx.Context())
	queryName := b.I18n().GetQueryName()
	lang := ctx.R.FormValue(queryName)
	if lang == "" {
		lang = b.i18nBuilder.GetCurrentLangFromCookie(ctx.R)
	}
	_, cur := language.MatchStrings(language.NewMatcher(supported), lang, ctx.R.Header.Get("Accept-Language"))

	var items h.HTMLComponents
	for i, tag := range supported {
		u, _ := url.Parse(ctx.R.RequestURI)
		qs := u.Query()
		qs.Del("__execute_event__")
		qs.Del("__reload__")
		qs.Set(queryName, tag.String())
		u.RawQuery = qs.Encode()
		items = append(items, VListItem(
			VListItemTitle(h.Text(bcp47.FlagLabel(tag.String()))),
		).Href(u.String()).Active(i == cur))
	}
	return VMenu(
		web.Slot(
			VListItem(
				web.Slot(VIcon("mdi-translate")).Name(VSlotPrepend),
				VListItemTitle(h.Text(msgr.Language)),
				web.Slot(
					h.Span(supported[cur].String()).Class("text-caption text-medium-emphasis mr-1"),
					VIcon("mdi-chevron-right").Size(SizeSmall),
				).Name(VSlotAppend),
			).Attr("v-bind", "props"),
		).Name("activator").Scope("{ props }"),
		VList(items...).Density(DensityCompact),
	).Location("start").OpenOnHover(false)
}

// userMenuTheme: light or dark, kept in the browser.
func (b *Builder) userMenuTheme(ctx *web.EventContext) h.HTMLComponent {
	msgr := MustGetMessages(ctx.Context())
	return VListItem(
		web.Slot(VIcon("mdi-theme-light-dark")).Name(VSlotPrepend),
		VListItemTitle(h.Text(msgr.Theme)),
		web.Slot(
			h.Tag("vx-theme-toggle").Attr("storage-key", UserMenuThemeStorageKey),
		).Name(VSlotAppend),
	)
}

// accountScript keeps a page of the account N (/admin/u/N/…) in it: its
// requests say it (the header X-Rvq-Account), its links and its history
// get its /u/N; the address of a page of it reached without one, too.
func accountScript(prefix string, n int) h.HTMLComponent {
	return h.Script(fmt.Sprintf(`(function(){
var P=%q,N=%d,H=%q,S=P+"/u/";
if(window.__rvqAccount)return;
function fix(u){try{var x=new URL(u,location.href);if(x.origin!==location.origin)return u;var p=x.pathname;
if((p===P||p.indexOf(P+"/")===0)&&p.indexOf(S)!==0){x.pathname=S+N+p.slice(P.length);return x.pathname+x.search+x.hash}}catch(e){}return u}
function same(u){try{return new URL(u,location.href).origin===location.origin}catch(e){return false}}
window.__rvqAccount={prefix:P,index:N,fix:fix};
var f=window.fetch;window.fetch=function(i,o){o=o||{};try{var u=typeof i==="string"?i:(i&&i.url)||String(i);
if(same(u)){var hd=new Headers(o.headers||(i instanceof Request?i.headers:undefined));if(!hd.has(H))hd.set(H,String(N));o.headers=hd}}catch(e){}
return f.call(this,i,o)};
var xo=XMLHttpRequest.prototype.open,xs=XMLHttpRequest.prototype.send;
XMLHttpRequest.prototype.open=function(m,u){this.__rvqSame=same(u);return xo.apply(this,arguments)};
XMLHttpRequest.prototype.send=function(){if(this.__rvqSame){try{this.setRequestHeader(H,String(N))}catch(e){}}return xs.apply(this,arguments)};
var ps=history.pushState,rs=history.replaceState;
history.pushState=function(s,t,u){return ps.call(this,s,t,u==null?u:fix(String(u)))};
history.replaceState=function(s,t,u){return rs.call(this,s,t,u==null?u:fix(String(u)))};
document.addEventListener("click",function(e){var a=e.target&&e.target.closest&&e.target.closest("a[href]");if(!a)return;
var v=a.getAttribute("href");if(!v||v.charAt(0)==="#")return;var g=fix(v);if(g!==v)a.setAttribute("href",g)},true);
document.addEventListener("submit",function(e){var fm=e.target;if(!fm||!fm.getAttribute)return;var v=fm.getAttribute("action");
if(v){var g=fix(v);if(g!==v)fm.setAttribute("action",g)}},true);
var here=fix(location.pathname+location.search+location.hash);if(here!==location.pathname+location.search+location.hash)rs.call(history,history.state,"",here);
})();`, prefix, n, account.Header))
}
