package gorm2op

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/data"
	"github.com/go-rvq/rvq/admin/utils/db_utils"
	"github.com/go-rvq/rvq/web"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var wildcardReg = regexp.MustCompile(`[%_]`)

type (
	Mode     uint8
	Preparer func(db *gorm.DB, mode Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB
	Deleter  func(db *gorm.DB, obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error)
	Updator  func(db *gorm.DB, obj interface{}, id model.ID, ctx *web.EventContext) (err error)
	Creator  func(db *gorm.DB, obj interface{}, ctx *web.EventContext) (err error)
	Finder   func(db *gorm.DB, obj interface{}, ctx *web.EventContext) (result any, err error)
)

const (
	Search Mode = 1 << iota
	Create
	Fetch
	FetchTitle
	Update
	Delete

	Read  = Search | Fetch
	Write = Create | Update
)

type deleteAssocEntry struct {
	db    *gorm.DB
	query string
	args  []interface{}
}

type deleteAssocStack struct {
	entries []*deleteAssocEntry
}

func DBCascade(id model.ID, db *gorm.DB) *gorm.DB {

	return nil
}

var Modes = []Mode{Search, Create, Fetch, FetchTitle, Update, Delete}

func (m Mode) Is(other ...Mode) bool {
	for _, mode := range other {
		if mode == m {
			return true
		}
	}
	return false
}

func (m Mode) Has(o Mode) bool {
	return m&o == o
}

func (m Mode) Split() (r []Mode) {
	for _, mode := range Modes {
		if m.Has(mode) {
			r = append(r, mode)
		}
	}
	return
}

var defaultPreparer Preparer = func(db *gorm.DB, mode Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
	return db
}

func DataOperator(db *gorm.DB) (r *DataOperatorBuilder) {
	db = db.Session(&gorm.Session{})
	r = NewCallbacks(&DataOperatorBuilder{
		db:       db,
		preparer: defaultPreparer,
	})
	db.Callback().Delete()
	return
}

type DataOperatorBuilder struct {
	db       *gorm.DB
	preparer Preparer
	deleter  Deleter
	creator  Creator
	updator  Updator
	finder   Finder

	CallbacksRegistrator[*DataOperatorBuilder]
	callbackMergers []CallbackMerger
}

func (b *DataOperatorBuilder) Updator() Updator {
	return b.updator
}

func (b *DataOperatorBuilder) SetUpdator(updator Updator) *DataOperatorBuilder {
	b.updator = updator
	return b
}

func (b *DataOperatorBuilder) Creator() Creator {
	return b.creator
}

func (b *DataOperatorBuilder) SetCreator(creator Creator) *DataOperatorBuilder {
	b.creator = creator
	return b
}

func (b *DataOperatorBuilder) DB() *gorm.DB {
	return b.db
}

func (b *DataOperatorBuilder) SetDB(db *gorm.DB) *DataOperatorBuilder {
	b.db = db
	return b
}

func (b *DataOperatorBuilder) Deleter() Deleter {
	return b.deleter
}

func (b *DataOperatorBuilder) SetDeleter(deleter Deleter) *DataOperatorBuilder {
	b.deleter = deleter
	return b
}

func (b *DataOperatorBuilder) Finder() Finder {
	return b.finder
}

func (b *DataOperatorBuilder) SetFinder(finder Finder) *DataOperatorBuilder {
	b.finder = finder
	return b
}

func (b DataOperatorBuilder) Clone() *DataOperatorBuilder {
	b.SetDot(&b)
	b.db = b.db.Session(&gorm.Session{})
	return &b
}

func (b *DataOperatorBuilder) CloneDataOperator() data.DataOperator {
	return b.Clone()
}

func (b *DataOperatorBuilder) Preparer() Preparer {
	return b.preparer
}

func (b *DataOperatorBuilder) SetPreparer(prepare Preparer) *DataOperatorBuilder {
	b.preparer = prepare
	return b
}

func (b *DataOperatorBuilder) WrapPrepare(f func(old Preparer) Preparer) *DataOperatorBuilder {
	old := b.preparer
	if old == nil {
		old = defaultPreparer
	}
	b.preparer = f(old)
	return b
}

func (b *DataOperatorBuilder) dbCopy() *gorm.DB {
	return b.db.Session(&gorm.Session{})
}

func (b DataOperatorBuilder) tx(f func(b *DataOperatorBuilder) error) error {
	return b.dbCopy().Transaction(func(tx *gorm.DB) error {
		b.db = tx
		return f(&b)
	})
}

func (b *DataOperatorBuilder) Prepare(mode Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
	db := b.dbCopy().Model(obj)

	if params == nil {
		params = &presets.SearchParams{}
	}

	if b.preparer != nil {
		db = b.preparer(db, mode, obj, id, params, ctx)
	}

	modelSchema, _ := schema.Parse(obj, &sync.Map{}, db.NamingStrategy)
	fmt := func(s string) string {
		return strings.ReplaceAll(s, "#TABLE#", modelSchema.Table)
	}

	if db.Dialector.Name() == "sqlite" {
		for _, cond := range params.SQLConditions {
			db = db.Where(strings.Replace(fmt(cond.Query), " ILIKE ", " LIKE ", -1), cond.Args...)
		}
	} else {
		for _, cond := range params.SQLConditions {
			db = db.Where(fmt(cond.Query), cond.Args...)
		}
	}
	return db
}

func (b *DataOperatorBuilder) Search(obj interface{}, params *presets.SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error) {
	if len(params.KeywordColumns) > 0 && len(params.Keyword) > 0 {
		var (
			segs []string
			args []interface{}
		)

		for _, c := range params.KeywordColumns {
			segs = append(segs, fmt.Sprintf("%s ILIKE ?", c))
			kw := wildcardReg.ReplaceAllString(params.Keyword, `\$0`)
			args = append(args, fmt.Sprintf("%%%s%%", kw))
		}

		params.SQLConditions = append(params.SQLConditions, &presets.SQLCondition{
			Query: strings.Join(segs, " OR "),
			Args:  args,
		})
	}

	var (
		db = b.Prepare(Search, obj, model.ID{}, params, ctx)
		do = func(state *CallbackState) (err error) {
			var (
				c   int64
				cdb = state.DB.Count(&c)
			)

			totalCount = int(c)

			if params.MustCount {
				return
			}

			if err = cdb.Error; err != nil {
				state.DB = cdb
				return
			} else if c == 0 {
				state.DB = cdb
				// reset result
				rv := reflect.ValueOf(obj)
				rv.Elem().Set(reflect.MakeSlice(rv.Elem().Type(), 0, 0))
				r = rv.Elem().Interface()
				return
			}

			if params.PerPage > 0 {
				db = db.Limit(int(params.PerPage))
				page := params.Page
				if page == 0 {
					page = 1
				}
				offset := (page - 1) * params.PerPage
				db = db.Offset(int(offset))
			}

			orderBy := params.OrderBy

			if len(orderBy) > 0 {
				db = db.Order(orderBy)
			}

			if b.finder != nil {
				state.Obj, err = b.finder(db, obj, ctx)
			} else {
				db = db.Find(state.Obj)
				err = db.Error
			}

			state.DB = db

			if err != nil {
				return
			}

			if state.Return != nil {
				r = state.Obj
			} else {
				r = reflect.ValueOf(obj).Elem().Interface()
			}

			return
		}
	)

	state := b.NewCallbackState(db, obj, ctx)
	state.SearchParams = params
	cbs := b.GetCallbacks(Search, ctx)

	for _, cbr := range GetContextCallbacks(params.Context) {
		cbs.Merge(&cbr.search)
	}
	err = cbs.
		Build(do).
		Execute(state)
	return
}

func (b *DataOperatorBuilder) Fetch(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
	db := db_utils.ModelIdWhere(b.Prepare(Fetch, obj, id, nil, ctx), obj, id)
	return b.GetCallbacks(Fetch, ctx).
		Build(func(state *CallbackState) (err error) {
			err = state.DB.First(state.Obj).Error
			if err == gorm.ErrRecordNotFound {
				err = presets.ErrRecordNotFound
			}
			return
		}).Execute(b.NewCallbackState(db, obj, ctx))
}

func (b *DataOperatorBuilder) FetchTitle(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
	err = db_utils.ModelIdWhere(b.Prepare(FetchTitle, obj, id, nil, ctx), obj, id).First(obj).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presets.ErrRecordNotFound
		}
		return
	}
	return
}

