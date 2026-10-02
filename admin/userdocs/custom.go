package userdocs

import (
	"path"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

// Custom are documents of the tree of no part of the admin — a guide, a
// tutorial —, put under the node Parent ("" the top of the tree), first among
// its children or last.
type Custom struct {
	Parent string
	First  bool
	Nodes  []*CustomNode
}

// CustomNode is a node of Custom: its document is the one of its ID
// (DocFile: ID/README.md), in any source, in the language of the request;
// its title, the heading of the document (# Title).
type CustomNode struct {
	ID       string
	Icon     string
	Children []*CustomNode
}

// Custom registers documents of no part of the admin in the tree (see
// Custom). Nodes whose parent is not in the tree of a request — a model it
// may not see — are left out of it.
func (b *Builder) Custom(c Custom) *Builder {
	b.customs = append(b.customs, c)
	return b
}

// withCustom is tree with the custom nodes, in the language of ctx.
func (b *Builder) withCustom(tree []*Node, ctx *web.EventContext) []*Node {
	locale := b.locale(ctx)
	var nodes func(cs []*CustomNode) []*Node
	nodes = func(cs []*CustomNode) (r []*Node) {
		for _, c := range cs {
			icon := c.Icon
			if icon == "" {
				icon = "mdi-file-document-outline"
			}
			r = append(r, &Node{ID: c.ID, Title: b.docTitle(locale, c.ID), Icon: icon, Children: nodes(c.Children)})
		}
		return
	}
	for _, c := range b.customs {
		add := nodes(c.Nodes)
		into := &tree
		if c.Parent != "" {
			p := findNode(tree, c.Parent)
			if p == nil {
				continue
			}
			into = &p.Children
		}
		if c.First {
			*into = slices.Concat(add, *into)
		} else {
			*into = append(*into, add...)
		}
	}
	return tree
}

// titleTTL is how long the title of a custom node is kept (titles): the tree
// is made several times a page, and a title is read from its document — a
// change of it (an edit, a translation of the boot) shows within it.
const titleTTL = 30 * time.Second

type cachedTitle struct {
	title string
	at    time.Time
}

// titles are the titles of the custom nodes read (docTitle), by locale and
// node; an edit of the documents empties it.
var titles sync.Map

// docTitle is the title of the document of node in locale: its first
// heading, else its last word.
func (b *Builder) docTitle(locale, node string) string {
	key := locale + "\x00" + node
	if c, ok := titles.Load(key); ok && time.Since(c.(cachedTitle).at) < titleTTL {
		return c.(cachedTitle).title
	}
	title := presets.HumanizeString(path.Base(node))
	if _, _, content, err := b.find(locale, node); err == nil {
		for _, l := range strings.Split(content, "\n") {
			if t, ok := strings.CutPrefix(strings.TrimSpace(l), "# "); ok {
				title = strings.TrimSpace(t)
				break
			}
		}
	}
	titles.Store(key, cachedTitle{title, time.Now()})
	return title
}

// forgetTitles empties the titles read: the documents changed.
func forgetTitles() { titles.Clear() }
