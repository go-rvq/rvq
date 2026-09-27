package docgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	. "github.com/go-rvq/htmlgo"
	"github.com/iancoleman/strcase"
	"golang.org/x/net/html"
)

type DocBuilder struct {
	title        string
	slug         string
	abstractText string
	children     []HTMLComponent
}

func Doc(vs ...HTMLComponent) (r *DocBuilder) {
	r = &DocBuilder{
		children: vs,
	}
	return
}

func (b *DocBuilder) Title(v string) (r *DocBuilder) {
	b.title = v
	return b
}

func (b *DocBuilder) Slug(v string) (r *DocBuilder) {
	b.slug = v
	return b
}

func (b *DocBuilder) AbstractText(v string) (r *DocBuilder) {
	b.abstractText = v
	return b
}

func (b *DocBuilder) GetPageURL() (r string) {
	slug := b.slug
	if slug == "" {
		slug = strcase.ToKebab(b.title)
	}
	u := path.Join("/", slug)
	if u == "/" {
		return "index.html"
	}
	return strings.TrimLeft(fmt.Sprintf("%s.html", u), "/")
}

func (b *DocBuilder) Write(ctx *Context) error {
	return Div(
		Div(
			Div(
				Button("").Children(
					Div(
						menuIcon,
					).Class("w-4 h-4 fill-current text-gray-300"),
				).Class("w-12 h-12 p-4").
					Attr("@click", "vars.hideAside = !vars.hideAside"),
			).Class("flex flex-row"),
			Div(
				H1(b.title).Class("mb-8"),
				If(len(b.abstractText) > 0,
					Div(
						Text(b.abstractText),
					).Class("mb-8 text-xl font-normal"),
				),
				Div(
					b.children...,
				).Class("border-t"),
			).Class("px-16 pb-12 pt-4 overflow-auto").
				ID("docMainBox"),
		).Class("flex flex-grow flex-col w-2/3"),
		Div(
			Div(
				Text("On This Page"),
				RawHTML("<toc></toc>"),
			).Class("sticky top-4 w-52"),
		).Class("font-medium text-base hidden xl:block text-gray-600 pt-4"),
	).
		Class("flex flex-row w-full").
		ID("docContentBox").
		Write(ctx)
}

func (b *DocBuilder) plainText() string {
	hb, err := Marshal(Div(
		Div(Text(b.abstractText)),
		Div(b.children...),
	), context.Background())
	if err != nil {
		panic(err)
	}

	return html2text(hb)
}

func html2text(in []byte) string {
	r := &bytes.Buffer{}

	tokenizer := html.NewTokenizer(bytes.NewReader(in))
	startToken := tokenizer.Token()

	for tt := tokenizer.Next(); tt != html.ErrorToken; tt = tokenizer.Next() {
		switch tt {
		case html.StartTagToken:
			startToken = tokenizer.Token()
			if startToken.Data == "highlightjs" {
				for _, v := range startToken.Attr {
					if v.Key == ":code" {
						code := v.Val
						if code != "" {
							if err := json.Unmarshal([]byte(code), &code); err != nil {
								panic(err)
							}
							code = strings.NewReplacer("\n", " ", "\t", " ").Replace(code)
							r.WriteString(code)
							r.WriteString(" ")
						}
					}
				}
			}
		case html.TextToken:
			if startToken.Data == "script" || startToken.Data == "style" {
				continue
			}
			txt := strings.TrimSpace(string(tokenizer.Text()))
			if txt != "" {
				r.WriteString(txt)
				r.WriteString(" ")
			}
		}
	}

	return r.String()
}

var menuIcon = RawHTML(`
<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="16px" height="16px" viewBox="0 0 16 16" version="1.1">
<g id="surface1">
<path style=" stroke:none;fill-rule:nonzero;fill:rgb(0%,0%,0%);fill-opacity:1;" d="M 2 12 L 2 11 L 14 11 L 14 12 Z M 2 8.5 L 2 7.5 L 14 7.5 L 14 8.5 Z M 2 5 L 2 4 L 14 4 L 14 5 Z M 2 5 "/>
</g>
</svg>
`)
