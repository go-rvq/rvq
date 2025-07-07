package perm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/ory/ladon"
	"github.com/qor5/web/v3/zeroer"
	"github.com/sunfmin/reflectutils"
)

var Verbose = false

type verReq struct {
	subjects       []string
	objs           []interface{}
	r              *http.Request
	req            *ladon.Request
	resourcesParts []string
}

type Verifier struct {
	builder *Builder
	module  string
	vr      *verReq
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
	}

	resourceParts := []string{b.module}
	if b.vr != nil {
		resourceParts = b.vr.resourcesParts
	}

	r.vr = &verReq{
		resourcesParts: append([]string{}, resourceParts...),
		req:            &ladon.Request{},
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
	return strings.Join(b.vr.resourcesParts, ":") + ":"
}

func (b *Verifier) ResourceWithModule() string {
	var m string
	if b.module != "" {
		m = b.module
	}
	return m + ":" + strings.Join(b.vr.resourcesParts, ":") + ":"
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
	if b.builder == nil {
		return nil
	}

	b.vr.req.Resource = b.Resource()

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

	var err error
	// any of the subjects have permission, then have permission
	for _, sub := range b.vr.subjects {
		b.vr.req.Subject = sub

		err = b.builder.ladon.IsAllowed(context.TODO(), b.vr.req)
		if Verbose {
			fmt.Printf("have permission: %+v, req: {%s}\n", err == nil, RequestToString(b.vr.req))
		}
		if err == nil {
			return nil
		}
	}

	return err
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
	if r.Context != nil && len(r.Context) > 0 {
		b, _ := json.Marshal(r.Context)
		s = append(s, fmt.Sprintf("context=%v", string(b)))
	}
	return strings.Join(s, ", ")
}
