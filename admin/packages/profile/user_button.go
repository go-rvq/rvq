package profile

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/account"
	"github.com/go-rvq/rvq/x/login"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// The button of the user, at the end of the top bar (presets_layout.go): the
// avatar, the name and the role; it opens the menu of the user — the account
// (its index when the browser holds several), the quick access to what is the
// user's (presets.UserMenuQuick), the preferences (presets.UserMenuPrefs), the
// other accounts of the browser, "Entrar com outra conta", "Sair desta
// conta" / "Sair de todas as contas".

// accountColors: the color of an account, by its user (the same on every tab).
// (CSS colors: the palette classes of Vuetify are not in every build)
var accountColors = []string{"#1976d2", "#388e3c", "#673ab7", "#ef6c00", "#00897b", "#c2185b", "#3f51b5", "#795548"}

func accountColor(id string) string {
	f := fnv.New32a()
	f.Write([]byte(id))
	return accountColors[int(f.Sum32()%uint32(len(accountColors)))]
}

func shortName(u user.User) string {
	name := u.GetName()
	if name == "" {
		name = u.GetAccountName()
	}
	if rs := []rune(name); len(rs) > 1 {
		name = string(rs[:1])
	}
	return strings.ToUpper(name)
}

func displayName(u user.User) string {
	if n := u.GetName(); n != "" {
		return n
	}
	return u.GetAccountName()
}

// UserButton is the button of the user and their menu.
func (c *Builder) UserButton(b *presets.Builder) func(ctx *web.EventContext) h.HTMLComponent {
	return func(ctx *web.EventContext) h.HTMLComponent {
		u := user.GetCurrentUser(ctx.R)
		if u == nil || u.Anonymous() {
			return VBtn(GetMessages(ctx.Context()).Login).Variant(VariantText).Href(c.lb.GetLoginPageURL())
		}
		accounts := c.lb.Accounts(ctx.R)
		multi := len(accounts) > 1
		color := ColorPrimary
		if multi {
			color = accountColor(u.GetID().String())
		}

		photo := avatarURL(u)
		avatar := func(size int) h.HTMLComponent {
			a := VAvatar().Color(color).Size(strconv.Itoa(size))
			if photo != "" {
				a.Image(photo)
			} else {
				a.Text(shortName(u))
			}
			if multi {
				a.Class("rvq-user-avatar--ring").Style("outline: 2px solid currentColor; outline-offset: 2px")
			}
			return a
		}

		others := 0
		for _, a := range accounts {
			if !a.Current {
				others++
			}
		}

		activator := VBtn("").
			Variant(VariantText).
			Class("rvq-user-button px-2").
			Height(44).
			Attr("v-bind", "props").
			Children(
				h.Div(
					VBadge(avatar(32)).
						Content(others).
						Color("success").
						Location("bottom start").
						Attr(":model-value", h.JSONString(others > 0)),
					h.Div(
						h.Div(h.Text(displayName(u))).Class("text-body-2 font-weight-medium text-truncate").Style("max-width: 160px"),
						h.Div(h.Text(u.GetRoles().FirstName())).Class("text-caption text-medium-emphasis text-truncate").Style("max-width: 160px"),
					).Class("d-none d-sm-flex flex-column align-start ml-2").Style("line-height: 1.15"),
					VIcon("mdi-chevron-down").Size(SizeSmall).Class("ml-1 d-none d-sm-flex"),
				).Class("d-flex align-center"),
			)

		menu := c.userMenu(ctx, b, u, accounts, avatar)

		// one menu, under the button: 320px, a phone's width too
		return VMenu(
			web.Slot(activator).Name("activator").Scope("{ props }"),
			menu,
		).Location(LocationBottom + " end").Offset(6).CloseOnContentClick(true).MaxHeight("calc(100vh - 72px)")
	}
}

