package script

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

func textOf(n *html.Node) string {
	var w bytes.Buffer
	for c := range n.Descendants() {
		if c.Type == html.TextNode {
			w.WriteString(c.Data)
		}
	}
	return w.String()
}

func nodeOfValueReplacer(n *html.Node) (r *html.Node) {
	var find func(n *html.Node) (r *html.Node)
	find = func(n *html.Node) (r *html.Node) {
		switch n.Type {
		case html.ElementNode:
			for c := range n.Descendants() {
				if r = find(c); r != nil {
					return
				}
			}
		case html.TextNode:
			if strings.Contains(n.Data, ResultReplacer) {
				return n
			}
		}
		return
	}

	for n.PrevSibling != nil {
		if r = find(n.PrevSibling); r != nil {
			return
		}
		n = n.PrevSibling
	}
	return
}
