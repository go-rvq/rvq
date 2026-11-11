package user

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type stackUser struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	Roles     []*role.Role `gorm:"many2many:user_role_join;joinForeignKey:UserID"`
	DeletedAt gorm.DeletedAt
	login.UserPass
}

func (u *stackUser) TableName() string               { return "users" }
func (u *stackUser) GetID() uuid.UUID                { return u.ID }
func (u *stackUser) SetID(v uuid.UUID)               { u.ID = v }
func (u *stackUser) Anonymous() bool                 { return IsAnonymous(u.ID) }
func (u *stackUser) GetName() string                 { return u.Name }
func (u *stackUser) SetName(v string)                { u.Name = v }
func (u *stackUser) SetEmail(v string)               { u.Account = v }
func (u *stackUser) SetRegistrationDate(v time.Time) {}
func (u *stackUser) GetStatus() string               { return "active" }
func (u *stackUser) GetRoles() role.Roles            { return u.Roles }
func (u *stackUser) SetRoles(roles role.Roles)       { u.Roles = roles }

var stackSeq int64

func stackDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:stack_%d?mode=memory&cache=shared", atomic.AddInt64(&stackSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&role.Role{}, &stackUser{}); err != nil {
		t.Fatal(err)
	}
	if err = EnsureRolePositions(db); err != nil {
		t.Fatal(err)
	}
	return db
}

// The roles of a user in its order: kept by their positions, read in it.
func TestRolePositions(t *testing.T) {
	db := stackDB(t)
	a := &role.Role{ID: uuid.New(), Name: "A"}
	z := &role.Role{ID: uuid.New(), Name: "Z"}
	u := &stackUser{ID: uuid.New(), Name: "Ana", Roles: []*role.Role{a, z}}
	if err := db.Create(u).Error; err != nil {
		t.Fatal(err)
	}
	if err := EnsureRolePositions(db); err != nil { // again: nothing
		t.Fatal(err)
	}
	names := func() string {
		roles, err := OrderedRoles(db, u.ID)
		if err != nil {
			t.Fatal(err)
		}
		return role.Roles(roles).FirstName() + "," + strings.Join(role.Roles(roles).Names()[1:], ",")
	}
	if got := names(); got != "A,Z" {
		t.Errorf("by name, no positions: %s", got)
	}
	if err := SetRolePositions(db, u.ID, []uuid.UUID{z.ID, a.ID}); err != nil {
		t.Fatal(err)
	}
	if got := names(); got != "Z,A" {
		t.Errorf("in its order: %s", got)
	}
}

// The source of the stacks: the user of a request — its subject, the
// version of its roles in its order (a role changed: another) —, its stack
// the code's first, then its roles' in its order; none for the anonymous.
func TestStackSource(t *testing.T) {
	pb := perm.New().Policies(perm.PolicyFor("Editor").WhoAre(perm.Denied).ToDo(perm.Anything).On("admin:secret:*"))
	b := &Builder{lb: login.New(i18n.New())}
	u := &stackUser{ID: uuid.New(), Name: "Ana", Roles: []*role.Role{
		{ID: uuid.New(), Name: "Editor", Permissions: pq.StringArray{"admin:users:{list,new}", "allow admin:secret:*"}},
		{ID: uuid.New(), Name: "Locked", Permissions: pq.StringArray{"deny admin:users:new"}},
	}}
	u.Account = "ana" // not the initial user's (its stack: the code's alone)
	src := b.StackSource(pb)
	r := httptest.NewRequest("GET", "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), login.UserKey, u))
	subject, version, build := src.StackOf(r)
	if subject != "user:"+u.ID.String() || build == nil {
		t.Fatalf("the source: %q, built: %v", subject, build != nil)
	}
	s, err := build()
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, l := range s.Lines {
		got = append(got, fmt.Sprintf("%d %s (%s) %d", l.Index, l.Text(), l.Origin, l.ReplacedBy))
	}
	want := "1 deny admin:secret:* (code: Editor) 4|2 admin:users:list (role: Editor) 0|3 admin:users:new (role: Editor) 5|" +
		"4 allow admin:secret:* (role: Editor) 0|5 deny admin:users:new (role: Locked) 0"
	if strings.Join(got, "|") != want {
		t.Errorf("the stack:\n%s\nwant\n%s", strings.Join(got, "|"), want)
	}
	role.Changed(nil)
	if _, v2, _ := src.StackOf(r); v2 == version {
		t.Error("a role changed: the same version")
	}
	u.Roles[0], u.Roles[1] = u.Roles[1], u.Roles[0]
	if _, v3, _ := src.StackOf(r); v3 == version {
		t.Error("the order changed: the same version")
	}
	anon := httptest.NewRequest("GET", "/", nil)
	if s, _, _ := src.StackOf(anon); s != "" {
		t.Errorf("the anonymous: %q", s)
	}
}

// The page of the permissions: Tudo, each line in its place, its origin and
// its state; Ativas, the active ones — an allow of no prefix —, and as a
// tree.
func TestPermissionsBody(t *testing.T) {
	b := &Builder{lb: login.New(i18n.New())}
	s := (&perm.Stack{}).Add("role: Editor", "admin:users:{list,new}", "deny admin:users:new", "allow admin:secret:*",
		"reset admin:pages:*").Resolve()
	body := h.MustString(b.permissionsBody(context.Background(), Messages_pt_BR, s, &stackUser{ID: uuid.New()}), context.Background())
	for _, want := range []string{
		"Tudo", "Ativas", "Índice", "Origem", "Situação",
		"substituída pela #3", "remoção", "role: Editor",
		"admin:secret:*", "Em árvore", "users:{", "Árvore",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("no %q", want)
		}
	}
	if strings.Contains(body, "allow admin:secret:*</code></div>") {
		t.Error("an allow active with its prefix")
	}
	// an allow under a wider deny: the deny but what it matches
	s = (&perm.Stack{}).Add("role: Editor", "admin:*", "deny admin:users:*", "allow admin:users:@list").Resolve()
	body = h.MustString(b.permissionsBody(context.Background(), Messages_pt_BR, s, &stackUser{ID: uuid.New()}), context.Background())
	for _, want := range []string{"ativa, menos o que #3 casa", "menos admin:users:@list"} {
		if !strings.Contains(body, want) {
			t.Errorf("no %q", want)
		}
	}
}
