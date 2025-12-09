package i18n

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/sunfmin/reflectutils"
)

type kv struct {
	key        string
	val        string
	allowEmpty bool
}

type moduleMissing struct {
	missingKeys []kv
	missingVals []kv
}
type DynaBuilder struct {
	lang    string
	missing map[ModuleKey]*moduleMissing
}

func DynaNew() (r *DynaBuilder) {
	return &DynaBuilder{
		missing: make(map[ModuleKey]*moduleMissing),
	}
}

func (d *DynaBuilder) Language(lang string) (r *DynaBuilder) {
	d.lang = lang
	return d
}

func (d *DynaBuilder) GetLanguage() string {
	return d.lang
}

func T(ctx context.Context, module ModuleKey, key string, args ...string) (r string) {
	return PT(ctx, module, "", key, args...)
}

func PT(ctx context.Context, module ModuleKey, prefix string, key string, args ...string) (r string) {
	return PTFk(ctx, module, func() string { return strcase.ToCamel(prefix + " " + key) }, key, args...)
}

func DynaFromContext(ctx context.Context) (r *DynaBuilder) {
	r, _ = ctx.Value(dynaBuilderKey).(*DynaBuilder)
	return
}

func PTFk(ctx context.Context, module ModuleKey, fFieldKey func() string, key string, args ...string) (r string) {
	defaultVal := strings.NewReplacer(args...).Replace(key)
	msgr := MustGetModuleMessages(ctx, module, nil)
	if msgr == nil {
		return defaultVal
	}

	builder := DynaFromContext(ctx)

	fieldKey := fFieldKey()
	val, err := reflectutils.Get(msgr, fieldKey)
	if err != nil {
		if builder != nil {
			builder.putMissingKey(module, fieldKey, key)
		}
		return defaultVal
	}

	if val.(string) == "" {
		if builder != nil {
			builder.putMissingVal(module, fieldKey, key)
		}
		val = defaultVal
	}

	return strings.NewReplacer(args...).Replace(val.(string))
}

func (d *DynaBuilder) putMissingKey(module ModuleKey, key, val string, allowEmpty ...bool) {
	if d.missing[module] == nil {
		d.missing[module] = &moduleMissing{}
	}
	mm := d.missing[module]

	for _, ck := range mm.missingKeys {
		if ck.key == key {
			return
		}
	}
	var ae bool
	for _, ae = range allowEmpty {
	}

	mm.missingKeys = append(mm.missingKeys, kv{key, val, ae})
}

func (d *DynaBuilder) putMissingVal(module ModuleKey, key, val string) {
	if d.missing[module] == nil {
		d.missing[module] = &moduleMissing{}
	}
	mm := d.missing[module]

	for _, ck := range mm.missingVals {
		if ck.key == key {
			return
		}
	}
	mm.missingVals = append(mm.missingVals, kv{key, val, true})
}

func (d *DynaBuilder) HaveMissingKeys() bool {
	return len(d.missing) > 0
}

func (d *DynaBuilder) PrettyMissingKeys() string {
	buf := new(bytes.Buffer)
	for module, missing := range d.missing {

		buf.WriteString(fmt.Sprintf("\nFor module %s, ", module))
		buf.WriteString("Missing the following translations\nCopy these to your Messages struct definition\n============================\n\n")

		for _, kv := range missing.missingKeys {
			_, _ = fmt.Fprintf(buf, "%s string\n", kv.key)
		}

		buf.WriteString("\n")
		buf.WriteString(fmt.Sprintf("\nCopy these to your Messages struct values for language: `%s`\n\n", d.lang))
		for _, kv := range missing.missingKeys {
			if !kv.allowEmpty {
				_, _ = fmt.Fprintf(buf, "%s: %#+v,\n", kv.key, kv.val)
			}
		}

		for _, kv := range missing.missingVals {
			_, _ = fmt.Fprintf(buf, "%s: %#+v,\n", kv.key, kv.val)
		}
	}

	return buf.String()
}
