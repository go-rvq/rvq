package l10n

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

func localizeSetModel(t *testing.T) *presets.ModelBuilder {
	t.Helper()
	pb := presets.New(i18n.New())
	m := pb.Model(&localeTestRecord{})
	b := New(nil).DefaultLocaleCode("en-US")
	b.RegisterLocale("en-US", "en-us", "English")
	b.RegisterLocale("pt-BR", "pt-br", "Português do Brasil")
	if err := b.ModelInstall(pb, m); err != nil {
		t.Fatal(err)
	}
	return m
}

func menuNames(rmb *presets.RowMenuBuilder) string {
	var ns []string
	for _, it := range rmb.Items() {
		ns = append(ns, it.Name())
	}
	return strings.Join(ns, ",")
}

// The action is in the menus of the listing and of the detail; named again
// (SetAction), it is the same item, under its name and title, in both.
func TestSetActionRenamesTheMenus(t *testing.T) {
	m := localizeSetModel(t)
	for _, rmb := range []*presets.RowMenuBuilder{m.Listing().RowMenu(), m.Detailing().RowMenu()} {
		if !strings.Contains(menuNames(rmb), "Localize") {
			t.Fatalf("no Localize in %s", menuNames(rmb))
		}
	}
	SetAction(m, "Translate", func(context.Context) string { return "Traduzir" })
	for _, rmb := range []*presets.RowMenuBuilder{m.Listing().RowMenu(), m.Detailing().RowMenu()} {
		names := menuNames(rmb)
		if strings.Contains(names, "Localize") || !strings.Contains(names, "Translate") {
			t.Errorf("the menu: %s", names)
		}
	}
	if got := ActionTitle(context.Background(), m); got != "Traduzir" {
		t.Errorf("the title: %q", got)
	}
}

// A localization takes what Expand says, with the switches of the request
// (their defaults when not sent), and is done by Do, in that order.
func TestLocalizeRecordsExpands(t *testing.T) {
	m := localizeSetModel(t)
	var (
		gotOn map[string]bool
		did   []string
	)
	ModelWithOptions(m, &ModelLocalizeOptions{
		Switches: []*LocalizeSwitch{{Name: "ancestors", Default: true}, {Name: "descendants"}},
		Expand: func(_ *web.EventContext, mids model.IDSlice, _ []string, on map[string]bool) (*Expansion, error) {
			gotOn = on
			return &Expansion{IDs: mids, Added: map[string]int{}}, nil
		},
		Do: func(_ *web.EventContext, mids model.IDSlice, to []string) error {
			for _, mid := range mids {
				did = append(did, mid.String()+">"+strings.Join(to, ","))
			}
			return nil
		},
	})
	r := httptest.NewRequest("POST", "/", strings.NewReader(url.Values{switchesField + "descendants": {"true"}}.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_ = r.ParseForm()
	ctx := &web.EventContext{R: r}
	n, err := LocalizeRecords(nil, m, nil, ctx, model.IDSlice{testMID()}, []string{"pt-BR"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 || len(did) != 1 || !strings.HasSuffix(did[0], ">pt-BR") {
		t.Errorf("done: %d %v", n, did)
	}
	if !gotOn["ancestors"] || !gotOn["descendants"] {
		t.Errorf("the switches: %v", gotOn)
	}
	// none chosen: nothing
	if n, _ := LocalizeRecords(nil, m, nil, ctx, model.IDSlice{testMID()}, nil); n != 0 {
		t.Errorf("localized to no locale: %d", n)
	}
}

func testMID() model.ID {
	return model.ID{Fields: model.Fields{model.SingleField("ID"), model.SingleField("LocaleCode")}, Values: []any{1, "en-US"}}
}

// What the switches take along is asked @edit once by record: its versions
// in each language one question.
func TestLocalizeRecordsAsksOnce(t *testing.T) {
	var asked []string
	pb := presets.New(i18n.New())
	pb.Permission(perm.New().AllowAll().OnAsk(func(a perm.Asked) { asked = append(asked, a.Resource) }).
		SubjectsFunc(func(*http.Request) []string { return []string{"editor"} }))
	m := pb.Model(&localeTestRecord{})
	b := New(nil).DefaultLocaleCode("en-US")
	b.RegisterLocale("en-US", "en-us", "English")
	if err := b.ModelInstall(pb, m); err != nil {
		t.Fatal(err)
	}
	id := func(n int, loc string) model.ID {
		return model.ID{Fields: model.Fields{model.SingleField("ID"), model.SingleField("LocaleCode")}, Values: []any{n, loc}}
	}
	ModelWithOptions(m, &ModelLocalizeOptions{
		Expand: func(_ *web.EventContext, mids model.IDSlice, _ []string, _ map[string]bool) (*Expansion, error) {
			// the parent 2, in two languages, and 3
			return &Expansion{IDs: append(model.IDSlice{id(2, "en-US"), id(2, "pt-BR"), id(3, "en-US")}, mids...)}, nil
		},
		Do: func(*web.EventContext, model.IDSlice, []string) error { return nil },
	})
	ctx := &web.EventContext{R: httptest.NewRequest("POST", "/", nil)}
	if _, err := LocalizeRecords(nil, m, nil, ctx, model.IDSlice{id(1, "en-US")}, []string{"zz"}); err != nil {
		t.Fatal(err)
	}
	count := func(suffix string) (n int) {
		for _, a := range asked {
			if strings.HasSuffix(a, suffix) {
				n++
			}
		}
		return
	}
	if n2, n3, n1 := count(":<2>:@edit"), count(":<3>:@edit"), count(":<1>:@edit"); n2 != 1 || n3 != 1 || n1 != 0 {
		t.Errorf("@edit asked: <2> %d, <3> %d, <1> %d — %v", n2, n3, n1, asked)
	}
}
