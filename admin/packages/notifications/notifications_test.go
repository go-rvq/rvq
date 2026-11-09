package notifications

import (
	"context"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/i18n"
	"github.com/google/uuid"
	"golang.org/x/text/language"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := AutoMigrate(db); err != nil {
		t.Fatal(err)
	}
	return db
}

// Told to each user once; unread counted, the newest first; read — one's
// own only —, all read.
func TestNotifications(t *testing.T) {
	db := testDB(t)
	ana, bia := uuid.New(), uuid.New()
	if err := Notify(db, Notification{Title: "A", Level: LevelWarning, Source: "x"}, ana, bia, ana, uuid.Nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Millisecond)
	if err := Notify(db, Notification{Title: "B"}, ana); err != nil {
		t.Fatal(err)
	}
	if n, _ := Unread(db, ana); n != 2 {
		t.Fatalf("ana's unread: %d", n)
	}
	if n, _ := Unread(db, bia); n != 1 {
		t.Fatalf("bia's unread: %d", n)
	}
	ns, _ := Latest(db, ana, 10)
	if len(ns) != 2 || ns[0].Title != "B" || ns[1].Level != LevelWarning || ns[0].Level != LevelInfo {
		t.Fatalf("ana's latest: %+v", ns)
	}
	bs, _ := Latest(db, bia, 10)
	// bia's id marked by ana: not read
	if err := MarkRead(db, ana, bs[0].ID); err != nil {
		t.Fatal(err)
	}
	if n, _ := Unread(db, bia); n != 1 {
		t.Error("another's read")
	}
	if err := MarkRead(db, ana, ns[0].ID); err != nil {
		t.Fatal(err)
	}
	if n, _ := Unread(db, ana); n != 1 {
		t.Errorf("one read: %d", n)
	}
	if err := MarkAllRead(db, ana); err != nil {
		t.Fatal(err)
	}
	if n, _ := Unread(db, ana); n != 0 {
		t.Errorf("all read: %d", n)
	}
}

// Every word in each language.
func TestMessagesFilled(t *testing.T) {
	for _, m := range []*Messages{Messages_en_US, Messages_pt_BR} {
		if m.BellTitle == "" || m.MarkAllRead == "" || m.SeeAll == "" || m.Notifications == "" || m.Open == "" {
			t.Errorf("%+v", m)
		}
	}
}

// Said by messages: each reader reads it in their language, with the
// arguments of the title and of the body; written in the language of who
// made it; without the module, as written.
type textMessages struct{ Below, Left string }

func TestNotificationText(t *testing.T) {
	const key i18n.ModuleKey = "test/notifications"
	b := i18n.New().SupportLanguages(language.English, language.BrazilianPortuguese).
		RegisterForModule(language.English, key, &textMessages{Below: "Below zero: %s", Left: "Left with %s."}).
		RegisterForModule(language.BrazilianPortuguese, key, &textMessages{Below: "Saldo negativo: %s", Left: "Ficou com %s."})
	en, pt := b.LanguageContext(context.Background(), "en"), b.LanguageContext(context.Background(), "pt-BR")
	n := Notification{}.Say(pt, key, Msg("Below", "Caixa"), Msg("Left", "-10.00"))
	if n.Title != "Saldo negativo: Caixa" || n.Body != "Ficou com -10.00." {
		t.Fatalf("written: %q, %q", n.Title, n.Body)
	}
	db := testDB(t)
	u := uuid.New()
	if err := Notify(db, n, u); err != nil {
		t.Fatal(err)
	}
	ns, _ := Latest(db, u, 1)
	if title, body := ns[0].Text(en); title != "Below zero: Caixa" || body != "Left with -10.00." {
		t.Fatalf("in english: %q, %q", title, body)
	}
	if title, _ := ns[0].Text(pt); title != "Saldo negativo: Caixa" {
		t.Fatalf("in portuguese: %q", title)
	}
	ns[0].Module = "unknown"
	if title, body := ns[0].Text(en); title != n.Title || body != n.Body {
		t.Fatalf("without the module: %q, %q", title, body)
	}
	if !ns[0].System || ns[0].Source != SourceSystem {
		t.Fatalf("said by messages, not the system's: %v %q", ns[0].System, ns[0].Source)
	}
	// a user's: their text, never translated, not the system's
	sender := uuid.New()
	mine := n
	mine.Source = SourceUser(sender)
	if err := Notify(db, mine, u); err != nil {
		t.Fatal(err)
	}
	var got Notification
	db.Where("source = ?", SourceUser(sender)).First(&got)
	if title, _ := got.Text(en); got.System || title != n.Title {
		t.Fatalf("a user's: system %v, %q", got.System, title)
	}
	// with no source, the system's
	if err := Notify(db, Notification{Title: "S"}, u); err != nil {
		t.Fatal(err)
	}
	got = Notification{}
	db.Where("title = ?", "S").First(&got)
	if !got.System || got.Source != SourceSystem {
		t.Fatalf("no source: %v %q", got.System, got.Source)
	}
	plain := Notification{Title: "T", Body: "B"}
	if title, body := plain.Text(en); title != "T" || body != "B" {
		t.Fatalf("plain: %q, %q", title, body)
	}
}

// The periods up to now begin at the start of a day: today, the last 7 days
// (today among them).
func TestPeriodSince(t *testing.T) {
	now := time.Date(2026, 11, 9, 15, 30, 0, 0, time.UTC)
	for val, want := range map[string]time.Time{
		"today": time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC),
		"7d":    time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC),
		"30d":   time.Date(2026, 10, 11, 0, 0, 0, 0, time.UTC),
	} {
		if got, ok := periodSince(val, now); !ok || !got.Equal(want) {
			t.Errorf("%s: %v %v", val, got, ok)
		}
	}
	if _, ok := periodSince("x", now); ok {
		t.Error("an unknown period")
	}
}
