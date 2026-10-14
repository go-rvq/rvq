// Package security is what an installation allows of the ways past its
// login: the `security:` of the app.yml (admin/app/env), each off unless it
// says so.
//
//	security:
//	  secure_key_login: true
package security

import "sync"

// Config is the `security:` of the app.yml.
type Config struct {
	// SecureKeyLogin: the secure key (x/login, `.secure_key`) logs in to the
	// admin as its administrator — the initial account —, no login page:
	// x/login's secure key login. Off, the key only passes the form
	// protection of the login page.
	SecureKeyLogin bool `yaml:"secure_key_login"`
}

var (
	mu  sync.RWMutex
	cfg Config
)

// Load sets the configuration: what admin/app/env reads of the app.yml.
func Load(c Config) {
	mu.Lock()
	defer mu.Unlock()
	cfg = c
}

// Get is the configuration.
func Get() Config {
	mu.RLock()
	defer mu.RUnlock()
	return cfg
}

// SecureKeyLogin says the secure key logs in as the administrator
// (Config.SecureKeyLogin).
func SecureKeyLogin() bool { return Get().SecureKeyLogin }
