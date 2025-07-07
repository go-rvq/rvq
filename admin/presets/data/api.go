package data

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/qor5/admin/v3/model"
	"github.com/qor5/web/v3"
)

type (
	Searcher func(obj interface{}, params *SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error)
	Fetcher  func(obj interface{}, id string, ctx *web.EventContext) (r interface{}, err error)
	Deleter  func(obj interface{}, id string, ctx *web.EventContext) (err error)
	Saver    func(obj interface{}, id string, ctx *web.EventContext) (err error)
)

// Data Layer
type DataOperator interface {
	Search(obj interface{}, params *SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error)
	// return ErrRecordNotFound if record not found
	Fetch(obj interface{}, id model.ID, ctx *web.EventContext) (err error)
	FetchTitle(obj interface{}, id model.ID, ctx *web.EventContext) (err error)
	Save(obj interface{}, id model.ID, ctx *web.EventContext) (err error)
	Create(obj interface{}, ctx *web.EventContext) (err error)
	Delete(obj interface{}, id model.ID, cascade bool, ctx *web.EventContext) (err error)
	CloneDataOperator() DataOperator
	Schema(model any) (schema model.Schema, err error)
}

type SQLConditions []*SQLCondition

type SearchParams struct {
	KeywordColumns []string
	Keyword        string
	SQLConditions  SQLConditions
	PerPage        int64
	Page           int64
	OrderBy        string
	Query          web.Query
	Context        context.Context
	MustCount      bool
}

func (p *SearchParams) ContextValue(key interface{}) (value any) {
	if p.Context == nil {
		return nil
	}
	return p.Context.Value(key)
}

func (p *SearchParams) SetContextValue(key, value interface{}) {
	if p.Context == nil {
		p.Context = context.Background()
	}
	p.Context = context.WithValue(p.Context, key, value)
}

func (p *SearchParams) WithContextValue(key, value interface{}) (deleter func()) {
	if p.Context == nil {
		p.Context = context.Background()
	}

	if ptr := web.GetContextValuer(p.Context, key); ptr != nil {
		return ptr.With(value)
	}

	p.Context = context.WithValue(p.Context, key, value)

	return func() {
		p.Context = web.GetContextValuer(p.Context, key).Top().Delete()
	}
}

func (p *SearchParams) Where(query string, args ...any) *SearchParams {
	p.SQLConditions = append(p.SQLConditions, &SQLCondition{
		Query: query,
		Args:  args,
	})
	return p
}

func (p *SearchParams) WhereModelID(id model.ID, withoutKeys ...string) *SearchParams {
	var (
		q []string
		v []any
	)

	for i, field := range id.Fields {
		if !slices.Contains(withoutKeys, field.Name()) {
			q = append(q, fmt.Sprintf("%s = ?", field.QuotedFullDBName()))
			v = append(v, id.Values[i])
		}
	}

	if len(q) > 0 {
		p.Where(strings.Join(q, " AND "), v...)
	}

	return p
}

func (p *SearchParams) WhereModelIDs(ids model.IDSlice, withoutKeys ...string) *SearchParams {
	var (
		allq []string
		args []any
	)

	for _, mid := range ids {
		var (
			q []string
		)

		for i, field := range mid.Fields {
			if !slices.Contains(withoutKeys, field.Name()) {
				q = append(q, fmt.Sprintf("%s = ?", field.QuotedFullDBName()))
				args = append(args, mid.Values[i])
			}
		}

		allq = append(allq, "("+strings.Join(q, " AND ")+")")
	}

	return p.Where("( "+strings.Join(allq, " OR ")+" )", args...)
}

type SQLCondition struct {
	Query string
	Args  []interface{}
}
