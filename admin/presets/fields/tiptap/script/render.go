package script

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gad-lang/gad"
	"github.com/gad-lang/gad/parser"
	"github.com/gad-lang/gad/parser/source"
	h "github.com/go-rvq/htmlgo"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const ResultReplacer = "{}"

type BuilderOption func(b *Builder)

func WithDelimiter(delimiter parser.MixedDelimiter) BuilderOption {
	return func(b *Builder) {
		b.delimiter = delimiter
	}
}

type Builder struct {
	delimiter parser.MixedDelimiter
}

func New(opts ...BuilderOption) *Builder {
	b := &Builder{
		delimiter: parser.MixedDelimiter{
			Start: []rune("{%"),
			End:   []rune("%}"),
		},
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *Builder) Delimiter() parser.MixedDelimiter {
	return b.delimiter
}

func (b *Builder) Build(htmlValue string) (result *ScriptBuilderResult, err error) {
	var (
		doc    *html.Node
		w      bytes.Buffer
		nodes  []*Node
		script = func(n *html.Node, val bool) (string, int, bool) {
			v := textOf(n)

			// no remove left spaces to prevent error on SourcePos
			if len(strings.TrimSpace(v)) == 0 {
				return "", 1, false
			}

			var (
				lc       = len(strings.Split(v, "\n"))
				prefix   = string(b.delimiter.Start)
				sufix    string
				hasEqual bool
			)

			if strings.HasPrefix(v, "- ") {
				prefix += "-"
				// no change length to prevent error on SourcePos
				v = " " + v[1:]
			} else if val {
				if strings.HasPrefix(v, "= ") {
					// no change length to prevent error on SourcePos
					v = " " + v[1:]
					hasEqual = true
				} else if strings.HasPrefix(v, "-= ") {
					prefix += "-"
					// no change length to prevent error on SourcePos
					v = "  " + v[2:]
					hasEqual = true
				}
			}

			if strings.HasSuffix(v, " -") {
				v = v[:len(v)-2]
				sufix = " -"
			}

			if val {
				prefix += "="
			}

			sufix += string(b.delimiter.End)

			v = fmt.Sprintf("%s\n//!!%d\n%s\n%s", prefix, len(nodes), v, sufix)
			return v, lc, hasEqual
		}
		addNode = func(n *html.Node, val bool) {
			var (
				lc       int
				hasEqual bool
				data     string
			)

			data, lc, hasEqual = script(n, val)

			old := *n
			reset := func() {
				n.Data = old.Data
				n.Type = old.Type
				n.FirstChild = old.FirstChild
				n.LastChild = old.LastChild
			}
			n.Type = html.RawNode
			n.FirstChild = nil
			n.LastChild = nil

			if hasEqual && n.PrevSibling != nil {
				n.Data = ""

				vNode := &html.Node{
					Type: html.RawNode,
					Data: data,
				}

				nor := nodeOfValueReplacer(n)

				if nor != nil {
					var (
						oldNor    = *nor
						parts     = strings.SplitN(nor.Data, ResultReplacer, 2)
						rightNode *html.Node
						nextOfNor = nor.NextSibling
					)

					nor.Data = parts[0]

					if len(parts) > 1 {
						rightNode = &html.Node{
							Type: html.TextNode,
							Data: parts[1],
						}

						nor.Parent.InsertBefore(rightNode, nextOfNor)
					}

					nor.Parent.InsertBefore(vNode, rightNode)
					oldReset := reset

					reset = func() {
						vNode.Parent.RemoveChild(vNode)
						if rightNode != nil {
							rightNode.Parent.RemoveChild(rightNode)
						}
						*nor = oldNor
						oldReset()
					}
				} else {
					ps := n.PrevSibling
					psc := &html.Node{
						Type: html.RawNode,
						Data: data,
					}

					ps.AppendChild(psc)
					oldReset := reset

					reset = func() {
						oldReset()
						ps.RemoveChild(psc)
					}
				}
				n.Data = ""
			} else {
				n.Data = data
			}

			nodes = append(nodes, &Node{
				Node:      n,
				Old:       &old,
				Index:     len(nodes),
				LineCount: lc,
				reset:     reset,
			})
		}
	)
	if doc, err = html.Parse(strings.NewReader("<body>" + htmlValue + "</body>")); err != nil {
		return
	}

	var body *html.Node
	for node := range doc.FirstChild.ChildNodes() {
		if node.Type == html.ElementNode && node.DataAtom == atom.Body {
			body = node
			break
		}
	}

	if body == nil {
		return
	}

nodes:
	for n := range body.Descendants() {
		if n.Type == html.ElementNode {
			if n.DataAtom == atom.Script {
				switch GetTypeFromNode(n) {
				case TypeBlock:
					addNode(n, false)
					continue nodes
				case TypeInline:
					addNode(n, false)
					continue nodes
				case TypeValue, TypeCodeBlockValue:
					addNode(n, true)
					continue nodes
				}
			}
		}
	}

	for n := range body.ChildNodes() {
		if err = html.Render(&w, n); err != nil {
			return
		}
	}

	result = &ScriptBuilderResult{
		Script:    w.String(),
		Root:      body,
		Nodes:     nodes,
		delimiter: b.delimiter,
	}
	return
}

type Node struct {
	Node,
	Old *html.Node
	Index     int
	LineCount int
	Messages  []*vx.Message
	reset     func()
	file      *source.SliceFile
}

func (n *Node) Reset() {
	n.reset()
}

type ScriptBuilderResult struct {
	Script       string
	Root         *html.Node
	Nodes        []*Node
	delimiter    parser.MixedDelimiter
	nodesFileSet *source.FileSet
	TraceLines   int
}

func (b *ScriptBuilderResult) Delimiter() parser.MixedDelimiter {
	return b.delimiter
}

func (r *ScriptBuilderResult) Render(out io.Writer) (err error) {
	for n := range r.Root.ChildNodes() {
		if err = html.Render(out, n); err != nil {
			return
		}
	}
	return
}

func (n *Node) AsMessages() (err error) {
	n.Node.Type = html.RawNode
	s, _ := h.MarshallString(h.Tag("vx-messages").Attr(":items", h.JSONString(n.Messages)), context.Background())
	n.Node.Data = s
	n.Node.Attr = nil
	return
}

func (n *Node) UpdateMessagesAttribute() (err error) {
	return n.updateMessagesAttribute("data-messages")
}

func (n *Node) updateMessagesAttribute(attrName string) (err error) {
	var (
		attrIndex int
		attr      *html.Attribute
		node      = n.Node
	)

	for i, attribute := range node.Attr {
		if attribute.Key == attrName {
			attr = &node.Attr[i]
			attrIndex = i
			break
		}
	}

	if len(n.Messages) == 0 {
		if attr != nil {
			// remove
			node.Attr = append(node.Attr[:attrIndex], node.Attr[attrIndex+1:]...)
		}
		return
	}

	if attr == nil {
		node.Attr = append(node.Attr, html.Attribute{Key: attrName, Val: ""})
		attr = &node.Attr[len(node.Attr)-1]
	}

	var b []byte
	if b, err = json.Marshal(n.Messages); err == nil {
		attr.Val = string(b)
	}
	return
}

func (r *ScriptBuilderResult) NodeSourceFileFromPos(pos source.FilePos) (n *Node, curPos source.FilePos) {
	var (
		nodeStartLine int
		nodeIndex     int
		file          = pos.File
	)

	for line := pos.Line; line >= 1; line-- {
		if d, err := file.Data.LineData(line); err == nil {
			if bytes.HasPrefix(d, []byte("//!!")) {
				nodeStartLine = line + 1
				nodeIndex, _ = strconv.Atoi(string(d[4:]))
				n = r.Nodes[nodeIndex]
				break
			}
		}
	}

	var err error
	if n.file == nil {
		n.file, err = pos.File.Slice(r.nodesFileSet, fmt.Sprintf("(node:%d)", nodeIndex), nodeStartLine, n.LineCount)
		if err != nil {
			return
		}
	}

	if curPos, err = n.file.CastPos(pos); err != nil {
		panic(err)
	}
	return
}

func (r *ScriptBuilderResult) AddNodeError(messages *Messages, err *ScriptError) bool {
	var (
		n, p  = r.NodeSourceFileFromPos(err.Pos)
		trace bytes.Buffer
		lines = r.TraceLines
	)

	if lines == 0 {
		lines = 10
	}

	if p.IsValid() {
		b := err.Pos.File.Data.Bytes()
		_ = b
		s := []rune(string(err.Pos.File.Data.Bytes()))
		_ = s
		p.File.Data.TraceLines(&trace, p.Line, p.Column, lines, lines)
		msg := messages.FormateTypeError(err.Type, p, err.Message)

		n.Messages = append(n.Messages, &vx.Message{
			Value: msg,
			Type:  "error",
			Detail: &vx.MessageDetails{
				Raw: fmt.Sprintf("<pre>%s\n\n%s</pre>", msg, html.EscapeString(trace.String())),
			},
		})
		return true
	}
	return false
}

func (r *ScriptBuilderResult) CheckError(messages *Messages, err error) (outError error, nodeFailed bool) {
	if err == nil {
		return
	}

	if r.nodesFileSet == nil {
		r.nodesFileSet = source.NewFileSet()
	}

	switch t := err.(type) {
	case *gad.RuntimeError:
		fs := t.FileSet()
		for i := 0; i < len(t.Trace); i++ {
			pos := t.Trace[i]
			fpos := fs.Position(pos)
			msg := t.Err.Error()
			e := NewTemplateScriptError(t, ScriptErrorTypeRun, fpos, msg)
			if r.AddNodeError(messages, e) {
				nodeFailed = true
				if outError == nil {
					outError = ErrScriptFailure
				}
			} else {
				outError = e
			}
		}
	case parser.ErrorList:
		for _, e := range t {
			if file := e.Pos.File; file.Name == gad.MainName {
				msg := e.Msg
				if sep := strings.Index(msg, "\n\tat (main)"); sep != -1 {
					msg = msg[:sep]
				}
				e := NewTemplateScriptError(t, ScriptErrorTypeParse, e.Pos, msg)

				if r.AddNodeError(messages, e) {
					nodeFailed = true
					if outError == nil {
						outError = &ScriptError{Type: ScriptErrorTypeParse, Pos: e.Pos, Message: msg}
					}
				} else {
					outError = e
				}
			}
		}
	case *gad.CompilerError:
		msg := t.Err.Error()
		pos := t.FileSet.Position(t.Node.Pos())
		e := NewTemplateScriptError(t, ScriptErrorTypeCompile, pos, msg)
		if nodeFailed = r.AddNodeError(messages, e); !nodeFailed {
			outError = ErrScriptFailure
		}
	}
	if outError == nil {
		outError = err
	}
	return
}

func (r *ScriptBuilderResult) UpdateNodes(ro bool) (err error) {
	for _, node := range r.Nodes {
		if ro && len(node.Messages) > 0 {
			node.AsMessages()
		} else {
			node.Reset()
			if err = node.UpdateMessagesAttribute(); err != nil {
				return
			}
		}
	}
	return
}

type ScriptRunBuilderOption func(b *ScriptRunBuilder)

func WithCompileOptions(co *gad.CompileOptions) ScriptRunBuilderOption {
	return func(b *ScriptRunBuilder) {
		b.co = co
	}
}

func WithConfigureVM(configure func(vm *gad.VM) error) ScriptRunBuilderOption {
	return func(b *ScriptRunBuilder) {
		b.configureVM = configure
	}
}

func WithVMRunOptions(opts *gad.RunOpts) ScriptRunBuilderOption {
	return func(b *ScriptRunBuilder) {
		b.runOpts = opts
	}
}

type ScriptRunBuilder struct {
	r           *ScriptBuilderResult
	co          *gad.CompileOptions
	runOpts     *gad.RunOpts
	configureVM func(vm *gad.VM) error
}

func (r *ScriptBuilderResult) Runner(opt ...ScriptRunBuilderOption) *ScriptRunBuilder {
	b := &ScriptRunBuilder{
		r: r,
	}

	for _, opt := range opt {
		opt(b)
	}

	if b.co == nil {
		b.co = &gad.CompileOptions{}
	}
	if b.runOpts == nil {
		b.runOpts = &gad.RunOpts{}
	}

	return b
}

func (b *ScriptRunBuilder) Run() (err error) {
	b.co.ParserOptions.Mode.Set(parser.ParseMixed | parser.ParseConfigDisabled).Clear(parser.ParseComments)
	b.co.ScannerOptions.MixedDelimiter = b.r.delimiter

	var bc *gad.Bytecode

	if bc, err = gad.Compile([]byte(b.r.Script), *b.co); err != nil {
		return
	}

	vm := gad.NewVM(bc)

	if b.configureVM != nil {
		if err = b.configureVM(vm); err != nil {
			return
		}
	}

	if _, err = vm.RunOpts(b.runOpts); err != nil {
		return
	}

	return
}

func Build(htmlValue string, opt ...BuilderOption) (result *ScriptBuilderResult, err error) {
	var r *ScriptBuilderResult
	if r, err = New(opt...).Build(htmlValue); err != nil {
		return
	}
	return r, nil
}
