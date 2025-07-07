package profile

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/helper/login_session"
	"github.com/go-rvq/rvq/admin/helper/user"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

const (
	signOutAllSession       = "signOutAllSession"
	loginSessionDialogEvent = "loginSessionDialogEvent"
	ModelID                 = "my-profile"
)

type Options struct {
	CurrentUser func(r *http.Request) user.User
	NewUser     func() user.User
}

type Builder struct {
	db                        *gorm.DB
	lb                        *login.Builder
	userMb                    *presets.ModelBuilder
	mgr                       *login_session.Manager
	notificationComponentFunc func(ctx *web.EventContext, u user.User) h.HTMLComponent
}

func New(db *gorm.DB, mgr *login_session.Manager, lb *login.Builder, userMb *presets.ModelBuilder) *Builder {
	ConfigureMessages(userMb.Builder().I18n())

	return &Builder{
		db:     db,
		lb:     lb,
		userMb: userMb,
		mgr:    mgr,
	}
}

func (c *Builder) NotificationComponentFunc() func(ctx *web.EventContext, u user.User) h.HTMLComponent {
	return c.notificationComponentFunc
}

func (c *Builder) SetNotificationComponentFunc(notificationComponentFunc func(ctx *web.EventContext, u user.User) h.HTMLComponent) *Builder {
	c.notificationComponentFunc = notificationComponentFunc
	return c
}

func (c *Builder) ProfilePage(p *presets.Builder) http.Handler {
	return p.Wrap(
		p.GetLayoutFunc()(func(ctx *web.EventContext) (r web.PageResponse, err error) {
			u := user.GetCurrentUser(ctx.R)
			r.PageTitle = "Meu Perfil"
			r.Actions = append(r.Actions,
				VBtn("Sessões de Login").
					Attr("@click", web.Plaid().URL("/profile").EventFunc(loginSessionDialogEvent).Query(presets.ParamID, u.GetID()).Go()),
				VBtn("Sair").
					Color("primary").
					Variant(VariantFlat).
					// Size(SizeSmall).
					Class("ml-2").
					Attr("@click", web.Plaid().URL(c.lb.GetLogoutURL()).Go()),
			)
			r.Body = h.HTMLComponents{}
			return
		}, nil))
}

func (c *Builder) Brand(b *presets.Builder) func(ctx *web.EventContext) h.HTMLComponent {
	return func(ctx *web.EventContext) h.HTMLComponent {
		u := user.GetCurrentUser(ctx.R)
		if u == nil {
			return VBtn("Entrar").Variant(VariantText).Href(c.lb.GetLoginPageURL())
		}

		var (
			notification h.HTMLComponent
			roles        = u.GetRoles()
		)

		if c.notificationComponentFunc != nil {
			notification = c.notificationComponentFunc(ctx, u)
		}

		prependProfile := VCard(
			web.Slot(
				VAvatar().Text(c.getAvatarShortName(u)).
					Color(ColorPrimaryLighten2).
					Size(SizeLarge).
					Class(fmt.Sprintf("rounded-lg text-%s", ColorPrimary)),
			).Name(VSlotPrepend),
			web.Slot(
				h.Div(
					VBtn(u.GetName()).
						Class("px-0").
						Style("min-width: 0").
						Variant(VariantText).
						Density(DensityCompact).
						Attr("@click", web.Plaid().URL(b.GetURIPrefix()+"/"+ModelID).PushState(true).Go()).
						Class(fmt.Sprintf(`text-subtitle-2 text-%s`, ColorSecondary)),
				).Class("d-flex justify-space-between align-center"),
			).Name(VSlotTitle),
			web.Slot(
				h.Div(h.Text(roles.FirstName())),
			).Name(VSlotSubtitle),
		).Class(W100).Flat(true)

		profileNewLook := VCard(
			VCardTitle(
				prependProfile,
				VCardText(
					h.Div(
						notification,
						VBtn("").
							Icon(true).
							Color("error").
							Density(DensityCompact).
							Variant(VariantText).
							Attr("@click", web.Plaid().URL(c.lb.GetLogoutURL()).Go()).
							Children(VIcon("mdi-logout")),
					).Class("border-s-md", "pl-4", H75)),
			).Class("d-inline-flex align-center justify-space-between  justify-center pa-0", W100),
		).Class("pa-0").Class(W100)
		return profileNewLook
	}
}

