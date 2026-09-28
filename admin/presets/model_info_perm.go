package presets

import (
	"fmt"
	"net/http"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/x/perm"
)

type ModelPermissioner struct {
	mb     *ModelBuilder
	parent *ModelPermissioner
}

func (i *ModelInfo) Permissioner() *ModelPermissioner {
	if i.parent != nil {
		return &ModelPermissioner{
			mb:     i.mb,
			parent: i.parent.Permissioner(),
		}
	}
	return i.mb.permissioner
}

func (p *ModelPermissioner) Verifier(id ID, parentID ...ID) (v *perm.Verifier) {
	if p.parent != nil {
		v = p.parent.Verifier(id, parentID...)
	} else {
		if p.mb.parent != nil {
			var pid ID
			if !p.mb.parent.singleton && len(parentID) > 0 {
				pid = parentID[0]
				parentID = parentID[1:]
			}
			v = p.mb.parent.Info().Permissioner().Verifier(pid, parentID...)
		} else {
			mb := p.mb.GetVerifierModel()
			v = mb.p.verifier.Spawn()
			for _, name := range mb.MenuGroup().PathNames() {
				v.SnakeOn(name)
			}
		}
	}

	v.SnakeOn(p.mb.id)

	if !p.mb.singleton && !id.IsZero() {
		if id := id.GetValue("ID"); id != nil {
			v.SnakeOn(fmt.Sprint(id))
		}
	}
	return
}

func (p *ModelPermissioner) ListVerifier(parentID ...ID) (v *perm.Verifier) {
	if p.mb.parent != nil {
		var pid ID
		if !p.mb.parent.singleton && len(parentID) > 0 {
			pid = parentID[0]
			parentID = parentID[1:]
		}
		v = p.mb.parent.Info().Permissioner().Verifier(pid, parentID...)
	} else {
		mb := p.mb.GetVerifierModel()
		v = mb.p.verifier.Spawn()
		for _, name := range mb.MenuGroup().PathNames() {
			v.SnakeOn(name)
		}
	}

	v.SnakeOn(p.mb.id)
	return
}

func (p *ModelPermissioner) ReqObjector(r *http.Request, obj any) *perm.Verifier {
	return p.Objector(r, p.mb.MustRecordID(obj), ParentsModelID(r)...)
}

func (p *ModelPermissioner) ReqObjectFielder(r *http.Request, obj any, field string) *perm.Verifier {
	return p.ReqObjector(r, obj).On(FieldPerm(field))
}

// ReqObjectFieldReader verifies reading a field of obj, and
// ReqObjectFieldWriter writing it: creating, when obj is a record not saved
// yet, else updating. field is the field's path from the record: "Title",
// "Tags.Name" (no indexes) — the same whether the field is written by the
// edit form or by a section of the detail page.
func (p *ModelPermissioner) ReqObjectFieldReader(r *http.Request, obj any, field string) *perm.Verifier {
	return p.ReqObjectFielder(r, obj, field).Do(PermGet)
}

func (p *ModelPermissioner) ReqObjectFieldWriter(r *http.Request, obj any, field string) *perm.Verifier {
	return p.ReqObjectFielder(r, obj, field).Do(p.writeAction(obj))
}

// ReqObjectSectionReader verifies seeing the section name of obj's detail
// page, and ReqObjectSectionWriter editing it in place (creating or updating,
// as for a field).
func (p *ModelPermissioner) ReqObjectSectionReader(r *http.Request, obj any, name string) *perm.Verifier {
	return p.ReqObjector(r, obj).On(SectionPerm(name)).Do(PermGet)
}

func (p *ModelPermissioner) ReqObjectSectionWriter(r *http.Request, obj any, name string) *perm.Verifier {
	return p.ReqObjector(r, obj).On(SectionPerm(name)).Do(p.writeAction(obj))
}

// writeAction is what writing obj is: creating a record not saved yet (no
// key), else updating it. A singleton is always updated.
func (p *ModelPermissioner) writeAction(obj any) string {
	if !p.mb.singleton && p.mb.MustRecordID(obj).IsZero() {
		return PermCreate
	}
	return PermUpdate
}

func (p *ModelPermissioner) ReqObjectFieldConder(r *http.Request, obj any, ok bool, field string) *perm.Verifier {
	v := p.ReqObjector(r, obj)
	if ok {
		return v.SnakeOn(FieldPerm(field))
	}
	return v
}

func (p *ModelPermissioner) ReqObjectUpdater(r *http.Request, obj any) *perm.Verifier {
	return p.Updater(r, p.mb.MustRecordID(obj), ParentsModelID(r)...)
}

func (p *ModelPermissioner) ReqObjectReader(r *http.Request, obj any) *perm.Verifier {
	return p.Reader(r, p.mb.MustRecordID(obj), ParentsModelID(r)...)
}

func (p *ModelPermissioner) ReqObjectDeleter(r *http.Request, obj any) *perm.Verifier {
	return p.Deleter(r, p.mb.MustRecordID(obj), ParentsModelID(r)...)
}

