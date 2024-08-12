package data

import (
	"net/url"

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
	Fetch(obj interface{}, id string, ctx *web.EventContext) (err error)
	FetchTitle(obj interface{}, id string, ctx *web.EventContext) (err error)
	Save(obj interface{}, id string, ctx *web.EventContext) (err error)
	Delete(obj interface{}, id string, ctx *web.EventContext) (err error)
	CloneDataOperator() DataOperator
}

type SearchParams struct {
	KeywordColumns []string
	Keyword        string
	SQLConditions  []*SQLCondition
	PerPage        int64
	Page           int64
	OrderBy        string
	PageQuery      url.Values
}

type SQLCondition struct {
	Query string
	Args  []interface{}
}
