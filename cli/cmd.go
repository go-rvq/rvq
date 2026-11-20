package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-rvq/rvq/x/lifecycle"
)

type CommandContext struct {
	Out,
	Err io.Writer
	parent    *CommandContext
	cmd       *Command
	flags     *flag.FlagSet
	name      string
	Context   context.Context
	InputArgs []string
	Args      Args
	NamedArgs any

	// plan is -plan of a command with a Lifecycle: its order, nothing run
	plan bool
	// lifecycle is the graph of this run (Lifecycle)
	lifecycle *lifecycle.Graph
}

// Lifecycle is the graph of this run of the command, made on the first
// call: the command's Lifecycle fills it, its Run started around by it.
// A sub-command has a graph of its own.
func (ctx *CommandContext) Lifecycle() *lifecycle.Graph {
	if ctx.lifecycle == nil {
		ctx.lifecycle = &lifecycle.Graph{}
	}
	return ctx.lifecycle
}

func (ctx *CommandContext) Name() string {
	return ctx.name
}

func (ctx *CommandContext) Path() []string {
	names := []string{ctx.name}
	p := ctx.parent
	for p != nil {
		names = append(names, p.Name())
		p = p.parent
	}

	slices.Reverse(names)
	return names
}

func (ctx *CommandContext) FullName() string {
	pth := ctx.Path()
	for i, s := range pth {
		pth[i] = strconv.Quote(s)
	}
	return strings.Join(pth, " 🠆 ")
}

func (ctx *CommandContext) Parent() *CommandContext {
	return ctx.parent
}

func (ctx *CommandContext) Cmd() *Command {
	return ctx.cmd
}

func (ctx *CommandContext) Flags() *flag.FlagSet {
	return ctx.flags
}

func (ctx *CommandContext) WithValue(name, value any) *CommandContext {
	ctx.Context = context.WithValue(ctx.Context, name, value)
	return ctx
}

func (ctx *CommandContext) Value(name any) any {
	return ctx.Context.Value(name)
}

func (ctx *CommandContext) Fork() *CommandContext {
	child := *ctx
	child.parent = ctx
	child.lifecycle = nil
	child.plan = false
	return &child
}

// Run runs the command: with a Lifecycle, its nodes started around its Run
// (runLifecycle).
func (ctx *CommandContext) Run() error {
	if ctx.cmd.Lifecycle != nil {
		if err := ctx.cmd.Lifecycle(ctx); err != nil {
			return err
		}
		return ctx.runLifecycle(ctx.Lifecycle())
	}
	if ctx.cmd.Run == nil {
		return nil
	}

	return ctx.cmd.Run(ctx)
}

// runLifecycle starts the nodes of g in order, runs the command's Run (nil:
// none), waits for the waiters and closes the nodes in the inverse order
// (lifecycle.Graph.Run); SIGINT and SIGTERM end it. With -plan, the order of
// g, nothing run.
func (ctx *CommandContext) runLifecycle(g *lifecycle.Graph) error {
	if ctx.plan {
		p, err := g.Plan()
		if err != nil {
			return err
		}
		_, err = io.WriteString(ctx.Out, p)
		return err
	}
	c, stop := signal.NotifyContext(ctx.Context, os.Interrupt, syscall.SIGTERM)
	defer stop()
	var run func(context.Context) error
	if ctx.cmd.Run != nil {
		run = func(c context.Context) error {
			rctx := *ctx
			rctx.Context = c
			return ctx.cmd.Run(&rctx)
		}
	}
	return g.Run(c, run)
}

func (ctx *CommandContext) Help() (err error) {
	return Help(ctx).Execute()
}

type Command struct {
	Name        string
	Usage       string
	Description string
	sub         map[string]*Command
	New         func(ctx *CommandContext) (err error)
	Help        func(ctx *CommandContext) (err error)
	Run         func(ctx *CommandContext) (err error)
	ParseArgs   func(ctx *CommandContext) (err error)
	// Lifecycle fills the graph of the run (ctx.Lifecycle()) with the nodes
	// the command starts around its Run and closes after it (x/lifecycle):
	// called only when the command runs — the application wired by then.
	// The command gets -plan (unless it has its own): the order, nothing
	// run.
	Lifecycle func(ctx *CommandContext) error
}

