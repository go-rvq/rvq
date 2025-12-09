package vue

import (
	"encoding/json"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web/tag"
)

type FormFieldBuilder struct {
	*UserComponentBuilder
}

func FormField(fieldComponent ...h.HTMLComponent) (b *FormFieldBuilder) {
	comp := UserComponent(fieldComponent...)
	b = &FormFieldBuilder{comp}
	if len(fieldComponent) == 1 {
		if tg, _ := fieldComponent[0].(tag.TagGetter); tg != nil {
			t := tg.GetHTMLTagBuilder()
			assign := t.GetAttr("v-assign")
			if assign != nil {
				t.RemoveAttr("v-assign")

				if v, ok := assign.Value.(string); ok {
					if strings.HasPrefix(v, "[form,") {
						// remove prefix '[form, ' and ']' sufix
						v = v[7 : len(v)-1]
						var vm = map[string]any{}
						json.Unmarshal([]byte(v), &vm)
						assigner := b.Assigner("form")
						for k, v := range vm {
							assigner.Set(k, v)
						}
					}
				}
			}
		}
	}
	return
}

func (b *FormFieldBuilder) Scope(name string, v ...any) *FormFieldBuilder {
	b.UserComponentBuilder.Scope(name, v...)
	return b
}

func (b *FormFieldBuilder) Setup(s string) *FormFieldBuilder {
	b.UserComponentBuilder.Setup(s)
	return b
}

func (b *FormFieldBuilder) Field() h.HTMLComponent {
	return b.Template().Childs[0]
}

func (b *FormFieldBuilder) Value(fieldName string, value any) *FormFieldBuilder {
	q := strconv.Quote(fieldName)

	b.Assign("form", fieldName, value).
		Scope("fieldValue").
		Setup(`({scope, computed}) => {
	scope.fieldValue = {
		value: computed({
			get: () => form[` + q + `],		
			set: (newValue) => {
				console.log("set", newValue)
				form[` + q + `] = newValue
			}
		})
	}
}`)
	return b
}
func (b *FormFieldBuilder) BindTo(comp h.TagGetter) *FormFieldBuilder {
	comp.GetHTMLTagBuilder().Attr("v-model", `fieldValue.value`)
	return b
}
func (b *FormFieldBuilder) Bind() *FormFieldBuilder {
	return b.BindTo(b.Field().(h.TagGetter))
}
