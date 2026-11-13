package presets

import (
	"context"
	"sync"
	"time"

	"github.com/go-rvq/rvq/web"
)

// The cards of pages of other sites (an external target: an address of http
// or https): read once (LinkPreview) by the admin's LinkCardsEvent — of who
// is in the admin, a few a request —, kept (ExternalLinkCardStore); what
// shows the texts to anybody else (the site's own cards) reads them of the
// store, never of the internet.

// externalFetchesMax is how many pages of other sites a LinkCardsEvent
// reads; the others, of no card yet, are read at a later one.
const externalFetchesMax = 4

// ExternalLinkCardStore keeps the cards of pages of other sites, by their
// address.
type ExternalLinkCardStore interface {
	// ExternalLinkCard is the card kept of href, and whether it is to be read
	// again (old); nil, none.
	ExternalLinkCard(ctx context.Context, href string) (card *LinkCard, stale bool)
	// SaveExternalLinkCard keeps the card of href.
	SaveExternalLinkCard(ctx context.Context, href string, card *LinkCard) error
}

// ExternalLinkCards sets where the cards of pages of other sites are kept
// (default: in memory, a day, MemoryExternalLinkCards).
func (b *Builder) ExternalLinkCards(s ExternalLinkCardStore) *Builder {
	b.externalCards = s
	return b
}

// GetExternalLinkCards is where the cards of pages of other sites are kept.
func (b *Builder) GetExternalLinkCards() ExternalLinkCardStore {
	if b.externalCards == nil {
		b.externalCards = NewMemoryExternalLinkCards(24 * time.Hour)
	}
	return b.externalCards
}

// externalLinkCard is the card of the page of another site href: kept,
// else — read allowed — read and kept; nil, none.
func (b *Builder) externalLinkCard(ctx *web.EventContext, href string, read bool) *LinkCard {
	store := b.GetExternalLinkCards()
	c, stale := store.ExternalLinkCard(ctx.R.Context(), href)
	if c != nil && !stale || !read {
		return c
	}
	got, err := LinkPreview(ctx.R.Context(), href)
	if err != nil {
		// read before: kept as it was
		return c
	}
	_ = store.SaveExternalLinkCard(ctx.R.Context(), href, got)
	return got
}

// ExternalTarget says target is an address of another site a card may link
// (http or https).
func ExternalTarget(target string) bool { return linkTitleURL(target) != nil }

// MemoryExternalLinkCards keeps the cards of pages of other sites in memory,
// each for ttl.
type MemoryExternalLinkCards struct {
	ttl   time.Duration
	mu    sync.Mutex
	cards map[string]memoryExternalCard
}

type memoryExternalCard struct {
	card *LinkCard
	at   time.Time
}

// NewMemoryExternalLinkCards is a store in memory of cards kept for ttl.
func NewMemoryExternalLinkCards(ttl time.Duration) *MemoryExternalLinkCards {
	return &MemoryExternalLinkCards{ttl: ttl, cards: map[string]memoryExternalCard{}}
}

// ExternalLinkCard is ExternalLinkCardStore's.
func (m *MemoryExternalLinkCards) ExternalLinkCard(_ context.Context, href string) (*LinkCard, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.cards[href]
	if !ok {
		return nil, false
	}
	return c.card, time.Since(c.at) > m.ttl
}

// SaveExternalLinkCard is ExternalLinkCardStore's.
func (m *MemoryExternalLinkCards) SaveExternalLinkCard(_ context.Context, href string, card *LinkCard) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cards[href] = memoryExternalCard{card, time.Now()}
	return nil
}
