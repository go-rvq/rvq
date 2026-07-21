package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// fakeUser is a minimal user.User (login.UserPass provides the UserPasser part).
type fakeUser struct {
	login.UserPass
	id uuid.UUID
}

func (f *fakeUser) GetID() uuid.UUID              { return f.id }
func (f *fakeUser) GetName() string               { return "Test" }
func (f *fakeUser) SetName(string)                {}
func (f *fakeUser) SetEmail(string)               {}
func (f *fakeUser) SetRegistrationDate(time.Time) {}
func (f *fakeUser) GetStatus() string             { return "active" }
func (f *fakeUser) GetRoles() role.Roles          { return nil }
func (f *fakeUser) SetRoles(role.Roles)           {}

func newTestAdmin(t *testing.T) (*presets.Builder, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	b := presets.New(i18n.New()).URIPrefix("/admin")
	b.I18n().SupportLanguages(language.BrazilianPortuguese, language.AmericanEnglish)
	b.DataOperator(gorm2op.DataOperator(db))
	b.Permission(
		perm.New().Policies(
			perm.PolicyFor("everyone").WhoAre(perm.Allowed).ToDo(perm.Anything).On(perm.Anything),
		).SubjectsFunc(func(r *http.Request) []string {
			return []string{"everyone"}
		}),
	)
	Configure(b, db)
	b.Build()
	return b, db
}

func listingAs(t *testing.T, b *presets.Builder, mb *presets.ModelBuilder, u login.UserPasser) (int, string) {
	t.Helper()
	r := multipartestutils.NewMultipartBuilder().
		PageURL(mb.Info().ListingHref()).
		EventFunc(actions.ReloadList).
		BuildEventFuncRequest()
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, u))
	w := httptest.NewRecorder()
	b.ServeHTTP(w, r)
	return w.Code, w.Body.String()
}

func TestOrganizacaoOwnerScoping(t *testing.T) {
	b, db := newTestAdmin(t)
	mb := b.GetModelByID("organizacoes")
	if mb == nil {
		t.Fatal("organizacoes model not registered")
	}

	userA := uuid.New()
	userB := uuid.New()
	if err := db.Create(&models.Organizacao{Nome: "Alpha SA", OwnerID: userA}).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&models.Organizacao{Nome: "Beta Ltda", OwnerID: userB}).Error; err != nil {
		t.Fatal(err)
	}

	// user A only sees their own organization
	code, body := listingAs(t, b, mb, &fakeUser{id: userA})
	if code != http.StatusOK {
		t.Fatalf("listing status = %d, want 200 (%s)", code, body)
	}
	if !strings.Contains(body, "Alpha SA") {
		t.Errorf("owner A should see Alpha SA")
	}
	if strings.Contains(body, "Beta Ltda") {
		t.Errorf("owner A must not see owner B's Beta Ltda")
	}
}
