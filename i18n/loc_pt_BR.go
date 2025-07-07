package i18n

import (
	"context"
	"time"

	"github.com/go-playground/locales/pt_BR"
	"github.com/go-playground/universal-translator"
)

const (
	PtDateDefaultFormat = "pt:date:default"
	PtDateShortFormat   = "pt:date:short"
	PtDateFullFormat    = "pt:date:full"

	PtTimeDefaultFormat = "pt:time:default"
	PtTimeShortFormat   = "pt:time:short"
	PtTimeFullFormat    = "pt:time:full"

	PtDateTimeDefaultFormat = "pt:datetime:default"
	PtDateTimeShortFormat   = "pt:datetime:short"
	PtDateTimeFullFormat    = "pt:datetime:full"
)

var Default_pt_BR = &DefaultMessages{
	True:  "Verdadeiro",
	False: "Falso",
	Yes:   "Sim",
	No:    "Não",
	DateLayout: TimeLayoutSpec{
		Default: PtDateDefaultFormat,
		Short:   PtDateShortFormat,
		Full:    PtDateFullFormat,
	},
	TimeLayout: TimeLayoutSpec{
		Default: PtTimeDefaultFormat,
		Short:   PtTimeShortFormat,
		Full:    PtTimeFullFormat,
	},
	DateTimeLayout: TimeLayoutSpec{
		Default: PtDateTimeDefaultFormat,
		Short:   PtDateTimeShortFormat,
		Full:    PtDateTimeFullFormat,
	},
}

func init() {
	p := pt_BR.New()
	uni := ut.New(p, p)
	pt, found := uni.GetTranslator("pt_BR")
	if !found {
		panic("translator not found")
	}

	TimeFormats.AddFuncs(func(add func(name TimeLayout, formatter TimeLayoutFunc)) {
		add(PtDateDefaultFormat, func(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
			return pt.FmtDateMedium
		})
		add(PtDateShortFormat, func(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
			return pt.FmtDateShort
		})
		add(PtDateFullFormat, func(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
			return pt.FmtDateFull
		})

		add(PtTimeDefaultFormat, func(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
			return pt.FmtTimeMedium
		})
		add(PtTimeShortFormat, func(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
			return pt.FmtTimeShort
		})
		add(PtTimeFullFormat, func(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
			return pt.FmtTimeFull
		})

		merge := func(t1, t2 func(time.Time) string, sep string) func(t time.Time) string {
			return func(t time.Time) string {
				a, b := t1(t), t2(t)
				return a + sep + b
			}
		}
		add(PtDateTimeDefaultFormat, func(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
			return merge(pt.FmtDateMedium, pt.FmtTimeMedium, " às ")
		})
		add(PtDateTimeShortFormat, func(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
			return merge(pt.FmtDateShort, pt.FmtTimeShort, " ")
		})
		add(PtDateTimeFullFormat, func(m *DefaultMessages, ctx context.Context) func(t time.Time) string {
			return merge(pt.FmtDateFull, pt.FmtTimeFull, ", às ")
		})
	})

	for i := time.January; i <= time.December; i++ {
		Default_pt_BR.MonthNames[i] = pt.MonthWide(i)
	}

	for i := time.January; i <= time.December; i++ {
		Default_pt_BR.AbbrMonthNames[i] = pt.MonthAbbreviated(i)
	}
}
