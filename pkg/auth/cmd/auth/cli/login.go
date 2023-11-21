package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"go.jetpack.io/pkg/auth"
	"go.jetpack.io/pkg/auth/session"
)

type loginFlags struct {
	client string
}

func LoginCmd() *cobra.Command {
	flags := &loginFlags{}

	command := &cobra.Command{
		Use:   "login <issuer> --client <client-id>",
		Short: "Login using OIDC",
		Args:  cobra.ExactArgs(1),
		RunE:  loginCmd(flags),
	}

	command.Flags().StringVar(&flags.client, "client", "", "client id")

	return command
}

func loginCmd(flags *loginFlags) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		issuer := args[0]
		clientID := flags.client

		if clientID == "" {
			return fmt.Errorf("please provide a client id")
		}
		return login(issuer, clientID)
	}
}

func login(issuer, clientID string) error {
	client, err := auth.NewClient(issuer, clientID)
	if err != nil {
		return err
	}

	tok, err := client.LoginFlow()
	if err != nil {
		return err
	}

	err = printToken(tok)
	if err != nil {
		return err
	}
	return nil
}

func printToken(tok *session.Token) error {
	fmt.Println("Tokens:")
	err := printJSON(tok)
	if err != nil {
		return err
	}

	fmt.Println("\nID Token Claims:")
	err = printJSON(tok.IDClaims())
	if err != nil {
		return err
	}

	return nil
}

func printJSON(v any) error {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	if !isTerminal() {
		color.NoColor = true
	}
	output, err := prettyjson.Format(bytes)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}

func isTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}
