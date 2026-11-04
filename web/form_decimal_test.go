package web_test

import (
	"net/url"
	"testing"

	"github.com/go-rvq/rvq/web"
	"github.com/shopspring/decimal"
)

// A decimal of the form is read from its text, "" zero (NULL when nullable);
// without its decoder it came out zero, with no error.
func TestUnmarshalFormDecimal(t *testing.T) {
	type forma struct{ Valor decimal.Decimal }
	var v struct {
		Valor   decimal.Decimal
		Vazio   decimal.Decimal
		Nulo    decimal.NullDecimal
		Fixado  decimal.NullDecimal
		Formas  []*forma
		Ponteir *decimal.Decimal
	}
	err := (&web.EventContext{}).UnmarshalFormValues(url.Values{
		"Valor":           {"-250.50"},
		"Vazio":           {""},
		"Nulo":            {""},
		"Fixado":          {"1234567890123.45678901"},
		"Formas[0].Valor": {"10"},
		"Ponteir":         {"3.3"},
	}, &v)
	if err != nil {
		t.Fatal(err)
	}
	if !v.Valor.Equal(decimal.RequireFromString("-250.5")) || !v.Vazio.IsZero() || v.Nulo.Valid ||
		!v.Fixado.Valid || v.Fixado.Decimal.String() != "1234567890123.45678901" ||
		len(v.Formas) != 1 || !v.Formas[0].Valor.Equal(decimal.NewFromInt(10)) ||
		v.Ponteir == nil || v.Ponteir.String() != "3.3" {
		t.Errorf("decoded: %+v", v)
	}
	if err = (&web.EventContext{}).UnmarshalFormValues(url.Values{"Valor": {"abc"}}, &v); err == nil {
		t.Error("a bad decimal: no error")
	}
}