func (c *Builder) userMenu(ctx *web.EventContext, b *presets.Builder, u user.User, accounts []*login.Account, avatar func(int) h.HTMLComponent) h.HTMLComponent {
	msgr := GetMessages(ctx.Context())
	cur := account.IndexOf(ctx.R)
	multi := len(accounts) > 1
	prefix := b.URIPrefixCtx(ctx.R.Context())

	var header h.HTMLComponents
	header = append(header, VListItem(
		web.Slot(h.A(avatar(48)).Href(prefix+"/"+ModelID).Class("mr-3")).Name(VSlotPrepend),
		VListItemTitle(h.Text(displayName(u))).Class("font-weight-bold"),
		VListItemSubtitle(h.Text(u.GetAccountName())),
		h.Div(
			h.Iff(u.GetRoles().FirstName() != "", func() h.HTMLComponent {
				return VChip(h.Text(u.GetRoles().FirstName())).Size(SizeXSmall).Color(ColorPrimary).Variant(VariantTonal)
			}),
			h.If(multi, VChip(h.Text("u/"+strconv.Itoa(cur))).Size(SizeXSmall).Color("orange-darken-3").Variant(VariantTonal).Class("ml-1")),
		).Class("mt-1"),
	).Class("py-3"))

	items := h.HTMLComponents{
		header,
		VDivider(),
	}
	if quick := b.UserMenuItems(ctx, presets.UserMenuQuick); len(quick) > 0 {
		items = append(items, quick, VDivider())
	}
	if prefs := b.UserMenuItems(ctx, presets.UserMenuPrefs); len(prefs) > 0 {
		items = append(items, prefs, VDivider())
	}

	// the other accounts of the browser
	if multi {
		items = append(items, VListSubheader(h.Text(msgr.OtherAccounts)))
		for _, a := range accounts {
			if a.Current {
				continue
			}
			au, _ := a.User.(user.User)
			if au == nil {
				continue
			}
			href := c.accountSwitchURL(ctx, a.Index)
			sub := au.GetAccountName()
			if r := au.GetRoles().FirstName(); r != "" {
				sub += " · " + r
			}
			cls := ""
			if !a.LoggedIn {
				href = c.lb.AccountLoginURL(a.Index)
				sub = msgr.SessionExpired
				cls = "rvq-user-menu__expired"
			}
			items = append(items, VListItem(
				web.Slot(otherAvatar(au)).Name(VSlotPrepend),
				VListItemTitle(h.Text(displayName(au))),
				VListItemSubtitle(h.Text(sub)),
				web.Slot(VChip(h.Text("u/"+strconv.Itoa(a.Index))).Size(SizeXSmall).Variant(VariantTonal)).Name(VSlotAppend),
			).Href(href).Class(cls))
		}
		items = append(items, VDivider())
	}

	// another account, the logouts
	if max := c.lb.GetMaxAccounts(); max > 1 {
		add := VListItem(
			web.Slot(VIcon("mdi-account-plus-outline").Color(ColorPrimary)).Name(VSlotPrepend),
			VListItemTitle(h.Text(msgr.AddAccount)).Class("text-primary font-weight-medium"),
		)
		if len(accounts) >= max {
			add.Disabled(true).Children(VListItemSubtitle(h.Text(fmt.Sprintf(msgr.AccountsLimit, max))))
		} else {
			add.Href(c.lb.AddAccountURL(ctx.R.RequestURI))
		}
		items = append(items, add)
	}
	logoutLabel := msgr.Logout
	if multi {
		logoutLabel = msgr.LogoutThisAccount
	}
	items = append(items, VListItem(
		web.Slot(VIcon("mdi-logout").Color("error")).Name(VSlotPrepend),
		VListItemTitle(h.Text(logoutLabel)).Class("text-error font-weight-medium"),
	).Href(c.lb.AccountLogoutURL(cur, false)))
	if multi {
		items = append(items, VListItem(
			web.Slot(VIcon("mdi-logout-variant")).Name(VSlotPrepend),
			VListItemTitle(h.Text(msgr.LogoutAll)),
			web.Slot(h.Span(strconv.Itoa(len(accounts))).Class("text-caption text-medium-emphasis")).Name(VSlotAppend),
		).Href(c.lb.AccountLogoutURL(cur, true)))
	}

	// the list as the admin's others: its size, its density
	return VCard(
		VList(items...),
	).Class("rvq-user-menu")
}

// accountSwitchURL is the page of the request in the account n: its index
// said (/admin/u/N/…, even the 0's), the page itself.
func (c *Builder) accountSwitchURL(ctx *web.EventContext, n int) string {
	p := ctx.R.URL.Path
	q := ctx.R.URL.Query()
	q.Del("__execute_event__")
	q.Del("__reload__")
	q.Del(presets.NoAccessQuery)
	// may that account not see it: its home instead (presets)
	q.Set(presets.SwitchAccountQuery, "1")
	for _, prefix := range c.lb.GetAccountPrefixes() {
		if p == prefix || len(p) > len(prefix) && p[:len(prefix)+1] == prefix+"/" {
			out := prefix + "/" + account.Segment + "/" + strconv.Itoa(n) + p[len(prefix):]
			if enc := q.Encode(); enc != "" {
				out += "?" + enc
			}
			return out
		}
	}
	return c.lb.AccountHomeURL(n)
}

// otherAvatar is the avatar of another account of the browser.
func otherAvatar(u user.User) h.HTMLComponent {
	a := VAvatar().Color(accountColor(u.GetID().String())).Size("32").Class("mr-3")
	if p := avatarURL(u); p != "" {
		return a.Image(p)
	}
	return a.Text(shortName(u))
}
