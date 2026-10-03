package perm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-rvq/rvq/web/zeroer"
	"github.com/go-rvq/rvq/x/osenv"
	"github.com/iancoleman/strcase"
	"github.com/ory/ladon"
	"github.com/sunfmin/reflectutils"
)

var ErrIsDanied = errors.New("verifier is denied")

var Verbose = osenv.GetBool("RVQ_PERM_VERBOSE", "Print all permissions verification", false)

type verReq struct {
	subjects       []string
	objs           []interface{}
	r              *http.Request
	req            *ladon.Request
	resourcesParts []string
	// preferredParts, when set, are the resource by the unique name of what
	// is verified (Prefer): its policies decide before the ones of
	// resourcesParts
	preferredParts []string
}

type VerifierMode uint8

func (m VerifierMode) Is(o VerifierMode) bool {
	return m == o
}

const (
	VerifierModeDefault VerifierMode = iota
	VerifierModeAllow
	VerifierModeDeny
	VerifierModeMustAllow
	VerifierModeMustDeny
)

type Verifier struct {
	builder *Builder
	module  string
	vr      *verReq
	mode    VerifierMode
}

func NewVerifier(module string, b *Builder) (r *Verifier) {
	r = &Verifier{
		module: module,
	}

	if b == nil {
		return r
	}

	r.builder = b
	return
}

func (b *Verifier) Deny() *Verifier {
	b.mode = VerifierModeDeny
	return b
}

func (b *Verifier) Allow() *Verifier {
	b.mode = VerifierModeAllow
	return b
}

func (b *Verifier) SetMode(v VerifierMode) *Verifier {
	b.mode = v
	return b
}

func (b *Verifier) Module() string {
	return b.module
}

func (b *Verifier) Spawn() (r *Verifier) {
	if b.builder == nil {
		return b
	}

	r = &Verifier{
		module:  b.module,
		builder: b.builder,
		mode:    b.mode,
	}

	resourceParts := []string{b.module}
	if b.vr != nil {
		resourceParts = b.vr.resourcesParts
	}

	r.vr = &verReq{
		resourcesParts: append([]string{}, resourceParts...),
		req:            &ladon.Request{},
	}
	if b.vr != nil && b.vr.preferredParts != nil {
		r.vr.preferredParts = append([]string{}, b.vr.preferredParts...)
	}

	if b.vr != nil {
		r.vr.r = b.vr.r
	}

	return
}

func (b *Verifier) Do(v string) (r *Verifier) {
	if b.builder == nil {
		return b
	}

	r = b.Spawn()
	r.vr.req.Action = v
	return
}

func (b *Verifier) Resource() string {
	return strings.Join(b.resourceParts(), ":") + ":"
}

// ResourceParts are a copy of the parts of the resource: what a resource by
// a unique name (Prefer) is made under.
func (b *Verifier) ResourceParts() []string {
	return append([]string{}, b.resourceParts()...)
}

// Prefer sets the resource by the unique name of what is verified — the
// parts under the root and its name, "posts" —: the parts added after
// (On, SnakeOn) go to both resources. Its policies decide first (IsAllowed):
// a deny of any subject by the unique name denies, an allow allows; only when
// none matches it does the resource of the ancestors decide.
func (b *Verifier) Prefer(parts ...string) *Verifier {
	if b.builder == nil {
		return b
	}
	b.vr.preferredParts = append([]string{}, parts...)
	return b
}

// PreferredResource is the resource by the unique name (Prefer), "" when
// there is none.
func (b *Verifier) PreferredResource() string {
	if b.vr == nil || b.vr.preferredParts == nil {
		return ""
	}
	return strings.Join(b.vr.preferredParts, ":") + ":"
}

// resourceParts is nil-safe: NewVerifier accepts a nil Builder (an app with no
// permission builder), and then there is no request to read parts from.
func (b *Verifier) resourceParts() []string {
	if b.vr == nil {
		return nil
	}
	return b.vr.resourcesParts
}

func (b *Verifier) ResourceWithModule() string {
	var m string
	if b.module != "" {
		m = b.module
	}
	return m + ":" + strings.Join(b.resourceParts(), ":") + ":"
}

// SnakeDo convert string to snake form.
// e.g. "SnakeDo" -> "snake_do"
func (b *Verifier) SnakeDo(actions ...string) (r *Verifier) {
	fixed := []string{b.module}
	for _, a := range actions {
		fixed = append(fixed, strcase.ToSnake(a))
	}
	return b.Do(strings.Join(fixed, ":"))
}

func (b *Verifier) On(vs ...string) (r *Verifier) {
	if b.builder == nil {
		return b
	}

	b.vr.resourcesParts = append(b.vr.resourcesParts, vs...)
	if b.vr.preferredParts != nil {
		b.vr.preferredParts = append(b.vr.preferredParts, vs...)
	}
	return b
}

