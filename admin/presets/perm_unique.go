package presets

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

// The parts of a resource of a permission are separated by colons: the
// module, the groups of the menu — each with a "/" at its end, so a group is
// never taken for a model —, the model, a record ("<7>"), a field ("#Title"),
// a section ("$Main"), a nested model…, and, last, the permission asked
// ("@edit") or the action ("!publish"):
//
//	:site/:seo/:seo_config:<7>:@edit   by the groups
//	:seo_config:<7>:@edit              by the unique name
//
// The resource by the unique name — the model's id, a page's path, under the
// module — decides before the one by the groups (perm.Verifier.Prefer).

// RecordPermPart is the part of a resource of a record of id: "<7>".
func RecordPermPart(id any) string {
	return "<" + fmt.Sprint(id) + ">"
}

// GroupPermPart is the part of a resource of the group name: "site/".
func GroupPermPart(name string) string {
	return strcase.ToSnakeWithIgnore(name, ".") + "/"
}

// ActionPerm is the permission of the action name: "!publish". A ":" of the
// name — "geo_ip:update", a kind of job — becomes ".": ":" separates the
// parts of a resource.
func ActionPerm(name string) string {
	return "!" + strings.ReplaceAll(strcase.ToSnake(name), ":", ".")
}

// UniquePermName is the unique name of the model in a permission: its id,
// "posts".
func (mb *ModelBuilder) UniquePermName() string {
	return strcase.ToSnake(mb.id)
}

// UniquePermName is the unique name of the page in a permission: its path,
// "/import".
func (b *HttpPageBuilder) UniquePermName() string {
	return b.path
}
