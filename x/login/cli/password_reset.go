package cli

import (
	"fmt"
	"os"

	"github.com/go-rvq/rvq/cli"
	"github.com/mattn/go-tty"
)

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

			fmt.Fprintf(os.Stdin, "Enter a new %q password: ", account)

			var TTY *tty.TTY
			if TTY, err = tty.Open(); err != nil {
				return
			}

			defer TTY.Close()

			var pwd string
			if pwd, err = TTY.ReadPassword(); err != nil {
				return
			}

			if len(pwd) < 6 {
				err = fmt.Errorf("expected at least %d chars, got %d", 6, len(pwd))
				return
			}

			fmt.Fprint(os.Stdin, "Confirm password: ")

			var pwd2 string
			if pwd2, err = TTY.ReadPassword(); err != nil {
				return
			}

			if pwd2 != pwd {
				err = fmt.Errorf("Passwords not equal")
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
