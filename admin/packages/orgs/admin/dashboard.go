package admin

import (
	"context"
	"sort"

	"github.com/go-rvq/rvq/admin/dashboard"
	"github.com/go-rvq/rvq/admin/packages/orgs/messages"
	"github.com/go-rvq/rvq/admin/packages/shared"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

const (
	// PanelOrg is the panel of an organization: a page of its detail
	// (PainelPath), its widgets of the organization it is of
	// (dashboard.ScopeOf).
	PanelOrg = "org"
	// PainelPath is the path of the panel in the detail of an organization.
	PainelPath = "/painel"
	// GroupOrg is the group of the widgets of an organization, of its panel
	// alone (Records) — other packages add theirs (d.GetGroup(GroupOrg)):
	// the people's.
	GroupOrg = "orgs"
	// WidgetMembers: who the organization is shared with, their roles.
	WidgetMembers = "members"
	// WidgetPending: the sharings waiting for the owner or for whom was
	// invited.
	WidgetPending = "pending"
	// WidgetActivity: the latest events of the sharing.
	WidgetActivity = "activity"
)

// activityRows are the events of the activity shown.
const activityRows = 8

// Dashboard makes the panel of the organizations of d (PanelOrg), a page of
// their detail, with the group GroupOrg — its widgets of the sharing
// (members, what waits, the activity) when the organizations are shared
// (Sharing). Other packages put theirs in the panel (d.GetPanel(PanelOrg))
// or in the group. After Install.
func (b *Builder) Dashboard(d *dashboard.Builder) *dashboard.Panel {
	mb := b.OrganizacaoModel()
	// of an organization alone: only on its panel (never the home's)
	g := d.Group(GroupOrg).Records(mb).Icon("mdi-domain").Order(-10).
		Title(func(ctx context.Context) string { return messages.Get(ctx).WidgetsGroup })
	if s := b.sharing; s != nil {
		orgWidgets(g, mb, s)
	}
	pn := d.Panel(PanelOrg).Group(g).
		Title(func(ctx context.Context) string { return messages.Get(ctx).Painel }).
		Record(mb, PainelPath)
	// in the menu, under the organization chosen: the first of its pages
	cfg := mb.MenuRecordSelectorConfig()
	cfg.Pages = append([]string{PainelPath}, cfg.Pages...)
	return pn
}

// orgWidgets are the widgets of the sharing of the organizations of mb.
func orgWidgets(g *dashboard.GroupBuilder, mb *presets.ModelBuilder, s *shared.Service) {
	detail := func(ctx *web.EventContext) string {
		return mb.Info().DetailingHref(dashboard.ScopeOf(ctx).Record)
	}
	g.Widget(WidgetMembers, func(w *dashboard.WidgetBuilder) {
		w.Order(1).Title(func(ctx context.Context) string { return messages.Get(ctx).Membros }).Link(detail).
			Load(func(ctx *web.EventContext, _ dashboard.Period) ([]dashboard.Block, error) {
				c, m := ctx.R.Context(), messages.Get(ctx.Context())
				acc, err := s.Subjects(c, OrganizacaoModelID, dashboard.ScopeOf(ctx).Record)
				if err != nil {
					return nil, err
				}
				list := &dashboard.List{Empty: m.SemMembros}
				for _, a := range acc {
					role := ""
					if a.Role != nil {
						role = a.Role.Name
					}
					list.Rows = append(list.Rows, dashboard.Row{Text: s.SubjectName(c, a.Subject), Right: role})
				}
				return []dashboard.Block{&dashboard.Metric{Value: float64(len(acc)), Caption: m.MembrosCaption}, list}, nil
			})
	}).Widget(WidgetPending, func(w *dashboard.WidgetBuilder) {
		w.Order(2).Title(func(ctx context.Context) string { return messages.Get(ctx).Pendentes }).Link(detail).
			Load(func(ctx *web.EventContext, _ dashboard.Period) ([]dashboard.Block, error) {
				c, m := ctx.R.Context(), messages.Get(ctx.Context())
				reqs, invs, err := s.Pending(c, OrganizacaoModelID, dashboard.ScopeOf(ctx).Record)
				if err != nil {
					return nil, err
				}
				loc := presets.TimeLocationOf(ctx)
				list := &dashboard.List{Empty: m.PendentesNenhum}
				for _, r := range reqs {
					list.Rows = append(list.Rows, dashboard.Row{Text: s.SubjectName(c, r.Subject), Sub: m.PedidoDono,
						Right: r.CreatedAt.In(loc).Format("02/01 15:04"), Level: dashboard.LevelWarning})
				}
				for _, i := range invs {
					list.Rows = append(list.Rows, dashboard.Row{Text: s.SubjectName(c, i.Subject), Sub: m.ConviteAceite,
						Right: i.CreatedAt.In(loc).Format("02/01 15:04"), Level: dashboard.LevelInfo})
				}
				return []dashboard.Block{list}, nil
			})
	}).Widget(WidgetActivity, func(w *dashboard.WidgetBuilder) {
		w.Order(3).Title(func(ctx context.Context) string { return messages.Get(ctx).Atividade }).Link(detail).
			Load(func(ctx *web.EventContext, _ dashboard.Period) ([]dashboard.Block, error) {
				c, m := ctx.R.Context(), messages.Get(ctx.Context())
				evs, err := s.History(c, OrganizacaoModelID, dashboard.ScopeOf(ctx).Record)
				if err != nil {
					return nil, err
				}
				sort.SliceStable(evs, func(i, j int) bool { return evs[i].At.After(evs[j].At) })
				if len(evs) > activityRows {
					evs = evs[:activityRows]
				}
				loc := presets.TimeLocationOf(ctx)
				list := &dashboard.List{Empty: m.SemAtividade}
				for _, e := range evs {
					sub := ""
					if e.Subject != "" {
						sub = s.SubjectName(c, e.Subject)
					}
					list.Rows = append(list.Rows, dashboard.Row{Text: shared.EventLabel(c, e.Event), Sub: sub,
						Right: e.At.In(loc).Format("02/01 15:04")})
				}
				return []dashboard.Block{list}, nil
			})
	})
}
