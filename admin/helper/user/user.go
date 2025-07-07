package user

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/qor5/admin/v3/model"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/gorm2op"
	"github.com/qor5/admin/v3/role"
	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/vue"
	"github.com/qor5/x/v3/login"
	"github.com/qor5/x/v3/perm"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"
)

type Builder struct {
	db *gorm.DB
	lb *login.Builder
	mb *presets.ModelBuilder

	ExpireAllSessionLogs  func(db *gorm.DB, userID uint) (err error)
	LoginInitialUserEmail string
	Roles                 []string
	UserManagerRoles      []string
}

func (b *Builder) UserRoles(u User) (roles []string) {
	if u == nil {
		return
	}
	if u.GetAccountName() == b.lb.GetInitialUserAccount() {
		return []string{RoleAdministrador}
	} else {
		return append(u.GetRoles().Names(), RoleLogged)
	}
}

func (b *Builder) AppendRoles(roles ...string) *Builder {
	b.Roles = append(b.Roles, roles...)
	return b
}

func FindRoles(db *gorm.DB, u User) (roles []*role.Role) {
	db.Model(u).Association("Roles").Find(&roles)
	return
}

func New(db *gorm.DB, lb *login.Builder, mb *presets.ModelBuilder, loginInitialUserEmail string) (c *Builder) {
	ConfigureMessages(mb.Builder().I18n())

	c = &Builder{
		db:                    db,
		lb:                    lb,
		mb:                    mb,
		LoginInitialUserEmail: loginInitialUserEmail,
		Roles:                 []string{RoleAdministrador},
		UserManagerRoles:      []string{},
	}

	mb.Listing().SearchFunc(func(model interface{}, params *presets.SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error) {
		qdb := db

		// If the current user doesn't has 'admin' role, do not allow them to view admin and manager users
		// We didn't do this on permission because of we are not supporting the permission on listing page
		if currentRoles := c.UserRoles(GetCurrentUser(ctx.R)); !slices.Contains(currentRoles, RoleAdministrador) {
			qdb = db.Joins("inner join user_role_join urj on users.id = urj.user_id inner join roles r on r.id = urj.role_id").
				Group("users.id").
				Having("COUNT(CASE WHEN r.name in (?) THEN 1 END) = 0", append(c.UserManagerRoles, RoleAdministrador))
		}

		return gorm2op.DataOperator(qdb).Search(model, params, ctx)
	})

	ed := mb.Editing(
		"Actions",
		"Name",
		"Account",
		"Roles",
		"Status",
	)

	ed.Validators.AppendFunc(func(obj interface{}, _ presets.FieldModeStack, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(User)
		if u.GetAccountName() == "" {
			err.FieldError("Account", GetMessages(ctx.Context()).ErrorAccountRequired)
		}
		return
	})

	getUserOrError := func(ctx *web.EventContext) (u User, err error) {
		id := ctx.R.FormValue("id")
		if id == "" {
			err = presets.MustGetMessages(ctx.Context()).ErrEmptyParamID
			return
		}
		u = c.mb.NewModel().(User)
		err = db.Where("id = ?", id).First(&u).Error
		return
	}

	mb.RegisterEventFunc("eventUnlockUser", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var u User
		if u, err = getUserOrError(ctx); err != nil {
			return
		}
		if err = u.UnlockUser(db, u); err != nil {
			return r, err
		}
		ctx.Flash = GetMessages(ctx.Context()).UserUnlockedSuccessfully
		ed.UpdateOverlayContent(ctx, &r, &u, "", nil)
		return r, nil
	})

	mb.RegisterEventFunc("eventSendResetPasswordEmail", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		var u User
		if u, err = getUserOrError(ctx); err != nil {
			return
		}
		if _, err = c.lb.SendResetPasswordLink(ctx.Request(), u); err != nil {
			return r, err
		}
		ctx.Flash = GetMessages(ctx.Context()).MailSentSuccessfully
		return
	})

	if c.ExpireAllSessionLogs != nil {
		if _, ok := c.mb.Model().(login.SessionSecureUserPasser); ok {
			mb.RegisterEventFunc("eventRevokeTOTP", func(ctx *web.EventContext) (r web.EventResponse, err error) {
				var u User
				if u, err = getUserOrError(ctx); err != nil {
					return
				}
				err = login.RevokeTOTP(u.(login.SessionSecureUserPasser), db, u, fmt.Sprint(u.GetID()))
				if err != nil {
					return r, err
				}
				err = c.ExpireAllSessionLogs(db, u.GetID())
				if err != nil {
					return r, err
				}
				ctx.Flash = GetMessages(ctx.Context()).AllSessionLogsExpiredSuccessfully
				ed.UpdateOverlayContent(ctx, &r, u, "", nil)
				return r, nil
			})
		}
	}

	ed.Field("Account").
		Label("Email").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		u, ok := field.Obj.(User)
		if !ok {
			return nil
		}
		if u.GetAccountName() == c.LoginInitialUserEmail {
			return nil
		}
		return VTextField().Attr(web.VField(field.Name, field.Value())...).Label(field.Label).ErrorMessages(field.Errors...)
	}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
		u := obj.(User)
		email := ctx.R.FormValue(field.Name)
		u.SetEmail(email)
		return nil
	})

	mb.DeletingRestriction.ObjHandler(presets.OkObjHandlerFuncT(func(u User, ctx *web.EventContext) (ok, handled bool) {
		if u.GetAccountName() == c.LoginInitialUserEmail {
			return true, true
		}
		return
	}))

	ed.Field("Roles").
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			u := field.Obj.(User)
			if u.GetAccountName() == c.LoginInitialUserEmail {
				return nil
			}
			values := []string{}
			var roles []*role.Role
			db.Model(u).Association("Roles").Find(&roles)
			for _, r := range roles {
				values = append(values, fmt.Sprint(r.ID))
			}

			roles = nil
			db.Find(&roles)
			allRoleItems := []DefaultOptionItem{}
			for _, r := range roles {
				allRoleItems = append(allRoleItems, DefaultOptionItem{
					Text:  r.Name,
					Value: fmt.Sprint(r.ID),
				})
			}

			return vue.FormField(
				vx.VXSelectMany().
					Label(field.Label).
					ItemText("text").ItemValue("value").
					Many(true).
					Chips(true).
					Items(allRoleItems).
					ErrorMessages(field.Errors...),
			).
				Value(field.Name, values).
				Bind()
		}).
		SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
			u := obj.(User)
			if u.GetAccountName() == c.LoginInitialUserEmail {
				return perm.PermissionDenied
			}
			rids := ctx.FormSliceValues(field.FormKey)
			var roles []*role.Role
			for _, id := range rids {
				uid, err1 := strconv.Atoi(id)
				if err1 != nil {
					continue
				}
				roles = append(roles, &role.Role{
					ID: uint(uid),
				})
			}

			if u.GetID() == 0 {
				err = reflectutils.Set(obj, field.Name, roles)
			} else {
				err = db.Model(u).Association(field.Name).Replace(roles)
			}
			if err != nil {
				return
			}
			return
		})

	ed.Field("Status").
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			u := field.Obj.(User)
			if u.GetAccountName() == c.LoginInitialUserEmail {
				return nil
			}
			return vx.VXSelectOne().
				Label(field.Label).
				ItemText("text").ItemValue("value").
				Attr(web.VField(field.Name, u.GetStatus())...).
				Chips(true).
				Items([]DefaultOptionItem{
					{Value: "active", Text: "Ativo"},
					{Value: "inactive", Text: "Inativo"},
				}).ErrorMessages(field.Errors...)
		})

	ed.WrapSaveFunc(func(in presets.SaveFunc) presets.SaveFunc {
		return func(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
			u := obj.(User)
			if u.GetAccountName() == c.LoginInitialUserEmail {
				return perm.PermissionDenied
			}
			u.SetRegistrationDate(time.Now())
			return in(obj, id, ctx)
		}
	})

	ed.Only("Name", "Status", "Account", "Roles", "CreatedAt", "UpdatedAt")

	d := mb.Detailing()
	d.Field("Roles").ComponentFunc(presets.FieldComponentWrapper(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		var (
			u     = field.Obj.(User)
			roles []*role.Role
		)
		db.Model(u).Association("Roles").Find(&roles)
		u.SetRoles(roles)
		return h.Text(strings.Join(c.UserRoles(u), ", "))
	}))
	d.Only("Name", "Status", "Account", "Roles", "CreatedAt", "UpdatedAt")

	ConfigureChangePasswordAction(c.lb, false, d, func(ctx *web.EventContext) (u User, err error) {
		fctx := presets.GetActionFormContext[*ChangePassword](ctx)
		u = c.mb.NewModel().(User)
		err = db.First(u, "id = ?", fctx.ID).Error
		return
	})

	d.
		Action("send_reset_password_email").
		Icon("mdi-email-arrow-right").
		SetEnabledObj(func(obj any, id string, ctx *web.EventContext) (ok bool, err error) {
			return strings.Contains(obj.(User).GetAccountName(), "@"), nil
		}).
		SetI18nLabel(func(ctx context.Context) string {
			return "Enviar email para alterar a senha"
		}).
		OnClick(func(ctx *web.EventContext, id string, obj any) string {
			return web.Plaid().EventFunc("eventSendResetPasswordEmail").
				Query("id", id).Go()
		})

	d.
		Action("unlock").
		SetEnabledObj(func(obj any, id string, ctx *web.EventContext) (ok bool, err error) {
			return obj.(User).GetLocked(), nil
		}).
		SetI18nLabel(func(ctx context.Context) string {
			return "Desbloquear"
		}).
		OnClick(func(ctx *web.EventContext, id string, obj any) string {
			return web.Plaid().EventFunc("eventUnlockUser").
				Query("id", id).Go()
		})

	d.
		Action("revoke_totp").
		SetEnabledObj(func(obj any, id string, ctx *web.EventContext) (ok bool, err error) {
			return obj.(User).GetIsTOTPSetup(), nil
		}).
		SetI18nLabel(func(ctx context.Context) string {
			return "Revogar TOTP"
		}).
		OnClick(func(ctx *web.EventContext, id string, obj any) string {
			return web.Plaid().EventFunc("eventRevokeTOTP").
				Query("id", id).Go()
		})

	cl := mb.Listing("ID", "Name", "Account", "Status", "Notes").PerPage(10)
	cl.Field("Account").Label("Email")
	cl.SearchColumns("users.Name", "Account")

	cl.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		msgr := GetMessages(ctx.Context())

		return []*vx.FilterItem{
			{
				Key:          "created",
				Label:        msgr.UserCreatedAt,
				ItemType:     vx.ItemTypeDatetimeRange,
				SQLCondition: `users.created_at %s ?`,
			},
			{
				Key:          "name",
				Label:        msgr.UserName,
				ItemType:     vx.ItemTypeString,
				SQLCondition: `users.name %s ?`,
			},
			{
				Key:          "status",
				Label:        msgr.UserStatus,
				ItemType:     vx.ItemTypeSelect,
				SQLCondition: `users.status %s ?`,
				Options: []*vx.SelectItem{
					{Text: msgr.Active, Value: "active"},
					{Text: msgr.Inactive, Value: "inactive"},
				},
			},
			{
				Key:          "registration_date",
				Label:        msgr.UserRegistrationDate,
				ItemType:     vx.ItemTypeDate,
				SQLCondition: `users.registration_date %s ?`,
				Folded:       true,
			},
			{
				Key:          "registration_date_range",
				Label:        msgr.UserRegistrationDateRange,
				ItemType:     vx.ItemTypeDateRange,
				SQLCondition: `users.registration_date %s ?`,
				Folded:       true,
			},
		}
	})

	return
}
