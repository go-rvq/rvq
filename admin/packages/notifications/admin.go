package notifications

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/comments"
	"github.com/go-rvq/rvq/admin/packages/perms"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelID is the id of the listing of the notifications of the user (under
// their profile: my-profile/notifications).
const ModelID = "notifications"

// filterRead is the filter of the tabs: "no" unread, "yes" read.
const filterRead = "read"

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
	// Users is the model of the users: given, who may list them sends a
	// notification to them (the action ActionSend of the listing)
	Users *presets.ModelBuilder
	// Profile is the model of the profile of the user (rvq profile): the
	// notifications go under it — the user browses all of theirs there;
	// nil, a model of their own, out of the menu.
	Profile *presets.ModelBuilder
	// Mailer, when given, mails a message of a user to who received it too
	// (Deliver): its answer by mail comes back to the conversation by the
	// app's inbox (Answer)
	Mailer func(ctx context.Context, m *Mail) error
	// ReplyFiles says whether the answer of the user of the request takes
	// files (nil: it does)
	ReplyFiles func(ctx *web.EventContext) bool
	// ValidateFile and ValidateMessage, when given, check the files of a
	// message sent or answered (rvq comments': MaxFileSize, MaxTotalSize);
	// the key, the notification answered ("" a new message)
	ValidateFile    func(ctx *web.EventContext, key string, f *comments.Attachment) error
	ValidateMessage func(ctx *web.EventContext, w *comments.Write) error
	p               *presets.Builder
	conv            *comments.Builder
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
	b.p = p
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
	mb.URIName("notifications").MenuIcon("mdi-bell-outline")
	if b.Profile != nil {
		b.Profile.AddChild(mb)
	}
	b.mb = mb
	mb.SetCreatingDisabled(true)
	mb.SetEditingDisabled(true)
	lb := mb.Listing("Level", "Title", "Body", "Sender", "CreatedAt", "ReadAt").SearchColumns("title", "body").OrderBy("created_at desc")
	lb.Field("CreatedAt").ComponentFunc(presets.TDDateTimeComponentFunc)
	// who sent it: the login of the user, or the system
	lb.Field("Sender").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(b.sender(ctx.Context(), field.Obj.(*Notification)))
	})
	// the importance: its icon in a column (its name the tooltip), the row
	// of its color
	lb.Field("Level").SetHiddenLabel(true)
	lb.Field("Level").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Td(levelIcon(ctx.Context(), field.Obj.(*Notification).Level)).Style("width: 1%; white-space: nowrap")
	})
	lb.RowWrapperFunc(func(row h.MutableAttrHTMLComponent, _ string, obj any, _ string, _ *web.EventContext) h.HTMLComponent {
		if n, ok := obj.(*Notification); ok {
			row.SetAttr("style", levelRowStyle(n.Level))
			row.SetAttr("data-notification-level", LevelColor(n.Level))
		}
		return row
	})
	lb.Field("ReadAt").ComponentFunc(presets.TDDateTimeComponentFunc)
	// the tabs: unread (the default), read, all
	lb.FilterTabsFunc(func(ctx *web.EventContext) []*presets.FilterTab {
		m := GetMessages(ctx.Context())
		return []*presets.FilterTab{
			{ID: "unread", Label: m.Unread, Query: url.Values{filterRead: {"no"}}, Default: true},
			{ID: "read", Label: m.Read, Query: url.Values{filterRead: {"yes"}}},
			{ID: "all", Label: m.All, Query: url.Values{filterRead: {}}},
		}
	})
	lb.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		m := GetMessages(ctx.Context())
		return []*vx.FilterItem{{Key: filterRead, Invisible: true,
			SQLConditionFunc: func(val, _ string) (bool, string, []any) {
				switch val {
				case "no":
					return true, "user_notifications.read_at IS NULL", nil
				case "yes":
					return true, "user_notifications.read_at IS NOT NULL", nil
				}
				return false, "", nil
			}},
			// the dates, a period up to now, the origin
			{Key: "created", Label: m.NotificationCreatedAt, ItemType: vx.ItemTypeDatetimeRange, SQLCondition: "user_notifications.created_at %s ?"},
			{Key: "period", Label: m.Period, ItemType: vx.ItemTypeSelect,
				Options: []*vx.SelectItem{
					{Text: m.PeriodToday, Value: "today"},
					{Text: m.Period7d, Value: "7d"},
					{Text: m.Period30d, Value: "30d"},
					{Text: m.Period90d, Value: "90d"},
					{Text: m.Period365d, Value: "365d"},
				},
				SQLConditionFunc: func(val, _ string) (bool, string, []any) {
					since, ok := periodSince(val, time.Now().In(presets.TimeLocationOf(ctx)))
					if !ok {
						return false, "", nil
					}
					return true, "user_notifications.created_at >= ?", []any{since}
				}},
			{Key: "level", Label: m.NotificationLevel, ItemType: vx.ItemTypeMultipleSelect, SQLCondition: "user_notifications.level %s ?",
				Options: levelOptions(ctx)},
			{Key: "origin", Label: m.Origin, ItemType: vx.ItemTypeSelect,
				Options: []*vx.SelectItem{
					{Text: m.OriginSystem, Value: "system"},
					{Text: m.OriginUsers, Value: "users"},
					{Text: m.OriginOthers, Value: "others"},
				},
				SQLConditionFunc: func(val, _ string) (bool, string, []any) {
					switch val {
					case "system":
						return true, "user_notifications.system", nil
					case "users":
						return true, "user_notifications.source LIKE ?", []any{sourceUserPrefix + "%"}
					case "others":
						return true, "NOT user_notifications.system AND user_notifications.source NOT LIKE ?", []any{sourceUserPrefix + "%"}
					}
					return false, "", nil
				}},
		}
	})
	// the title and the text in the language of the user (Text)
	text := func(field *presets.FieldContext, ctx *web.EventContext) string {
		n := field.Obj.(*Notification)
		title, body := n.Text(ctx.Context())
		if field.Name == "Title" {
			return title
		}
		return body
	}
	for _, name := range []string{"Title", "Body"} {
		lb.Field(name).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			s := text(field, ctx)
			// rich text: its words, the first
			if field.Name == "Body" && field.Obj.(*Notification).HTML {
				s = comments.PlainText(s)
				if r := []rune(s); len(r) > 140 {
					s = string(r[:140]) + "…"
				}
			}
			return h.Td(h.Text(s))
		})
	}
	dt := mb.Detailing("Title", "Body", "Link", "Level", "Sender", "System", "Source", "CreatedAt", "ReadAt", "Conversation")
	// the conversation of a user's notification: its messages, its answer
	b.conv = b.configureConversation()
	dt.Field("Conversation").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		n := field.Obj.(*Notification)
		if n.System || n.ID == uuid.Nil {
			return nil
		}
		return h.Div(
			h.Div(h.Text(field.Label)).Class("text-subtitle-1 mb-2"),
			b.conv.Inline(ctx, n.ID.String()),
		).Class("mt-4").Attr("data-notification-conversation", true)
	})
	// the attachments, to who received or sent their message
	p.MuxSetup(func(prefix string, mux *http.ServeMux) {
		mux.Handle("GET "+strings.TrimSuffix(prefix, "/")+filesPath+"{id}", FileHandler(b.db))
	})
	// the source: the user as a chip of their login, their id copied by its
	// icon; the system, a chip; another, as it is
	dt.Field("Source").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return vx.VXReadonlyField(b.sourceChip(ctx.Context(), field.Obj.(*Notification))).Label(field.Label)
	})
	dt.Field("Sender").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return vx.VXReadonlyField(b.sender(ctx.Context(), field.Obj.(*Notification))).Label(field.Label)
	})
	dt.Field("Level").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return vx.VXReadonlyField(levelChip(ctx.Context(), field.Obj.(*Notification).Level)).Label(field.Label)
	})
	for _, name := range []string{"Title", "Body"} {
		dt.Field(name).ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			s := text(field, ctx)
			if s == "" {
				return nil
			}
			// rich text (a user's): sanitized
			if field.Name == "Body" && field.Obj.(*Notification).HTML {
				return vx.VXReadonlyField(h.Div(h.RawHTML(comments.SanitizeHTML(s))).Class(web.StdContentClass)).Label(field.Label)
			}
			return vx.VXReadonlyField().Label(field.Label).Value(s)
		})
	}
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
	// deleted softly: the trash of the listing, restored
	perms.SetupTrash(mb, b.db, perms.TrashOptions{TableName: Notification{}.TableName()})
	// read, unread again: an icon of the row (at the left of its menu),
	// confirmed
	mark := func(name, icon string, read bool) {
		label := func(ctx context.Context) string {
			if read {
				return GetMessages(ctx).MarkRead
			}
			return GetMessages(ctx).MarkUnread
		}
		lb.ItemAction(name).Inline(true).Icon(icon).SetI18nLabel(label).
			SetEnabledObj(func(obj any, _ string, _ *web.EventContext) (bool, error) {
				n, ok := obj.(*Notification)
				return ok && (n.ReadAt == nil) == read, nil
			}).
			ComponentFunc(func(_ string, ctx *web.EventContext) (h.HTMLComponent, error) {
				m := GetMessages(ctx.Context())
				q := m.MarkReadConfirm
				if !read {
					q = m.MarkUnreadConfirm
				}
				return h.Div(h.Text(q)), nil
			}).
			UpdateFunc(func(id string, ctx *web.EventContext) error {
				nid, err := uuid.Parse(id)
				if err != nil {
					return err
				}
				if read {
					return MarkRead(b.db, currentUser(ctx), nid)
				}
				return MarkUnread(b.db, currentUser(ctx), nid)
			})
	}
	if b.Users != nil {
		b.configureSend(lb)
	}
	mark("mark_read", "mdi-email-open-outline", true)
	mark("mark_unread", "mdi-email-mark-as-unread", false)

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

