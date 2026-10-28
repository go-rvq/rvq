package mail_sender

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// The subject of a message is made of parts (MessageBuilder.Topic, Ref,
// Reply, Subject), put together as the settings say (MailSender.SubjectRules):
//
//	Re: [PAC] [ISSUE - Chamados] [#12] Title
//	│   │     │                  │     └ the subject (Subject)
//	│   │     │                  └ the record (Ref) — SubjectRef
//	│   │     └ the thing it is of (Topic) — SubjectTopic
//	│   └ the site's name (MailSender.SubjectPrefix, else SiteName) — SubjectSite
//	└ an answer of a conversation (Reply): always first
//
// Which parts go is told by the family of the message (the Kind before its
// dot: "issue.reply" is "issue") and to whom it goes (External: outside the
// installation; else the team).

// SubjectParts are the parts of a subject a rule puts in it, as bits.
type SubjectParts uint8

const (
	// SubjectSite: the site's name
	SubjectSite SubjectParts = 1 << iota
	// SubjectTopic: the thing the message is of (Topic)
	SubjectTopic
	// SubjectRef: the record (Ref)
	SubjectRef
)

// SubjectAll are all the parts.
const SubjectAll = SubjectSite | SubjectTopic | SubjectRef

// Has says p has every part of q.
func (p SubjectParts) Has(q SubjectParts) bool { return p&q == q }

// SubjectFamilyOther is the family of the rule of the messages of a family
// with none of its own.
const SubjectFamilyOther = "*"

// SubjectRule are the parts of the subjects of a family of messages: of
// those to the team (Team) and of those outside the installation (External).
type SubjectRule struct {
	Team     SubjectParts `json:"team"`
	External SubjectParts `json:"external"`
}

// DefaultSubjectRule is the rule of a family no rule names: everything to
// the team; to outside, the site's name and the record.
var DefaultSubjectRule = SubjectRule{Team: SubjectAll, External: SubjectSite | SubjectRef}

// SubjectRules are the rules by family (SubjectFamilyOther: the others),
// kept as JSON.
type SubjectRules map[string]SubjectRule

// Of is the rule of family: its own, else the others', else the default.
func (r SubjectRules) Of(family string) SubjectRule {
	if rule, ok := r[family]; ok {
		return rule
	}
	if rule, ok := r[SubjectFamilyOther]; ok {
		return rule
	}
	return DefaultSubjectRule
}

// Scan reads them from their JSON.
func (r *SubjectRules) Scan(src any) error {
	*r = SubjectRules{}
	var b []byte
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("mail sender: subject rules of %T", src)
	}
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, r)
}

// Value writes them as JSON.
func (r SubjectRules) Value() (driver.Value, error) {
	if r == nil {
		return "{}", nil
	}
	b, err := json.Marshal(r)
	return string(b), err
}

// SubjectFamily is a family of messages an application sends (its Kinds'
// prefix): its name, and what the admin calls it.
type SubjectFamily struct {
	Name  string
	Label func(ctx context.Context) string
}

var (
	subjectFamiliesMu sync.RWMutex
	subjectFamilies   = map[string]SubjectFamily{}
)

// RegisterSubjectFamily registers the family name of the messages an
// application sends — its rule in the settings' table —, label what the
// admin calls it.
func RegisterSubjectFamily(name string, label func(ctx context.Context) string) {
	subjectFamiliesMu.Lock()
	defer subjectFamiliesMu.Unlock()
	subjectFamilies[name] = SubjectFamily{Name: name, Label: label}
}

// SubjectFamilies are the families registered, by name.
func SubjectFamilies() []SubjectFamily {
	subjectFamiliesMu.RLock()
	defer subjectFamiliesMu.RUnlock()
	out := make([]SubjectFamily, 0, len(subjectFamilies))
	for _, f := range subjectFamilies {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// SiteName is the site's name, when the settings name none
// (MailSender.SubjectPrefix): the application gives it.
var SiteName = func() string { return "" }

// subjectSite is the site's name of the subjects: the settings', else the
// application's (SiteName); its brackets out (they are the subject's).
func (c *MailSender) subjectSite() string {
	clean := func(s string) string { return strings.TrimSpace(strings.Trim(strings.TrimSpace(s), "[]")) }
	if n := clean(c.SubjectPrefix); n != "" {
		return n
	}
	return clean(SiteName())
}

// Family is the family of the message: its Kind before the dot.
func (b *MessageBuilder) Family() string {
	f, _, _ := strings.Cut(b.kind, ".")
	return f
}

// composeSubject is the subject of b with the parts of parts, the site's
// name site.
func (b *MessageBuilder) composeSubject(site string, parts SubjectParts) string {
	var s []string
	if b.reply {
		s = append(s, "Re:")
	}
	bracket := func(v string) {
		if v = strings.TrimSpace(v); v != "" {
			s = append(s, "["+v+"]")
		}
	}
	if parts.Has(SubjectSite) {
		bracket(site)
	}
	for _, l := range b.labels {
		bracket(l)
	}
	if parts.Has(SubjectTopic) {
		bracket(b.topic)
	}
	if parts.Has(SubjectRef) {
		bracket(b.ref)
	}
	if v := strings.TrimSpace(b.subject); v != "" {
		s = append(s, v)
	}
	return strings.Join(s, " ")
}

// composeWith makes the subject of b as the settings c say — its parts by its
// family and whether it goes outside —, once: then it goes as it is.
func (c *MailSender) composeWith(b *MessageBuilder) {
	if b.verbatim {
		b.idFooter = c.IDFooter
		return
	}
	rule := c.SubjectRules.Of(b.Family())
	parts := rule.Team
	if b.external {
		parts = rule.External
	}
	b.subject = b.composeSubject(c.subjectSite(), parts)
	b.labels = nil
	b.verbatim = true
	b.idFooter = c.IDFooter
}
