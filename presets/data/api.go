package data

import "github.com/qor5/web/v3"

// Data Layer
type DataOperator interface {
	Search(obj interface{}, params *SearchParams, ctx *web.EventContext) (r interface{}, totalCount int, err error)
	// return ErrRecordNotFound if record not found
	Fetch(obj interface{}, id string, ctx *web.EventContext) (r interface{}, err error)
	Save(obj interface{}, id string, ctx *web.EventContext) (err error)
	Delete(obj interface{}, id string, ctx *web.EventContext) (err error)
	Clone() DataOperator
}