func (c *Builder) Install(p *presets.Builder) (err error) {
	if permB := p.GetPermission(); permB != nil {
		permB.CreatePolicies(perm.PolicyFor(user.RoleLogged).WhoAre(perm.Allowed).ToDo(perm.Anything).On("*:my_profile:*"))
	}

	m := Model(p, &Profile{}, presets.ModelWithID(ModelID)).
		Singleton(true).
		MenuIcon("mdi-account").
		InMenu(false)

	m.Singleton(true)
	m.DeletingRestriction.ObjHandler(presets.OkObjHandlerFunc(func(obj any, ctx *web.EventContext) (ok, handled bool) {
		return true, true
	}))

	fetch := func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
		u := user.GetCurrentUser(ctx.R)
		if u == nil {
			return presets.ErrRecordNotFound
		}
		profile := obj.(*Profile)
		profile.ID = u.GetID()
		profile.Status = u.GetStatus()
		profile.Name = u.GetName()
		profile.AccountName = u.GetAccountName()
		profile.Roles = strings.Join(u.GetRoles().Names(), ", ")
		return
	}

	e := m.Editing()
	e.FetchFunc(fetch)
	e.SaveFunc(func(obj interface{}, id presets.ID, ctx *web.EventContext) (err error) {
		u := user.GetCurrentUser(ctx.R)

		if u == nil {
			return presets.ErrRecordNotFound
		}

		if err = c.db.First(u).Error; err != nil {
			return
		}

		if err = c.db.Save(u).Error; err != nil {
			return
		}

		return nil
	})

	d := m.Detailing()
	d.FetchFunc(fetch)

	d.Action("login_sessions").
		Icon("mdi-account-network").
		SetI18nLabel(func(ctx context.Context) string {
			return login_session.GetMessages(ctx).LoginSessions
		}).
		BuildComponentFunc(c.loginSession)

	d.Action(signOutAllSession).
		Icon("mdi-logout-variant").
		SetI18nLabel(func(ctx context.Context) string {
			return login_session.GetMessages(ctx).SignOutAllOtherSessions
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) (err error) {
			u := user.GetCurrentUser(ctx.R)

			if err = c.mgr.ExpireOtherSessionLogs(c.db, ctx.R, u.GetID()); err != nil {
				return err
			}

			ctx.FlashMessage(login_session.GetMessages(ctx.Context()).SignOutAllSuccessfullyTips)
			return
		})

	user.ConfigureChangePasswordAction(c.lb, true, d, func(ctx *web.EventContext) (user.User, error) {
		return user.GetCurrentUser(ctx.R), nil
	})
	return
}

func (c *Builder) getAvatarShortName(u user.User) string {
	name := u.GetName()
	if name == "" {
		name = u.GetAccountName()
	}
	if rs := []rune(name); len(rs) > 1 {
		name = string(rs[:1])
	}

	return strings.ToUpper(name)
}

func (c *Builder) loginSession(_ string, ctx *web.EventContext, cb *presets.ContentComponentBuilder) (err error) {
	u := user.GetCurrentUser(ctx.R)
	if u == nil {
		return errors.New("no user found")
	}
	cb.Body, err = c.mgr.Sessions(c.db, ctx, u.GetID())

	msgr := login_session.GetMessages(ctx.Context())

	cb.AddMenu(VList(
		VListItem(h.Text(msgr.SignOutAllOtherSessions)).
			Attr("@click",
				web.Plaid().
					EventFunc("presets_DoAction").
					Query(presets.ParamAction, signOutAllSession).
					Go()),
	))

	return
}

type Profile struct {
	ID          uint `admin:"ro"`
	Name        string
	AccountName string `admin:"ro"`
	Status      string `admin:"ro"`
	Roles       string `admin:"ro"`
}
