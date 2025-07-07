package web

import (
	"bytes"
	"context"
	"strings"

	h "github.com/go-rvq/htmlgo"
)

type RunScriptBuilder struct {
	scripts       []string
	setups        []string
	beforeUnmount []string
	unmounted     []string
}

func RunScript(script ...string) (r *RunScriptBuilder) {
	r = &RunScriptBuilder{
		scripts: script,
	}
	return
}

func (r *RunScriptBuilder) BeforeUnmount(s ...string) *RunScriptBuilder {
	r.beforeUnmount = append(r.beforeUnmount, s...)
	return r
}

func (r *RunScriptBuilder) Unmount(s ...string) *RunScriptBuilder {
	r.unmounted = append(r.unmounted, s...)
	return r
}

func (r *RunScriptBuilder) Script(s ...string) *RunScriptBuilder {
	r.scripts = append(r.scripts, s...)
	return r
}

func (r *RunScriptBuilder) Setup(s ...string) *RunScriptBuilder {
	r.setups = append(r.setups, s...)
	return r
}

func (r *RunScriptBuilder) MarshalHTML(ctx context.Context) ([]byte, error) {
	var (
		sb  bytes.Buffer
		tag = h.Tag("go-plaid-run-script")

		toString = func(name string, v []string) {
			if len(v) == 0 {
				return
			}

			sb.WriteString("(scope) => {\n")
			for i, s := range v {
				if i > 0 {
					sb.WriteString("\n\n")
				}
				s = strings.TrimSpace(s)
				sb.WriteString(s)
				if s[0] == '(' {
					// is function, add call args
					sb.WriteString("(scope)")
				}
			}

			sb.WriteString("\n}")
			tag.Attr(":"+name, sb.Bytes())
			sb.Reset()
		}
	)

	toString("script", r.scripts)
	toString("setup", r.setups)
	toString("before-unmount", r.beforeUnmount)
	toString("unmounted", r.unmounted)

	return tag.MarshalHTML(ctx)
}
