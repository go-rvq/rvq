// Package integration provides helpers for exercising preset apps as real HTTP
// servers in integration tests — notably the bun/TypeScript tests under
// js/integration_tests, which boot a Go app (in-memory SQLite) and drive its
// event funcs over HTTP.
package integration

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

// ServeEnv listens on 127.0.0.1:$PORT (PORT unset or "0" → an ephemeral port),
// prints "LISTENING http://<addr>" to stdout once bound, and serves handler.
// It blocks until the server stops. The printed line lets a test runner wait for
// readiness and learn the bound address.
func ServeEnv(handler http.Handler) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "0"
	}
	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	fmt.Printf("LISTENING http://%s\n", ln.Addr().String())
	_ = os.Stdout.Sync()
	return http.Serve(ln, handler)
}

// ServeEnvMain runs ServeEnv and exits the process non-zero on error, so a
// package main entry point is a one-liner:
//
//	func main() { integration.ServeEnvMain(handler) }
func ServeEnvMain(handler http.Handler) {
	if err := ServeEnv(handler); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