func (b *Verifier) SnakeOn(vs ...string) (r *Verifier) {
	if b.builder == nil {
		return b
	}

	var fixed []string
	for _, v := range vs {
		if v == "" {
			continue
		}
		fixed = append(fixed, strcase.ToSnakeWithIgnore(v, "."))
	}

	b.On(fixed...)
	return b
}

func (b *Verifier) ObjectOn(v interface{}) (r *Verifier) {
	if b.builder == nil || v == nil {
		return b
	}

	id, err := reflectutils.Get(v, "ID")
	if err == nil && !zeroer.IsZero(id) {
		b.vr.objs = append(b.vr.objs, v)
		b.SnakeOn(fmt.Sprint(id))
	}

	return b
}

func (b *Verifier) RemoveOn(length int) (r *Verifier) {
	if b.builder == nil {
		return b
	}
	if len(b.vr.resourcesParts) >= length {
		b.vr.resourcesParts = b.vr.resourcesParts[:len(b.vr.resourcesParts)-length]
	}
	if len(b.vr.preferredParts) >= length {
		b.vr.preferredParts = b.vr.preferredParts[:len(b.vr.preferredParts)-length]
	}
	return b
}

func (b *Verifier) WithReq(v *http.Request) *Verifier {
	if b.builder == nil {
		return b
	}
	b.vr.r = v
	return b
}

func (b *Verifier) From(v string) (r *Verifier) {
	if b.builder == nil {
		return b
	}

	b.vr.subjects = append(b.vr.subjects, v)
	return b
}

func (b *Verifier) Given(v ladon.Context) (r *Verifier) {
	if b.builder == nil {
		return b
	}
	b.vr.req.Context = v
	return b
}

func (b *Verifier) Allowed() bool {
	return b.IsAllowed() == nil
}

func (b *Verifier) Denied() bool {
	return b.IsAllowed() != nil
}

func (b *Verifier) IsAllowed() error {
	if b.builder == nil || b.mode >= VerifierModeMustAllow {
		return nil
	}

	// the permission asked is the resource with its verb at the end:
	// "presets:site/:seo/:seo_config:7:@edit", "…:7:!publish"
	b.vr.req.Resource = b.Resource() + b.vr.req.Action

	if len(b.vr.subjects) == 0 && b.builder.subjectsFunc != nil {
		b.vr.subjects = b.builder.subjectsFunc(b.vr.r)
	}

	if len(b.vr.subjects) == 0 {
		b.vr.subjects = []string{Anonymous}
	}

	if b.builder.contextFunc != nil {
		newContext := b.builder.contextFunc(b.vr.r, b.vr.objs)
		if newContext != nil {
			for k, v := range b.vr.req.Context {
				newContext[k] = v
			}
			b.vr.req.Context = newContext
		}
	}

	// the unique name first (Prefer): a deny of any subject denies, an allow
	// allows; when no policy of it matches, the ancestors decide
	if b.mode == VerifierModeDefault && b.vr.preferredParts != nil {
		if decided, err := b.preferredDecision(); decided {
			return err
		}
	}

	var err error
	// any of the subjects have permission, then have permission
	for _, sub := range b.vr.subjects {
		b.vr.req.Subject = sub

		switch b.mode {
		case VerifierModeDeny:
			err = ErrIsDanied
		case VerifierModeDefault:
			err = b.builder.ladon.IsAllowed(context.TODO(), b.vr.req)
		}

		if Verbose {
			fmt.Printf("have permission: %+v, req: {%s}\n", err == nil, RequestToString(b.vr.req))
		}

		if err == nil {
			return nil
		}
	}

	return err
}

// preferredDecision is the decision of the policies of the resource by the
// unique name (Prefer), for all the subjects: denied when one of them is
// denied by a policy (deny wins), allowed when one is allowed; undecided when
// no policy matches it for any subject.
func (b *Verifier) preferredDecision() (decided bool, err error) {
	req := *b.vr.req
	req.Resource = b.PreferredResource() + req.Action
	allowed := false
	for _, sub := range b.vr.subjects {
		req.Subject = sub
		e := b.builder.ladon.IsAllowed(context.TODO(), &req)
		if Verbose {
			fmt.Printf("by the unique name: %v, req: {%s}\n", e, RequestToString(&req))
		}
		switch {
		case e == nil:
			allowed = true
		case errors.Is(e, ladon.ErrRequestForcefullyDenied):
			return true, e
		}
	}
	return allowed, nil
}

func RequestToString(r *ladon.Request) string {
	var s []string
	if r.Resource != "" {
		s = append(s, fmt.Sprintf("resurce=%q", r.Resource))
	}
	if r.Action != "" {
		s = append(s, fmt.Sprintf("action=%q", r.Action))
	}
	if r.Subject != "" {
		s = append(s, fmt.Sprintf("subject=%q", r.Subject))
	}
	if len(r.Context) > 0 {
		b, _ := json.Marshal(r.Context)
		s = append(s, fmt.Sprintf("context=%v", string(b)))
	}
	return strings.Join(s, ", ")
}
