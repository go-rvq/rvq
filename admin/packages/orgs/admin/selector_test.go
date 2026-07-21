package admin

import (
	"bytes"
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	"github.com/google/uuid"
)

func TestOrgSelectorComponent(t *testing.T) {
	b, db := newTestAdmin(t)

	me := uuid.New()
	db.Create(&models.Organizacao{Nome: "Alpha SA", OwnerID: me})
	db.Create(&models.Organizacao{Nome: "Beta Ltda", OwnerID: me})
	db.Create(&models.Organizacao{Nome: "Foreign Inc", OwnerID: uuid.New()})

	r := httptest.NewRequest("GET", "/admin", nil)
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &fakeUser{id: me}))
	ctx := &web.EventContext{R: r}

	comp := OrgSelectorComponent(b, db)(ctx)
	if comp == nil {
		t.Fatal("selector should render for a user with organizations")
	}
	var buf bytes.Buffer
	if err := h.Fprint(&buf, comp, ctx.Context()); err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	// lists the user's own organizations, not others'
	if !strings.Contains(html, "Alpha SA") || !strings.Contains(html, "Beta Ltda") {
		t.Errorf("selector should list the user's organizations; got:\n%s", html)
	}
	if strings.Contains(html, "Foreign Inc") {
		t.Errorf("selector must not list organizations owned by others")
	}
	// binds the selection to the global var
	if !strings.Contains(html, "vars."+OrganizationVar) {
		t.Errorf("selector should bind to vars.%s; got:\n%s", OrganizationVar, html)
	}

	// no user -> no selector
	r2 := httptest.NewRequest("GET", "/admin", nil)
	if OrgSelectorComponent(b, db)(&web.EventContext{R: r2}) != nil {
		t.Errorf("selector should not render without a current user")
	}
}
