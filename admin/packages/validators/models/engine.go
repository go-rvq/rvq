package models

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/gad-lang/gad"
	"github.com/gad-lang/gad/encoder"
)

// The validation contract (GAD): the script receives the value to analyze as its
// `value` parameter. It signals success by returning a truthy value (typically
// `return true`) and signals a validation failure by throwing a message key —
// `throw "MESSAGE_KEY"` — which the engine translates to the user's language via
// the validator's Messages. A plain falsy return (no throw) is a generic
// failure. Any other runtime/compile error is a real script error.

// compileScript compiles a GAD validation script (self-contained, using only the
// default builtins and the `value` param) into bytecode.
func compileScript(script string) (*gad.Bytecode, error) {
	st := gad.NewSymbolTable(gad.NewBuiltins().Build().Builtins().NameSet)
	_, bc, err := gad.Compile(st, []byte(script), gad.CompileOptions{})
	if err != nil {
		return nil, fmt.Errorf("validator script: %w", err)
	}
	return bc, nil
}

// EncodeScript compiles script and returns its serialized bytecode, for caching
// in the Validator.Bytecode column. Our scripts use no external modules, so the
// bytecode is self-contained.
func EncodeScript(script string) ([]byte, error) {
	bc, err := compileScript(script)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err = encoder.EncodeBytecodeTo(encoder.NewWriteContext(context.Background(), encoder.NewWriter(&buf)), bc); err != nil {
		return nil, fmt.Errorf("encode bytecode: %w", err)
	}
	return buf.Bytes(), nil
}

// decodeBytecode deserializes cached bytecode.
func decodeBytecode(data []byte) (*gad.Bytecode, error) {
	return encoder.DecodeBytecodeFrom(encoder.NewReadContext(encoder.NewReader(bytes.NewReader(data))))
}

// evalBytecode runs compiled bytecode with value as the `value` param. It returns
// ok=true when the script returns truthy; when the script throws a message key,
// ok=false and msgKey is the thrown key; a plain falsy return yields ok=false with
// an empty msgKey. A non-thrown error is a real script error.
func evalBytecode(bc *gad.Bytecode, value string) (ok bool, msgKey string, err error) {
	vm := gad.NewVM(gad.NewBuiltins().Build(), bc)
	ret, rerr := vm.Run(gad.Str(value))
	if rerr != nil {
		// a `throw "KEY"` surfaces as *gad.Error with an empty Name (system
		// errors carry a Name); treat its Message as the message key.
		var ge *gad.Error
		if errors.As(rerr, &ge) && ge.Name == "" && ge.Message != "" {
			return false, ge.Message, nil
		}
		return false, "", fmt.Errorf("validator script: %w", rerr)
	}
	if ret == nil {
		return false, "", nil
	}
	if b, isBool := ret.(gad.Bool); isBool {
		return bool(b), "", nil
	}
	return !ret.IsFalsy(), "", nil
}

// Eval runs the validator against value, using the cached Bytecode when available
// (falling back to compiling the effective script). It returns ok, the thrown
// message key (when the script threw one) and any real script error.
func (v *Validator) Eval(value string) (ok bool, msgKey string, err error) {
	var bc *gad.Bytecode
	if len(v.Bytecode) > 0 {
		if decoded, derr := decodeBytecode(v.Bytecode); derr == nil {
			bc = decoded
		}
	}
	if bc == nil {
		if bc, err = compileScript(v.EffectiveValue()); err != nil {
			return false, "", err
		}
	}
	return evalBytecode(bc, value)
}

// Validate reports whether value is valid (ignoring the specific message key).
func (v *Validator) Validate(value string) (bool, error) {
	ok, _, err := v.Eval(value)
	return ok, err
}

// ErrInvalid is returned by Check when the value fails validation and no specific
// message is available.
var ErrInvalid = errors.New("invalid value")

// MsgKeyInvalid is the fallback message key used when the script signals failure
// without throwing a specific key.
const MsgKeyInvalid = "invalid"

// Check runs the validator and, when value is invalid, returns an error whose
// message is the validator's translation (for lang, falling back to DefaultLang)
// of the thrown message key — or of MsgKeyInvalid for a plain falsy return.
// Real script errors are returned as-is.
func (v *Validator) Check(value, lang string) error {
	ok, msgKey, err := v.Eval(value)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	if msgKey == "" {
		msgKey = MsgKeyInvalid
	}
	if msg := v.Message(lang, msgKey); msg != msgKey {
		return errors.New(msg)
	}
	return ErrInvalid
}
