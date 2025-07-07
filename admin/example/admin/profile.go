package admin

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/go-rvq/rvq/admin/example/models"
	plogin "github.com/go-rvq/rvq/admin/login"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"
)

const (
	signOutAllSessionEvent = "signOutAllSessionEvent"
)

func profile(db *gorm.DB) presets.ComponentFunc {
	return func(ctx *web.EventContext) h.HTMLComponent {
		u := getCurrentUser(ctx.R)
		if u == nil {
			return VBtn("Login").Variant(VariantText).Href("/auth/login")
		}

		var roles []string
		for _, role := range u.Roles {
			roles = append(roles, role.Name)
		}

		total := notifierCount(db)(ctx)
		content := notifierComponent(db)(ctx)
		icon := VIcon("mdi-bell-outline").Size(20).Color("grey-darken-1")
		notification := VMenu().Children(
			h.Template().Attr("v-slot:activator", "{ props }").Children(
				VBtn("").Icon(true).Children(
					h.If(total > 0,
						VBadge(
							icon,
						).Content(total).Dot(true).Color("error"),
					).Else(icon),
				).Attr("v-bind", "props").
					Density(DensityCompact).
					Variant(VariantText),
			),
			VCard(content),
		)
		_ = notification
		profileNewLook := VCard(
			web.Slot(
				VAvatar().Text(getAvatarShortName(u)).Color(ColorPrimaryLighten2).Size(SizeLarge).Class(fmt.Sprintf("rounded-lg text-%s", ColorPrimary)),
			).Name(VSlotPrepend),
			web.Slot(
				h.Div(
					h.Div(h.Text(u.Name)).Class(fmt.Sprintf(`text-subtitle-2 text-%s`, ColorSecondary)),
					VBtn("").Attr("@click", web.Plaid().URL(logoutURL).Go()).
						Icon(true).Density(DensityCompact).Variant(VariantText).Children(
						VIcon("mdi-chevron-right").Size(SizeSmall),
					).Class("mr-8"),
				).Class("d-flex justify-space-between align-center"),
			).Name(VSlotTitle),
			web.Slot(
				h.Div(h.Text(roles[0])),
			).Name(VSlotSubtitle),
			web.Slot(
				VRow(
					VCol(
						notification,
					),
				).Class("border-s"),
			).Name(VSlotAppend),
		).Class(W100)
		return profileNewLook
	}
}

type Profile struct{}

