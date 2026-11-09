package mail_sender

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// The attachments go with the message, built more than once (sent again).
func TestBuildAttachments(t *testing.T) {
	b := Message().To("a@example.com").Subject("s").Body("<p>b</p>").
		Attachment(&Attachment{Name: "pauta.txt", Body: io.NopCloser(strings.NewReader("1. abertura"))})
	for i := 0; i < 2; i++ {
		m, err := b.Build()
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		if _, err := m.WriteTo(&buf); err != nil {
			t.Fatal(err)
		}
		if out := buf.String(); !strings.Contains(out, "pauta.txt") || !strings.Contains(out, "multipart/mixed") {
			t.Fatalf("build %d: no attachment\n%s", i, out)
		}
	}
}
