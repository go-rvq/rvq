package script

import "golang.org/x/net/html"

type Type uint8

const (
	TypeBlock = iota + 1
	TypeInline
	TypeValue
	TypeCodeBlockValue

	AttrTypeScript         = "text/script"
	AttrTypeInline         = AttrTypeScript + "Block"
	AttrTypeValue          = AttrTypeScript + "Value"
	AttrTypeCodeBlockValue = "text/codeBlockScriptValue"
)

func GetTypeFromNode(n *html.Node) Type {
	for _, a := range n.Attr {
		if a.Key == "type" {
			switch a.Val {
			case AttrTypeScript:
				return TypeBlock
			case AttrTypeInline:
				return TypeInline
			case AttrTypeValue:
				return TypeValue
			case AttrTypeCodeBlockValue:
				return TypeCodeBlockValue
			}
			return 0
		}
	}
	return 0
}
