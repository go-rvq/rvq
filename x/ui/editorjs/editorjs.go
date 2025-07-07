package editorjs

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"strings"
	"unicode"
)

type BlockFile struct {
	Url string
}

type BlockData struct {
	Title          string        `json:"title,omitempty"`
	Message        string        `json:"message,omitempty"`
	Text           string        `json:"text,omitempty"`
	Level          uint          `json:"level,omitempty"`
	Caption        string        `json:"caption,omitempty"`
	Url            string        `json:"url,omitempty"`
	File           *BlockFile    `json:"file,omitempty"`
	Stretched      bool          `json:"stretched,omitempty"`
	WithBackground bool          `json:"withBackground,omitempty"`
	WithBorder     bool          `json:"withBorder,omitempty"`
	WithHeadings   bool          `json:"withHeadings,omitempty"`
	Items          BlockItems    `json:"items,omitempty"`
	Content        [][]string    `json:"content,omitempty"` // table content
	Style          string        `json:"style,omitempty"`
	Code           string        `json:"code,omitempty"`
	Service        string        `json:"service,omitempty"`
	Source         string        `json:"source,omitempty"`
	Embed          string        `json:"embed,omitempty"`
	Width          uint16        `json:"width,omitempty"`
	Height         uint16        `json:"height,omitempty"`
	HTML           template.HTML `json:"html,omitempty"`
}

type BlockItemValue struct {
	Content string       `json:"content,omitempty"`
	Items   []*BlockItem `json:"items,omitempty"`
}

type BlockItem struct {
	BlockItemValue
}

func (i *BlockItem) UnmarshalJSON(bytes []byte) error {
	if IsJSON(bytes) {
		return json.Unmarshal(bytes, &i.BlockItemValue)
	}
	i.BlockItemValue = BlockItemValue{Content: string(bytes)}
	return nil
}

type BlockItems []*BlockItem

type Block struct {
	Type string     `json:"type,omitempty"`
	Data *BlockData `json:"data,omitempty"`
}

func FirstChar(b []byte) (_ byte, ok bool) {
	for _, b := range b {
		if !unicode.IsSpace(rune(b)) {
			return b, true
		}
	}
	return
}

func LastChar(bytes []byte) (_ byte, ok bool) {
	for i := len(bytes) - 1; i >= 0; i-- {
		b := bytes[i]
		if !unicode.IsSpace(rune(b)) {
			return b, true
		}
	}
	return
}

func IsJSON(bytes []byte) bool {
	if len(bytes) > 0 {
		if b, ok := FirstChar(bytes); ok && b == '{' || b == '[' {
			var b2 byte
			if b2, ok = LastChar(bytes); ok && b == '{' && b2 == '}' || b == '[' && b2 == ']' {
				return true
			}
		}
	}
	return false
}

func Parse(bytes []byte) (blocks Blocks, err error) {
	if len(bytes) == 0 {
		return
	}

	switch bytes[0] {
	case '{':
		var o struct {
			Blocks Blocks `json:"blocks,omitempty"`
		}
		err = json.Unmarshal(bytes, &o)
		blocks = o.Blocks
	case '[':
		err = json.Unmarshal(bytes, &blocks)
	default:
		err = errors.New("Unsupported JSON type")
	}
	return
}

func MustParse(bytes []byte) (blocks Blocks) {
	var err error
	if blocks, err = Parse(bytes); err != nil {
		panic(err)
	}
	return
}

type Blocks []*Block

func (b Blocks) Html(w io.Writer) (err error) {
	var t Transformer
	for _, b := range b {
		if t = HTMLTransformers[b.Type]; t != nil {
			if err = t(b.Data, w); err != nil {
				return
			}
		} else {
			return &BlockTypeNotSupported{b.Type}
		}
	}
	return
}

func (b Blocks) Markdown(w io.Writer) (err error) {
	var t Transformer
	for _, b := range b {
		if t = MarkdownTransformers[b.Type]; t != nil {
			if err = t(b.Data, w); err != nil {
				return
			}
		} else {
			return &BlockTypeNotSupported{b.Type}
		}
	}
	return
}

func Htmlify(rawValue string) (s string, err error) {
	if len(rawValue) > 0 {
		if rawValue[0] != '{' {
			return rawValue, nil
		}

		var blocks Blocks
		if blocks, err = Parse([]byte(rawValue)); err != nil {
			return
		}
		var out strings.Builder
		err = blocks.Html(&out)
		s = out.String()
	}
	return
}
