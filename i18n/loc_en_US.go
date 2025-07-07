package i18n

import (
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

var Default_en = &DefaultMessages{
	True:  "True",
	False: "False",
	Yes:   "Yes",
	No:    "No",
}

func init() {
	e := en.New()
	uni := ut.New(e, e)
	pt, found := uni.GetTranslator("en")
	if !found {
		panic("translator not found")
	}

	for i := time.January; i <= time.December; i++ {
		Default_en.MonthNames[i] = pt.MonthWide(i)
	}

	for i := time.January; i <= time.December; i++ {
		Default_en.AbbrMonthNames[i] = pt.MonthAbbreviated(i)
	}
}
