package comments

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/web"
)

// The text of rich text: its words; what a comment shows of it: no scripts,
// no handlers; how many a tree has.
func TestTextAndCount(t *testing.T) {
	if got := PlainText("<p>Olá <b>mundo</b></p><p>dois</p>"); got != "Olá mundo dois" {
		t.Errorf("PlainText: %q", got)
	}
	if PlainText("<p><br></p>") != "" {
		t.Error("an empty editor is empty")
	}
	s := SanitizeHTML(`<p onclick="x()">a<script>alert(1)</script><a href="javascript:x">l</a></p>`)
	if strings.Contains(s, "script") || strings.Contains(s, "onclick") || strings.Contains(s, "javascript") {
		t.Errorf("SanitizeHTML: %q", s)
	}
	tree := []*Comment{{ID: "a", Replies: []*Comment{{ID: "b", Replies: []*Comment{{ID: "c"}}}}}, {ID: "d"}}
	if Count(tree) != 4 {
		t.Errorf("Count: %d", Count(tree))
	}
	if Initials("Ana Maria Souza") != "AS" || Initials("") != "?" {
		t.Error("Initials")
	}
	if HumanSize(2048) != "2 KB" || HumanSize(3<<20) != "3.0 MB" {
		t.Error("HumanSize")
	}
}

// The limits of the files: each one, all of them together; none, any.
func TestValidate(t *testing.T) {
	b := &Builder{
		ValidateFile:    MaxFileSize(func(*web.EventContext, string) int64 { return 10 }),
		ValidateMessage: MaxTotalSize(func(*web.EventContext, string) int64 { return 15 }),
	}
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	ok := &Write{Files: []Attachment{{Name: "a", Data: make([]byte, 8)}, {Name: "b", Data: make([]byte, 6)}}}
	if err := b.Validate(ctx, ok); err != nil {
		t.Errorf("within the limits: %v", err)
	}
	if err := b.Validate(ctx, &Write{Files: []Attachment{{Name: "big", Data: make([]byte, 11)}}}); err == nil || !strings.Contains(err.Error(), "big") {
		t.Errorf("a file too large: %v", err)
	}
	if err := b.Validate(ctx, &Write{Files: []Attachment{{Name: "a", Data: make([]byte, 9)}, {Name: "b", Data: make([]byte, 9)}}}); err == nil {
		t.Error("too large together")
	}
	if err := (&Builder{}).Validate(ctx, &Write{Files: []Attachment{{Data: make([]byte, 1<<20)}}}); err != nil {
		t.Errorf("no limits: %v", err)
	}
}
