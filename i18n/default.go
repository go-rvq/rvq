package i18n

import "context"

type DefaultMessages struct {
	MonthNames     [13]string
	AbbrMonthNames [13]string

	TimeLayout     TimeLayoutSpec
	DateLayout     TimeLayoutSpec
	DateTimeLayout TimeLayoutSpec

	True  string
	False string

	Yes string
	No  string
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
