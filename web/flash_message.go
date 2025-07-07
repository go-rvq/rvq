package web

import (
	"fmt"

	h "github.com/theplant/htmlgo"
)

type FlashMessage struct {
	Text     string `json:"text,omitempty"`
	HtmlText string `json:"htmlText,omitempty"`
	Color    string `json:"color,omitempty"`
	Duration int    `json:"duration,omitempty"`
	// Detail this comp open in Dialog
	Detail h.HTMLComponent `json:"-"`
	// HtmlDetail this comp is Detail wraped into Dialog
	HtmlDetail  string `json:"rawDetail,omitempty"`
	NotClosable bool   `json:"notClosable,omitempty"`
}

func (m *FlashMessage) Error() string {
	return m.Text
}

func NewFlashMessages(msg any, color ...string) (m FlashMessages) {
	if msg == nil {
		return
	}
	switch t := msg.(type) {
	case *FlashMessage:
		return FlashMessages{t}
	case FlashMessages:
		return t
	case []*FlashMessage:
		return FlashMessages(t)
	default:
		return FlashMessages{NewFlashMessage(msg, color...)}
	}
}

func NewFlashMessage(msg any, color ...string) (m *FlashMessage) {
	if m, _ := msg.(*FlashMessage); m != nil {
		return m
	}

	m = &FlashMessage{}
	for _, m.Color = range color {
	}

	switch t := msg.(type) {
	case *ValidationErrors:
		gErr := t.GetGlobalError()
		if len(gErr) > 0 {
			m.Text = gErr
			if m.Color == "" {
				m.Color = "error"
			}
		}
	case error:
		if m.Color == "" {
			m.Color = "error"
		}
		m.Text = t.Error()
	case string:
		m.Text = t
	case h.HTMLComponent:
		m.Text = "Bad flash message type. Contact your administrator."
		m.Color = "error"
	case *FlashMessage:
		m = t
	case FlashMessage:
		*m = t
	default:
		m.Text = fmt.Sprintf("%v", t)
	}
	return
}

func (m *FlashMessage) SetDuration(v int) *FlashMessage {
	m.Duration = v
	return m
}

func (m *FlashMessage) SetDetail(c h.HTMLComponent) *FlashMessage {
	m.Detail = c
	return m
}

func (m *FlashMessage) SetClosable(v bool) *FlashMessage {
	m.NotClosable = !v
	return m
}

func (e *EventContext) FlashMessage(msg string, color ...string) *FlashMessage {
	fm := &FlashMessage{
		Text: msg,
	}
	for _, s := range color {
		fm.Color = s
	}
	e.Flash = fm
	return fm
}

type FlashMessages []*FlashMessage

func (m *FlashMessages) Append(message ...*FlashMessage) {
	*m = append(*m, message...)
}
