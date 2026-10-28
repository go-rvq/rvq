package mail_sender

import (
	"strings"
	"testing"
)

// The subject in parts, as the settings say: the answer first; the site's
// name (the settings', brackets out; else the application's), the thing and
// the record by the rule of the family and of to whom it goes — the family's
// own, else the others', else the default —; as it is when verbatim.
func TestSubjectParts(t *testing.T) {
	defer func(f func() string) { SiteName = f }(SiteName)
	SiteName = func() string { return "Site" }
	msg := func() *MessageBuilder {
		return Message().Kind("issue.reply").Topic("ISSUE - Chamados").Ref("#12").Subject("Title").Reply()
	}
	for name, c := range map[string]struct {
		c    MailSender
		b    *MessageBuilder
		want string
	}{
		"default, to the team":    {MailSender{SubjectPrefix: "[PAC] "}, msg(), "Re: [PAC] [ISSUE - Chamados] [#12] Title"},
		"default, outside":        {MailSender{SubjectPrefix: "PAC"}, msg().External(true), "Re: [PAC] [#12] Title"},
		"the application's name":  {MailSender{}, msg(), "Re: [Site] [ISSUE - Chamados] [#12] Title"},
		"the family's own rule":   {MailSender{SubjectPrefix: "PAC", SubjectRules: SubjectRules{"issue": {Team: SubjectRef}}}, msg(), "Re: [#12] Title"},
		"the others' rule":        {MailSender{SubjectPrefix: "PAC", SubjectRules: SubjectRules{SubjectFamilyOther: {Team: SubjectTopic}}}, msg(), "Re: [ISSUE - Chamados] Title"},
		"nothing but the subject": {MailSender{SubjectRules: SubjectRules{"issue": {}}}, msg(), "Re: Title"},
		"verbatim":                {MailSender{SubjectPrefix: "PAC"}, Message().Kind("issue.reply").Subject("Re: [X] kept").Verbatim(), "Re: [X] kept"},
		"labels after the site":   {MailSender{SubjectPrefix: "PAC"}, Message().Kind("x.y").Label("L").Subject("S"), "[PAC] [L] S"},
	} {
		c.c.composeWith(c.b)
		if got := c.b.FormattedSubject(); got != c.want {
			t.Errorf("%s: %q, want %q", name, got, c.want)
		}
	}
	// the rules kept as JSON
	r := SubjectRules{"issue": {Team: SubjectAll, External: SubjectRef}}
	v, _ := r.Value()
	var back SubjectRules
	if err := back.Scan(v); err != nil || back.Of("issue") != r["issue"] || back.Of("contact") != DefaultSubjectRule {
		t.Errorf("JSON: %v %v", back, err)
	}
}

// The ids of a message of a record: its token and kind in its headers
// always; at the end of its body — the reference and the thing, the token,
// the Message-ID and the kind — when the settings say, once.
func TestMessageIDs(t *testing.T) {
	mk := func(footer bool) *MessageBuilder {
		b := Message().To("a@example.com").Kind("issue.reply").Ref("#12").Topic("ISSUE - Chamados").
			Identify("hermon:issue:01a1").MessageID("m.1@example.com").Subject("T").Body("<p>x</p>")
		(&MailSender{IDFooter: footer}).composeWith(b)
		return b
	}
	b := mk(true)
	m, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}
	if h := m.GetGenHeader("X-Mail-Record"); len(h) != 1 || h[0] != "hermon:issue:01a1" {
		t.Errorf("X-Mail-Record: %v", h)
	}
	if h := m.GetGenHeader("X-Mail-Kind"); len(h) != 1 || h[0] != "issue.reply" {
		t.Errorf("X-Mail-Kind: %v", h)
	}
	for _, w := range []string{"Ref.: #12 · ISSUE - Chamados", "ID: hermon:issue:01a1", "Message-ID: &lt;m.1@example.com&gt; · issue.reply"} {
		if !strings.Contains(b.body, w) {
			t.Errorf("the footer: no %q in %s", w, b.body)
		}
	}
	if _, err := b.Build(); err != nil || strings.Count(b.body, "data-mail-ids") != 1 {
		t.Errorf("the footer twice: %s", b.body)
	}
	if b := mk(false); func() bool { b.Build(); return strings.Contains(b.body, "data-mail-ids") }() {
		t.Error("the footer, off")
	}
}