func configProfile(b *presets.Builder, db *gorm.DB) {
	m := b.Model(&Profile{}).URIName("profile").
		MenuIcon("mdi-account").Label("Profile").Singleton(true)

	eb := m.Editing("Info", "Actions", "Sessions")

	m.RegisterEventHandler(signOutAllSessionEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)

		u := getCurrentUser(ctx.R)

		if u.GetAccountName() == loginInitialUserEmail {
			return r, perm.PermissionDenied
		}

		if err = expireOtherSessionLogs(db, ctx.R, u.ID); err != nil {
			return r, err
		}

		presets.ShowMessage(&r, msgr.SignOutAllSuccessfullyTips, "")
		r.Reload = true
		return
	})

	eb.FetchFunc(func(obj interface{}, id string, ctx *web.EventContext) (r interface{}, err error) {
		u := getCurrentUser(ctx.R)
		if u == nil {
			return nil, errors.New("cannot get current user")
		}
		return u, nil
	})

	eb.SetterFunc(func(obj interface{}, ctx *web.EventContext) {
		u := obj.(*models.User)
		u.Name = ctx.R.FormValue("name")
		return
	})

	eb.Field("Info").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)

		u := obj.(*models.User)
		var roles []string
		for _, v := range u.Roles {
			roles = append(roles, v.Name)
		}

		return h.Div(
			VRow(
				VCol(
					VTextField().Label(msgr.Name).Attr(web.VField("name", u.Name)...),
				),
			).Class("my-n6"),
			VRow(
				VCol(
					vx.VXReadonlyField().Label(msgr.Email).Value(u.Account),
				),
			).Class("my-n6"),
			VRow(
				VCol(
					vx.VXReadonlyField().Label(msgr.Company).Value(u.Company),
				),
			).Class("my-n6"),
			VRow(
				VCol(
					vx.VXReadonlyField().Label(msgr.Role).Value(strings.Join(roles, ", ")),
				),
			).Class("my-n6"),
			VRow(
				VCol(
					vx.VXReadonlyField().Label(msgr.Status).Value(u.Status),
				),
			),
		).Class("mx-2 mt-4")
	})

	eb.Field("Actions").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)

		// We don't allow public user to change its password
		u := getCurrentUser(ctx.R)
		if u.GetAccountName() == loginInitialUserEmail {
			return h.RawHTML("")
		}

		var actionBtns h.HTMLComponents
		if u.OAuthProvider == "" && u.Account != "" {
			actionBtns = append(actionBtns,
				VBtn("").
					Variant(VariantOutlined).Color("primary").
					Children(VIcon("lock_outline").Size(SizeSmall), h.Text(msgr.ChangePassword)).
					Class("mr-2").
					OnClick(plogin.OpenChangePasswordDialogEvent),
			)
		}

		return h.Div(
			actionBtns...,
		).Class("mx-2 mt-4 text-left")
	})

	eb.Field("Sessions").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		msgr := i18n.MustGetModuleMessages(ctx.R, I18nExampleKey, Messages_en_US).(*Messages)

		u := obj.(*models.User)
		items := []*models.LoginSession{}
		if err := db.Where("user_id = ?", u.ID).Find(&items).Error; err != nil {
			panic(err)
		}

		isPublicUser := false
		if u.GetAccountName() == loginInitialUserEmail {
			isPublicUser = true
		}

		currentTokenHash := getStringHash(login.GetSessionToken(loginBuilder, ctx.R), LoginTokenHashLen)

		var (
			expired        = msgr.Expired
			active         = msgr.Active
			currentSession = msgr.CurrentSession
		)

		activeDevices := make(map[string]struct{})
		for _, item := range items {
			if isPublicUser {
				item.IP = msgr.HideIPTips
			}

			if isTokenValid(*item) {
				item.Status = expired
			} else {
				item.Status = active
				activeDevices[fmt.Sprintf("%s#%s", item.Device, item.IP)] = struct{}{}
			}
			if item.TokenHash == currentTokenHash {
				item.Status = currentSession
			}

			item.Time = humanize.Time(item.CreatedAt)
		}

		{
			newItems := make([]*models.LoginSession, 0, len(items))
			for _, item := range items {
				if item.Status == expired {
					_, ok := activeDevices[fmt.Sprintf("%s#%s", item.Device, item.IP)]
					if ok {
						continue
					}
				}
				newItems = append(newItems, item)
			}
			items = newItems
		}

		if isPublicUser {
			if len(items) > 10 {
				items = items[:10]
			}
		}

		sort.Slice(items, func(i, j int) bool {
			if items[j].Status == currentSession {
				return false
			}
			if items[i].Status == expired &&
				items[j].Status == active {
				return false
			}
			if items[i].CreatedAt.Sub(items[j].CreatedAt) < 0 {
				return false
			}
			return true
		})

		sessionTableHeaders := []DataTableHeader{
			{msgr.Time, "Time", "25%", false},
			{msgr.Device, "Device", "25%", false},
			{msgr.IPAddress, "IP", "25%", false},
			{"", "Status", "25%", true},
		}

		return h.Div(
			VCard(
				VRow(
					VCol(
						VCardTitle(h.Text(msgr.LoginSessions)),
						VCardSubtitle(h.Text(msgr.LoginSessionsTips)),
					),
					VCol(
						h.If(!isPublicUser,
							VBtn("").Attr("@click", web.Plaid().EventFunc(signOutAllSessionEvent).Go()).
								Variant(VariantOutlined).Color("primary").
								Children(VIcon("warning").Size(SizeSmall), h.Text(msgr.SignOutAllOtherSessions))),
					).Class("text-right mt-6 mr-4"),
				),
				VDataTable().Headers(sessionTableHeaders).
					Items(items).
					ItemsPerPage(-1),
				// TODO fix it .HideDefaultFooter(true),
			),
		).Class("mx-2 mt-12 mb-4")
	})
}

func getAvatarShortName(u *models.User) string {
	name := u.Name
	if name == "" {
		name = u.Account
	}
	if rs := []rune(name); len(rs) > 1 {
		name = string(rs[:1])
	}

	return strings.ToUpper(name)
}
