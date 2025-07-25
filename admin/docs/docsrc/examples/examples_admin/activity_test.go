package examples_admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/theplant/gofixtures"
)

var activityData = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.activity_logs (id, user_id, created_at, creator, action, model_keys, model_name, model_label, model_link, model_diffs) VALUES (1, 0, '2024-05-30 07:02:53.393836 +00:00', 'smile', 'Create', '1:xxx', 'WithActivityProduct', 'with-activity-products', '', '');

INSERT INTO public.with_activity_products (title, code, price, id, created_at, updated_at, deleted_at) VALUES ('P11111111111', 'code11111111', 0, 1, '2024-05-30 07:02:53.389781 +00:00', '2024-05-30 07:15:38.585837 +00:00', null);

`, []string{"with_activity_products", "activity_logs"}))

func TestActivity(t *testing.T) {
	pb := presets.New().DataOperator(gorm2op.DataOperator(TestDB))
	ActivityExample(pb, TestDB)

	dbr, _ := TestDB.DB()

	cases := []multipartestutils.TestCase{
		{
			Name:  "Index Page",
			Debug: true,
			ReqFunc: func() *http.Request {
				activityData.TruncatePut(dbr)
				return httptest.NewRequest("GET", "/with-activity-products", nil)
			},
			ExpectPageBodyContainsInOrder: []string{"P11111111111"},
		},
		{
			Name:  "Activity Model details should have timeline",
			Debug: true,
			ReqFunc: func() *http.Request {
				activityData.TruncatePut(dbr)
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/with-activity-products?__execute_event__=presets_DetailingDrawer&id=1").
					BuildEventFuncRequest()
				return req
			},
			ExpectPortalUpdate0ContainsInOrder: []string{"WithActivityProduct 1"},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			multipartestutils.RunCase(t, c, pb)
		})
	}
}
