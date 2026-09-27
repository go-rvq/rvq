package examples_admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/theplant/gofixtures"
	"github.com/theplant/testingutils"
)

var l10nData = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.l10n_models (id, created_at, updated_at, deleted_at, title, locale_code) VALUES (1, 
'2024-06-04 23:27:40.442281 +00:00', '2024-06-04 23:27:40.442281 +00:00', null, 'My model title', 'International');


`, []string{"l10n_models"}))

var l10nDataWithChina = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.l10n_models (id, created_at, updated_at, deleted_at, title, locale_code) VALUES (1, 
'2024-06-04 23:27:40.442281 +00:00', '2024-06-04 23:27:40.442281 +00:00', null, 'My model title', 'International');
INSERT INTO public.l10n_models (id, created_at, updated_at, deleted_at, title, locale_code) VALUES (1, '2024-06-04 23:50:28.847833 +00:00', '2024-06-04 23:50:28.844069 +00:00', null, '中文标题', 'China');


`, []string{"l10n_models"}))

func TestLocalization(t *testing.T) {
	pb := presets.New(i18n.New()).DataOperator(gorm2op.DataOperator(TestDB))
	LocalizationExample(pb, TestDB)

	cases := []multipartestutils.TestCase{
		{
			Name:  "Index Page",
			Debug: true,
			ReqFunc: func() *http.Request {
				l10nData.TruncatePut(SqlDB)
				return httptest.NewRequest("GET", "/l10n-models", nil)
			},
			ExpectPageBodyContainsInOrder: []string{"My model title", "International"},
		},
		{
			Name:  "Index Page with locale code",
			Debug: true,
			ReqFunc: func() *http.Request {
				l10nDataWithChina.TruncatePut(SqlDB)
				return httptest.NewRequest("GET", "/l10n-models?locale=China", nil)
			},
			// a page's body travels as the escaped payload of the app portal, so the
			// markup is matched the way it is written there
			ExpectPageBodyContainsInOrder: []string{"中文标题", `v-chip color=&#39;success&#39; :variant=&#39;\"flat\"&#39; :label=&#39;true&#39; :size=&#39;\"small\"&#39;\u003eChina`},
		},
		{
			Name:  "Localize dialog",
			Debug: true,
			ReqFunc: func() *http.Request {
				l10nData.TruncatePut(SqlDB)
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/l10n-models?__execute_event__=l10n_LocalizeEvent&id=1_International").
					BuildEventFuncRequest()
				return req
			},
			ExpectPortalUpdate0ContainsInOrder: []string{"China", "Japan"},
		},
		{
			Name:  "Show detail",
			Debug: true,
			ReqFunc: func() *http.Request {
				l10nDataWithChina.TruncatePut(SqlDB)
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/l10n-models?__execute_event__=presets_Edit&overlay=RightDrawer&id=1_China").
					BuildEventFuncRequest()
				return req
			},
			ExpectPortalUpdate0ContainsInOrder: []string{"中文标题"},
		},
		{
			Name:  "Update detail",
			Debug: true,
			ReqFunc: func() *http.Request {
				l10nDataWithChina.TruncatePut(SqlDB)
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/l10n-models?__execute_event__=presets_Update&id=1_China").
					AddField("Title", "Updated Title").
					AddField("LocaleCode", "China").
					BuildEventFuncRequest()
				// L10nModel has an UpdatedAt, so the update requires the stamp
				// the rendered form would have carried (see SignForm).
				var stored L10nModel
				TestDB.First(&stored, "id = ? AND locale_code = ?", 1, "China")
				req, err := pb.SignForm(req, &stored)
				if err != nil {
					panic(err)
				}
				return req
			},
			EventResponseMatch: func(t *testing.T, er *multipartestutils.TestEventResponse) {
				var m L10nModel
				TestDB.Find(&m, "id = ? AND locale_code = ?", 1, "China")
				if m.Title != "Updated Title" {
					t.Errorf("title is wrong %#+v", m)
				}
			},
		},
		{
			Name:  "Delete China locale",
			Debug: true,
			ReqFunc: func() *http.Request {
				l10nDataWithChina.TruncatePut(SqlDB)
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/l10n-models?__execute_event__=presets_DoDelete&id=1_China").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *multipartestutils.TestEventResponse) {
				// deleted softly, in its own locale: the locale is not renamed
				// (a later translation to China frees the key, see l10n/events.go)
				var m L10nModel
				TestDB.Unscoped().Find(&m, "id = ? AND deleted_at IS NOT NULL", 1)
				if m.LocaleCode != "China" {
					t.Errorf("delete is wrong %#+v", m)
				}
			},
		},
		{
			Name:  "Localize to China and Japan",
			Debug: true,
			ReqFunc: func() *http.Request {
				l10nData.TruncatePut(SqlDB)
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/l10n-models?__execute_event__=l10n_DoLocalizeEvent&id=1_International&localize_from=International").
					AddField("LocalizeTo", "China").
					AddField("LocalizeTo", "Japan").
					BuildEventFuncRequest()
				return req
			},
			EventResponseMatch: func(t *testing.T, er *multipartestutils.TestEventResponse) {
				var localeCodes []string
				TestDB.Raw("SELECT locale_code FROM l10n_models ORDER BY locale_code").Scan(&localeCodes)
				if diff := testingutils.PrettyJsonDiff(
					[]string{"China", "International", "Japan"},
					localeCodes); diff != "" {
					t.Error(diff)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			multipartestutils.RunCase(t, c, pb)
		})
	}
}
