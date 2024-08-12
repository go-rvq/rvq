package presets

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
	_ "unsafe"

	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"
)

func RecoverPrimaryColumnValuesBySlug(dec SlugDecoder, slug string) (r map[string]string, err error) {
	defer func() {
		if e := recover(); e != nil {
			r = nil
			err = fmt.Errorf("wrong slug: %v", slug)
		}
	}()
	r = dec.PrimaryColumnValuesBySlug(slug)
	return r, nil
}

func ShowMessage(r *web.EventResponse, msg string, color string) {
	if msg == "" {
		return
	}

	if color == "" {
		color = "success"
	}

	web.AppendRunScripts(r, fmt.Sprintf(
		`vars.presetsMessage = { show: true, message: %s, color: %s}`,
		h.JSONString(msg), h.JSONString(color)))
}

func copyURLWithQueriesRemoved(u *url.URL, qs ...string) *url.URL {
	newU, _ := url.Parse(u.String())
	newQuery := newU.Query()
	for _, k := range qs {
		newQuery.Del(k)
	}
	newU.RawQuery = newQuery.Encode()
	return newU
}

func GetOverlay(ctx *web.EventContext) actions.OverlayMode {
	return actions.OverlayMode(ctx.R.FormValue(ParamOverlay))
}

func isInDialogFromQuery(ctx *web.EventContext) bool {
	return actions.OverlayMode(ctx.R.FormValue(ParamOverlay)) == actions.Dialog
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func OkOrError(ok bool, err error) error {
	if !ok {
		return err
	}
	return nil
}

type dotToken struct {
	Field            string
	Left             string
	IsArray          bool
	ArrayIndex       int
	IsAppendingArray bool
}

//go:linkname nextDot github.com/sunfmin/reflectutils.nextDot
func nextDot(name string) (t *dotToken, err error)

func GetFieldStruct(i interface{}, name string) (_ *reflect.StructField) {
	var err error

	t := reflect.TypeOf(i)

	if name == "" {
		return
	}

	var token *dotToken
	token, err = nextDot(name)
	if err != nil {
		return nil
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		sf, ok := t.FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(name, token.Field)
		})

		if !ok {
			return nil
		}

		if token.Left == "" {
			return &sf
		}

		return GetFieldStruct(reflect.Zero(sf.Type).Interface(), token.Left)
	}

	return
}
