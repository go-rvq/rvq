package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-rvq/rvq/cli"
	"github.com/mattn/go-isatty"
	"github.com/mattn/go-tty"
)

// MinPasswordLen is the least length of a password set by passwd.
const MinPasswordLen = 6

// stdinIsTerminal says whether the standard input is a terminal: then the
// password is read from it, not echoed; else (a pipe, a file) a line each.
var stdinIsTerminal = func() bool {
	fd := os.Stdin.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

// readNewPassword reads a new password and its confirmation by read, each
// asked on prompt (none when nil): at least MinPasswordLen chars, the two
// equal.
func readNewPassword(read func() (string, error), prompt io.Writer, account string) (pwd string, err error) {
	ask := func(s string) {
		if prompt != nil {
			fmt.Fprint(prompt, s)
		}
	}
	ask(fmt.Sprintf("Enter a new %q password: ", account))
	if pwd, err = read(); err != nil {
		return
	}
	if len(pwd) < MinPasswordLen {
		return "", fmt.Errorf("expected at least %d chars, got %d", MinPasswordLen, len(pwd))
	}
	ask("Confirm password: ")
	var pwd2 string
	if pwd2, err = read(); err != nil {
		return
	}
	if pwd2 != pwd {
		//lint:ignore ST1005 shown to the user as a sentence
		return "", fmt.Errorf("Passwords not equal")
	}
	return
}

// errNoPassword: the input ended before a password.
var errNoPassword = errors.New("no password in the input: the password and its confirmation, a line each")

// linesReader reads the lines of r, one a call (the end of line out).
func linesReader(r io.Reader) func() (string, error) {
	sc := bufio.NewScanner(r)
	return func() (string, error) {
		if !sc.Scan() {
			if err := sc.Err(); err != nil {
				return "", err
			}
			return "", errNoPassword
		}
		return strings.TrimRight(sc.Text(), "\r"), nil
	}
}

type PasswordResetContext[T any] struct {
	findUser    func(account string) (user T, err error)
	setPassword func(model T, password string) error
}

func NewPasswordResetContext[T any](
	findUser func(account string) (user T, err error),
	setPassword func(model T, password string) error,
) *PasswordResetContext[T] {
	return &PasswordResetContext[T]{findUser: findUser, setPassword: setPassword}

}

func PasswordReset[T any](context func(ctx *cli.CommandContext) (*PasswordResetContext[T], error)) *cli.Command {
	return &cli.Command{
		Name:        "passwd",
		Description: "changes user password",
		Usage:       "USER_ACCOUNT",
		ParseArgs: func(ctx *cli.CommandContext) (err error) {
			if err = ctx.Args.Eq(1); err != nil {
				return err
			}
			return nil
		},
		Run: func(ctx *cli.CommandContext) (err error) {
			var (
				account = ctx.Args[0]
				prc     *PasswordResetContext[T]
			)

			if prc, err = context(ctx); err != nil {
				return err
			}

			var user T
			user, err = prc.findUser(ctx.Args[0])
			if err != nil {
				return err
			}

			var pwd string
			if stdinIsTerminal() {
				// the terminal: not echoed
				var TTY *tty.TTY
				if TTY, err = tty.Open(); err != nil {
					return
				}
				defer TTY.Close()
				pwd, err = readNewPassword(TTY.ReadPassword, ctx.Err, account)
			} else {
				// a script: the password and its confirmation, a line each
				pwd, err = readNewPassword(linesReader(os.Stdin), nil, account)
			}
			if err != nil {
				return
			}

			if err = prc.setPassword(user, pwd); err != nil {
				return
			}

			fmt.Fprintf(ctx.Out, "password for user %q changed.\n", account)
			return nil
		},
	}
}
