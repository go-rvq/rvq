package mail_sender

import (
	"strings"
	"testing"
)

// The files of a reply: its attachment, a picture pasted (inline), not its
// text nor its HTML.
func TestParseInboundAttachments(t *testing.T) {
	raw := strings.ReplaceAll(`From: Zz <zz@example.com>
Subject: Re: x
Message-ID: <r1@example.com>
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=B1

--B1
Content-Type: multipart/related; boundary=B2

--B2
Content-Type: text/plain; charset=utf-8

Segue o arquivo.
--B2
Content-Type: image/png; name=foto.png
Content-Disposition: inline
Content-ID: <img1>
Content-Transfer-Encoding: base64

iVBORw0KGgo=
--B2--
--B1
Content-Type: application/pdf
Content-Disposition: attachment; filename=orcamento.pdf
Content-Transfer-Encoding: base64

JVBERi16eg==
--B1--
`, "\n", "\r\n")
	m, err := ParseInbound([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(m.Text) != "Segue o arquivo." {
		t.Errorf("the text: %q", m.Text)
	}
	if len(m.Attachments) != 2 {
		t.Fatalf("the files: %+v", m.Attachments)
	}
	names := map[string]string{}
	for _, a := range m.Attachments {
		names[a.Name] = string(a.Data)
	}
	if names["orcamento.pdf"] != "%PDF-zz" || names["foto.png"] == "" {
		t.Errorf("the files: %v", names)
	}
}
