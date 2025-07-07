package presets

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
	_ "unsafe"

	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/str_utils"
	"github.com/sunfmin/reflectutils"
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

func ShowMessage(r *web.EventResponse, msg any, color ...string) {
	var (
		m       = NewFlashMessage(msg, color...)
		text    = m.Text
		textKey = "message"
	)
	if len(text) == 0 {
		if m.HtmlText != "" {
			textKey = "htmlText"
			text = m.HtmlText
		} else {
			return
		}
	}

	if m.Color == "" {
		m.Color = "success"
	}

	web.AppendRunScripts(r, fmt.Sprintf(
		`vars.presetsMessage = { show: true, %s: %s, color: %s}`,
		textKey, h.JSONString(text), h.JSONString(m.Color)))
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

func ToStringContext(ctx *web.EventContext, v any) string {
	if v == nil || reflectutils.IsNil(v) {
		return ""
	}

	switch t := v.(type) {
	case []rune:
		return string(t)
	case []byte:
		return string(t)
	case time.Time:
		return t.Format("2006-01-02 15:04:05")
	case *time.Time:
		if t == nil {
			return ""
		}
		return t.Format("2006-01-02 15:04:05")
	case sql.NullInt32:
		return fmt.Sprint(t.Int32)
	case sql.NullInt64:
		return fmt.Sprint(t.Int64)
	case sql.NullFloat64:
		return fmt.Sprint(t.Float64)
	case sql.Null[uint]:
		return fmt.Sprint(t.V)
	case sql.Null[uint8]:
		return fmt.Sprint(t.V)
	case sql.Null[uint16]:
		return fmt.Sprint(t.V)
	case sql.Null[uint32]:
		return fmt.Sprint(t.V)
	case sql.Null[uint64]:
		return fmt.Sprint(t.V)

	case ContextStringer:
		return t.ContextString(ctx)
	}
	return fmt.Sprint(v)
}

var (
	HumanizeString = str_utils.HumanizeString
	NamifyString   = str_utils.NamifyString
	SplitString    = str_utils.SplitString
)
