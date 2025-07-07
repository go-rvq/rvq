package presets

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/web"
)

type ModelInfo struct {
	mb        *ModelBuilder
	parent    *ModelInfo
	parentObj interface{}

	slice interface{}
	index int
	p     *ModelPermissioner
}

func (i *ModelInfo) Root() *ModelInfo {
	for i.parent != nil {
		i = i.parent
	}
	return i
}

func (i *ModelInfo) Builder() *ModelBuilder {
	return i.mb
}

func (i ModelInfo) ChildOf(parent *ModelInfo, obj any) *ModelInfo {
	i.parent = parent
	i.parentObj = obj
	return &i
}

func (i ModelInfo) ItemOf(slice interface{}, index int) *ModelInfo {
	i.slice = slice
	i.index = index
	return &i
}

func (i *ModelInfo) Parent() (*ModelInfo, interface{}) {
	return i.parent, i.parentObj
}

func (i *ModelInfo) Slice() (interface{}, int) {
	return i.slice, i.index
}

func (i *ModelInfo) GetID(obj interface{}) (id ID, err error) {
	return i.mb.RecordID(obj)
}

func (i *ModelInfo) LookupID(obj interface{}) (id ID, level int, err error) {
	if id, err = i.mb.RecordID(obj); err != nil {
		return
	}

	if len(id.Schema.PrimaryFields()) == 0 && i.parent != nil {
		if i.parent != nil {
			return i.parent.LookupID(i.parentObj)
		}
	}
	return
}

func (i *ModelInfo) MustID(obj interface{}) (id ID) {
	id = i.mb.MustRecordID(obj)

	if id.IsZero() && i.parent != nil {
		if i.parent != nil {
			return i.parent.MustID(obj)
		}
	}
	return
}

func (i *ModelInfo) Schema() model.Schema {
	return i.mb.Schema()
}

func (i ModelInfo) ListingHref(parentID ...ID) string {
	s := i.mb.p.prefix + "/" + i.mb.URI()
	for i, id := range parentID {
		s = strings.Replace(s, "{parent_"+strconv.Itoa(i)+"_id}", id.String(), 1)
	}
	return s
}

func (i ModelInfo) ListingHrefParts() (r []any) {
	return append([]any{i.mb.p.prefix}, i.mb.SplitedURI()...)
}

func (i ModelInfo) ListingHrefCtx(ctx *web.EventContext) string {
	return i.ListingHref(ParentsModelID(ctx.R)...)
}

func (i ModelInfo) UpdateHref(parentID ...ID) string {
	s := i.mb.p.prefix + "/" + i.mb.SaveURI()
	for i, id := range parentID {
		s = strings.Replace(s, "{parent_"+strconv.Itoa(i)+"_id}", id.String(), 1)
	}
	return s
}

func (i ModelInfo) UpdateHrefCtx(ctx *web.EventContext) string {
	return i.UpdateHref(ParentsModelID(ctx.R)...)
}

func (i ModelInfo) EditingHref(id any, parentID ...ID) string {
	return path.Join(i.UpdateHref(parentID...), fmt.Sprint(id), "edit")
}

func (i ModelInfo) EditingHrefCtx(ctx *web.EventContext, id any) string {
	return i.EditingHref(id, ParentsModelID(ctx.R)...)
}

func (i ModelInfo) DetailingHref(id any, parentID ...ID) string {
	s := i.ListingHref(parentID...)
	if id != nil {
		s = path.Join(s, fmt.Sprint(id))
	}
	return s
}

func (i ModelInfo) DetailingHrefCtx(ctx *web.EventContext, id any) string {
	return i.DetailingHref(id, ParentsModelID(ctx.R)...)
}

func (i ModelInfo) HasDetailing() bool {
	return i.mb.hasDetailing
}

func (i ModelInfo) PresetsPrefix() string {
	return i.mb.p.prefix
}

func (i ModelInfo) URIName() string {
	return i.mb.uriName
}

func (i ModelInfo) URI() string {
	return i.mb.URI()
}

func (i ModelInfo) Label() string {
	return i.mb.label
}
