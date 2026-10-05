package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-rvq/rvq/cli"
)

// RunRemoteIfAsked runs `NAME http_api …` (args, as os.Args) when it goes to a
// remote server — --remote, or $HTTP_API_REMOTE —, before the application is
// built: no database, no boot, only the request. Called first in main;
// handled false when args are not such a run.
func RunRemoteIfAsked(args []string) (handled bool, code int) {
	if len(args) < 2 || args[1] != "http_api" || !remoteAsked(args[2:]) {
		return false, 0
	}
	ctx, err := (&cli.Command{Name: args[0]}).Sub(HttpApiCommand(Config{})).Parse(&cli.CommandContext{
		InputArgs: args[1:],
	})
	if err == nil {
		err = ctx.Run()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		return true, 1
	}
	return true, 0
}

// remoteAsked says whether the flags of http_api (args) send it to a remote
// server.
func remoteAsked(args []string) bool {
	for _, a := range args {
		if a == "--" {
			break
		}
		name, value, hasValue := strings.Cut(strings.TrimLeft(a, "-"), "=")
		if strings.HasPrefix(a, "-") && name == "remote" {
			return !hasValue || value != ""
		}
	}
	return os.Getenv("HTTP_API_REMOTE") != ""
}
