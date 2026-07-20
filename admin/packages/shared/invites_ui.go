package shared

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/helper/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// invitesPortal holds the current user's pending share invites.
	invitesPortal      = "sharedInvites"
	invitesLoadEvent   = "shared_invites_load"
	invitesAcceptEvent = "shared_invites_accept"
	invitesRejectEvent = "shared_invites_reject"
	paramInviteID      = "invite_id"
)

// InstallInvites mounts the "my share invites" surface: a portal (menu-top item)
// listing the current user's pending invites with accept/decline buttons, and
// the events backing them. Accepting grants the actual access (see AcceptFor).
// notifier, when given, notifies the inviter on acceptance.
func InstallInvites(b *presets.Builder, db *gorm.DB, notifier ...Notifier) {
	if err := AutoMigrateInvites(db); err != nil {
		panic(err)
	}
	registerMessages(b.I18n())
	var n Notifier
	if len(notifier) > 0 {
		n = notifier[0]
	}

	reloadInvites := func(ctx *web.EventContext) *web.PortalUpdate {
		return &web.PortalUpdate{Name: invitesPortal, Body: invitesBody(db, ctx)}
	}

	wb := b.GetWebBuilder()
	wb.RegisterEventFunc(invitesLoadEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.UpdatePortals = append(r.UpdatePortals, reloadInvites(ctx))
		return
	})
	wb.RegisterEventFunc(invitesAcceptEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		if subject, ok := currentSubject(ctx); ok {
			if id, e := uuid.Parse(ctx.R.FormValue(paramInviteID)); e == nil {
				_ = AcceptFor(db, id, subject, n)
			}
		}
		r.UpdatePortals = append(r.UpdatePortals, reloadInvites(ctx))
		return
	})
	wb.RegisterEventFunc(invitesRejectEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		if subject, ok := currentSubject(ctx); ok {
			if id, e := uuid.Parse(ctx.R.FormValue(paramInviteID)); e == nil {
				_ = RejectFor(db, id, subject)
			}
		}
		r.UpdatePortals = append(r.UpdatePortals, reloadInvites(ctx))
		return
	})

	b.AddMenuTopItemFunc("shared_invites", func(ctx *web.EventContext) h.HTMLComponent {
		if _, ok := currentSubject(ctx); !ok {
			return nil
		}
		return web.Portal(invitesBody(db, ctx)).Name(invitesPortal)
	})
}

// currentSubject is the current user's identifier (the value matched against a
// share invite's Subject).
func currentSubject(ctx *web.EventContext) (string, bool) {
	u := user.GetCurrentUser(ctx.R)
	if u == nil {
		return "", false
	}
	return u.GetID().String(), true
}

// invitesBody renders the current user's pending invites with accept/decline
// actions. It renders nothing when there are none.
func invitesBody(db *gorm.DB, ctx *web.EventContext) h.HTMLComponent {
	subject, ok := currentSubject(ctx)
	if !ok {
		return h.Div()
	}
	invites, err := PendingFor(db, subject)
	if err != nil || len(invites) == 0 {
		return h.Div()
	}
	m := msgs(ctx.Context())

	var rows []h.HTMLComponent
	for _, inv := range invites {
		id := inv.ID.String()
		rows = append(rows, v.VListItem(
			v.VListItemTitle(h.Text(inv.Resource)),
			v.VListItemSubtitle(h.Text(m.InviteFrom+": "+inv.InvitedBy)),
			web.Slot(
				v.VBtn(m.AcceptInvite).Size(v.SizeSmall).Color("primary").Variant(v.VariantTonal).Class("mr-1").
					Attr("@click", web.Plaid().EventFunc(invitesAcceptEvent).Query(paramInviteID, id).Go()),
				v.VBtn(m.RejectInvite).Size(v.SizeSmall).Variant(v.VariantText).
					Attr("@click", web.Plaid().EventFunc(invitesRejectEvent).Query(paramInviteID, id).Go()),
			).Name("append"),
		))
	}
	return v.VList(
		append([]h.HTMLComponent{v.VListSubheader(h.Text(m.MyInvites))}, rows...)...,
	).Density(v.DensityCompact)
}
