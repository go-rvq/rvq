package admin

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/orgs/messages"
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

// OrganizationVar is the global UI variable that holds the selected organization
// id. Because vars is global in the UI, the selected value is accessible
// anywhere as vars.currentOrganization (e.g. to scope other modules' UI).
const OrganizationVar = "currentOrganization"

// OrgSelectorComponent renders the organization selector shown at the top of the
// admin menu: a dropdown of the current user's organizations bound to
// vars.currentOrganization. Selecting one sets the global var and navigates to
// that organization's area (/orgs/{id}).
func OrgSelectorComponent(b *presets.Builder, db *gorm.DB) presets.ComponentFunc {
	return func(ctx *web.EventContext) h.HTMLComponent {
		u := user.GetCurrentUser(ctx.R)
		if u == nil {
			return nil
		}
		var orgs []models.Organizacao
		db.Where("owner_id = ?", u.GetID()).Order("nome").Find(&orgs)
		if len(orgs) == 0 {
			return nil
		}
		items := make([]map[string]string, len(orgs))
		for i, o := range orgs {
			items[i] = map[string]string{"title": o.Nome, "value": o.ID.String()}
		}
		return v.VSelect().
			Label(messages.Get(ctx.Context()).SelecioneOrganizacao).
			Items(items).
			ItemTitle("title").
			ItemValue("value").
			Attr("v-model", "vars."+OrganizationVar).
			Density(v.DensityCompact).
			HideDetails(true).
			// on change, (re)load the org's permitted resource menu into the portal
			Attr("@update:modelValue",
				web.Plaid().EventFunc(nestedMenuEvent).Query(paramSelectedOrg, web.Var("$event")).Go())
	}
}

// configureSelector installs the organization selector at the top of the menu.
func configureSelector(b *presets.Builder, db *gorm.DB) {
	b.AddMenuTopItemFunc("orgs_selector", OrgSelectorComponent(b, db))
}
