package portal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

// The server of the portals' content (Server): the page names no address —
// only a token (Builder.Source) —, and the runtime of the site (the bundle of
// js/portal) asks for the content, and for what the content does (an
// action: send a form…), at one fixed path, the token and the action in its
// headers.
const (
	// DefaultPath is where the content of a portal is asked for.
	DefaultPath = "/_rvq/portal"
	// HeaderToken carries the token of the portal (Builder.Source).
	HeaderToken = "X-Rvq-Portal"
	// HeaderAction carries the action asked for: none, the content.
	HeaderAction = "X-Rvq-Portal-Action"
	// DefaultTTL is how long a token is good for.
	DefaultTTL = 24 * time.Hour
)

var (
	// ErrBadToken is a token the server did not sign, or not a token.
	ErrBadToken = errors.New("portal: a token not valid")
	// ErrExpired is a token past its time: the page is older than it.
	ErrExpired = errors.New("portal: the token expired")
)

// Key is what a token names: a content of a kind (a Source's), its id, in a
// language, until a moment.
type Key struct {
	Kind    string `json:"k"`
	ID      string `json:"i"`
	Locale  string `json:"l,omitempty"`
	Expires int64  `json:"e"`
}

// Signer signs a value and checks it (presets.FormSigner is one).
type Signer interface {
	Sign(value string) string
	Unsign(signed string) (value string, err error)
}

// Content is the content of a portal: its template (a Vue template, compiled
// in the browser) and the variables of it (the portal's content options'
// data); the scripts it needs (loaded once for the page, before it is
// shown), and its kind — what makes it in the browser (RvqPortal.kinds,
// registered by those scripts): its behavior, its template made anew —, none
// for the template as it is.
type Content struct {
	Template string   `json:"template"`
	Data     any      `json:"data,omitempty"`
	Scripts  []string `json:"scripts,omitempty"`
	Kind     string   `json:"kind,omitempty"`
}

// Source answers the content of the portals of its kind.
type Source interface {
	Content(r *http.Request, key Key) (*Content, error)
}

// ActionSource is a Source whose content does things — sends a form, saves a
// draft —: each an action, at the same path, the request as the content's
// script made it (its body untouched: a multipart of files too).
type ActionSource interface {
	Source
	ServeAction(w http.ResponseWriter, r *http.Request, key Key, action string)
}

// SourceFunc is a Source of a function.
type SourceFunc func(r *http.Request, key Key) (*Content, error)

func (f SourceFunc) Content(r *http.Request, key Key) (*Content, error) { return f(r, key) }

// Server serves the content of the portals and their actions at its Path:
// a POST with the token of the portal (HeaderToken) and the action
// (HeaderAction; none, the content).
type Server struct {
	// Path is where it is mounted (DefaultPath).
	Path    string
	signer  Signer
	ttl     time.Duration
	mu      sync.RWMutex
	sources map[string]Source

	// the bundle (script)
	scriptOnce    sync.Once
	scriptBody    []byte
	scriptVersion string
	scriptErr     error
}

// NewServer is a server whose tokens signer signs.
func NewServer(signer Signer) *Server {
	return &Server{Path: DefaultPath, signer: signer, ttl: DefaultTTL, sources: map[string]Source{}}
}

// TTL sets how long a token is good for.
func (s *Server) TTL(d time.Duration) *Server {
	s.ttl = d
	return s
}

// Register makes src the source of the portals of kind.
func (s *Server) Register(kind string, src Source) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sources[kind] = src
	return s
}

func (s *Server) source(kind string) Source {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sources[kind]
}

// Token is the token of the content id of kind, in locale: opaque to the
// page, signed, good for the server's TTL.
func (s *Server) Token(kind, id, locale string) string {
	b, _ := json.Marshal(Key{Kind: kind, ID: id, Locale: locale, Expires: time.Now().Add(s.ttl).Unix()})
	return s.signer.Sign(base64.RawURLEncoding.EncodeToString(b))
}

// Key is what token names: ErrBadToken when the server did not sign it,
// ErrExpired when it is past its time.
func (s *Server) Key(token string) (k Key, err error) {
	v, err := s.signer.Unsign(strings.TrimSpace(token))
	if err != nil {
		return k, ErrBadToken
	}
	b, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil || json.Unmarshal(b, &k) != nil || k.Kind == "" {
		return Key{}, ErrBadToken
	}
	if time.Now().Unix() > k.Expires {
		return k, ErrExpired
	}
	return k, nil
}

// errorReply is an error the runtime understands: expired (the page is
// older than its token: to load again), or not (nothing to do).
func errorReply(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": code})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		errorReply(w, http.StatusMethodNotAllowed, "method")
		return
	}
	key, err := s.Key(r.Header.Get(HeaderToken))
	switch {
	case errors.Is(err, ErrExpired):
		errorReply(w, http.StatusGone, "expired")
		return
	case err != nil:
		errorReply(w, http.StatusForbidden, "token")
		return
	}
	src := s.source(key.Kind)
	if src == nil {
		errorReply(w, http.StatusNotFound, "kind")
		return
	}
	if action := r.Header.Get(HeaderAction); action != "" {
		as, ok := src.(ActionSource)
		if !ok {
			errorReply(w, http.StatusNotFound, "action")
			return
		}
		as.ServeAction(w, r, key, action)
		return
	}
	c, err := src.Content(r, key)
	if err != nil || c == nil {
		errorReply(w, http.StatusNotFound, "content")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(c)
}