func (p *ModelPermissioner) ReqObjectWithRelatedDeleter(r *http.Request, obj any) *perm.Verifier {
	return p.WithRelatedDeleter(r, p.mb.MustRecordID(obj), ParentsModelID(r)...)
}

func (p *ModelPermissioner) ReqObjectActioner(r *http.Request, obj any, action string) *perm.Verifier {
	return p.ReqObjectReader(r, obj).Do(action)
}

func (p *ModelPermissioner) ReqCreator(r *http.Request) *perm.Verifier {
	return p.Creator(r, ParentsModelID(r)...)
}

func (p *ModelPermissioner) Actioner(r *http.Request, action string, id ID, parentID ...ID) *perm.Verifier {
	return p.Reader(r, id, parentID...).SnakeDo(action)
}

func (p *ModelPermissioner) ReqList(r *http.Request) *perm.Verifier {
	return p.List(r, ParentsModelID(r)...)
}

func (p *ModelPermissioner) ReqListDo(r *http.Request, action string) *perm.Verifier {
	return p.ReqList(r).Do(action)
}

func (p *ModelPermissioner) ReqLister(r *http.Request) *perm.Verifier {
	return p.Lister(r, ParentsModelID(r)...)
}

func (p *ModelPermissioner) ReqListActioner(r *http.Request, action string) *perm.Verifier {
	return p.ReqLister(r).Do(action)
}

func (p *ModelPermissioner) Objector(r *http.Request, id ID, parentID ...ID) *perm.Verifier {
	return p.Verifier(id, parentID...).WithReq(r)
}

func (p *ModelPermissioner) ObjectDo(r *http.Request, action string, id ID, parentID ...ID) *perm.Verifier {
	return p.Objector(r, id, parentID...).Do(action)
}

func (p *ModelPermissioner) List(r *http.Request, parentID ...ID) *perm.Verifier {
	return p.ListVerifier(parentID...).WithReq(r)
}

func (p *ModelPermissioner) ListDo(r *http.Request, action string, parentID ...ID) *perm.Verifier {
	return p.List(r, parentID...).Do(action)
}

func (p *ModelPermissioner) ObjectFielder(r *http.Request, field string, id ID, parentID ...ID) *perm.Verifier {
	return p.ObjectDo(r, FieldPerm(field), id, parentID...)
}

func (p *ModelPermissioner) ObjectReadActioner(r *http.Request, action string, id model.ID, parentID ...ID) *perm.Verifier {
	return p.Reader(r, id, parentID...).Do(action)
}

func (p *ModelPermissioner) Updater(r *http.Request, id ID, parentID ...ID) *perm.Verifier {
	return p.ObjectDo(r, PermUpdate, id, parentID...)
}

func (p *ModelPermissioner) Reader(r *http.Request, id ID, parentID ...ID) *perm.Verifier {
	return p.ObjectDo(r, PermGet, id, parentID...)
}

func (p *ModelPermissioner) Deleter(r *http.Request, id ID, parentID ...ID) *perm.Verifier {
	return p.ObjectDo(r, PermDelete, id, parentID...)
}

func (p *ModelPermissioner) WithRelatedDeleter(r *http.Request, id ID, parentID ...ID) *perm.Verifier {
	return p.ObjectDo(r, PermDeleteWithRelated, id, parentID...)
}

func (p *ModelPermissioner) Creator(r *http.Request, parentID ...ID) *perm.Verifier {
	return p.ListDo(r, PermCreate, parentID...)
}

func (p *ModelPermissioner) Lister(r *http.Request, parentID ...ID) *perm.Verifier {
	return p.ListDo(r, PermList, parentID...)
}

func (p *ModelPermissioner) Default() *perm.Verifier {
	var (
		listing = p.mb.singleton
		parents []ID
		id      = ID{
			Fields: []model.Field{model.SingleField("ID")},
			Values: []any{"*"},
		}
	)

	mb := p.mb
	for mb.parent != nil {
		parents = append(parents, id)
		mb = mb.parent
	}

	if listing {
		id = ID{}
	}

	return p.Verifier(id, parents...)
}

// CanReadField and CanWriteField say whether the request may read, or write,
// the field of obj (a record of this model; field its path from it). See
// ReqObjectFieldReader and ReqObjectFieldWriter.
func (i *ModelInfo) CanReadField(r *http.Request, obj any, field string) bool {
	return i.Permissioner().ReqObjectFieldReader(r, obj, field).Allowed()
}

func (i *ModelInfo) CanWriteField(r *http.Request, obj any, field string) bool {
	return i.Permissioner().ReqObjectFieldWriter(r, obj, field).Allowed()
}

// CanReadSection and CanWriteSection say whether the request may see, or edit
// in place, the section name of obj's detail page. See ReqObjectSectionReader.
func (i *ModelInfo) CanReadSection(r *http.Request, obj any, name string) bool {
	return i.Permissioner().ReqObjectSectionReader(r, obj, name).Allowed()
}

func (i *ModelInfo) CanWriteSection(r *http.Request, obj any, name string) bool {
	return i.Permissioner().ReqObjectSectionWriter(r, obj, name).Allowed()
}
