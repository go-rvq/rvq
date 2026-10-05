package accesskeys

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type keyUser struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string
}

// keysApp is an admin with users and their keys; requests as user (by a
// session: no key), the handler and the builder.
func keysApp(t *testing.T) (*gorm.DB, *presets.Builder, *Builder, func(user uuid.UUID, r *http.Request) *httptest.ResponseRecorder) {
	t.Helper()
	db := testDB(t)
	if err := db.AutoMigrate(&keyUser{}, &perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	p := presets.New(i18n.New()).URIPrefix("/admin")
	p.DataOperator(gorm2op.DataOperator(db))
	pb := perm.New().AllowAll()
	p.Permission(pb.SubjectsFunc(func(*http.Request) []string { return []string{"admin"} }))
	users := p.Model(&keyUser{}, presets.ModelWithID("users"))
	b := New(db)
	b.Authenticator().FindUser = func(_ context.Context, id uuid.UUID) (any, error) {
		u := &keyUser{}
		return u, db.First(u, "id = ?", id).Error
	}
	if err := b.Install(p, users); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	p.Build(mux)
	do := func(user uuid.UUID, r *http.Request) *httptest.ResponseRecorder {
		r = r.WithContext(context.WithValue(r.Context(), login.UserKey, &keyUser{ID: user}))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}
	return db, p, b, do
}

func createKey(href string, form ...[2]string) *http.Request {
	mb := multipartestutils.NewMultipartBuilder().PageURL(href).EventFunc(actions.Create)
	for _, f := range form {
		mb = mb.AddField(f[0], f[1])
	}
	return mb.BuildEventFuncRequest()
}

// My keys: a user makes their own key — its code shown once, only its hash
// kept, its permissions allows of its subject —, sees only their own; the
// keys of a user under the users; a key lasts at most MaxValidity; no request
// by a key reaches any key.
func TestMyKeys(t *testing.T) {
	db, p, _, do := keysApp(t)
	ana, bob := uuid.New(), uuid.New()
	db.Create(&keyUser{ID: ana, Name: "Ana"})
	db.Create(&keyUser{ID: bob, Name: "Bob"})
	my := p.GetModelByID(MyKeysModelID).Info().ListingHref()

	w := do(ana, createKey(my, [2]string{"Name", "deploy"}, [2]string{"Enabled", "true"},
		[2]string{"Permissions[0].Actions[0]", "*"}, [2]string{"Permissions[0].Resources[0]", "admin:/site-files:*"}))
	if w.Code != http.StatusOK {
		t.Fatalf("creating: %d %.300s", w.Code, w.Body.String())
	}
	var keys []AccessKey
	db.Preload("Permissions").Find(&keys)
	if len(keys) != 1 {
		t.Fatalf("keys: %+v\n%.500s", keys, w.Body.String())
	}
	k := keys[0]
	if k.UserID != ana || k.CreatedByID == nil || *k.CreatedByID != ana || !k.Enabled || len(k.Prefix) != 10 || k.CodeHash == "" {
		t.Errorf("the key: %+v", k)
	}
	if !strings.Contains(w.Body.String(), CodePrefix+k.Prefix+"_") {
		i := strings.Index(w.Body.String(), "created")
		t.Errorf("the code not shown: %d %.600s", i, w.Body.String()[max(0, i-300):])
	}
	if d := time.Until(k.ExpiresAt); d < 89*24*time.Hour || d > 91*24*time.Hour {
		t.Errorf("the default validity: %v", d)
	}
	if len(k.Permissions) != 1 || k.Permissions[0].Subject != k.Subject() || k.Permissions[0].Effect != perm.Allowed ||
		k.Permissions[0].ReferID != k.ID.String() || strings.Join(k.Permissions[0].Resources, ",") != "admin:/site-files:*" {
		t.Errorf("its permissions: %+v", k.Permissions)
	}

	// made disabled: kept so
	do(ana, createKey(my, [2]string{"Name", "off"}, [2]string{"Enabled", "false"}))
	var off AccessKey
	if db.First(&off, "name = ?", "off").Error != nil || off.Enabled {
		t.Errorf("a key made disabled: %+v", off)
	}

	// too long: refused
	far := time.Now().Add(400 * 24 * time.Hour).Format("2006-01-02 15:04")
	do(ana, createKey(my, [2]string{"Name", "far"}, [2]string{"ExpiresAt", far}))
	if db.First(&AccessKey{}, "name = ?", "far").Error == nil {
		t.Error("a key of 400 days made")
	}

	// bob's: not in ana's keys; under bob, in the users
	do(bob, createKey(my, [2]string{"Name", "bobs"}))
	if body := do(ana, httptest.NewRequest("GET", my, nil)).Body.String(); strings.Contains(body, "bobs") || !strings.Contains(body, "deploy") {
		t.Errorf("ana's keys:\n%.500s", body)
	}
	users := p.GetModelByID("users").Info().ListingHref()
	if body := do(ana, httptest.NewRequest("GET", users+"/"+bob.String()+"/access-keys", nil)).Body.String(); !strings.Contains(body, "bobs") || strings.Contains(body, "deploy") {
		t.Errorf("bob's keys under the users:\n%.500s", body)
	}

	// a request by a key: no key reached
	r := httptest.NewRequest("GET", my, nil)
	r = r.WithContext(perm.WithRestriction(r.Context(), k.Subject()))
	if w := do(ana, r); w.Code == http.StatusOK && strings.Contains(w.Body.String(), "deploy") {
		t.Errorf("a key reached the keys: %d", w.Code)
	}
}

// The permissions as the form posts them — each list item its __value, the
// editor's flags, no plain field of the item —: saved; with no actions, any
// (*): the resource says all.
func TestKeyPermissionsAsTheFormPosts(t *testing.T) {
	db, p, _, do := keysApp(t)
	ana := uuid.New()
	db.Create(&keyUser{ID: ana, Name: "Ana"})
	my := p.GetModelByID(MyKeysModelID).Info().ListingHref()
	for _, c := range []struct {
		name, actions string
		form          [][2]string
	}{
		{"only-resources", "*", [][2]string{
			{"Permissions[0].Resources[0].__value", "admin:/site-files:{@get,!git}"}, {"Permissions[0].Resources[0].$id", "1"}}},
		{"both", "@get", [][2]string{
			{"Permissions[0].Actions[0].__value", "@get"}, {"Permissions[0].Actions[0].$id", "1"},
			{"Permissions[0].Resources[0].__value", "admin:/site-files:{@get,!git}"}, {"Permissions[0].Resources[0].$id", "2"}}},
	} {
		form := append([][2]string{{"Name", c.name}, {"Enabled", "true"}, {"Permissions.__present", "1"},
			{"Permissions[0].__pos", "0"}, {"Permissions[0].__index", "0"}, {"Permissions[0].__new", "true"}}, c.form...)
		do(ana, createKey(my, form...))
		var k AccessKey
		db.Preload("Permissions").First(&k, "name = ?", c.name)
		if len(k.Permissions) != 1 || strings.Join(k.Permissions[0].Resources, ",") != "admin:/site-files:{@get,!git}" ||
			strings.Join(k.Permissions[0].Actions, ",") != c.actions {
			t.Errorf("%s: %+v", c.name, k.Permissions)
		}
	}
}
