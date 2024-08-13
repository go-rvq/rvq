package gorm2op

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/data"
	"github.com/qor5/web/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var wildcardReg = regexp.MustCompile(`[%_]`)

type (
	Mode     uint8
	Preparer func(db *gorm.DB, mode Mode, obj interface{}, id string, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB
	Deleter  func(db *gorm.DB, obj interface{}, id string, ctx *web.EventContext) (err error)
	Creator  func(db *gorm.DB, obj interface{}, ctx *web.EventContext) (err error)
)

const (
	Search Mode = iota
	Create
	Fetch
	FetchTitle
	Update
	Delete
)

func (m Mode) Is(other ...Mode) bool {
	for _, mode := range other {
		if mode == m {
			return true
		}
	}
	return false
}

var defaultPreparer Preparer = func(db *gorm.DB, mode Mode, obj interface{}, id string, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
	return db
}

func DataOperator(db *gorm.DB) (r *DataOperatorBuilder) {
	r = &DataOperatorBuilder{db: db, preparer: defaultPreparer}
	return
}

type DataOperatorBuilder struct {
	db         *gorm.DB
	preparer   Preparer
	deleter    Deleter
	creator    Creator
	updator    Deleter
	postCreate []func(db *gorm.DB, obj interface{}, ctx *web.EventContext) (err error)
}

func (b *DataOperatorBuilder) Updator() Deleter {
	return b.updator
}

func (b *DataOperatorBuilder) SetUpdator(updator Deleter) *DataOperatorBuilder {
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

func (b DataOperatorBuilder) Clone() *DataOperatorBuilder {
	return &b
}

func (b DataOperatorBuilder) CloneDataOperator() data.DataOperator {
	return &b
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

func (b *DataOperatorBuilder) PostCreate(cb func(db *gorm.DB, obj interface{}, ctx *web.EventContext) (err error)) *DataOperatorBuilder {
	b.postCreate = append(b.postCreate, cb)
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

func (b *DataOperatorBuilder) Prepare(mode Mode, obj interface{}, id string, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
	db := b.dbCopy().Model(obj)

	if params == nil {
		params = &presets.SearchParams{}
	}

	if b.preparer != nil {
		db = b.preparer(db, mode, obj, id, params, ctx)
	}

	if db.Dialector.Name() == "sqlite" {
		for _, cond := range params.SQLConditions {
			db = db.Where(strings.Replace(cond.Query, " ILIKE ", " LIKE ", -1), cond.Args...)
		}
	} else {
		for _, cond := range params.SQLConditions {
			db = db.Where(cond.Query, cond.Args...)
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
		c   int64
		wh  = b.Prepare(Search, obj, "", params, ctx)
		cdb = wh.Session(&gorm.Session{}).Count(&c)
	)

	if err = cdb.Error; err != nil {
		return
	}

	totalCount = int(c)

	if params.PerPage > 0 {
		wh = wh.Limit(int(params.PerPage))
		page := params.Page
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * params.PerPage
		wh = wh.Offset(int(offset))
	}

	orderBy := params.OrderBy
	if len(orderBy) > 0 {
		wh = wh.Order(orderBy)
	}

	err = wh.Find(obj).Error
	if err != nil {
		return
	}
	r = reflect.ValueOf(obj).Elem().Interface()
	return
}

func (b *DataOperatorBuilder) primarySluggerWhere(wh *gorm.DB, obj interface{}, id string) *gorm.DB {
	if id == "" {
		return wh
	}

	if slugger, ok := obj.(presets.SlugDecoder); ok {
		cs := slugger.PrimaryColumnValuesBySlug(id)
		for key, value := range cs {
			wh = wh.Where(fmt.Sprintf("%s = ?", key), value)
		}
	} else {
		cond := "id = ?"
		if tb, ok := obj.(schema.Tabler); ok {
			cond = fmt.Sprintf(`%q.%s`, tb.TableName(), cond)
		}
		wh = wh.Where(cond, id)
	}

	return wh
}

func (b *DataOperatorBuilder) Fetch(obj interface{}, id string, ctx *web.EventContext) (err error) {
	db := b.primarySluggerWhere(b.Prepare(Fetch, obj, id, nil, ctx), obj, id)
	db = db.First(obj)
	err = db.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presets.ErrRecordNotFound
		}
		return
	}
	return
}

func (b *DataOperatorBuilder) FetchTitle(obj interface{}, id string, ctx *web.EventContext) (err error) {
	err = b.primarySluggerWhere(b.Prepare(FetchTitle, obj, id, nil, ctx), obj, id).First(obj).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return presets.ErrRecordNotFound
		}
		return
	}
	return
}

func (b *DataOperatorBuilder) Save(obj interface{}, id string, ctx *web.EventContext) (err error) {
	return b.tx(func(b *DataOperatorBuilder) (err error) {
		if id == "" {
			if b.creator != nil {
				db := b.dbCopy()
				if err = b.creator(db, obj, ctx); err != nil {
					return
				}
				for _, f := range b.postCreate {
					if err = f(db, obj, ctx); err != nil {
						break
					}
				}
			} else {
				db := b.Prepare(Create, obj, id, nil, ctx)
				if err = db.Create(obj).Error; err != nil {
					return
				}
				for _, f := range b.postCreate {
					if err = f(db, obj, ctx); err != nil {
						break
					}
				}
			}
			return
		}

		if b.updator != nil {
			return b.updator(b.db, obj, id, ctx)
		}
		err = b.primarySluggerWhere(b.Prepare(Update, obj, id, nil, ctx), obj, id).Save(obj).Error
		return
	})
}

func (b *DataOperatorBuilder) Delete(obj interface{}, id string, ctx *web.EventContext) (err error) {
	return b.tx(func(b *DataOperatorBuilder) (err error) {
		db := b.Prepare(Delete, obj, id, nil, ctx)
		if b.deleter != nil {
			return b.deleter(db, obj, id, ctx)
		}
		err = b.primarySluggerWhere(db, obj, id).Delete(obj).Error
		return
	})
}