func (b *DataOperatorBuilder) Save(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
	if id.IsZero() {
		return b.Create(obj, ctx)
	}
	return b.Update(obj, id, ctx)
}

func (b *DataOperatorBuilder) Create(obj interface{}, ctx *web.EventContext) (err error) {
	return b.tx(func(b *DataOperatorBuilder) (err error) {
		var (
			db = b.db.Session(&gorm.Session{})
			do func(state *CallbackState) error
		)
		if b.creator == nil {
			db = b.Prepare(Create, obj, model.ID{}, nil, ctx)
			do = func(state *CallbackState) error {
				return state.DB.Create(state.Obj).Error
			}
		} else {
			do = func(state *CallbackState) error {
				return b.creator(state.DB, state.Obj, state.Ctx)
			}
		}
		return b.GetCallbacks(Create, ctx).
			Build(do).
			Execute(b.NewCallbackState(db, obj, ctx))
	})
}

func (b *DataOperatorBuilder) Update(obj interface{}, id model.ID, ctx *web.EventContext) (err error) {
	return b.tx(func(b *DataOperatorBuilder) (err error) {
		var (
			db = b.db.Session(&gorm.Session{})
			do func(state *CallbackState) error
		)

		if b.updator == nil {
			db = db_utils.ModelIdWhere(b.Prepare(Update, obj, id, nil, ctx), obj, id)
			do = func(state *CallbackState) error {
				return state.DB.Select("*").Updates(state.Obj).Error
			}
		} else {
			do = func(state *CallbackState) error {
				return b.updator(state.DB, state.Obj, id, state.Ctx)
			}
		}

		return b.GetCallbacks(Update, ctx).
			Build(do).
			Execute(b.NewCallbackState(db, obj, ctx))
	})
}

