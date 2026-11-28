package admin

import (
	"context"

	"github.com/go-rvq/rvq/admin/dashboard"
	orgsadmin "github.com/go-rvq/rvq/admin/packages/orgs/admin"
	"github.com/go-rvq/rvq/admin/packages/people/messages"
	"github.com/go-rvq/rvq/admin/packages/people/models"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
)

// WidgetPeople is the widget of the people of an organization, in the group
// of the organizations' (orgsadmin.GroupOrg).
const WidgetPeople = "people"

// peopleRows are the people the widget lists, the newest.
const peopleRows = 5

// Widgets puts the widget of the people (WidgetPeople) in the group of the
// organizations of d: how many are active, the newest, the link to their
// listing. After Install.
func (b *Builder) Widgets(d *dashboard.Builder) *Builder {
	mb, org := b.personMB, b.OrganizationBuilder.OrganizacaoModel()
	// the href of the person id ("" the listing) of the organization of ctx
	href := func(ctx *web.EventContext, id any) string {
		oid, err := org.ParseRecordID(dashboard.ScopeOf(ctx).Record)
		if err != nil {
			return ""
		}
		return mb.Info().DetailingHref(id, oid)
	}
	listing := func(ctx *web.EventContext) string { return href(ctx, nil) }
	d.GetGroup(orgsadmin.GroupOrg).Widget(WidgetPeople, func(w *dashboard.WidgetBuilder) {
		w.Order(10).Title(func(ctx context.Context) string { return messages.Get(ctx).Persons }).
			Allowed(dashboard.ListAllowed(mb)).Link(listing).
			Load(func(ctx *web.EventContext, _ dashboard.Period) ([]dashboard.Block, error) {
				m := messages.Get(ctx.Context())
				q := func() *gorm.DB {
					return b.db.WithContext(ctx.R.Context()).Model(&models.Person{}).
						Where("organizacao_id = ?", dashboard.ScopeOf(ctx).Record)
				}
				var active int64
				if err := q().Where("active").Count(&active).Error; err != nil {
					return nil, err
				}
				var ps []models.Person
				if err := q().Order("created_at DESC").Limit(peopleRows).Find(&ps).Error; err != nil {
					return nil, err
				}
				list := &dashboard.List{Empty: m.NoPersons, Link: listing(ctx)}
				for _, p := range ps {
					list.Rows = append(list.Rows, dashboard.Row{Text: p.Name, Sub: p.Document, Link: href(ctx, p.ID)})
				}
				return []dashboard.Block{&dashboard.Metric{Value: float64(active), Caption: m.PersonsCaption}, list}, nil
			})
	})
	return b
}
