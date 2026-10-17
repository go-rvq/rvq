// Package featurescli is the command `features` of an application: what its
// configuration turns on, of what its binary was compiled with (x/features).
//
//	app features      the features on
//	app features all  every feature compiled: asked in app.yml, on, what it depends on
//
// A feature the binary has not is none: never shown. It needs no server:
// RunIfAsked runs it before the application boots, the `features:` of the
// app.yml read already (admin/app/env).
package featurescli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/go-rvq/rvq/cli"
	"github.com/go-rvq/rvq/x/features"
)

// Command is `features`, the features on, and its sub-command `all`.
func Command() *cli.Command {
	return (&cli.Command{
		Name:        "features",
		Description: "the features on (`all`: every feature compiled)",
		Run:         func(ctx *cli.CommandContext) error { return PrintOn(ctx.Out) },
	}).Sub(&cli.Command{
		Name:        "all",
		Description: "every feature compiled: asked in app.yml, on, what it depends on",
		Run:         func(ctx *cli.CommandContext) error { return PrintAll(ctx.Out) },
	})
}

// RunIfAsked runs `features …` when args (os.Args) ask for it: handled, and
// the exit code.
func RunIfAsked(args []string) (handled bool, code int) {
	if len(args) < 2 || args[1] != "features" {
		return false, 0
	}
	root := (&cli.Command{Name: filepath.Base(args[0])}).Sub(Command())
	ctx, err := root.Parse(&cli.CommandContext{InputArgs: args[1:]})
	if err == nil {
		err = ctx.Run()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		return true, 1
	}
	return true, 0
}

// PrintOn writes the features on: name and label.
func PrintOn(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, s := range features.All() {
		if s.On {
			fmt.Fprintf(tw, "%s\t%s\n", s.Name, s.Label)
		}
	}
	return tw.Flush()
}

// PrintAll writes every feature compiled: asked, on — or why not —, what it
// depends on; then the keys of app.yml of no feature.
func PrintAll(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tASKED\tON\tDEPENDS ON\tLABEL")
	for _, s := range features.All() {
		var deps []string
		for _, d := range s.Dependencies() {
			deps = append(deps, d.Name)
		}
		on := yes(s.On)
		if !s.On && s.Asked && len(s.Missing) > 0 {
			on = "no (needs " + strings.Join(s.Missing, ", ") + ")"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", s.Name, yes(s.Asked), on, dash(strings.Join(deps, ", ")), s.Label)
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	for _, k := range features.Unknown() {
		fmt.Fprintf(w, "app.yml: features.%s: no such feature\n", k)
	}
	return nil
}

func yes(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func dash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}
