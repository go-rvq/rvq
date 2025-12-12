package script

import (
	"errors"
	"fmt"

	"github.com/gad-lang/gad/parser/source"
)

var (
	ErrScriptFailure = errors.New("script failure")
)

type ScriptErrorType string

const (
	ScriptErrorTypeParse   ScriptErrorType = "PARSE"
	ScriptErrorTypeCompile ScriptErrorType = "COMPILE"
	ScriptErrorTypeRun     ScriptErrorType = "RUN"
)

type ScriptError struct {
	CausedBy error
	Type     ScriptErrorType
	Pos      source.FilePos
	Message  string
}

func NewTemplateScriptError(causedBy error, Type ScriptErrorType, pos source.FilePos, message string) *ScriptError {
	return &ScriptError{CausedBy: causedBy, Type: Type, Pos: pos, Message: message}
}

func (e *ScriptError) Error() string {
	return fmt.Sprintf("Script Error [%s] at [%d:%d]: %s", e.Type, e.Pos.Line, e.Pos.Column, e.Message)
}

func (e *ScriptError) Cause() error {
	return e.CausedBy
}

type ScriptTypedError struct {
	CausedBy *ScriptError
	Message  string
}

func NewScriptTypedError(causedBy *ScriptError, message string) *ScriptTypedError {
	return &ScriptTypedError{CausedBy: causedBy, Message: message}
}

func (e *ScriptTypedError) Error() string {
	return e.Message
}

func (e *ScriptTypedError) Cause() error {
	return e.CausedBy
}
