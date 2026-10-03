package presets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

type descPost struct {
	ID    uint
	Title string
	Body  string
}

type descMessages struct {
	DescPost                  string
	DescPost_Desc             string
	DescPostTitle             string
	DescPostTitle_Desc        string
	DescPostTitle_Hint        string
	DescPostBody_Hint         string
	DescPost_Action_Sync      string
	DescPost_Action_Sync_Desc string
}

// What a field is (_Desc) and how to fill it (_Hint): the text under it in a
// form is both, in this order; the description alone is what it is. A model
// and an action have theirs as _Desc too.
func TestFieldDescription(t *testing.T) {
	ib := i18n.New().SupportLanguages(language.English).
		RegisterForModule(language.English, ModelsI18nModuleKey, &descMessages{
			DescPost: "Post", DescPost_Desc: "A text of the site",
			DescPostTitle: "Title", DescPostTitle_Desc: "The name of the post.", DescPostTitle_Hint: "Up to 80 letters.",
			DescPostBody_Hint:    "Markdown.",
			DescPost_Action_Sync: "Sync", DescPost_Action_Sync_Desc: "Sends it again",
		})
	b := New(ib)
	mb := b.Model(&descPost{})
	mb.Editing("Title", "Body")
	sync := mb.Detailing().Action("Sync")

	var ctx context.Context
	ib.EnsureLanguage(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { ctx = r.Context() })).
		ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	info := mb.Info()
	title, body := mb.editing.GetField("Title"), mb.editing.GetField("Body")
	for _, c := range []struct{ what, got, want string }{
		{"the description", title.ContextDescription(info, ctx), "The name of the post."},
		{"the hint", title.ContextHint(info, ctx), "Up to 80 letters."},
		{"the form", title.ContextFormHint(info, ctx), "The name of the post. Up to 80 letters."},
		{"a hint alone", body.ContextFormHint(info, ctx), "Markdown."},
		{"no description", body.ContextDescription(info, ctx), ""},
		{"the model", mb.TDescription(ctx), "A text of the site"},
		{"an action", sync.RequestDescription(mb, ctx), "Sends it again"},
	} {
		if c.got != c.want {
			t.Errorf("%s: %q, want %q", c.what, c.got, c.want)
		}
	}
}
