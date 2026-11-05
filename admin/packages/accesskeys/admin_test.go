package accesskeys

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/model"
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
// kept, its permissions its denies —, sees only their own; the
// keys of a user under the users; a key lasts at most MaxValidity; no request
// by a key reaches any key.
func TestMyKeys(t *testing.T) {
	db, p, _, do := keysApp(t)
	ana, bob := uuid.New(), uuid.New()
	db.Create(&keyUser{ID: ana, Name: "Ana"})
	db.Create(&keyUser{ID: bob, Name: "Bob"})
	my := p.GetModelByID(MyKeysModelID).Info().ListingHref()

	w := do(ana, createKey(my, [2]string{"Name", "deploy"}, [2]string{"Enabled", "true"},
		[2]string{"Permissions.__present", "1"},
		[2]string{"Permissions[0].__value", "deny admin:users:*"}, [2]string{"Permissions[0].$id", "1"}))
	if w.Code != http.StatusOK {
		t.Fatalf("creating: %d %.300s", w.Code, w.Body.String())
	}
	var keys []AccessKey
	db.Find(&keys)
	if len(keys) != 1 {
		t.Fatalf("keys: %+v\n%.500s", keys, w.Body.String())
	}
	k := keys[0]
	if k.UserID != model.UUID(ana) || k.CreatedByID != model.UUID(ana) || !k.Enabled || len(k.Prefix) != 10 || k.CodeHash == "" {
		t.Errorf("the key: %+v", k)
	}
	if !strings.Contains(w.Body.String(), CodePrefix+k.Prefix+"_") {
		i := strings.Index(w.Body.String(), "created")
		t.Errorf("the code not shown: %d %.600s", i, w.Body.String()[max(0, i-300):])
	}
	if d := time.Until(k.ExpiresAt); d < 89*24*time.Hour || d > 91*24*time.Hour {
		t.Errorf("the default validity: %v", d)
	}
	if strings.Join(k.Permissions, " ") != "deny admin:users:*" {
		t.Errorf("its permissions: %q", k.Permissions)
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
	all := &AccessKey{ID: uuid.New(), Name: "all"}
	if err := p.GetPermission().Register(all.Subject(), "1", func() (*perm.Stack, error) { return all.Stack(), nil }); err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", my, nil)
	r = r.WithContext(perm.WithRestriction(r.Context(), all.Subject()))
	if w := do(ana, r); w.Code == http.StatusOK && strings.Contains(w.Body.String(), "deploy") {
		t.Errorf("a key reached the keys: %d", w.Code)
	}
}

// The permissions as the form posts them — a list of texts, each item its
// __value —: kept as written, in their order, the blank ones out; only
// denies — a line that allows, resets, or a group of one that does not deny,
// refused.
func TestKeyPermissionsAsTheFormPosts(t *testing.T) {
	db, p, _, do := keysApp(t)
	ana := uuid.New()
	db.Create(&keyUser{ID: ana, Name: "Ana"})
	my := p.GetModelByID(MyKeysModelID).Info().ListingHref()
	post := func(name string, lines ...string) {
		form := [][2]string{{"Name", name}, {"Enabled", "true"}, {"Permissions.__present", "1"}}
		for i, v := range lines {
			form = append(form, [2]string{fmt.Sprintf("Permissions[%d].__value", i), v},
				[2]string{fmt.Sprintf("Permissions[%d].$id", i), fmt.Sprint(i + 1)})
		}
		do(ana, createKey(my, form...))
	}
	post("denies", "deny admin:users:*", " ", " deny admin:pages:{@edit,@delete} ")
	var k AccessKey
	if err := db.First(&k, "name = ?", "denies").Error; err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(k.Permissions, "|"); got != "deny admin:users:*|deny admin:pages:{@edit,@delete}" {
		t.Errorf("the denies: %q", got)
	}
	for i, bad := range []string{"admin:*", "allow admin:*", "reset admin:users:*", "deny", "denyadmin:users:*",
		"admin:users:{deny @edit, @list}"} {
		name := fmt.Sprint("bad", i)
		post(name, "deny admin:users:*", bad)
		if db.First(&AccessKey{}, "name = ?", name).Error == nil {
			t.Errorf("a key of a line %q made", bad)
		}
	}
}

// userSource is the stack of the user of a request: its header X-User.
type userSource map[string][]string

func (s userSource) StackOf(r *http.Request) (string, string, func() (*perm.Stack, error)) {
	u := r.Header.Get("X-User")
	if u == "" {
		return "", "", nil
	}
	return "user:" + u, "1", func() (*perm.Stack, error) { return (&perm.Stack{}).Add("role: R", s[u]...), nil }
}

// A request by a key: allowed what its user may, but what the key denies;
// never a key, whatever the user may (ProtectedLines); a key of before the
// denies, what its lines of no prefix reach, and its user may.
func TestKeyStackOfARequest(t *testing.T) {
	pb := perm.New().Stacks(userSource{"all": {"*"}, "pages": {"admin:pages:*"}})
	n := 0
	allowed := func(user, resource string, lines ...string) bool {
		t.Helper()
		n++
		k := &AccessKey{ID: uuid.New(), Name: fmt.Sprint("k", n), Permissions: lines}
		if err := pb.Register(k.Subject(), "1", func() (*perm.Stack, error) { return k.Stack(), nil }); err != nil {
			t.Fatal(err)
		}
		parts := strings.Split(resource, ":")
		r := httptest.NewRequest("GET", "/admin", nil)
		r.Header.Set("X-User", user)
		r = r.WithContext(perm.WithRestriction(r.Context(), k.Subject()))
		v := perm.NewVerifier("admin", pb).Spawn().Do(parts[len(parts)-1]).On(parts[:len(parts)-1]...).WithReq(r)
		return v.IsAllowed() == nil
	}
	for _, c := range []struct {
		user, resource string
		lines          []string
		want           bool
	}{
		{"all", "users:@list", nil, true},
		{"all", "users:@list", []string{"deny admin:users:*"}, false},
		{"all", "pages:@list", []string{"deny admin:users:*"}, true},
		{"all", "users:@edit", []string{"deny admin:users:{@edit,@delete}"}, false},
		{"all", "users:@list", []string{"deny admin:users:{@edit,@delete}"}, true},
		{"pages", "users:@list", nil, false},
		{"pages", "pages:@list", nil, true},
		{"all", MyKeysModelID + ":@list", nil, false},
		{"all", "users:" + UserKeysModelID + ":@list", nil, false},
		{"all", UsesModelID + ":@list", nil, false},
		// of before the denies: its lines are what it reaches
		{"all", "users:@list", []string{"admin:pages:*"}, false},
		{"all", "pages:@list", []string{"admin:pages:*"}, true},
		{"pages", "users:@list", []string{"admin:*"}, false},
	} {
		if got := allowed(c.user, c.resource, c.lines...); got != c.want {
			t.Errorf("%s, %s by %q: %v, want %v", c.user, c.resource, c.lines, got, c.want)
		}
	}
}

// The keys of before the lines: their policies made lines ("deny " before a
// deny's), once — the policies gone; a key of lines left as it is.
func TestMigrateLines(t *testing.T) {
	db := testDB(t)
	if err := db.AutoMigrate(&perm.DefaultDBPolicy{}); err != nil {
		t.Fatal(err)
	}
	old, _ := newKey(t, db, uuid.New(), nil)
	for _, pol := range []*perm.DefaultDBPolicy{
		{ReferID: old.ID.String(), Subject: old.Subject(), Effect: perm.Allowed, Resources: []string{"admin:*"}},
		{ReferID: old.ID.String(), Subject: old.Subject(), Effect: perm.Denied, Resources: []string{"admin:users:*"}},
	} {
		pol.ID = uuid.New()
		if err := db.Create(pol).Error; err != nil {
			t.Fatal(err)
		}
	}
	lined, _ := newKey(t, db, uuid.New(), func(k *AccessKey) { k.Permissions = []string{"deny admin:pages:*"} })
	for range 2 {
		if err := MigrateLines(db); err != nil {
			t.Fatal(err)
		}
	}
	var k AccessKey
	db.First(&k, "id = ?", old.ID)
	if got := strings.Join(k.Permissions, "|"); got != "admin:*|deny admin:users:*" && got != "deny admin:users:*|admin:*" {
		t.Errorf("the lines of the key of before: %q", got)
	}
	if !k.legacy() {
		t.Error("the key of before not of before")
	}
	var n int64
	db.Model(&perm.DefaultDBPolicy{}).Where("refer_id = ?", old.ID.String()).Count(&n)
	if n != 0 {
		t.Errorf("its policies left: %d", n)
	}
	var r AccessKey
	db.First(&r, "id = ?", lined.ID)
	if strings.Join(r.Permissions, " ") != "deny admin:pages:*" {
		t.Errorf("a key of lines changed: %q", r.Permissions)
	}
}
