package docgo

import (
	"strings"

	"github.com/go-rvq/htmlgo"
	"github.com/shurcooL/github_flavored_markdown"
)

func Markdown(body string) htmlgo.HTMLComponent {
	return htmlgo.ComponentFunc(func(ctx *htmlgo.Context) (err error) {
		body = strings.Replace(body, "~", "`", -1)
		return htmlgo.RawHTML(github_flavored_markdown.Markdown([]byte(body))).Write(ctx)
	})
}
