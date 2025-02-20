package str_utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// HumanizeString humanize separates string based on capitalizd letters
// e.g. "OrderItem" -> "Order Item, CNNName to CNN Name"
func HumanizeString(str string) string {
	var human []rune
	input := []rune(str)
	for i, l := range input {
		if i > 0 && unicode.IsUpper(l) {
			if (!unicode.IsUpper(input[i-1]) && input[i-1] != ' ') || (i+1 < len(input) && !unicode.IsUpper(input[i+1]) && input[i+1] != ' ' && input[i-1] != ' ') {
				human = append(human, rune(' '))
			}
		}
		human = append(human, l)
	}
	return strings.Title(string(human))
}

// NamifyString Joins string parts based on special separeted chars
// e.g. "order_item-data" -> "OrderItemData"
func NamifyString(s string) string {
	s = RemoveAccents(s)

	var (
		human   []rune
		toUpper bool
	)

	s = "_" + s
	for _, c := range s {
		if c == '_' || c == '-' {
			toUpper = true
			continue
		} else if c == '!' {
			toUpper = true
		} else if toUpper {
			toUpper = false
			if c >= 'a' && c <= 'z' {
				c -= 'a' - 'A'
			}
		}
		human = append(human, c)
	}
	return string(human)
}

// NamifyString Joins string parts based on special separeted chars
// e.g. "order_item-data" -> []string{"order", "item", "data"}
func SplitString(s string) (parts []string) {
	var buf []rune
	var split bool
	s = "_" + s

loop:
	for _, c := range s {
		switch c {
		case '_', '-', ' ', '.', '!', '@':
			split = true
			continue loop
		}

		if split {
			split = false
			if len(buf) > 0 {
				parts = append(parts, string(buf))
				buf = nil
			}
		}
		buf = append(buf, c)
	}

	if len(buf) > 0 {
		parts = append(parts, string(buf))
	}

	return
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}
