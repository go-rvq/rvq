package admin

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/orgs/messages"
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

// The organization in the side menu: a root item, "Organization", with a chip
// — the organization chosen, or "Select" — that opens the dialog of the
// user's organizations; once one is chosen, the item has, opened, its menu:
// the resources nested under it the user may list (by group).
//
// The one chosen is the one of the page (/orgs/{id}/…), else the one kept in
// the browser (the cookie OrgCookie, set by the dialog), else the only one of
// the user, else the default of the user (Organizacao.Padrao).

const (
	// OrgCookie keeps, in the browser, the organization chosen.
	OrgCookie = "rvq_org"
	// orgMenuValue is the value of the root item (a group: not a link).
	orgMenuValue = "g:orgs"
)

// userOrgs are the organizations of the current user (the ones they own).
func userOrgs(db *gorm.DB, r *http.Request) (orgs []models.Organizacao) {
	u := user.GetCurrentUser(r)
	if u == nil || u.Anonymous() {
		return nil
	}
	db.Where("owner_id = ?", u.GetID()).Order("nome").Find(&orgs)
	return
}

// chosenOrg is the organization chosen, among orgs: the one of the page,
// else the cookie's, else the only one, else the default; nil none.
func chosenOrg(b *presets.Builder, r *http.Request, orgs []models.Organizacao) *models.Organizacao {
	id := ""
	base := b.GetURIPrefix() + "/" + OrganizacaoURIName + "/"
	if rest, ok := strings.CutPrefix(r.URL.Path, base); ok {
		id, _, _ = strings.Cut(rest, "/")
	}
	if id == "" {
		if c, err := r.Cookie(OrgCookie); err == nil {
			id = c.Value
		}
	}
	for i := range orgs {
		if orgs[i].ID.String() == id {
			return &orgs[i]
		}
	}
	if len(orgs) == 1 {
		return &orgs[0]
	}
	for i := range orgs {
		if orgs[i].Padrao {
			return &orgs[i]
		}
	}
	return nil
}

// orgMenuRoot is the root item of the organization (presets.MenuRootFunc).
func orgMenuRoot(b *presets.Builder, db *gorm.DB) presets.MenuRootFunc {
	return func(ctx *web.EventContext) (nodes []*presets.MenuNode, opened []string, active string) {
		orgs := userOrgs(db, ctx.R)
		if len(orgs) == 0 {
			return nil, nil, ""
		}
		msgs := messages.Get(ctx.Context())
		org := chosenOrg(b, ctx.R, orgs)
		root := &presets.MenuNode{
			Title: msgs.Organizacao,
			Value: orgMenuValue,
			Props: map[string]any{"prependIcon": "mdi-domain"},
		}
		chip := &presets.MenuChip{Title: msgs.SelecionarOrganizacao, Icon: "mdi-swap-horizontal",
			Hint: msgs.TrocarOrganizacao, Color: "grey",
			// the action ActionChoose of the Organizacao, its dialog
			Event: actions.OpenActionDialog,
			Query: url.Values{
				presets.ParamAction:  {ActionChoose},
				presets.ParamOverlay: {string(actions.Dialog)},
			}.Encode()}
		if mb := OrganizacaoModel(b); mb != nil {
			chip.URL = mb.Info().ListingHref()
		}
		if org != nil {
			chip.Title, chip.Color = org.Nome, "primary"
			root.Children, active = orgResourceNodes(b, ctx, org)
			// open, with the group of the page shown
			opened = append(opened, orgMenuValue)
			for _, n := range root.Children {
				if active != "" && strings.HasPrefix(n.Value, "g:") && hasValue(n.Children, active) {
					opened = append(opened, n.Value)
				}
			}
		}
		root.SetChip(chip)
		return []*presets.MenuNode{root}, opened, active
	}
}

func hasValue(nodes []*presets.MenuNode, value string) bool {
	for _, n := range nodes {
		if n.Value == value || hasValue(n.Children, value) {
			return true
		}
	}
	return false
}

