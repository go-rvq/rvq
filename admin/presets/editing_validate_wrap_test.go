package presets

import (
	"errors"
	"testing"

	"github.com/go-rvq/rvq/web"
)

// WrapPostValidate wraps the POST validation — it used to wrap the pre one,
// dropping a post validation set before and running the pre one twice.
func TestWrapPostValidateWrapsThePostValidation(t *testing.T) {
	var b EditingBuilder
	var calls []string
	b.PreValidate(func(*web.EventContext, any) error { calls = append(calls, "pre"); return nil })
	b.PostValidate(func(*web.EventContext, any) error { calls = append(calls, "post"); return errors.New("post") })
	b.WrapPostValidate(func(old func(*web.EventContext, any) error) func(*web.EventContext, any) error {
		return func(ctx *web.EventContext, obj any) error {
			calls = append(calls, "wrapped")
			return old(ctx, obj)
		}
	})

	if err := b.postValidate(nil, nil); err == nil || err.Error() != "post" {
		t.Fatalf("err = %v, want the original post validation's", err)
	}
	if got := len(calls); got != 2 || calls[0] != "wrapped" || calls[1] != "post" {
		t.Errorf("calls = %v, want [wrapped post]", calls)
	}
}
