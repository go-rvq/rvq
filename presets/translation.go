package presets

import (
	"context"
	"fmt"
	"strings"

	"github.com/qor5/x/v3/i18n"
)

var CommonMessagesTranslator i18n.TranslateFunc = func(ctx context.Context, key string, args ...string) (v string, ok bool) {
	if v = MustGetMessages(ctx).Common.Get(key); v != "" {
		v = strings.NewReplacer(args...).Replace(v)
		ok = true
	}
	return
}

func KeyFormatTranslatorD(prefix, keyFormat string, allowEmpty bool, module ...i18n.ModuleKey) i18n.Translator {
	if len(module) == 0 {
		module = append(module, ModelsI18nModuleKey)
	}

	t := make(i18n.Translators, len(module))

	for i, key := range module {
		t[i] = i18n.ModuleTranslator(key, allowEmpty, func(key string) string {
			return fmt.Sprintf(keyFormat, prefix, key)
		})
	}

	t[0] = i18n.WrapFallback(t[0], func(old i18n.TranslateFunc) i18n.TranslateFunc {
		return func(ctx context.Context, key string, args ...string) (v string, ok bool) {
			v, ok = CommonMessagesTranslator.Translate(ctx, key, args...)
			if !ok {
				return old(ctx, key)
			}
			return
		}
	})
	return t
}

func KeyFormatTranslator(prefix, keyFormat string, module ...i18n.ModuleKey) i18n.Translator {
	return KeyFormatTranslatorD(prefix, keyFormat, false, module...)
}

func (mb *ModelBuilder) KeyFormatTranslator(keyFormat string, module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslatorD(keyFormat, false, module...)
}

func (mb *ModelBuilder) KeyFormatTranslatorD(keyFormat string, allowEmpty bool, module ...i18n.ModuleKey) i18n.Translator {
	if len(module) == 0 {
		module = append(module, mb.I18nModuleKeyOrDefault())
	}
	return KeyFormatTranslatorD(mb.label, keyFormat, allowEmpty, module...)
}

func (mb *ModelBuilder) Translator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%[2]s", module...)
}

func (mb *ModelBuilder) FilterTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%s_Filter_%s", module...)
}

func (mb *ModelBuilder) FieldTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%s%s", module...)
}

func (mb *ModelBuilder) FieldHintTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslatorD("%s%s", true, module...)
}

func (mb *ModelBuilder) ActionTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%s_Action_%s", module...)
}

func (mb *ModelBuilder) BulkActionTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%s_BulkAction_%s", module...)
}

func (mb *ModelBuilder) T(ctx context.Context, key string, args ...string) string {
	return i18n.Translate(mb.Translator(), ctx, key, args...)
}

func (mb *ModelBuilder) TFormat(ctx context.Context, fmt, key string, args ...string) string {
	return i18n.Translate(mb.KeyFormatTranslator(fmt), ctx, key, args...)
}
