package notifications

import (
	"context"
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelID is the id of the listing of the notifications of the user.
const ModelID = "my_notifications"

// bellLimit is how many unread notifications the bell shows.
const bellLimit = 10

const (
	evOpen    = "notifications_open"
	evReadAll = "notifications_read_all"
	paramID   = "notification"
)

// Builder puts the notification center in the admin (Install).
type Builder struct {
	db *gorm.DB
	mb *presets.ModelBuilder
	// Profile is the model of the profile of the user (rvq profile): the
	// notifications go under it — the user browses all of theirs there;
	// nil, a model of their own, out of the menu.
	Profile *presets.ModelBuilder
}

// New is the notification center of the notifications of db.
func New(db *gorm.DB) *Builder { return &Builder{db: db} }

// Model is the listing of the notifications of the user, once installed.
func (b *Builder) Model() *presets.ModelBuilder { return b.mb }

// currentUser is the user of the request (uuid.Nil: none).
func currentUser(ctx *web.EventContext) uuid.UUID {
	if u := user.GetCurrentUser(ctx.R); u != nil && !u.Anonymous() {
		return u.GetID()
	}
	return uuid.Nil
}

// Install puts in p the notifications of the user: their unread ones in the
// bell (a source of the notification center: opened, read), all of them in
// their listing (My notifications: theirs only, read only).
func (b *Builder) Install(p *presets.Builder) error {
	if err := AutoMigrate(b.db); err != nil {
		return err
	}
	ConfigureMessages(p.I18n())
	if permB := p.GetPermission(); permB != nil {
		permB.CreatePolicies(perm.PolicyFor(user.RoleLogged).WhoAre(perm.Allowed).ToDo(perm.Anything).On("*:" + ModelID + ":*"))
	}
	config := presets.ModelConfig().SetId(ModelID).SetModuleKey(MessagesKey)
	var mb *presets.ModelBuilder
	if b.Profile != nil {
		mb = presets.NewModelBuilder(p, &Notification{}, config)
	} else {
		mb = p.Model(&Notification{}, config).InMenu(false)
	}
	mb.URIName("my-notifications").MenuIcon("mdi-bell-outline")
	if b.Profile != nil {
		b.Profile.AddChild(mb)
	}
	b.mb = mb
	mb.SetCreatingDisabled(true)
	mb.SetEditingDisabled(true)
	mb.Listing("CreatedAt", "Title", "Body", "ReadAt").SearchColumns("title", "body").OrderBy("created_at desc")
	mb.Detailing("Title", "Body", "Link", "Level", "Source", "CreatedAt", "ReadAt")
	// theirs only
	mb.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).
			WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
				cb.Pre(func(state *gorm2op.CallbackState) error {
					state.DB = state.DB.Where("user_notifications.user_id = ?", currentUser(state.Ctx))
					return nil
				})
			})
	})
	// opened: read, where it takes
	mb.Detailing().Action("open_notification").Icon("mdi-open-in-new").
		SetI18nLabel(func(ctx context.Context) string { return GetMessages(ctx).Open }).
		SetEnabledObj(func(obj any, _ string, _ *web.EventContext) (bool, error) {
			n, ok := obj.(*Notification)
			return ok && n.Link != "", nil
		}).
		OnClick(func(ctx *web.EventContext, id string, _ any) string {
			return web.Plaid().EventFunc(evOpen).Query(paramID, id).Go()
		})

	wb := p.GetWebBuilder()
	wb.RegisterEventFunc(evOpen, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		id, err := uuid.Parse(ctx.R.FormValue(paramID))
		if err != nil {
			return r, err
		}
		uid := currentUser(ctx)
		var n Notification
		if err := b.db.Where("id = ? AND user_id = ?", id, uid).First(&n).Error; err != nil {
			return r, err
		}
		if err := MarkRead(b.db, uid, id); err != nil {
			return r, err
		}
		if n.Link != "" {
			r.RedirectURL = n.Link
		} else {
			r.Reload = true
		}
		return
	})
	wb.RegisterEventFunc(evReadAll, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		if err := MarkAllRead(b.db, currentUser(ctx)); err != nil {
			return r, err
		}
		r.Reload = true
		return
	})

	p.AddNotificationSource(presets.NotificationSource{
		Name: "notifications",
		Count: func(ctx *web.EventContext) int {
			n, _ := Unread(b.db, currentUser(ctx))
			return int(n)
		},
		Content: func(ctx *web.EventContext) h.HTMLComponent { return b.bell(p, ctx) },
	})
	return nil
}

// bell is the content of the notifications in the bell: the newest unread,
// each opened (read) by a click; all marked read; the listing (nil: none
// unread).
func (b *Builder) bell(p *presets.Builder, ctx *web.EventContext) h.HTMLComponent {
	uid := currentUser(ctx)
	total, _ := Unread(b.db, uid)
	if total == 0 {
		return nil
	}
	ns, _ := Latest(b.db, uid, bellLimit)
	m := GetMessages(ctx.Context())
	items := []h.HTMLComponent{v.VListSubheader(h.Text(fmt.Sprintf(m.BellTitle, total)))}
	for _, n := range ns {
		children := []h.HTMLComponent{v.VListItemTitle(h.Text(n.Title)), v.VListItemSubtitle(h.Text(n.Body))}
		if color := map[string]string{LevelWarning: "warning", LevelError: "error"}[n.Level]; color != "" {
			children = append(children, web.Slot(v.VIcon("mdi-alert-circle-outline").Color(color)).Name(v.VSlotPrepend))
		}
		item := v.VListItem(children...).Lines(2).
			Attr("@click", web.Plaid().EventFunc(evOpen).Query(paramID, n.ID.String()).Go()).
			Attr("data-notification", n.ID.String())
		items = append(items, item)
	}
	items = append(items, v.VDivider(),
		v.VListItem(v.VListItemTitle(h.Text(m.MarkAllRead))).Attr("@click", web.Plaid().EventFunc(evReadAll).Go()).
			Attr("data-notifications-read-all", true))
	if b.mb != nil {
		items = append(items, v.VListItem(v.VListItemTitle(h.Text(m.SeeAll))).Href(b.mb.Info().ListingHref()))
	}
	return v.VList(items...).Density(v.DensityCompact)
}