// sender is who sent n: the login of the user (SourceUser), the system, or
// its source.
func (b *Builder) sender(ctx context.Context, n *Notification) h.HTMLComponent {
	if n.System {
		return h.Span(GetMessages(ctx).OriginSystem).Class("text-medium-emphasis font-italic")
	}
	id, ok := strings.CutPrefix(n.Source, sourceUserPrefix)
	if !ok {
		return h.Text(n.Source)
	}
	u := b.user(id)
	if u == nil {
		return h.Text(n.Source)
	}
	return h.Span(u.GetAccountName()).Attr("title", u.GetName()).Attr("data-notification-sender", true)
}

// sourceChip is the source of n: the login of the user (SourceUser) as a
// chip, an icon at its right copying their id; the system's, a chip.
func (b *Builder) sourceChip(ctx context.Context, n *Notification) h.HTMLComponent {
	if n.System {
		return v.VChip(h.Text(GetMessages(ctx).OriginSystem)).Size(v.SizeSmall).PrependIcon("mdi-cog-outline")
	}
	id, ok := strings.CutPrefix(n.Source, sourceUserPrefix)
	if !ok {
		return h.Text(n.Source)
	}
	login := id
	if u := b.user(id); u != nil {
		login = u.GetAccountName()
	}
	m := GetMessages(ctx)
	return v.VChip(
		h.Text(login),
		v.VIcon("mdi-content-copy").Size(v.SizeSmall).Class("ms-2").
			Attr("@click.stop", "copyToClipboard("+h.JSONString(id)+")").
			Attr("role", "button", "aria-label", m.CopyUserID, "data-copy-user-id", id).
			Children(v.VTooltip().Activator("parent").Location("top").Text(m.CopyUserID)),
	).Size(v.SizeSmall).PrependIcon("mdi-account-outline").Attr("title", id).Attr("data-notification-source", true)
}