// orgResourceNodes are the nodes of the resources nested under org the user
// may list (in it), by their menu group, and the one of the page shown.
func orgResourceNodes(b *presets.Builder, ctx *web.EventContext, org *models.Organizacao) (nodes []*presets.MenuNode, active string) {
	orgMb := OrganizacaoModel(b)
	if orgMb == nil {
		return nil, ""
	}
	orgID := org.ID.String()
	mid, err := orgMb.ParseRecordID(orgID)
	if err != nil || mid.IsZero() {
		return nil, ""
	}
	prefix := b.GetURIPrefix()
	var order []string
	groups := map[string][]*presets.MenuNode{}
	for _, c := range orgMb.Children() {
		if c.Permissioner().Lister(ctx.R, mid).Denied() {
			continue
		}
		g := c.MenuGroupName()
		seg := c.UriName()
		if g != "" {
			seg = g + "/" + seg
		}
		href := prefix + "/" + OrganizacaoURIName + "/" + orgID + "/" + seg
		icon := c.GetMenuIcon()
		if icon == "" {
			icon = "mdi-circle-small"
		}
		n := &presets.MenuNode{Title: c.TTitlePlural(ctx.Context()), Value: href,
			Props: map[string]any{"href": href, "prependIcon": icon}}
		if ctx.R.URL.Path == href || strings.HasPrefix(ctx.R.URL.Path, href+"/") {
			active = href
		}
		if _, ok := groups[g]; !ok {
			order = append(order, g)
		}
		groups[g] = append(groups[g], n)
	}
	for _, g := range order {
		if g == "" {
			nodes = append(nodes, groups[g]...)
			continue
		}
		nodes = append(nodes, &presets.MenuNode{
			Title:    groupTitle(b, ctx, g),
			Value:    "g:orgs:" + g,
			Props:    map[string]any{"prependIcon": groupIcon(b, g)},
			Children: groups[g],
		})
	}
	return
}

// groupTitle is the title of the menu group g (the one the module named), its
// name capitalized when it has none.
func groupTitle(b *presets.Builder, ctx *web.EventContext, g string) string {
	if mg := b.GetMenuGroup(g); mg != nil {
		if t := mg.TTitle(ctx.Context()); t != "" && t != g {
			return t
		}
	}
	return capitalize(g)
}

func groupIcon(b *presets.Builder, g string) string {
	if mg := b.GetMenuGroup(g); mg != nil && mg.GetIcon() != "" {
		return mg.GetIcon()
	}
	return "mdi-folder-outline"
}

// ActionChoose is the listing action of the Organizacao that chooses the
// organization of the menu: its dialog lists the user's organizations —
// searched, by pages —; one clicked is kept (OrgCookie) and its page opened.
const ActionChoose = "choose_organization"

const (
	// orgChooseListEvent draws the list of the dialog again: the search
	// (paramChooseQ) and the page (paramChoosePage).
	orgChooseListEvent  = "orgs_choose_list"
	orgChooseListPortal = "orgsChooseList"
	paramChooseQ        = "orgs_q"
	paramChoosePage     = "orgs_page"
	paramChooseOrg      = "orgs_org"
	// choosePageSize is how many organizations a page of the dialog has.
	choosePageSize = 10
)

// chooseList is the list of the dialog: the organizations of the user whose
// name or description has q, the page page of them, the one chosen marked —
// a click on one does the action (ActionChoose) with it.
func chooseList(b *presets.Builder, db *gorm.DB, mb *presets.ModelBuilder, ctx *web.EventContext, q string, page int) h.HTMLComponent {
	msgs := messages.Get(ctx.Context())
	u := user.GetCurrentUser(ctx.R)
	if u == nil || u.Anonymous() {
		return h.Div()
	}
	query := db.Model(&models.Organizacao{}).Where("owner_id = ?", u.GetID())
	if q = strings.TrimSpace(q); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("lower(nome) LIKE ? OR lower(descricao) LIKE ?", like, like)
	}
	var total int64
	query.Count(&total)
	pages := int((total + choosePageSize - 1) / choosePageSize)
	if page < 1 {
		page = 1
	}
	if pages > 0 && page > pages {
		page = pages
	}
	var orgs []models.Organizacao
	query.Order("nome").Offset((page - 1) * choosePageSize).Limit(choosePageSize).Find(&orgs)
	chosen := chosenOrg(b, refererRequest(ctx.R), userOrgs(db, ctx.R))

	href := mb.Info().ListingHref()
	var items h.HTMLComponents
	for _, o := range orgs {
		body := h.HTMLComponents{v.VListItemTitle(h.Text(o.Nome))}
		if o.Descricao != "" {
			body = append(body, v.VListItemSubtitle(h.Text(o.Descricao)))
		}
		item := v.VListItem(body...).
			PrependIcon("mdi-domain").
			Attr("@click", web.Plaid().URL(href).EventFunc(actions.DoListingAction).
				Query(presets.ParamAction, ActionChoose).
				Query(paramChooseOrg, o.ID.String()).Go())
		if chosen != nil && chosen.ID == o.ID {
			item.Active(true).Color("primary").AppendIcon("mdi-check")
		}
		items = append(items, item)
	}
	var list h.HTMLComponent = v.VList(items...).Density(v.DensityCompact)
	if len(orgs) == 0 {
		list = h.Div(h.Text(msgs.NenhumaOrganizacao)).Class("text-medium-emphasis pa-4")
	}
	var pager h.HTMLComponent
	if pages > 1 {
		pager = v.VPagination().
			Length(pages).
			ModelValue(page).
			Density(v.DensityCompact).
			TotalVisible(7).
			Attr("@update:model-value", web.Plaid().URL(href).EventFunc(orgChooseListEvent).
				Query(paramChooseQ, q).Query(paramChoosePage, web.Var("$event")).Go())
	}
	return h.Div(list, pager)
}

