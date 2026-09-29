// Package osenv reads the configuration from the environment, with the API of
// github.com/theplant/osenv — Get, GetBool, GetInt64: a key, what it is for and
// its default — but quiet: that one writes every variable to the standard
// output as it is read (at the packages' init, before main), which a command
// like `app version` cannot keep out of what it prints.
//
// What was read is recorded (Vars), for a command to list. OSENV_PRINT=true
// writes each one, as it is read, to the standard error — as the other did.
package osenv

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// Var is a variable read: its key, what it is for, its default and where it
// was read.
type Var struct {
	Key, Message, Default, At string
}

var (
	mu    sync.Mutex
	vars  []Var
	print = func() bool { v, _ := strconv.ParseBool(os.Getenv("OSENV_PRINT")); return v }()
)

// Vars are the variables read so far, in the order they were read.
func Vars() []Var {
	mu.Lock()
	defer mu.Unlock()
	return append([]Var(nil), vars...)
}

// String is the variable as the other package printed it.
func (v Var) String() string {
	def := ""
	if v.Default != "" {
		def = fmt.Sprintf(" (default: %s)", v.Default)
	}
	return fmt.Sprintf("%s: %s,%s, at: %s", v.Key, v.Message, def, v.At)
}

// Get is the value of key, or defaultValue when it is empty.
func Get(key string, msg string, defaultValue string) string {
	return get(key, msg, defaultValue)
}

// GetBool is the value of key as a bool, or defaultValue when it is empty.
func GetBool(key string, msg string, defaultValue bool) (r bool) {
	sv := get(key, msg, fmt.Sprint(defaultValue))
	if sv == "" {
		return defaultValue
	}
	r, _ = strconv.ParseBool(sv)
	return
}

// GetInt64 is the value of key as an int64, or defaultValue when it is empty.
func GetInt64(key string, msg string, defaultValue int64) (r int64) {
	sv := get(key, msg, fmt.Sprint(defaultValue))
	if sv == "" {
		return defaultValue
	}
	r, _ = strconv.ParseInt(sv, 10, 64)
	return
}

func get(key, msg, defaultValue string) string {
	at := ""
	if pc, _, _, ok := runtime.Caller(2); ok {
		at = runtime.FuncForPC(pc).Name()
	}
	v := Var{Key: key, Message: msg, Default: defaultValue, At: at}
	mu.Lock()
	vars = append(vars, v)
	mu.Unlock()
	if print {
		fmt.Fprintln(os.Stderr, v)
	}

	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