func (b *DataOperatorBuilder) Delete(obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error) {
	return b.tx(func(b *DataOperatorBuilder) (err error) {
		var (
			db = b.Prepare(Delete, obj, id, nil, ctx)
			do = func(state *CallbackState) error {
				return state.DB.Delete(state.Obj).Error
			}
		)
		if b.deleter != nil {
			do = func(state *CallbackState) error {
				return b.deleter(state.DB, state.Obj, id, cascade, ctx)
			}
		}

		db = db_utils.ModelIdWhere(db, obj, id)
		return b.GetCallbacks(Delete, ctx).
			Build(do).
			Execute(b.NewCallbackState(db, obj, ctx))
	})
}

func (b *DataOperatorBuilder) NewCallbackState(db *gorm.DB, obj any, ctx *web.EventContext) *CallbackState {
	sharedDB := db.Session(&gorm.Session{})
	sharedDB.Statement = &gorm.Statement{
		DB:       sharedDB,
		ConnPool: db.Statement.ConnPool,
		Settings: db.Statement.Settings,
		Context:  db.Statement.Context,
	}

	return &CallbackState{
		CommonDB: b.db.Session(&gorm.Session{}),
		SharedDB: sharedDB,
		DB:       db,
		Obj:      obj,
		Ctx:      ctx,
		data:     make(map[any]any),
	}
}

type CallbackState struct {
	CommonDB,
	SharedDB,
	DB *gorm.DB
	SearchParams *presets.SearchParams
	Return       any
	Obj          interface{}
	Ctx          *web.EventContext
	data         map[any]any
	dones        []func() error
}

func (s *CallbackState) Done(f func() error) *CallbackState {
	s.dones = append(s.dones, f)
	return s
}

func (s *CallbackState) Set(key any, value any) *CallbackState {
	s.data[key] = value
	return s
}

func (s *CallbackState) GetOk(key any) (v any, ok bool) {
	v, ok = s.data[key]
	return
}

func (s *CallbackState) Get(key any) any {
	return s.data[key]
}
