package i18n

import (
	"context"
	"strings"

	"github.com/qor5/web/v3/str_utils"
	"github.com/sunfmin/reflectutils"
)

type (
	Translator interface {
		Translate(ctx context.Context, key string, args ...string) (v string, ok bool)
	}
	TranslateFunc func(ctx context.Context, key string, args ...string) (v string, ok bool)
	Translators   []Translator

	FallbackTranslator interface {
		Translator
		Fallback(ctx context.Context, key string, args ...string) (v string, ok bool)
	}

	Fallback func(ctx context.Context, key string, args ...string) (v string, ok bool)
)

type translatorWithFallback struct {
	t TranslateFunc
	f TranslateFunc
}

func TranslatorWithFallback(t TranslateFunc, f TranslateFunc) FallbackTranslator {
	return &translatorWithFallback{t: t, f: f}
}

func (t *translatorWithFallback) Translate(ctx context.Context, key string, args ...string) (string, bool) {
	return t.t(ctx, key, args...)
}

func (t *translatorWithFallback) Fallback(ctx context.Context, key string, args ...string) (string, bool) {
	return t.f(ctx, key, args...)
}

func (f TranslateFunc) Translate(ctx context.Context, key string, args ...string) (string, bool) {
	return f(ctx, key, args...)
}

func (ts Translators) Translate(ctx context.Context, key string, args ...string) (v string, ok bool) {
	for _, t := range ts {
		if v, ok = t.Translate(ctx, key, args...); ok {
			return
		}
	}
	return
}

func (ts Translators) Fallback(ctx context.Context, key string, args ...string) (s string, ok bool) {
	for _, t := range ts {
		if f, _ := t.(FallbackTranslator); f != nil {
			if s, ok = f.Fallback(ctx, key, args...); ok {
				return
			}
		}
	}
	return
}

func (ts Translators) Append(translators ...Translator) Translators {
	return append(ts, translators...)
}

func (f Fallback) Translate(ctx context.Context, key string, args ...string) (string, bool) {
	return "", false
}

func (f Fallback) Fallback(ctx context.Context, key string, args ...string) (string, bool) {
	return f(ctx, key, args...)
}

func ModuleTranslator(module ModuleKey, fFieldKey func(key string) string) Translator {
	return TranslatorWithFallback(
		func(ctx context.Context, key string, args ...string) (v string, ok bool) {
			msgr := MustGetModuleMessages(ctx, module, nil)

			if msgr == nil {
				return
			}

			fieldKey := fFieldKey(key)
			val, err := reflectutils.Get(msgr, fieldKey)
			if err == nil {
				vs, _ := val.(string)
				if vs != "" {
					v = strings.NewReplacer(args...).Replace(vs)
					ok = true
				}
			}
			return
		},
		func(ctx context.Context, key string, args ...string) (v string, ok bool) {
			b, _ := ctx.Value(dynaBuilderKey).(*DynaBuilder)
			if b != nil {
				msgr := MustGetModuleMessages(ctx, module, nil)
				fk := fFieldKey(key)
				val, err := reflectutils.Get(msgr, fk)
				if err != nil {
					b.putMissingKey(module, fk, key)
				} else if val.(string) == "" {
					b.putMissingVal(module, fk, key)
				}
			}

			key = str_utils.HumanizeString(key)
			return strings.NewReplacer(args...).Replace(key), true
		},
	)

}

func Translate(t Translator, ctx context.Context, key string, args ...string) (s string) {
	var ok bool
	if s, ok = t.Translate(ctx, key, args...); ok {
		return
	}

	if f, _ := t.(FallbackTranslator); f != nil {
		if s, ok = f.Fallback(ctx, key, args...); ok {
			return
		}
	}
	return str_utils.HumanizeString(strings.NewReplacer(args...).Replace(key))
}

func TranslateHandler(t Translator, ctx context.Context) func(key string, args ...string) string {
	return func(key string, args ...string) string {
		return Translate(t, ctx, key, args...)
	}
}

func WrapFallback(t Translator, f func(old TranslateFunc) TranslateFunc) Translator {
	tf, _ := t.(FallbackTranslator)
	if tf == nil {
		return TranslatorWithFallback(t.Translate, f(nil))
	}
	return TranslatorWithFallback(t.Translate, f(tf.Fallback))
}
