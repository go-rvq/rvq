package helper

import "errors"

var (
	//lint:ignore ST1005 shown to the user as a sentence
	ErrFieldRequired = errors.New("This field is required")
)
