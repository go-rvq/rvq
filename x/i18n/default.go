package i18n

import "context"

type DefaultMessages struct {
	MonthNames     [13]string `i18n:"hint='The names of the months, January first (the item 0 is unused).'"`
	AbbrMonthNames [13]string `i18n:"label='Abbreviated month names', hint='The short names of the months, January first (the item 0 is unused).'"`

	TimeLayout     TimeLayoutSpec `i18n:"hint='How a time is written: Go layouts (the reference time 15:04:05) for the default, short and full forms.'"`
	DateLayout     TimeLayoutSpec `i18n:"hint='How a date is written: Go layouts (the reference date 2006-01-02) for the default, short and full forms.'"`
	DateTimeLayout TimeLayoutSpec `i18n:"label='Date and time layout', hint='How a date with its time is written: Go layouts (reference 2006-01-02 15:04:05) for the default, short and full forms.'"`

	True  string `i18n:"hint='How a true value is shown.'"`
	False string `i18n:"hint='How a false value is shown.'"`

	Yes string `i18n:"hint='The answer yes.'"`
	No  string `i18n:"hint='The answer no.'"`
}

func (m *DefaultMessages) TrueOrFalse(v bool) string {
	if v {
		return m.True
	}
	return m.False
}

func (m *DefaultMessages) YesOrNo(v bool) string {
	if v {
		return m.Yes
	}
	return m.No
}

const DefaultKey ModuleKey = "i18n:default"

func GetMessages(ctx context.Context) *DefaultMessages {
	return MustGetModuleMessages(ctx, DefaultKey, Default_en).(*DefaultMessages)
}
