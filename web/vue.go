package web

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-rvq/rvq/web/vue"
	h "github.com/theplant/htmlgo"
)

type Var = vue.Var

type RawVar = vue.RawVar

func VAssign(varName string, v interface{}) []interface{} {
	var varVal string
	switch t := v.(type) {
	case string:
		varVal = t
	case []byte:
		varVal = string(t)
	case map[string]interface{}:
		l := len(t)
		if l == 0 {
			varVal = "{}"
		} else {
			var b strings.Builder
			b.WriteString("{")
			for k, v := range t {
				b.WriteString(strconv.Quote(k))
				b.WriteString(": ")

				switch t := v.(type) {
				case Var:
					b.WriteString(string(t))
				case []byte:
					b.WriteString(string(t))
				default:
					b.WriteString(h.JSONString(t))
				}
				b.WriteString(",")
			}
			varVal = b.String()
			varVal = varVal[:len(varVal)-1] + "}"
		}
	default:
		varVal = h.JSONString(t)
	}
	return []interface{}{
		"v-assign",
		fmt.Sprintf("[%s, %s]", varName, varVal),
	}
}

func VField(name string, value interface{}) []interface{} {
	objValue := map[string]interface{}{name: value}
	return append([]interface{}{
		"v-model",
		fmt.Sprintf("form[%s]", h.JSONString(name)),
	}, VAssign("form", objValue)...)
}

func VModel(name string) []interface{} {
	return append([]interface{}{
		"v-model",
		name,
	})
}

func GlobalEvents() *h.HTMLTagBuilder {
	return h.Tag("global-events")
}

var objectScriptRe = regexp.MustCompile(`":(\w+)"(\s*):(\s*)"(\w+)"(?ms)`)

func EvaluatedJSONObject(s string) string {
	s = objectScriptRe.ReplaceAllString(s, `"$1"$2:$3$4`)
	return s
}