func (b *Command) Sub(cmd *Command) *Command {
	if b.sub == nil {
		b.sub = make(map[string]*Command)
	}
	b.sub[cmd.Name] = cmd
	return b
}

func (b *Command) Parse(ctx *CommandContext) (_ *CommandContext, err error) {
	if ctx == nil {
		ctx = &CommandContext{
			name:      filepath.Base(os.Args[0]),
			InputArgs: os.Args[1:],
		}
	}

	if ctx.Context == nil {
		ctx.Context = context.Background()
	}

	if ctx.Out == nil {
		ctx.Out = os.Stdout
	}

	if ctx.Err == nil {
		ctx.Err = os.Stderr
	}

	dot := b

parse:
	ctx.cmd = dot
	ctx.flags = flag.NewFlagSet(ctx.name, flag.ContinueOnError)

	var defaultFlags struct {
		Stdout string
		Stderr string
	}

	ctx.flags.StringVar(&defaultFlags.Stdout, "rvq:cli:out", "-", "program stdout file name")
	ctx.flags.StringVar(&defaultFlags.Stderr, "rvq:cli:err", "-", "program stderr file name")

	if dot.New != nil {
		if err = dot.New(ctx); err != nil {
			err = ToContextParserError(ctx, "new", err)
			return
		}
	}

	if dot.Lifecycle != nil && ctx.flags.Lookup("plan") == nil {
		ctx.flags.BoolVar(&ctx.plan, "plan", false, "print the order the command starts and closes, run nothing")
	}

	if err = ctx.flags.Parse(ctx.InputArgs); err != nil {
		if err == flag.ErrHelp {
			helpCmd := *dot
			helpCmd.Run = func(ctx *CommandContext) (err error) {
				return Help(ctx).WithSubCommands().Execute()
			}
			ctx.cmd = &helpCmd
			return ctx, nil
		}
		err = ToContextParserError(ctx, "parse", err)
		return
	}

	ctx.Args = ctx.flags.Args()

	if dot.ParseArgs != nil {
		if err = dot.ParseArgs(ctx); err != nil {
			err = ToContextParserError(ctx, "parseArgs", err)
			return
		}
	}

	// the writers the flags name; without them, the ones given (the
	// parent's, os.Stdout and os.Stderr by default)
	if ctx.Out, err = openWriterOrDefault(defaultFlags.Stdout, ctx.Out); err != nil {
		return
	}

	if ctx.Err, err = openWriterOrDefault(defaultFlags.Stderr, ctx.Err); err != nil {
		return
	}

	if len(ctx.Args) > 0 {
		subName := ctx.Args[0]
		switch subName {
		case "--":
			ctx.Args = ctx.Args[1:]
		case "help", "--help":
			helpCmd := *dot
			helpCmd.Run = func(ctx *CommandContext) (err error) {
				return Help(ctx).WithSubCommands().Execute()
			}
			ctx.cmd = &helpCmd
			return ctx, nil
		default:
			if len(dot.sub) > 0 {
				sub := dot.sub[subName]
				if sub == nil {
					err = ToContextParserError(ctx, "sub", fmt.Errorf("unknown command: %s", subName))
					return
				}
				ctx.Args = ctx.Args[1:]
				ctx = ctx.Fork()
				ctx.name = subName
				ctx.InputArgs = ctx.Args
				ctx.Args = nil
				dot = sub
				goto parse
			}
		}
	}

	return ctx, nil
}

func openWriter(name string) (f *os.File, err error) {
	if _, err = os.Stat(name); os.IsNotExist(err) {
		if f, err = os.Create(name); err != nil {
			return
		}
	} else if err == nil {
		return os.OpenFile(name, os.O_WRONLY, 0666)
	}
	return
}

func openWriterOrDefault(name string, defaultw io.Writer) (io.Writer, error) {
	switch name {
	case "", "-":
		return defaultw, nil
	default:
		return openWriter(name)
	}
}
