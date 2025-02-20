package vue

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/qor5/web/v3/tag"
	h "github.com/theplant/htmlgo"
)

type FormFieldBuilder struct {
	*UserComponentBuilder
}

func FormField(fieldComponent h.HTMLComponent) (b *FormFieldBuilder) {
	t := fieldComponent.(tag.TagGetter).GetHTMLTagBuilder()
	comp := UserComponent(fieldComponent)
	b = &FormFieldBuilder{comp}

	assign := tag.GetAttr(t, "v-assign")
	if assign != nil {
		tag.RemoveAttr(t, "v-assign")

		if v, ok := assign.Get().(string); ok {
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
	return
}

func (b *FormFieldBuilder) Field() h.HTMLComponent {
	return (*tag.Children(b.Template()))[0]
}

func (b *FormFieldBuilder) Value(fieldName string, value any) *FormFieldBuilder {
	q := strconv.Quote(fieldName)

	b.Assign("form", fieldName, value).
		Scope("fieldValue").
		Setup(`({scope, computed}) => {
	scope.fieldValue = {
		value: computed({
			get: () => form[` + q + `],		
			set: (newValue) => form[` + q + `] = newValue
		})
	}
}`)
	return b
}
