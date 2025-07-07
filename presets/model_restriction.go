package presets

import "github.com/qor5/web/v3"

type RestrictionsGetter interface {
	Handlers() OkHandlers
}

type Restrictor interface {
	Can(ctx *web.EventContext) bool
}

type Restriction[T any] struct {
	dot      T
	handlers OkHandlers
}

func NewRestriction[T any](obj T, f ...func(r *Restriction[T])) (r *Restriction[T]) {
	r = new(Restriction[T])
	r.dot = obj
	for _, f := range f {
		f(r)
	}
	return
}

func (r *Restriction[T]) Dot() T {
	return r.dot
}

func (r *Restriction[T]) Handlers() OkHandlers {
	return r.handlers
}

func (r *Restriction[T]) Handler(f ...OkHandler) *Restriction[T] {
	r.handlers = append(r.handlers, f...)
	return r
}

func (r *Restriction[T]) SetHandlers(handler OkHandlers) *Restriction[T] {
	r.handlers = handler
	return r
}

func (r *Restriction[T]) Can(ctx *web.EventContext) bool {
	return !CallOkHandler(r.handlers, ctx)
}

func (r *Restriction[T]) Insert(g ...RestrictionsGetter) *Restriction[T] {
	for _, g := range g {
		r.handlers = append(r.handlers, OkHandlerFunc(func(ctx *web.EventContext) (ok bool, handled bool) {
			return g.Handlers().Handle(ctx)
		}))
	}
	return r
}

func (r Restriction[T]) Clone(to T) *Restriction[T] {
	r.dot = to
	return &r
}

type ObjRestrictionsGetter interface {
	ObjHandlers() OkObjHandlers
}

type ObjRestrictor interface {
	CanObj(obj any, ctx *web.EventContext) bool
}

type ObjRestriction[T any] struct {
	dot         T
	objHandlers OkObjHandlers
}

func NewObjRestriction[T any](obj T, f ...func(r *ObjRestriction[T])) (r *ObjRestriction[T]) {
	r = new(ObjRestriction[T])
	r.dot = obj
	for _, f := range f {
		f(r)
	}
	return
}

func (r *ObjRestriction[T]) Insert(g ...ObjRestrictionsGetter) *ObjRestriction[T] {
	for _, g := range g {
		r.objHandlers = append(r.objHandlers, OkObjHandlerFunc(func(obj any, ctx *web.EventContext) (ok bool, handled bool) {
			return g.ObjHandlers().Handle(obj, ctx)
		}))
	}
	return r
}

func (r *ObjRestriction[T]) ObjHandler(f ...OkObjHandler) *ObjRestriction[T] {
	r.objHandlers = append(r.objHandlers, f...)
	return r
}

func (r *ObjRestriction[T]) ObjHandlers() OkObjHandlers {
	return r.objHandlers
}

func (r *ObjRestriction[T]) SetObjHandlers(objHandlers OkObjHandlers) *ObjRestriction[T] {
	r.objHandlers = objHandlers
	return r
}

func (r *ObjRestriction[T]) CanObj(obj any, ctx *web.EventContext) bool {
	return !CallOkObjHandlers(r.objHandlers, obj, ctx)
}

func (r ObjRestriction[T]) Clone(to T) *ObjRestriction[T] {
	r.dot = to
	return &r
}

type ListingRestrictionField[T any] struct {
	ListingRestriction *Restriction[T]
}

func (r *ListingRestrictionField[T]) CanList(ctx *web.EventContext) bool {
	return r.ListingRestriction.Can(ctx)
}

type CreatingRestrictionField[T any] struct {
	CreatingRestriction *Restriction[T]
}

func (r *CreatingRestrictionField[T]) CanCreate(ctx *web.EventContext) bool {
	return r.CreatingRestriction.Can(ctx)
}

type EditingRestrictionField[T any] struct {
	EditingRestriction *ObjRestriction[T]
}

func (r *EditingRestrictionField[T]) CanEditObj(obj any, ctx *web.EventContext) bool {
	return r.EditingRestriction.CanObj(obj, ctx)
}

type DetailingRestrictionField[T any] struct {
	DetailingRestriction *ObjRestriction[T]
}

func (r *DetailingRestrictionField[T]) CanDetailObj(obj any, ctx *web.EventContext) bool {
	return r.DetailingRestriction.CanObj(obj, ctx)
}

type DeletingRestrictionField[T any] struct {
	DeletingRestriction *ObjRestriction[T]
}

func (r *DeletingRestrictionField[T]) CanDeleteObj(obj any, ctx *web.EventContext) bool {
	return r.DeletingRestriction.CanObj(obj, ctx)
}
