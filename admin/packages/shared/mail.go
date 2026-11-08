package shared

import (
	"embed"
	"io/fs"

	"github.com/go-rvq/rvq/x/mailtpl"
)

// mails are the default templates of the mails of the sharing (mailtpl):
// templates/mail/shared/confirm.<language>.gadx — the owner asked to
// confirm a sharing, a role, a new owner.
//
//go:embed mail
var mails embed.FS

func init() {
	sub, err := fs.Sub(mails, "mail")
	if err != nil {
		panic(err)
	}
	mailtpl.Register("rvq.shared", sub)
}
