package perm

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gobwas/glob"
	"github.com/ory/ladon"
	"github.com/ory/ladon/manager/memory"
)

// StackSource gives the stack of the permissions of the user of a request
// (Builder.Stacks).
type StackSource interface {
	// StackOf is, of the user of r, the subject its resolved stack is
	// registered by ("user:<id>"; "" when r has none: the policies decide,
	// as before), the version of what the stack is made of (its roles, their
	// lines, their order: when it changes, the stack is made again) and how
	// it is made — its lines in order (Stack.Add), the code's first
	// (Builder.CodeLines), then the database's (Builder.DBLines). Not
	// resolved: Register resolves it.
	StackOf(r *http.Request) (subject, version string, build func() (*Stack, error))
}

// stacks are the stacks registered: their own ladon — the policies of their
// active lines, of their subjects —, apart from the policies of the code and
// of the database, which only make them.
type stacks struct {
	m        sync.Mutex
	ladon    *ladon.Ladon
	versions map[string]string
	ids      map[string][]string
	// code is the version of the policies (of the code and of the
	// database): changed, every stack is made again
	code int
}

func newStacks() *stacks {
	return &stacks{
		ladon:    &ladon.Ladon{Manager: memory.NewMemoryManager(), Matcher: &PathMatcher{}},
		versions: map[string]string{},
		ids:      map[string][]string{},
	}
}

// Stacks makes the requests of a user decided by the active lines of the
// user's stack (StackSource): resolved, registered as policies of the user's
// subject — once a version — in a ladon of their own; the request asks as
// that subject alone. The policies of the code and of the database (the
// roles') no longer decide a user's request: they are what the stacks are
// made of. A request of no user's subject is decided as before.
func (b *Builder) Stacks(src StackSource) *Builder {
	b.stackSource = src
	if b.stacks == nil {
		b.stacks = newStacks()
	}
	return b
}

// GetStacks is the source of the stacks, when set.
func (b *Builder) GetStacks() StackSource { return b.stackSource }

// Register registers the stack of subject, made by build and resolved: the
// policies of its active lines, in its own ladon (StackPolicies) — those of
// the version before taken out. Nothing when the version registered is
// version already (and no policy changed since).
func (b *Builder) Register(subject, version string, build func() (*Stack, error)) error {
	if b.stacks == nil {
		b.stacks = newStacks()
	}
	st := b.stacks
	st.m.Lock()
	full := version + "#" + strconv.Itoa(st.code)
	current := st.versions[subject] == full
	st.m.Unlock()
	if current {
		return nil
	}
	// made out of the lock: building reads the policies of the code, under
	// the builder's lock, which a policy made takes before this one
	s, err := build()
	if err != nil {
		return err
	}
	s.Resolve()
	policies := StackPolicies(subject, s.Active())

	st.m.Lock()
	defer st.m.Unlock()
	for _, id := range st.ids[subject] {
		_ = st.ladon.Manager.Delete(context.TODO(), id)
	}
	st.ids[subject] = nil
	for _, p := range policies {
		if err := st.ladon.Manager.Create(context.TODO(), p.policy); err != nil {
			return err
		}
		st.ids[subject] = append(st.ids[subject], p.GetID())
	}
	st.versions[subject] = full
	return nil
}

// Registered says whether the stack of subject is registered.
func (b *Builder) Registered(subject string) bool {
	if b.stacks == nil {
		return false
	}
	b.stacks.m.Lock()
	defer b.stacks.m.Unlock()
	_, ok := b.stacks.versions[subject]
	return ok
}

// Forget takes the stacks of subjects out — all of them when none —: made
// again at their next request.
func (b *Builder) Forget(subjects ...string) {
	if b.stacks == nil {
		return
	}
	st := b.stacks
	st.m.Lock()
	defer st.m.Unlock()
	if len(subjects) == 0 {
		for s := range st.versions {
			subjects = append(subjects, s)
		}
	}
	for _, s := range subjects {
		for _, id := range st.ids[s] {
			_ = st.ladon.Manager.Delete(context.TODO(), id)
		}
		delete(st.ids, s)
		delete(st.versions, s)
	}
}

// policiesChanged makes every stack be made again: a policy — of the code or
// of the database — changed.
func (b *Builder) policiesChanged() {
	if b.stacks == nil {
		return
	}
	b.stacks.m.Lock()
	b.stacks.code++
	b.stacks.m.Unlock()
}

// StackPolicies are the policies of the active lines of a stack, of subject:
// a line of no prefix (an allow is one) allowed, a deny denied — each its
// pattern, its actions (any, when it names none), its conditions and, when
// it has exceptions, an ExceptCondition of them.
func StackPolicies(subject string, active []*Line) (out []*PolicyBuilder) {
	for _, l := range active {
		effect := Allowed
		if l.Denies() {
			effect = Denied
		}
		actions := l.Actions
		if anyAction(actions) {
			actions = []string{Anything}
		}
		conditions := l.Conditions
		if len(l.Except) > 0 {
			conditions = ladon.Conditions{}
			for k, c := range l.Conditions {
				conditions[k] = c
			}
			conditions[exceptKey] = NewExceptCondition(l.Except...)
		}
		p := &PolicyBuilder{policy: &ladon.DefaultPolicy{
			ID:         fmt.Sprintf("stack:%s:%d", subject, l.Index),
			Subjects:   []string{subject},
			Effect:     effect,
			Actions:    actions,
			Resources:  []string{l.Pattern},
			Conditions: conditions,
		}}
		out = append(out, p)
	}
	return
}

// stackSubject registers the stack of the user of r (Stacks) and is its
// subject: "" when there is none, or when the stack could not be made.
func (b *Builder) stackSubject(r *http.Request) string {
	if b.stackSource == nil || r == nil {
		return ""
	}
	subject, version, build := b.stackSource.StackOf(r)
	if subject == "" || build == nil {
		return ""
	}
	if err := b.Register(subject, version, build); err != nil {
		if Verbose {
			fmt.Printf("the stack of %s: %v\n", subject, err)
		}
		return ""
	}
	return subject
}

// ladonOf is the ladon that decides for subjects: the stacks', when the
// first is a stack registered; else the policies'.
func (b *Builder) ladonOf(subjects []string) *ladon.Ladon {
	if len(subjects) > 0 && b.Registered(subjects[0]) {
		return b.stacks.ladon
	}
	return b.ladon
}

// exceptKey is the key of the ExceptCondition of a policy of a line.
const exceptKey = "stack:except"

// ExceptCondition is fulfilled by a request whose resource matches none of
// its patterns (globs): the exceptions of a line (Line.Except), what it no
// longer counts for.
type ExceptCondition struct {
	Patterns []string
	globs    []glob.Glob
}

// NewExceptCondition is the condition of the exceptions patterns.
func NewExceptCondition(patterns ...string) *ExceptCondition {
	c := &ExceptCondition{Patterns: patterns}
	for _, p := range patterns {
		g, err := glob.Compile(p)
		if err != nil {
			g = glob.MustCompile(glob.QuoteMeta(p))
		}
		c.globs = append(c.globs, g)
	}
	return c
}

// GetName is the name of the condition.
func (c *ExceptCondition) GetName() string { return "StackExceptCondition" }

// Fulfills says whether the resource of r matches none of the patterns.
func (c *ExceptCondition) Fulfills(_ context.Context, _ interface{}, r *ladon.Request) bool {
	for _, g := range c.globs {
		if g.Match(r.Resource) {
			return false
		}
	}
	return true
}
