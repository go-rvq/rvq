package presets

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/perm"
	"github.com/sunfmin/reflectutils"
)

type ModelInfo struct {
	mb        *ModelBuilder
	parent    *ModelInfo
	parentObj interface{}

	slice interface{}
	index int
}

func (b ModelInfo) ChildOf(parent *ModelInfo, obj any) *ModelInfo {
	b.parent = parent
	b.parentObj = obj
	return &b
}

func (b ModelInfo) ItemOf(slice interface{}, index int) *ModelInfo {
	b.slice = slice
	b.index = index
	return &b
}

func (b *ModelInfo) Parent() (*ModelInfo, interface{}) {
	return b.parent, b.parentObj
}

func (b *ModelInfo) Slice() (interface{}, int) {
	return b.slice, b.index
}

func (b *ModelInfo) GetID(obj interface{}) (id ID, level int, err error) {
	id.Value, err = reflectutils.Get(obj, "ID")
	for err == reflectutils.NoSuchFieldError && b.parent != nil {
		if b.parent != nil {
			return b.parent.GetID(obj)
		}
	}
	return
}

func (b ModelInfo) ListingHref(parentID ...ID) string {
	s := b.mb.p.prefix + "/" + b.mb.URI()
	for i, id := range parentID {
		s = strings.Replace(s, "{parent_"+strconv.Itoa(i)+"_id}", id.String(), 1)
	}
	return s
}

func (b ModelInfo) ListingHrefCtx(ctx *web.EventContext) string {
	return b.ListingHref(ParentsModelID(ctx.R)...)
}

func (b ModelInfo) UpdateHref(parentID ...ID) string {
	s := b.mb.p.prefix + "/" + b.mb.SaveURI()
	for i, id := range parentID {
		s = strings.Replace(s, "{parent_"+strconv.Itoa(i)+"_id}", id.String(), 1)
	}
	return s
}

func (b ModelInfo) UpdateHrefCtx(ctx *web.EventContext) string {
	return b.UpdateHref(ParentsModelID(ctx.R)...)
}

func (b ModelInfo) EditingHref(id any, parentID ...ID) string {
	return path.Join(b.UpdateHref(parentID...), fmt.Sprint(id), "edit")
}

func (b ModelInfo) EditingHrefCtx(ctx *web.EventContext, id any) string {
	return b.EditingHref(id, ParentsModelID(ctx.R)...)
}

func (b ModelInfo) DetailingHref(id any, parentID ...ID) string {
	s := b.ListingHref(parentID...)
	if id != nil {
		s = path.Join(s, fmt.Sprint(id))
	}
	return s
}

func (b ModelInfo) DetailingHrefCtx(ctx *web.EventContext, id any) string {
	return b.DetailingHref(id, ParentsModelID(ctx.R)...)
}

func (b ModelInfo) HasDetailing() bool {
	return b.mb.hasDetailing
}

func (b ModelInfo) PresetsPrefix() string {
	return b.mb.p.prefix
}

func (b ModelInfo) URIName() string {
	return b.mb.uriName
}

func (b ModelInfo) URI() string {
	return b.mb.URI()
}

func (b ModelInfo) Label() string {
	return b.mb.label
}

func (b ModelInfo) Verifier() *perm.Verifier {
	mb := b.mb.GetVerifierModel()
	v := mb.p.verifier.Spawn()
	if mb.menuGroupName != "" {
		v.SnakeOn("mg_" + mb.menuGroupName)
	}
	return v.SnakeOn(mb.uriName)
}

func (b ModelInfo) CanUpdate(r *http.Request, obj interface{}) bool {
	return b.Verifier().Do(PermUpdate).ObjectOn(obj).WithReq(r).IsAllowed() == nil
}

func (b ModelInfo) CanRead(r *http.Request, obj interface{}) bool {
	return b.Verifier().Do(PermGet).ObjectOn(obj).WithReq(r).IsAllowed() == nil
}

func (b ModelInfo) CanDelete(r *http.Request) bool {
	return b.Verifier().Do(PermDelete).WithReq(r).IsAllowed() == nil
}

func (b ModelInfo) CanCreate(r *http.Request) bool {
	return b.Verifier().Do(PermCreate).WithReq(r).IsAllowed() == nil
}
