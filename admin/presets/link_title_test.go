package presets

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"
	"time"
)

// The title of a page: of its <title>, its spaces one, its entities read, of
// its charset; none past the head, of no HTML.
func TestTitleOf(t *testing.T) {
	for in, want := range map[string]string{
		"<html><head><title>\n  Plan  Ahead &amp; Co \n</title></head></html>": "Plan Ahead & Co",
		"<title>Só isto</title><title>outro</title>":                           "Só isto",
		"<html><body><title>no body</title></body></html>":                     "",
		"no html at all": "",
	} {
		if got := titleOf(strings.NewReader(in)); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
	if got := cleanTitle(strings.Repeat("a", 400)); len([]rune(got)) != linkTitleMaxRunes+1 {
		t.Errorf("a long title: %d", len([]rune(got)))
	}
}

// The addresses read: of http and https, of a host, no user; of the internet
// only — not of the machine, of its network, of the cloud's metadata.
func TestLinkTitleAddresses(t *testing.T) {
	for raw, ok := range map[string]bool{
		"https://example.com/a": true, "http://example.com": true, "ftp://example.com": false,
		"/a/b": false, "javascript:alert(1)": false, "https://u:p@example.com": false, "https://": false,
	} {
		if (linkTitleURL(raw) != nil) != ok {
			t.Errorf("%q: %v", raw, !ok)
		}
	}
	for ip, ok := range map[string]bool{
		"93.184.216.34": true, "2606:2800:220:1::248": true,
		"127.0.0.1": false, "::1": false, "10.1.2.3": false, "192.168.0.1": false, "172.16.0.1": false,
		"169.254.169.254": false, "100.64.0.1": false, "0.0.0.0": false, "fc00::1": false, "fe80::1": false,
		"::ffff:127.0.0.1": false, "224.0.0.1": false,
	} {
		if publicAddr(netip.MustParseAddr(ip)) != ok {
			t.Errorf("%s: %v", ip, !ok)
		}
	}
}

// A page read: its title; an address of the machine refused (as any of the
// network); a page that does not answer in time, none — once, no second try.
func TestLinkTitle(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch r.URL.Path {
		case "/slow":
			time.Sleep(LinkTitleTimeout + 500*time.Millisecond)
		case "/json":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"title":"no"}`)
			return
		case "/latin1":
			w.Header().Set("Content-Type", "text/html; charset=iso-8859-1")
			w.Write([]byte("<title>Ol\xe1 mundo</title>"))
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<html><head><title>Contato</title></head><body></body></html>")
	}))
	defer srv.Close()

	// the machine's: refused, never asked
	if got := LinkTitle(context.Background(), srv.URL+"/"); got != "" || calls != 0 {
		t.Fatalf("the machine's read: %q, %d", got, calls)
	}

	linkTitleAllowed = func(netip.Addr) bool { return true }
	t.Cleanup(func() { linkTitleAllowed = publicAddr })
	if got := LinkTitle(context.Background(), srv.URL+"/"); got != "Contato" {
		t.Errorf("the title: %q", got)
	}
	if got := LinkTitle(context.Background(), srv.URL+"/latin1"); got != "Olá mundo" {
		t.Errorf("of its charset: %q", got)
	}
	if got := LinkTitle(context.Background(), srv.URL+"/json"); got != "" {
		t.Errorf("of no HTML: %q", got)
	}
	calls = 0
	start := time.Now()
	if got := LinkTitle(context.Background(), srv.URL+"/slow"); got != "" {
		t.Errorf("slow: %q", got)
	}
	if d := time.Since(start); d > LinkTitleTimeout+time.Second {
		t.Errorf("waited %v", d)
	}
	if calls != 1 {
		t.Errorf("asked %d times, want once", calls)
	}
}

// The card of another site's page: what it says of itself to be shared first
// (Open Graph, then Twitter's), else its title and description; its cover of
// https, of its address; kept once read, read again old.
func TestLinkPreview(t *testing.T) {
	pages := map[string]string{
		"/og": `<head><title>Plain</title><meta name="description" content="Plain desc">
<meta property="og:title" content="OG title"><meta property="og:description" content="OG desc">
<meta property="og:image" content="https://cdn.example/i.png"><meta name="twitter:title" content="TW title"></head>`,
		"/tw": `<head><title>Plain</title><meta name="twitter:title" content="TW title">
<meta name="twitter:description" content="TW desc"><meta name="twitter:image" content="/rel.png"></head>`,
		"/plain": `<head><title> Plain </title><meta name="description" content="Plain desc"></head>`,
		"/none":  `<head></head><body>no title</body>`,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, pages[r.URL.Path])
	}))
	defer srv.Close()
	linkTitleAllowed = func(netip.Addr) bool { return true }
	t.Cleanup(func() { linkTitleAllowed = publicAddr })

	c, err := LinkPreview(context.Background(), srv.URL+"/og")
	if err != nil || c.Title != "OG title" || c.Summary != "OG desc" || c.Cover != "https://cdn.example/i.png" {
		t.Errorf("og: %+v %v", c, err)
	}
	// twitter's; a cover of http (the test's server) left out
	c, err = LinkPreview(context.Background(), srv.URL+"/tw")
	if err != nil || c.Title != "TW title" || c.Summary != "TW desc" || c.Cover != "" {
		t.Errorf("twitter: %+v %v", c, err)
	}
	c, err = LinkPreview(context.Background(), srv.URL+"/plain")
	if err != nil || c.Title != "Plain" || c.Summary != "Plain desc" || c.Href != srv.URL+"/plain" {
		t.Errorf("plain: %+v %v", c, err)
	}
	if _, err = LinkPreview(context.Background(), srv.URL+"/none"); err != ErrLinkPreview {
		t.Errorf("no title: %v", err)
	}

	// kept: the store's, not read again; old, read again
	store := NewMemoryExternalLinkCards(time.Hour)
	ctx := context.Background()
	_ = store.SaveExternalLinkCard(ctx, "https://x.test/a", &LinkCard{Href: "https://x.test/a", Title: "Kept"})
	if got, stale := store.ExternalLinkCard(ctx, "https://x.test/a"); got == nil || got.Title != "Kept" || stale {
		t.Errorf("kept: %+v %v", got, stale)
	}
	old := NewMemoryExternalLinkCards(-time.Second)
	_ = old.SaveExternalLinkCard(ctx, "https://x.test/a", &LinkCard{Title: "Old"})
	if _, stale := old.ExternalLinkCard(ctx, "https://x.test/a"); !stale {
		t.Error("an old card not stale")
	}
	if !ExternalTarget("https://example.com/a") || ExternalTarget("/a") || ExternalTarget("admin:posts:") {
		t.Error("ExternalTarget")
	}
}
