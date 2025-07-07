package integration_test

import (
	"net/http"
	"testing"

	"github.com/go-rvq/rvq/admin/example/admin"
	"github.com/go-rvq/rvq/web/multipartestutils"
)

func TestOrders(t *testing.T) {
	h := admin.TestHandler(TestDB)
	dbr, _ := TestDB.DB()

	cases := []multipartestutils.TestCase{
		{
			Name:  "Show order detail",
			Debug: true,
			ReqFunc: func() *http.Request {
				admin.OrdersExampleData.TruncatePut(dbr)
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/orders?__execute_event__=presets_DetailingDrawer&id=11").
					BuildEventFuncRequest()
				return req
			},
			ExpectPortalUpdate0ContainsInOrder: []string{`Basic Information`},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			multipartestutils.RunCase(t, c, h)
		})
	}
}
