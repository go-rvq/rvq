package accesskeys

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type testUser struct{ ID uuid.UUID }

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&AccessKey{}, &AccessKeyUse{}); err != nil {
		t.Fatal(err)
	}
	return db
}

// newKey makes a key of user, its code returned.
func newKey(t *testing.T, db *gorm.DB, user uuid.UUID, mod func(*AccessKey)) (*AccessKey, string) {
	t.Helper()
	code, prefix := NewCode()
	k := &AccessKey{UserID: user, Name: "deploy", Prefix: prefix, CodeHash: CodeHash(code),
		ExpiresAt: time.Now().Add(time.Hour), Enabled: true}
	if mod != nil {
		mod(k)
	}
	if err := db.Create(k).Error; err != nil {
		t.Fatal(err)
	}
	return k, code
}

// A code: "hck_<prefix>_<secret>", its prefix read back, only its hash kept,
// matching only itself.
func TestCode(t *testing.T) {
	code, prefix := NewCode()
	if !strings.HasPrefix(code, CodePrefix+prefix+"_") || len(prefix) != 10 {
		t.Fatalf("code %q, prefix %q", code, prefix)
	}
	if p, err := ParseCode(code); err != nil || p != prefix {
		t.Errorf("parsed %q %v", p, err)
	}
	for _, bad := range []string{"", "hck_", "hck_short_x", "abc_" + prefix + "_" + strings.Repeat("a", 32), "hck_" + prefix} {
		if _, err := ParseCode(bad); err == nil {
			t.Errorf("%q parsed", bad)
		}
	}
	k := &AccessKey{CodeHash: CodeHash(code)}
	if !k.Matches(code) || k.Matches(code+"x") || strings.Contains(k.CodeHash, code) {
		t.Error("the hash")
	}
	other, _ := NewCode()
	if other == code {
		t.Error("two codes alike")
	}
}

// A key authenticates its user — the request restricted to it, the key told —;
// a wrong code, a key disabled, expired, or of no user, does not; too many
// wrong codes from an address, nothing for a while. Each request kept.
func TestAuthenticator(t *testing.T) {
	db := testDB(t)
	ana := uuid.New()
	now := time.Now()
	var told []string
	a := &Authenticator{DB: db,
		FindUser: func(_ context.Context, id uuid.UUID) (any, error) {
			if id != ana {
				return nil, errors.New("no user")
			}
			return &testUser{ID: id}, nil
		},
		Locate: func(string) string { return "Viçosa, BR" },
		KindOf: func(r *http.Request) string { return "git" },
		OnUse:  func(_ *http.Request, k *AccessKey, use *AccessKeyUse) { told = append(told, k.Name+" "+use.Path) },
		Now:    func() time.Time { return now },
	}
	k, code := newKey(t, db, ana, nil)

	r := httptest.NewRequest("GET", "/admin/site-files.git/info/refs", nil)
	r.RemoteAddr = "10.0.0.1:1234"
	user, ctx, err := a.AuthKey(r, code)
	if err != nil || user.(*testUser).ID != ana {
		t.Fatalf("the key: %v %v", user, err)
	}
	if got := perm.RestrictionOf(ctx); len(got) != 1 || got[0] != k.Subject() {
		t.Errorf("the restriction: %v", got)
	}
	if KeyOf(ctx) == nil || KeyOf(ctx).ID != k.ID {
		t.Error("the key not told")
	}
	a.KeyServed(r.WithContext(ctx), 200)
	var uses []AccessKeyUse
	db.Find(&uses)
	if len(uses) != 1 || uses[0].Kind != "git" || uses[0].IP != "10.0.0.1" || uses[0].Place != "Viçosa, BR" ||
		uses[0].Status != 200 || uses[0].Path != "/admin/site-files.git/info/refs" || uses[0].UserID != ana {
		t.Errorf("the history: %+v", uses)
	}
	var kept AccessKey
	db.First(&kept, "id = ?", k.ID)
	if kept.LastUsedAt == nil || kept.LastUsedIP != "10.0.0.1" {
		t.Errorf("the last use: %v %q", kept.LastUsedAt, kept.LastUsedIP)
	}
	if len(told) != 1 {
		t.Errorf("the user's log told %v", told)
	}

	// refused
	_, disabled := newKey(t, db, ana, func(k *AccessKey) { k.Enabled = false })
	_, expired := newKey(t, db, ana, func(k *AccessKey) { k.ExpiresAt = now.Add(-time.Minute) })
	_, orphan := newKey(t, db, uuid.New(), nil)
	for name, c := range map[string]struct {
		code string
		want error
	}{
		"wrong":    {code[:len(code)-1] + "x", ErrKeyNotFound},
		"unknown":  {func() string { c, _ := NewCode(); return c }(), ErrKeyNotFound},
		"disabled": {disabled, ErrKeyUnusable},
		"expired":  {expired, ErrKeyUnusable},
		"no user":  {orphan, ErrKeyNoUser},
	} {
		r := httptest.NewRequest("GET", "/", nil)
		r.RemoteAddr = "10.0.0.2:1"
		if _, _, err := a.AuthKey(r, c.code); !errors.Is(err, c.want) {
			t.Errorf("%s: %v", name, err)
		}
	}
	// deleted: gone
	db.Delete(&AccessKey{}, "id = ?", k.ID)
	if _, _, err := a.AuthKey(r, code); !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("a key deleted: %v", err)
	}

	// too many wrong codes from an address: even a right one refused there
	k2, code2 := newKey(t, db, ana, nil)
	_ = k2
	bad := httptest.NewRequest("GET", "/", nil)
	bad.RemoteAddr = "10.0.0.9:1"
	for i := 0; i < maxFailures; i++ {
		a.AuthKey(bad, code2+"x")
	}
	if _, _, err := a.AuthKey(bad, code2); !errors.Is(err, ErrTooManyErrors) {
		t.Errorf("blocked: %v", err)
	}
	other := httptest.NewRequest("GET", "/", nil)
	other.RemoteAddr = "10.0.0.8:1"
	if _, _, err := a.AuthKey(other, code2); err != nil {
		t.Errorf("another address blocked: %v", err)
	}
	now = now.Add(failureWindow + time.Second)
	if _, _, err := a.AuthKey(bad, code2); err != nil {
		t.Errorf("still blocked after the window: %v", err)
	}
}
