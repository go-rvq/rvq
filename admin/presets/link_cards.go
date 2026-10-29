package presets

import (
	"strings"

	"github.com/go-rvq/rvq/web"
)

// The link cards of a text (x/markup: a page of the site, a record of the
// admin, linked on a line of its own): what each is — its title, its
// address, its cover and its summary —, answered by LinkCardsEvent to the
// script that draws them (rvq mdkit link-cards), for who may see it; a
// target not known, or not theirs to see, has none (its link stays).

// LinkCardsEvent answers the cards of the targets posted (target, many): a
// map of target to LinkCard.
const LinkCardsEvent = "presets_LinkCards"

// linkCardsMax is how many targets an event answers.
const linkCardsMax = 50

// LinkCard is what a card shows.
type LinkCard struct {
	Href    string `json:"href"`
	Title   string `json:"title"`
	Cover   string `json:"cover,omitempty"`
	Summary string `json:"summary,omitempty"`
}

// LinkCardCoverer is a record with a cover for its card (an address).
type LinkCardCoverer interface {
	LinkCardCover() string
}

// LinkCardSummarizer is a record with a summary for its card.
type LinkCardSummarizer interface {
	LinkCardSummary() string
}

// LinkCardResolver is the card of target for the request; nil, it does not
// know it (or the user may not see it).
type LinkCardResolver func(ctx *web.EventContext, target string) *LinkCard

// LinkCardResolver adds a resolver of the application's targets (a path of
// its site), asked before the admin's record paths.
func (b *Builder) LinkCardResolver(f LinkCardResolver) *Builder {
	b.linkCardResolvers = append(b.linkCardResolvers, f)
	return b
}

// LinkCard is the card of target: the application's resolvers', else — a
// record path (admin:…:) — the record's, the listing's, by who may see it.
func (b *Builder) LinkCard(ctx *web.EventContext, target string) *LinkCard {
	for _, f := range b.linkCardResolvers {
		if c := f(ctx, target); c != nil {
			return c
		}
	}
	if strings.HasPrefix(target, PermModule+":") {
		return b.recordLinkCard(ctx, target)
	}
	return nil
}

// recordLinkCard is the card of the record path p: a record — its title,
// its cover and summary when it has them (LinkCardCoverer,
// LinkCardSummarizer) —, a listing, a singleton; who may not see it, none.
func (b *Builder) recordLinkCard(ctx *web.EventContext, p string) *LinkCard {
	t, err := b.ParseRecordPath(p)
	if err != nil || len(t.Models) == 0 || t.Action != nil {
		return nil
	}
	last := t.Models[len(t.Models)-1]
	parents := make(IDSlice, len(t.Models)-1)
	for n, s := range t.Models[:len(t.Models)-1] {
		parents[n] = s.ID
	}
	mb := last.Model
	perms := mb.Info().Permissioner()
	if t.Listing || mb.singleton {
		if perms.Lister(ctx.R, parents...).Denied() {
			return nil
		}
		return &LinkCard{Href: t.Href(b), Title: mb.TTitlePlural(ctx.Context())}
	}
	if perms.Reader(ctx.R, last.ID, parents...).Denied() {
		return nil
	}
	obj := mb.NewModel()
	if err := mb.Fetcher(obj, last.ID, ctx); err != nil {
		return nil
	}
	c := &LinkCard{Href: t.Href(b), Title: mb.RecordTitle(obj, ctx)}
	if v, ok := obj.(LinkCardCoverer); ok {
		c.Cover = v.LinkCardCover()
	}
	if v, ok := obj.(LinkCardSummarizer); ok {
		c.Summary = v.LinkCardSummary()
	}
	return c
}

// linkCards is LinkCardsEvent.
func (b *Builder) linkCards(ctx *web.EventContext) (r web.EventResponse, err error) {
	out := map[string]*LinkCard{}
	r.Data = out
	_ = ctx.R.ParseMultipartForm(1 << 20)
	for _, target := range ctx.R.Form["target"] {
		if len(out) >= linkCardsMax {
			break
		}
		if _, done := out[target]; done || target == "" {
			continue
		}
		if c := b.LinkCard(ctx, target); c != nil {
			out[target] = c
		}
	}
	return
}

// LinkCardSearchEvent answers what a link card may link matching q — the
// application's (LinkCardSearcher: the pages of its site), then the records
// of the search of the top bar (by their record paths) —, in groups.
const LinkCardSearchEvent = "presets_LinkCardSearch"

// LinkCardOption is what a link card may link: its title, its target (a path
// of the site, a record path), where it is (sub).
type LinkCardOption struct {
	Title  string `json:"title"`
	Target string `json:"target"`
	Sub    string `json:"sub,omitempty"`
}

// LinkCardOptionGroup is a group of options: its title, or its kind (the
// picker's word of it: pages, posts).
type LinkCardOptionGroup struct {
	Kind  string           `json:"kind,omitempty"`
	Title string           `json:"title"`
	Icon  string           `json:"icon,omitempty"`
	Items []LinkCardOption `json:"items"`
}

// LinkCardSearcher is the application's options matching q.
type LinkCardSearcher func(ctx *web.EventContext, q string) []*LinkCardOptionGroup

// LinkCardSearcher adds a searcher of the application's options, before the
// admin's records.
func (b *Builder) LinkCardSearcher(f LinkCardSearcher) *Builder {
	b.linkCardSearchers = append(b.linkCardSearchers, f)
	return b
}

// linkCardSearch is LinkCardSearchEvent.
func (b *Builder) linkCardSearch(ctx *web.EventContext) (r web.EventResponse, err error) {
	q := strings.TrimSpace(ctx.R.FormValue("q"))
	groups := []*LinkCardOptionGroup{}
	r.Data = groups
	if len([]rune(q)) < 2 {
		return
	}
	for _, f := range b.linkCardSearchers {
		for _, g := range f(ctx, q) {
			if g != nil && len(g.Items) > 0 {
				groups = append(groups, g)
			}
		}
	}
	gs, err := b.globalSearch(ctx)
	if err != nil {
		return
	}
	found, _ := gs.Data.([]*globalSearchGroup)
	for _, g := range found {
		og := &LinkCardOptionGroup{Title: g.Title, Icon: g.Icon}
		for _, it := range g.Items {
			if it.Path != "" {
				og.Items = append(og.Items, LinkCardOption{Title: it.Title, Target: it.Path, Sub: g.Title})
			}
		}
		if len(og.Items) > 0 {
			groups = append(groups, og)
		}
	}
	r.Data = groups
	return
}
