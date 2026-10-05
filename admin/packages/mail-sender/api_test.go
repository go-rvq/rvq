package mail_sender

import (
	"strings"
	"testing"
)

// A mail answered goes to its Reply-To — who wrote what it carries —, when
// the message says one.
func TestMessageReplyTo(t *testing.T) {
	m, err := Message().To("site@example.com").ReplyTo("visitor@example.com").Subject("Hi").Body("<p>x</p>").Build()
	if err != nil {
		t.Fatal(err)
	}
	if got := m.GetGenHeader("Reply-To"); len(got) != 1 || !strings.Contains(got[0], "visitor@example.com") {
		t.Errorf("Reply-To = %v", got)
	}
	m, _ = Message().To("site@example.com").Build()
	if got := m.GetGenHeader("Reply-To"); len(got) != 0 {
		t.Errorf("Reply-To without one = %v", got)
	}
}