// refererRequest is r at the page the dialog was opened from (its
// Referer): the organization of that page is the one chosen.
func refererRequest(r *http.Request) *http.Request {
	ref, err := url.Parse(r.Referer())
	if err != nil || ref.Path == "" {
		return r
	}
	rr := r.Clone(r.Context())
	rr.URL = ref
	return rr
}

// configureChoose installs the action ActionChoose on the Organizacao (its
// listing's), and the event of its list.
func configureChoose(b *presets.Builder, db *gorm.DB, mb *presets.ModelBuilder) {
	mb.RegisterEventFunc(orgChooseListEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		page, _ := strconv.Atoi(ctx.R.FormValue(paramChoosePage))
		r.UpdatePortals = append(r.UpdatePortals, &web.PortalUpdate{
			Name: orgChooseListPortal,
			Body: chooseList(b, db, mb, ctx, ctx.R.FormValue(paramChooseQ), page),
		})
		return
	})
	href := func() string { return mb.Info().ListingHref() }
	mb.Listing().Action(ActionChoose).
		Icon("mdi-swap-horizontal").
		SetI18nLabel(func(ctx context.Context) string { return messages.Get(ctx).TrocarOrganizacao }).
		SetDialogWidth("520").
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			msgs := messages.Get(ctx.Context())
			search := v.VTextField().
				Placeholder(msgs.PesquisarOrganizacao).
				PrependInnerIcon("mdi-magnify").
				Clearable(true).
				Density(v.DensityCompact).
				Variant(v.FieldVariantOutlined).
				HideDetails(true).
				Autofocus(true).
				Class("mb-2").
				Attr("@update:model-value", web.Plaid().URL(href()).EventFunc(orgChooseListEvent).
					Query(paramChooseQ, web.Var("$event || ''")).Query(paramChoosePage, 1).Go()).
				Attr(":debounce", 300)
			return h.Div(
				search,
				web.Portal(chooseList(b, db, mb, ctx, "", 1)).Name(orgChooseListPortal),
			), nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			msgs := messages.Get(ctx.Context())
			orgID := ctx.R.FormValue(paramChooseOrg)
			var org *models.Organizacao
			for _, o := range userOrgs(db, ctx.R) {
				if o.ID.String() == orgID {
					org = &o
					break
				}
			}
			if org == nil {
				return errors.New(msgs.SelecioneOrganizacao)
			}
			path := b.GetURIPrefix()
			if path == "" {
				path = "/"
			}
			http.SetCookie(ctx.W, &http.Cookie{Name: OrgCookie, Value: org.ID.String(), Path: path,
				MaxAge: 365 * 24 * 3600, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: ctx.R.TLS != nil})
			// its page, the menu of it drawn there
			ctx.Resp.RedirectURL = b.GetURIPrefix() + "/" + OrganizacaoURIName + "/" + org.ID.String()
			return nil
		})
}

// configureMenuRoot installs the root item of the organization in the side
// menu, and the action of its chip (ActionChoose).
func configureMenuRoot(b *presets.Builder, db *gorm.DB) {
	if mb := OrganizacaoModel(b); mb != nil {
		configureChoose(b, db, mb)
	}
	b.AddMenuRootFunc(orgMenuRoot(b, db))
}
