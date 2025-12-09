package printer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	h "github.com/go-rvq/htmlgo"
	"gopkg.in/yaml.v3"
)

// TagAttibutes this is a map of attributes
type Tag struct {
	// Attributes is a list of attributes. The attribute is pairs of attribute name and attribute value.
	// For boolean attributes, must set key of then (`len(attribute) == 1`).
	Attributes [][]string
	Content    string
}

func (t *Tag) Build(b *h.HTMLTagBuilder) *h.HTMLTagBuilder {
	for _, attr := range t.Attributes {
		switch len(attr) {
		case 1:
			b.SetAttr(attr[0], true)
		case 2:
			b.SetAttr(attr[0], attr[1])
		}
	}

	if len(t.Content) > 0 {
		b.Append(h.RawHTML(t.Content))
	}

	return b
}

type NamedTag struct {
	Name string
	Tag
}

func (t *NamedTag) Build() *h.HTMLTagBuilder {
	b := h.Tag(t.Name)
	return t.Tag.Build(b)
}

type Printer struct {
	PrintButton string
	Title       string
	Lang        string
	Metas       []*Tag
	Styles      []*Tag
	Scripts     []*Tag
	PreBody     []*NamedTag
	Body        string
	PostBody    []*NamedTag
}

func (p *Printer) ParseValues(values url.Values) (err error) {
	if options := values.Get("options"); len(options) > 0 {
		if err = json.NewDecoder(bytes.NewBufferString(options)).Decode(p); err != nil {
			err = fmt.Errorf("Error parsing options: %v", err)
			return
		}
	}
	p.Body = values.Get("body")
	return
}

func (p *Printer) Component() h.HTMLComponent {
	var (
		head = h.Head()
		body = h.Body()
	)

	for _, t := range p.Metas {
		head.Append(t.Build(h.Meta()))
	}

	if p.Title != "" {
		head.Append(h.Title(p.Title))
	}

	for _, t := range p.Styles {
		if len(t.Content) > 0 {
			head.Append(t.Build(h.Style(t.Content)))
		} else {
			head.Append(t.Build(h.Tag("link").Attr("rel", "stylesheet")))
		}
	}

	if p.PrintButton != "" {
		body.Append(h.Div(h.RawHTML(p.PrintButton)).Class("noprint-media"))
	}

	for _, t := range p.PreBody {
		body.Append(t.Build())
	}

	if len(p.Body) > 0 {
		body.Append(h.RawHTML(p.Body))
	}

	for _, t := range p.PostBody {
		body.Append(t.Build())
	}

	body.Append(h.Style(`
    .noprint { display: none !important; }
    .print { display: block !important; }
	@media print {
 		.noprint-media { display: none !important; }
	}
`),
		h.Script(`
document.addEventListener("DOMContentLoaded", function(e) {
	document.querySelectorAll('[data-print-script]').forEach((e) => {
		const s = e.getAttribute("data-print-script")
		if (s) {
			const f = new Function(s)
			f.call(e)
		}
	})
})
`))

	r := h.HTML(head, body).(h.HTMLComponents)
	if p.Lang != "" {
		r[1].(*h.HTMLTagBuilder).SetAttr("lang", p.Lang)
	}

	return r
}

func ParseRequest(r *http.Request) (p *Printer, err error) {
	return ParseRequestType(r.Header.Get("Content-Type"), r.Form, r.Body)
}

func ParseRequestType(contentType string, values url.Values, r io.Reader) (p *Printer, err error) {
	switch contentType {
	case "application/json":
		return ParseJSON(r)
	case "application/x-yaml", "text/yaml":
		return ParseYAML(r)
	default:
		return ParseValues(values)
	}
}

func ParseJSON(r io.Reader) (p *Printer, err error) {
	p = new(Printer)
	err = json.NewDecoder(r).Decode(p)
	return
}

func ParseYAML(r io.Reader) (p *Printer, err error) {
	p = new(Printer)
	err = yaml.NewDecoder(r).Decode(p)
	return
}

func ParseValues(values url.Values) (p *Printer, err error) {
	p = new(Printer)
	err = p.ParseValues(values)
	return
}
