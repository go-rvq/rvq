package tiptap

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-rvq/rvq/utils"
	"golang.org/x/net/html"
)

var (
	imgDataImagePrefix = utils.UnsafeStringToBytes("data:image/")
	base64Bytes        = utils.UnsafeStringToBytes("base64,")
)

func getAttr(n *html.Node, name string) *html.Attribute {
	for i, attr := range n.Attr {
		if attr.Key == name {
			return &n.Attr[i]
		}
	}
	return nil
}

func traverse(n *html.Node, do func(n *html.Node) error) (err error) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err = do(c); err != nil {
			return
		}
		if err = traverse(c, do); err != nil {
			return
		}
	}
	return
}

func traverseTag(tagName string, n *html.Node, do func(n *html.Node) error) (err error) {
	return traverseType(html.ElementNode, n, func(n *html.Node) error {
		if n.Data == tagName {
			return do(n)
		}
		return nil
	})
}

func traverseImg(n *html.Node, do func(n *html.Node, src *html.Attribute, format string, data []byte) error) (err error) {
	return traverseTag("img", n, func(n *html.Node) (err error) {
		if src := getAttr(n, "src"); src != nil {
			if len(src.Val) == 0 {
				return
			}

			b := utils.UnsafeStringToBytes(src.Val)
			if bytes.HasPrefix(b, imgDataImagePrefix) {
				b := b[len(imgDataImagePrefix):]
				pos := bytes.IndexByte(b, ';')
				if pos == -1 || pos > 5 {
					return
				}

				format := string(b[:pos])
				switch format {
				case "jpeg", "png", "jpg", "gif", "webp":
					b = b[pos+1:]
					if !bytes.HasPrefix(b, base64Bytes) {
						return
					}
					b = b[len(base64Bytes):]
				default:
					return
				}

				var bf bytes.Buffer
				for _, c := range b {
					switch c {
					case ' ', '\t', '\r', '\n':
					default:
						bf.WriteByte(c)
					}
				}

				b = bf.Bytes()
				var (
					dbuf = make([]byte, base64.StdEncoding.DecodedLen(len(b)))
					bc   int
				)

				if bc, err = base64.StdEncoding.Decode(dbuf, b); err != nil {
					fmt.Println("Error decoding Base64:", err)
					return
				}
				b = dbuf[:bc]
				dbuf = nil
				if err = do(n, src, format, b); err != nil {
					return
				}
			}
		}
		return
	})
}

func traverseType(typ html.NodeType, n *html.Node, do func(n *html.Node) error) (err error) {
	return traverse(n, func(n *html.Node) error {
		if n.Type == typ {
			return do(n)
		}
		return nil
	})
}

var spacesRe = regexp.MustCompile(`\s+`)

func text(n *html.Node) string {
	var buf bytes.Buffer
	_ = traverseType(html.TextNode, n, func(n *html.Node) error {
		if len(n.Data) > 0 {
			buf.WriteString(n.Data)
		}
		return nil
	})
	return spacesRe.ReplaceAllString(strings.TrimSpace(buf.String()), " ")
}