// user is the user of id (nil: none, no Users).
func (b *Builder) user(id string) user.User {
	if b.Users == nil {
		return nil
	}
	u, _ := b.Users.NewModel().(user.User)
	if u == nil || b.db.Where("id = ?", id).Limit(1).Find(u).Error != nil || u.GetID() == uuid.Nil {
		return nil
	}
	return u
}

// levelOptions are the levels of importance, the options of their filter.
func levelOptions(ctx *web.EventContext) (out []*vx.SelectItem) {
	for _, l := range Levels {
		out = append(out, &vx.SelectItem{Text: LevelLabel(ctx.Context(), l), Value: l})
	}
	return
}

// periodSince is when the period val (today, 7d, 30d, 90d, 365d) up to now
// begins.
func periodSince(val string, now time.Time) (time.Time, bool) {
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	switch val {
	case "today":
		return day, true
	case "7d":
		return day.AddDate(0, 0, -6), true
	case "30d":
		return day.AddDate(0, 0, -29), true
	case "90d":
		return day.AddDate(0, 0, -89), true
	case "365d":
		return day.AddDate(0, 0, -364), true
	}
	return time.Time{}, false
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
		title, body := n.Text(ctx.Context())
		children := []h.HTMLComponent{v.VListItemTitle(h.Text(title)), v.VListItemSubtitle(h.Text(body))}
		// its importance: the icon, the color of the theme
		children = append(children, web.Slot(v.VIcon(LevelIcon(n.Level)).Color(LevelColor(n.Level))).Name(v.VSlotPrepend))
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
