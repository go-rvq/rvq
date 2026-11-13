package perm

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// The sets of permissions: a name for permissions given to a resource
// together — reading, editing, publishing —, written at the end of a line,
// "%name": "admin:posts:%readers", "deny admin:users:%editors". A line keeps
// the name; what it gives comes of the set each time a stack is made
// (Stack.Add, by the expander set by SetSetExpander — presets', by the tree
// of the permissions of the admin and the sets registered): a set changed
// changes what every line of it gives; a set that is no more, or has nothing
// for its resource, gives nothing — the line is there, ignored (Line.Ignored).

// SetPrefix is the beginning of the segment of a set: "%readers".
const SetPrefix = "%"

// setNameRe is a name of a set.
var setNameRe = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// ValidSetName says name is a name a set may have: lower case letters,
// digits and "_", a letter first.
func ValidSetName(name string) bool { return setNameRe.MatchString(name) }

// SetRef is the resource and the set of a pattern of one line that ends by
// a set ("admin:posts:%readers": "admin:posts:", "readers"); ok false for any
// other.
func SetRef(p string) (resource, name string, ok bool) {
	segs := splitTop(p)
	if len(segs) < 2 {
		return "", "", false
	}
	last := segs[len(segs)-1]
	if !strings.HasPrefix(last, SetPrefix) || strings.HasSuffix(last, ":") {
		return "", "", false
	}
	name = strings.TrimPrefix(last, SetPrefix)
	if !ValidSetName(name) {
		return "", "", false
	}
	return strings.Join(segs[:len(segs)-1], ""), name, true
}

// SetExpander is the lines a set gives to the resource (its pattern, with
// its ":" at its end), each with its prefix when it has one ("deny …"); ok
// false when it gives nothing there (no such set, nothing for that kind of
// resource).
type SetExpander func(resource, name string) (lines []string, ok bool)

var (
	setExpanderMu sync.RWMutex
	setExpander   SetExpander
)

// SetSetExpander sets the expander of the sets of the process (presets'
// Builder, when it is built); nil, none: every line of a set is ignored.
func SetSetExpander(f SetExpander) {
	setExpanderMu.Lock()
	defer setExpanderMu.Unlock()
	setExpander = f
}

// ExpandSet is what the set name gives to resource (SetSetExpander).
func ExpandSet(resource, name string) ([]string, bool) {
	setExpanderMu.RLock()
	f := setExpander
	setExpanderMu.RUnlock()
	if f == nil {
		return nil, false
	}
	return f(resource, name)
}

// SetLineKind is the kind of a line a set gives under a line of kind outer:
// its own when it has one (a set's "deny" denies whatever the line says),
// else the line's.
func SetLineKind(outer, inner LineKind) LineKind {
	if inner != LineGrant {
		return inner
	}
	return outer
}

// SetBlocks are the blocks of a definition of a set by kind of resource —
// "{model: {@list, <*>:@get}, singleton: @get}": "model" its lines
// "{@list, <*>:@get}", "singleton" "@get" —, each checked (a pattern
// relative to its resource: no scope). A definition that cannot be read,
// its error.
func SetBlocks(def string) (map[string]string, error) {
	def = strings.TrimSpace(def)
	if !strings.HasPrefix(def, "{") || !strings.HasSuffix(def, "}") {
		return nil, fmt.Errorf("a set is {kind: permissions, …}: %q", def)
	}
	out := map[string]string{}
	for _, item := range SplitResources(def[1 : len(def)-1]) {
		i := strings.Index(item, ":")
		if i <= 0 {
			return nil, fmt.Errorf("%q: a kind of resource and its permissions, kind: …", item)
		}
		kind, body := strings.TrimSpace(item[:i]), strings.TrimSpace(item[i+1:])
		if !ValidSetName(kind) {
			return nil, fmt.Errorf("%q: not a kind of resource", kind)
		}
		if _, dup := out[kind]; dup {
			return nil, fmt.Errorf("%q: the kind twice", kind)
		}
		if err := CheckSetBlock(body); err != nil {
			return nil, fmt.Errorf("%s: %w", kind, err)
		}
		out[kind] = body
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("a set of no kind")
	}
	return out, nil
}

// CheckSetBlock says what is wrong with the permissions of a block of a set
// (relative to its resource: "{@list, <*>:@get}"): nil when nothing is.
func CheckSetBlock(body string) error {
	lines, err := LinesOf(body)
	if err != nil {
		return err
	}
	if len(lines) == 0 {
		return fmt.Errorf("no permission")
	}
	for _, l := range lines {
		k, p := ParseLine(l)
		if k == LineUnset {
			return ErrUnset
		}
		if p == "" {
			return fmt.Errorf("%q: no permission", l)
		}
	}
	return nil
}

// SetBlockLines are the lines of a block of a set given to resource (its
// pattern, its ":" at its end): each of the block's after it, with its
// prefix; a set of the block ("%authors": of the same resource) by sub, the
// lines of the other set, nil for none.
func SetBlockLines(resource, body string, sub func(name string) []string) ([]string, error) {
	lines, err := LinesOf(body)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, l := range lines {
		k, p := ParseLine(l)
		if strings.HasPrefix(p, SetPrefix) && ValidSetName(strings.TrimPrefix(p, SetPrefix)) {
			if sub == nil {
				continue
			}
			for _, sl := range sub(strings.TrimPrefix(p, SetPrefix)) {
				sk, sp := ParseLine(sl)
				out = append(out, SetLineKind(k, sk).Prefix()+sp)
			}
			continue
		}
		out = append(out, k.Prefix()+resource+p)
	}
	return out, nil
}
