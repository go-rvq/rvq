package notifications

import (
	"context"
	"net/http"
	"time"

	"github.com/go-rvq/rvq/admin/dashboard"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/google/uuid"
)

const (
	// GroupNotifications is the group of the widgets of the notifications
	// (Widgets).
	GroupNotifications = "notifications"
	// WidgetUnread is the widget of the unread notifications of the user.
	WidgetUnread = "unread"
)

// widgetAction is the action of a row of the widget: the notification read.
const widgetAction = "read"

// Widgets makes the group GroupNotifications of d, its widget WidgetUnread:
// the unread notifications of the user, the newest, each one read by its
// button (it leaves the list), opened by its link; the listing of all. A
// panel takes it as any group (the profile's: rvq profile Dashboard). After
// Install.
func (b *Builder) Widgets(d *dashboard.Builder) *dashboard.GroupBuilder {
	title := func(ctx context.Context) string { return GetMessages(ctx).Notifications }
	return d.Group(GroupNotifications).Title(title).Icon("mdi-bell-outline").Widget(WidgetUnread, func(w *dashboard.WidgetBuilder) {
		w.Title(title).Size(dashboard.Half).TTL(time.Second).
			Allowed(func(r *http.Request) bool {
				u := user.GetCurrentUser(r)
				return u != nil && !u.Anonymous()
			}).
			Link(func(*web.EventContext) string { return b.listingHref() }).
			Action(widgetAction, func(ctx *web.EventContext, arg string) error {
				id, err := uuid.Parse(arg)
				if err != nil {
					return err
				}
				return MarkRead(b.db, currentUser(ctx), id)
			}).
			Load(func(ctx *web.EventContext, _ dashboard.Period) ([]dashboard.Block, error) {
				uid := currentUser(ctx)
				ns, err := Latest(b.db, uid, bellLimit)
				if err != nil {
					return nil, err
				}
				total, err := Unread(b.db, uid)
				if err != nil {
					return nil, err
				}
				m := GetMessages(ctx.Context())
				list := &dashboard.List{Empty: m.NoUnread, Link: b.listingHref(), LinkText: m.SeeAll}
				loc := presets.TimeLocationOf(ctx)
				for _, n := range ns {
					title, body := n.Text(ctx.Context())
					list.Rows = append(list.Rows, dashboard.Row{
						Text:   title,
						Sub:    body,
						Right:  n.CreatedAt.In(loc).Format("02/01 15:04"),
						Link:   n.Link,
						Level:  dashboardLevel(n.Level),
						Action: &dashboard.Action{Name: widgetAction, Arg: n.ID.String(), Label: m.MarkRead, Icon: "mdi-check", Removes: true},
					})
				}
				return []dashboard.Block{&dashboard.Metric{Value: float64(total), Caption: m.Unread}, list}, nil
			})
	})
}

// listingHref is the listing of the notifications of the user ("": not
// installed).
func (b *Builder) listingHref() string {
	if b.mb == nil {
		return ""
	}
	return b.mb.Info().ListingHref()
}
