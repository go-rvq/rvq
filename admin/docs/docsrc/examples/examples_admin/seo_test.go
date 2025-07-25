package examples_admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/theplant/gofixtures"
)

var seoData = gofixtures.Data(gofixtures.Sql(`
INSERT INTO public.seo_posts (title, seo, id, created_at, updated_at, deleted_at) VALUES ('The seo post 1', 
'{"OpenGraphImageFromMediaLibrary":{"ID":0,"Url":"","VideoLink":"","FileName":"","Description":""}}', 1, '2024-05-31 10:02:13.114089 +00:00', '2024-05-31 10:02:13.114089 +00:00', null);

`, []string{"seo_posts"}))

func TestSEOExampleBasic(t *testing.T) {
	pb := presets.New()
	SEOExampleBasic(pb, TestDB)

	cases := []multipartestutils.TestCase{
		{
			Name:  "Index Page",
			Debug: true,
			ReqFunc: func() *http.Request {
				seoData.TruncatePut(SqlDB)
				return httptest.NewRequest("GET", "/seo-posts", nil)
			},
			ExpectPageBodyContainsInOrder: []string{"The seo post 1"},
		},
		{
			Name:  "Edit SEO Title",
			Debug: true,
			ReqFunc: func() *http.Request {
				seoData.TruncatePut(SqlDB)
				req := multipartestutils.NewMultipartBuilder().
					PageURL("/seo-posts?__execute_event__=presets_Detailing_Field_Save&detailField=SEO&id=1").
					AddField("Seo.EnabledCustomize", "true").
					AddField("Seo.Title", "My seo title").
					BuildEventFuncRequest()
				return req
			},
			// TODO: Not assert correct, should be "My seo title"
			ExpectPortalUpdate0ContainsInOrder: []string{`Open Graph Preview`},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			multipartestutils.RunCase(t, c, pb)
		})
	}
}
