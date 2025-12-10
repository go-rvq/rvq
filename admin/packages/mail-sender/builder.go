package mail_sender

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/packages/helper/nested"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/markbates/goth"
	google2 "github.com/markbates/goth/providers/google"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

const (
	oauthCallbackUri           = "EnvioDeEmailPorGmail/google/callback"
	oauthLoginUri              = "EnvioDeEmailPorGmail/google/connect"
	oauthLogoutUri             = "EnvioDeEmailPorGmail/google/logout"
	oauthCallbackSuccess       = "EnvioDeEmailPorGmail/google/success"
	gmailClearCredentialsEvent = "gmail:clear-credentials"
)

type Builder struct {
	db *gorm.DB
}

func New(db *gorm.DB, i18nB *i18n.Builder) *Builder {
	ConfigureMessages(i18nB)
	return &Builder{
		db: db,
	}
}

func (b *Builder) Install(p *presets.Builder) (err error) {
	if err = b.db.AutoMigrate(
		&MailSender{},
	); err != nil {
		return
	}

	mb := Model(p,
		&MailSender{},
		presets.ModelConfig().SetSingleton(true),
	).
		InMenu(true).
		MenuIcon("mdi-email-fast")

	nested.New(mb).
		Field("SMTP").
		Editing(func(b *presets.FieldsBuilder) *presets.FieldsBuilder {
			b.Field("Password").ComponentFunc(presets.PasswordFieldComponentFunc)
			return b
		})

	nested.New(mb).
		Field("Gmail").
		Editing(func(b *presets.FieldsBuilder) *presets.FieldsBuilder {
			b.
				Field("CredentialsFile").
				ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
					return presets.FileFieldComponentFunc(field, ctx).(*v.VFileInputBuilder).
						PrependIcon("mdi-code-json").
						Hint(field.HintLoader()).
						Attr("accept", "application/json").
						PersistentHint(true)
				}).
				SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
					file := ctx.R.MultipartForm.File[field.FormKey]
					if len(file) > 0 {
						o := obj.(*GmailSender)
						var f multipart.File
						if f, err = file[0].Open(); err != nil {
							return
						}
						defer f.Close()

						var b []byte
						if b, err = io.ReadAll(f); err != nil {
							return
						}
						if len(b) > 0 {
							if err = o.Credentials.Scan(b); err == Messages_en_US.ErrGmailSenderCredentialsInvalid {
								err = GetMessages(ctx.Context()).ErrGmailSenderCredentialsInvalid
							}
						} else {
							err = GetMessages(ctx.Context()).ErrGmailSenderCredentialsInvalid
						}
					}
					return
				})
			b.Append("CredentialsFile")
			return b
		})

	presets.ConfigureSelectField(mb, "Sender", 0, &presets.SelectConfig{
		AvailableKeysFunc: func(ctx *presets.FieldContext) (keys []string) {
			return []string{"GMAIL", "SMTP"}
		},
	})

	d := mb.Detailing()

	var getConfig = func(c *GmailSender, r *http.Request) (config *oauth2.Config, err error) {
		if config, err = c.Config(); err != nil {
			return
		}

		scheme := "https"
		if r.TLS == nil {
			scheme = "http"
		}

		callbackUrl := fmt.Sprintf("%s://%s%s/%s", scheme, r.Host, mb.Info().DetailingHref(model.ID{}), oauthCallbackUri)
		config.RedirectURL = callbackUrl
		return
	}

	{
		type sendMailTestForm struct {
			Sender  string `admin:"required"`
			To      string `admin:"required"`
			Subject string `admin:"required"`
			Message string `admin:"required"`
		}

		mb := NewModel(p, &sendMailTestForm{})
		e := mb.Editing()
		presets.ConfigureSelectField(mb, "Sender", presets.WRITE, &presets.SelectConfig{
			AvailableKeysFunc: func(ctx *presets.FieldContext) (keys []string) {
				s := presets.GetActionFormContext[*sendMailTestForm](ctx.EventContext).Obj.(*MailSender)
				if s.Gmail.IsValid() {
					keys = append(keys, "GMAIL")
				}
				if s.SMTP.IsValid() {
					keys = append(keys, "SMTP")
				}
				return
			},
			KeyLabelsFunc: func(ctx *presets.FieldContext, key []string) []string {
				return key
			},
		})

		presets.ActionForm[*sendMailTestForm](
			d.Action("TestSendMail"),
			e,
			func(ctx *presets.ActionFormContext[*sendMailTestForm]) (err error) {
				var c = ctx.Obj.(*MailSender)

				err = c.SendByMethod(ctx.Form.Sender, Message().
					To(ctx.Form.To).
					Label(ctx.Context.R.Host).
					Subject(ctx.Form.Subject).
					Body(ctx.Form.Message))

				if err == nil {
					ctx.Context.Flash = GetMessages(ctx.Context.Context()).SendMailSuccessfully
				}
				return
			}).
			InitForm(func(fctx *presets.ActionFormContext[*sendMailTestForm]) error {
				msgr := GetMessages(fctx.Context.Context())
				if fctx.Form.Subject == "" {
					fctx.Form.Subject = msgr.TestMessageSubject
				}
				if fctx.Form.Message == "" {
					fctx.Form.Message = msgr.TestMessageBody
				}
				return nil
			}).
			Fetch(true).
			Build()
	}

	reset := func(crendentials bool) error {
		var (
			obj MailSender
			m   = map[string]any{
				"gmail__callback_uri": nil,
				"gmail__user":         nil,
				"gmail__token":        nil,
			}
		)

		if crendentials {
			m["gmail__credentials"] = nil
		}

		return b.db.Session(&gorm.Session{}).Model(&obj).Where("TRUE").UpdateColumns(m).Error
	}

	mb.RegisterEventFunc(gmailClearCredentialsEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		if err = reset(true); err == nil {
			r.AppendRunScript("window.location.href = window.location.href")
		}
		return
	})

	isSenderEnabled := func(f *presets.FieldBuilder) *presets.FieldBuilder {
		return f.WrapEnabled(func(old func(ctx *presets.FieldContext) bool) func(ctx *presets.FieldContext) bool {
			return func(ctx *presets.FieldContext) bool {
				return old(ctx) && ctx.Obj.(*MailSender).Sender == strings.ToUpper(f.Name())
			}
		})
	}

	e := mb.Editing()
	isSenderEnabled(e.Field("Gmail"))
	isSenderEnabled(e.Field("SMTP"))

	d.Field("Gmail").ComponentFunc(
		presets.FieldComponentWrapper(
			func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
				var (
					comps h.HTMLComponents
					o     = field.Obj.(*MailSender).Gmail
					lmsgr = login.GetMessages(ctx.Context())
				)

				if !o.Credentials.IsZero() {
					comps = append(comps, v.VChip().
						Class("me-2").
						Closable(true).
						Text(o.Credentials.Data.ProjectID).
						Attr("@click:close", web.Plaid().URL(ctx.R.URL.Path).EventFunc(gmailClearCredentialsEvent).Go()),
					)

					if o.Token.IsZero() {
						comps = append(comps, v.VBtn(lmsgr.SignInBtn).
							Class("me-2").
							Attr("@click",
								fmt.Sprintf(`(e) => e.view.location.href = "%s/%s"`, ctx.R.URL.Path, oauthLoginUri)))
					} else {
						comps = append(comps,
							h.Text(fmt.Sprintf("%s <%s>", o.User.Data.Name, o.User.Data.Email)),
							v.VBtn(lmsgr.SignOutBtn).Class("ms-4").
								Attr("@click", fmt.Sprintf(`(e) => e.view.location.href = "%s/%s"`, ctx.R.URL.Path, oauthLogoutUri)),
						)
					}
				}

				return h.SimplifyComponent(comps)
			}))

	d.AddPageFunc(perm.PermVerifier(), oauthLoginUri, func(ctx *web.EventContext, obj any, mid model.ID, r *web.PageResponse) (err error) {
		c := obj.(*MailSender).Gmail
		var config *oauth2.Config
		if config, err = getConfig(&c, ctx.R); err != nil {
			return err
		}
		r.RedirectURL = config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		return
	})

	d.AddPageFunc(perm.PermVerifier(), oauthLogoutUri, func(ctx *web.EventContext, obj any, mid model.ID, r *web.PageResponse) (err error) {
		if err = reset(false); err == nil {
			ctx.Flash = GetMessages(ctx.Context()).GmailSenderLogOutSuccessfully
			r.RedirectURL = mb.Info().DetailingHref("")
		}
		return
	})

	d.AddPageFunc(perm.PermVerifier(), oauthCallbackUri, func(ctx *web.EventContext, obj any, mid model.ID, r *web.PageResponse) (err error) {
		var (
			c      = obj.(*MailSender).Gmail
			config *oauth2.Config
		)

		if config, err = getConfig(&c, ctx.R); err != nil {
			return err
		}

		tok, err := config.Exchange(context.TODO(), ctx.R.FormValue("code"))

		if err != nil {
			return fmt.Errorf("Unable to retrieve token from web: %v", err)
		}

		token := &GmailToken{
			AccessToken:  tok.AccessToken,
			TokenType:    tok.TokenType,
			RefreshToken: tok.RefreshToken,
			Expiry:       tok.Expiry,
			ExpiresIn:    tok.ExpiresIn,
			IDToken:      tok.Extra("id_token").(string),
		}

		s := &google2.Session{
			AuthURL:      "",
			AccessToken:  tok.AccessToken,
			RefreshToken: tok.RefreshToken,
			ExpiresAt:    tok.Expiry,
			IDToken:      token.IDToken,
		}

		p := google2.New(config.ClientID, config.ClientSecret, config.RedirectURL, config.Scopes...)

		var user goth.User
		if user, err = p.FetchUser(s); err != nil {
			return err
		}

		var tb, _ = json.Marshal(token)
		var ub, _ = json.Marshal(user)

		if err = b.db.Session(&gorm.Session{}).Model(obj).UpdateColumns(map[string]any{
			"gmail__callback_uri": config.RedirectURL,
			"gmail__user":         string(ub),
			"gmail__token":        string(tb),
		}).Error; err == nil {
			r.RedirectURL = path.Join(ctx.R.URL.Path, "../success")
		}
		return
	})

	d.AddPageFunc(perm.PermVerifier(), oauthCallbackSuccess, func(ctx *web.EventContext, obj any, mid model.ID, r *web.PageResponse) (err error) {
		ctx.Flash = GetMessages(ctx.Context()).GmailSenderConfiguredSuccessfully
		r.RedirectURL = path.Join(ctx.R.URL.Path, "../../..")
		return
	})

	return
}
