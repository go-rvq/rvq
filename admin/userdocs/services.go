package userdocs

import (
	"context"
	"strings"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// StablePathsNode is the guide of the stable paths: the rule of the
// admin's paths and its services (admin.services()).
const StablePathsNode = "guides/stable-paths"

// servicesFunc is admin.services(): the services of the admin
// (presets.Service), from the registry — never a copy —, each with its
// path, its routes, how they are reached; the guide of the stable paths
// (guides/stable-paths) shows them.
func (b *Builder) servicesFunc(ctx *web.EventContext, _ ...string) (string, error) {
	return ServicesMarkdown(ctx.Context(), b.p), nil
}

// ServicesMarkdown is the services of p in the language of ctx, in
// Markdown: a section of each, its routes in a table.
func ServicesMarkdown(ctx context.Context, p *presets.Builder) string {
	m := GetMessages(ctx)
	svcs := p.Services()
	if len(svcs) == 0 {
		return m.ServicesNone + "\n"
	}
	var sb strings.Builder
	for _, s := range svcs {
		sb.WriteString("### " + mdText(s.GetTitle(ctx)))
		if !s.IsMountOnly() {
			sb.WriteString(" — `" + s.Path() + "`")
		}
		sb.WriteString("\n\n")
		if d := s.GetDescription(ctx); d != "" {
			sb.WriteString(mdText(d) + "\n\n")
		}
		if s.IsOuter() {
			sb.WriteString("*" + m.ServicesOuter + "*\n\n")
		}
		if rs := s.Routes(); len(rs) > 0 {
			sb.WriteString("| " + m.ServicesMethod + " | " + m.ServicesPath + " | " + m.ServicesAuth + " | " + m.ServicesWhat + " |\n| --- | --- | --- | --- |\n")
			for _, r := range rs {
				var what string
				if r.Description != nil {
					what = r.Description(ctx)
				}
				sb.WriteString("| " + cell(r.Method) + " | `" + cell(s.RoutePath(r)) + "` | " + cell(mdText(r.GetAuth(ctx))) + " | " + cell(mdText(what)) + " |\n")
			}
			sb.WriteString("\n")
		}
		if ms := s.Mounts(); len(ms) > 0 {
			sb.WriteString(m.ServicesMounts + ":\n\n")
			for _, mp := range ms {
				sb.WriteString("- `" + mp + "`\n")
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// mdText is v as text of Markdown: a "<" outside its code spans is not taken
// for a tag ("<name>.git").
func mdText(v string) string {
	parts := strings.Split(v, "`")
	for i := 0; i < len(parts); i += 2 {
		parts[i] = strings.ReplaceAll(parts[i], "<", "&lt;")
	}
	return strings.Join(parts, "`")
}
