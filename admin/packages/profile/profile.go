package profile

import (
	"context"
	"errors"
	"net/http"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/packages/login_session"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
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
	// mb is the model of the profile (Model), once installed
	mb *presets.ModelBuilder
}

// Model is the model of the profile of the user (my_profile), once
// installed: what is the user's own goes under it (AddChild) — their keys,
// their credentials.
func (c *Builder) Model() *presets.ModelBuilder { return c.mb }

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

// Brand is what the top bar shows of the user (presets.ProfileFunc): their
// notifications and their button (UserButton), with its menu.
func (c *Builder) Brand(b *presets.Builder) func(ctx *web.EventContext) h.HTMLComponent {
	button := c.UserButton(b)
	return func(ctx *web.EventContext) h.HTMLComponent {
		var notification h.HTMLComponent
		if u := user.GetCurrentUser(ctx.R); u != nil && !u.Anonymous() && c.notificationComponentFunc != nil {
			notification = c.notificationComponentFunc(ctx, u)
		}
		return h.HTMLComponents{notification, button(ctx)}
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
	c.mb = m

	// the menu of the user: their profile first
	p.AddUserMenuItem(&presets.UserMenuItem{
		Name:    "MyProfile",
		Group:   presets.UserMenuQuick,
		Order:   10,
		Enabled: true,
		Handler: func(ctx *web.EventContext) h.HTMLComponent {
			return VListItem(
				web.Slot(VIcon("mdi-account-outline")).Name(VSlotPrepend),
				VListItemTitle(h.Text(GetMessages(ctx.Context()).MyProfile)),
			).Href(p.URIPrefixCtx(ctx.R.Context()) + "/" + ModelID)
		},
	})
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
		if au, ok := u.(AvatarUser); ok {
			profile.Avatar = au.GetAvatar()
		}
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
		// what the profile edits: the name, the picture
		profile := obj.(*Profile)
		u.SetName(profile.Name)
		if au, ok := u.(AvatarUser); ok {
			au.SetAvatar(profile.Avatar)
		}

		if err = c.db.Save(u).Error; err != nil {
			return
		}

		return nil
	})

	d := m.Detailing()
	d.FetchFunc(fetch)
	configureAvatarField(c.db, e, d)

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
	ID   uuid.UUID `admin:"ro" gorm:"type:uuid"`
	Name string
	// Avatar is the picture of the user, when their model has one
	// (AvatarUser)
	Avatar      media_library.MediaBox
	AccountName string `admin:"ro"`
	Status      string `admin:"ro"`
	Roles       string `admin:"ro"`
}
