// Command mailtest runs the mail server of the local development (package
// mailtest): SMTP and IMAP in memory, nothing leaving the machine. An
// application whose RVQ_MAIL_SMTP is its SMTP address sends every message
// there; each one is logged as it arrives, and read by IMAP — any address,
// the password "test".
//
//	go run github.com/go-rvq/rvq/cmd/mailtest -smtp 127.0.0.1:2525 -imap 127.0.0.1:1143
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-rvq/rvq/admin/packages/mail-sender/mailtest"
)

func main() {
	smtpAddr := flag.String("smtp", "127.0.0.1:2525", "where the SMTP server listens")
	imapAddr := flag.String("imap", "127.0.0.1:1143", "where the IMAP server listens")
	verbose := flag.Bool("v", false, "log what the servers say")
	flag.Parse()

	opts := mailtest.Options{SMTPAddr: *smtpAddr, IMAPAddr: *imapAddr}
	if *verbose {
		opts.Logger = log.New(os.Stderr, "mailtest: ", log.LstdFlags)
	}
	s, err := mailtest.Start(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	fmt.Printf("mailtest: SMTP %s, IMAP %s (any address, password %q)\n", s.SMTPAddr, s.IMAPAddr, mailtest.Password)
	fmt.Printf("mailtest: the application sends here with RVQ_MAIL_SMTP=%s\n", s.SMTPAddr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	seen := 0
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-stop:
			return
		case <-tick.C:
			msgs := s.Messages()
			for _, m := range msgs[seen:] {
				fmt.Printf("%s  %s → %s  %q  %s\n", m.Received.Format("15:04:05"), m.From,
					strings.Join(m.To, ", "), m.Subject, m.MessageID)
			}
			seen = len(msgs)
		}
	}
}
