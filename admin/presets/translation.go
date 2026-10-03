package presets

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-rvq/rvq/x/i18n"
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

// KeyFormatTranslator builds a translator that composes its keys as
// fmt.Sprintf(keyFormat, prefix, key).
//
// Every label the admin shows for a model is a field on a messages struct,
// found by name. There is no table of translations: the key is built from the
// model and the thing being named, and read off the struct by reflection. A key
// with no field is reported as missing rather than guessed at.
//
// # The prefix
//
// Each key starts with the model's label, which NewModelBuilder derives from
// the struct name in CamelCase — models.ContactMessage becomes ContactMessage —
// unless the builder was given one. The plural, used where a set is named, is
// the inflected form: ContactMessages. Renaming the struct therefore renames
// every key of that model, and the labels fall back to their English field
// names until the messages follow.
//
// # The formats
//
// Each translator below composes prefix and key differently, and the format is
// the whole contract:
//
//	Translator            <Model>                      the title, keyed by the
//	                                                   label itself: TTitle
//	                                                   passes the singular,
//	                                                   TTitlePlural the plural
//	FieldTranslator       <Model><Field>               a field's label
//	HintTranslator        <Model><Field>_Desc          what a field is, and
//	                      <Model><Field>_Hint          how to fill it: the text
//	                                                   under it in a form is
//	                                                   both, _Desc then _Hint;
//	                                                   the suffix is added by
//	                                                   the caller, not by the
//	                                                   format
//	FilterTranslator      <Model>_Filter_<Key>         a filter's label, and the
//	                                                   text of each of its
//	                                                   options
//	ActionTranslator      <Model>_Action_<Name>        an action's label
//	BulkActionTranslator  <Model>_BulkAction_<Name>    a bulk action's label
//
// The descriptions — what a model, an action, a section is — are their keys
// with "_Desc", as a field's, none required: <Model>_Desc (TDescription),
// <Model>_Action_<Name>_Desc, <Model>_BulkAction_<Name>_Desc
// (RequestDescription), <Model><Section>_Desc; a job's,
// WorkerJob<Name>_Desc. A group, a page and a verifier of permissions take
// theirs as their title, by a function (DescriptionFunc, Description).
//
// So a Subject field on ContactService is ContactServiceSubject, what it is
// ContactServiceSubject_Desc, its hint ContactServiceSubject_Hint, and the
// "enabled" tab of its listing
// ContactService_Filter_Enabled.
//
// A format that is needed only once is passed directly to TFormat, as the row
// menu does with "%sRowMenuItem%s".
//
// # Which module, and what happens when the key is not there
//
// The messages struct comes from the module the model was given with
// SetModuleKey, defaulting to ModelsI18nModuleKey. An application registers its
// own — hermon-cms sets I18nHermonModelsKey on every model — so its keys live
// beside its models rather than in this package.
//
// The first module also falls back to the common messages, which is what makes
// a plain "Save" or "Delete" resolve without every model repeating it.
//
// Missing keys are recorded through the dyna builder when one is in the
// context, which is how they surface for translation instead of silently
// reading in English. Passing allowEmpty says the absence is expected — hints
// are the case: most fields have none, and demanding one for each would report
// every field as untranslated.
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

// Translator keys by the label alone: <Model>. TTitle passes the singular
// and TTitlePlural the plural.
func (mb *ModelBuilder) Translator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%[2]s", module...)
}

// FilterTranslator keys as <Model>_Filter_<Key>, for a filter's label and for
// the text of each of its options.
func (mb *ModelBuilder) FilterTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%s_Filter_%s", module...)
}

// FieldTranslator keys as <Model><Field>.
func (mb *ModelBuilder) FieldTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%s%s", module...)
}

// HintTranslator keys as <Model><Field>, and the caller appends the _Hint —
// the format does not. It allows an empty result: most fields have no hint, and
// requiring one would report every field as untranslated.
func (mb *ModelBuilder) HintTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslatorD("%s%s", true, module...)
}

// ActionTranslator keys as <Model>_Action_<Name>.
func (mb *ModelBuilder) ActionTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%s_Action_%s", module...)
}

// BulkActionTranslator keys as <Model>_BulkAction_<Name>.
func (mb *ModelBuilder) BulkActionTranslator(module ...i18n.ModuleKey) i18n.Translator {
	return mb.KeyFormatTranslator("%s_BulkAction_%s", module...)
}

func (mb *ModelBuilder) T(ctx context.Context, key string, args ...string) string {
	return i18n.Translate(mb.Translator(), ctx, key, args...)
}

func (mb *ModelBuilder) TFormat(ctx context.Context, fmt, key string, args ...string) string {
	return i18n.Translate(mb.KeyFormatTranslator(fmt), ctx, key, args...)
}
