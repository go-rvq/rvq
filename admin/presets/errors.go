package presets

import "errors"

var (
	ErrRecordNotFound         = errors.New("record not found")
	ErrDeleteRecordNotAllowed = errors.New("delete record not allowed")
	ErrUpdateRecordNotAllowed = errors.New("update record not allowed")
	ErrReadRecordNotAllowed   = errors.New("read record not allowed")
	ErrCreateRecordNotAllowed = errors.New("create record not allowed")
	ErrActionNotAllowed       = errors.New("action not allowed")
)
