package dashboard

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// A widget Live: a selector of its intervals in its header, reloading its
// portal; its footer with the seconds; its blocks kept at most half the
// shortest interval.
func TestLive(t *testing.T) {
	w := &Widget{ID: "sys/cpu", Group: "sys", Name: "cpu", Title: func(context.Context) string { return "CPU" }, Live: []time.Duration{2 * time.Second, 5 * time.Second, 10 * time.Second},
		Load: func(*web.EventContext, Period) ([]Block, error) { return []Block{&Note{Text: "x"}}, nil }}
	sel := h.MustString(liveSelector(Messages_en_US, w), context.Background())
	for _, want := range []string{`data-dashboard-live='sys/cpu'`, `data-dashboard-live-choice='2000'`, `data-dashboard-live-choice='5000'`,
		`data-dashboard-live-choice='10000'`, `every: 2000`, PortalName("sys/cpu"), "rvq-dash-live:", Messages_en_US.LiveEvery} {
		if !strings.Contains(sel, want) {
			t.Errorf("the selector has no %q:\n%s", want, sel)
		}
	}

	// on the dashboard and on a record's panel
	_, b0 := newTest()
	ctx0 := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	p0 := NewPeriod(Period7d, now, "", "")
	for name, c := range map[string]h.HTMLComponent{
		"card":  b0.card(ctx0, Messages_en_US, w, &Prefs{}, p0, false, true, true),
		"panel": b0.panelCard(ctx0, Messages_en_US, w, p0, false),
	} {
		s := h.MustString(c, context.Background())
		if !strings.Contains(s, `data-dashboard-live='sys/cpu'`) {
			t.Errorf("%s: no selector", name)
		}
		// the button of its page
		if !strings.Contains(s, `data-dashboard-page='sys/cpu'`) || !strings.Contains(s, WidgetPagePath+"?id=sys%2Fcpu") {
			t.Errorf("%s: no button of its page", name)
		}
	}

	clock := now
	bd, b := newTest()
	bd.Now = func() time.Time { return clock }
	n := 0
	w.Load = func(*web.EventContext, Period) ([]Block, error) { n++; return nil, nil }
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	p := NewPeriod(Period7d, now, "", "")
	b.blocks(ctx, w, p, false)
	clock = clock.Add(1500 * time.Millisecond)
	b.blocks(ctx, w, p, false)
	if n != 2 {
		t.Errorf("kept past half the shortest interval: loads %d", n)
	}
	ctx.R = ctx.R.WithContext(context.Background())
	body := h.MustString(b.Render(ctx, w, p, true), context.Background())
	if !strings.Contains(body, clock.Format("15:04:05")) {
		t.Errorf("the footer without the seconds: %s", body)
	}
}
