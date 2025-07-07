package editorjs

import (
	"fmt"
	"io"
	"strings"
)

type Transformer func(data *BlockData, out io.Writer) error

var HTMLTransformers = map[string]Transformer{
	"delimiter": func(data *BlockData, out io.Writer) (err error) {
		out.Write([]byte(`<br/>`))
		return
	},

	"header": func(data *BlockData, out io.Writer) (err error) {
		fmt.Fprintf(out, `<h%d>%s</h%[1]d>`, data.Level, data.Text)
		return
	},

	"paragraph": func(data *BlockData, out io.Writer) (err error) {
		fmt.Fprintf(out, `<p>%s</p>`, data.Text)
		return
	},

	"table": func(data *BlockData, out io.Writer) (err error) {
		out.Write([]byte(`<div class="mdl-grid"><table class="mdl-data-table mdl-js-data-table mdl-shadow--2dp">`))
		rows := data.Content
		var numericColumns []bool
		if len(rows) > 0 {
			numericColumns = make([]bool, len(rows[0]))

			for i, h := range rows[0] {
				if strings.ToLower(h[0:5]) == "[[n]]" {
					rows[0][i] = rows[0][i][5:]
					numericColumns[i] = true
				}
			}

			if data.WithHeadings {
				out.Write([]byte("<theader><tr>"))
				for i, h := range rows[0] {
					out.Write([]byte("<th"))
					if !numericColumns[i] {
						out.Write([]byte(` class="mdl-data-table__cell--non-numeric"`))
					}
					out.Write([]byte(">"))
					out.Write([]byte(h))
					out.Write([]byte("</th>"))
				}
				out.Write([]byte("</tr></theader>"))
				rows = rows[1:]
			}
		}

		if len(rows) > 0 {
			out.Write([]byte("<tbody>"))
			for _, row := range rows {
				out.Write([]byte("<tr>"))
				for i, v := range row {
					out.Write([]byte("<td"))
					if !numericColumns[i] {
						out.Write([]byte(` class="mdl-data-table__cell--non-numeric"`))
					}
					out.Write([]byte(">"))
					out.Write([]byte(v))
					out.Write([]byte("</td>"))
				}
				out.Write([]byte("</tr>"))
			}
			out.Write([]byte("</tbody>"))
		}

		out.Write([]byte("</table></div>"))
		return
	},

	"list": func(data *BlockData, out io.Writer) (err error) {
		var (
			listStyle = "ol"
			recursor  func(items []*BlockItem, listStyle string)
		)

		out.Write([]byte(`<div>`))
		recursor = func(items []*BlockItem, listStyle string) {
			if err != nil {
				return
			}
			out.Write([]byte{'<'})
			out.Write([]byte(listStyle))
			out.Write([]byte{'>'})

			for _, item := range items {
				if err != nil {
					return
				}

				if item.Content == "" && len(item.Items) == 0 {
					continue
				}

				out.Write([]byte("<li>"))
				out.Write(StrToBytes(&item.Content))
				recursor(item.Items, listStyle)
				out.Write([]byte("</li>"))
			}

			out.Write([]byte("</"))
			out.Write([]byte(listStyle))
			out.Write([]byte{'>'})
		}

		if data.Style == "unordered" {
			listStyle = "ul"
		}

		recursor(data.Items, listStyle)

		out.Write([]byte(`</div>`))
		return
	},

	"image": func(data *BlockData, out io.Writer) (err error) {
		var caption = data.Caption

		out.Write([]byte(`<div class="mdl-grid">`))
		if caption == "" {
			caption = "Image"
		}

		out.Write([]byte(`<img src="`))
		if data.File != nil && data.File.Url != "" {
			out.Write([]byte(data.File.Url))
		} else {
			out.Write([]byte(data.Url))
		}
		out.Write([]byte(`" alt="`))
		out.Write([]byte(caption))
		out.Write([]byte(`" />`))

		out.Write([]byte(`</div>`))
		return
	},

	"quote": func(data *BlockData, out io.Writer) (err error) {
		out.Write([]byte(`<blockquote>`))
		out.Write(StrToBytes(&data.Text))
		out.Write([]byte(`"</blockquote>`))

		if data.Caption != "" {
			out.Write([]byte(" - "))
			out.Write([]byte(data.Caption))
		}

		return
	},

	"code": func(data *BlockData, out io.Writer) (err error) {
		out.Write([]byte(`<pre><code>`))
		out.Write([]byte(data.Code))
		out.Write([]byte(`"</code></pre>`))

		return
	},

	"embed": func(data *BlockData, out io.Writer) (err error) {
		out.Write([]byte(`<div class="mdl-grid">`))
		switch data.Service {
		case "vimeo":
			fmt.Fprintf(out, `<div><iframe src="%s" height="%d" frameborder="0" allow="autoplay; fullscreen; picture-in-picture" allowfullscreen></iframe></div>`,
				data.Embed, data.Height)
		case "youtube":
			fmt.Fprintf(out, `<div><iframe src="%s" width="%d" height="%d" title="YouTube video player" frameborder="0" allow="autoplay; fullscreen; picture-in-picture" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe></div>`,
				data.Embed, data.Width, data.Height)
		default:
			err = ErrEmbServiceNotSupported
		}
		out.Write([]byte(`</div>`))
		return
	},

	"warning": func(data *BlockData, out io.Writer) (err error) {
		out.Write([]byte(`<div class="qor-alert qor-alert--danger" role="alert" data-type="warning">`))
		if data.Title != "" {
			out.Write([]byte(`<span class="qor-alert-title">` + data.Title + `</span> `))
		}
		if data.Message != "" {
			out.Write([]byte(`<span class="qor-alert-message">` + data.Message + `</span>`))
		}
		out.Write([]byte(`</div>`))
		return
	},

	"raw": func(data *BlockData, out io.Writer) (err error) {
		out.Write([]byte(data.HTML))
		return
	},
}

var MarkdownTransformers = map[string]Transformer{
	"delimiter": func(data *BlockData, out io.Writer) (err error) {
		out.Write([]byte("\n--\n"))
		return
	},

	"header": func(data *BlockData, out io.Writer) (err error) {
		fmt.Fprintf(out, strings.Repeat("#", int(data.Level))+" "+data.Text)
		return
	},

	"paragraph": func(data *BlockData, out io.Writer) (err error) {
		fmt.Fprintf(out, "\n%s\n\n", data.Text)
		return
	},

	"table": func(data *BlockData, out io.Writer) (err error) {
		return
	},

	"list": func(data *BlockData, out io.Writer) (err error) {
		return
	},

	"image": func(data *BlockData, out io.Writer) (err error) {
		return
	},

	"quote": func(data *BlockData, out io.Writer) (err error) {
		return
	},

	"code": func(data *BlockData, out io.Writer) (err error) {
		out.Write([]byte(data.Code))
		return
	},

	"embed": func(data *BlockData, out io.Writer) (err error) {
		return
	},

	"warning": func(data *BlockData, out io.Writer) (err error) {
		return
	},
}
