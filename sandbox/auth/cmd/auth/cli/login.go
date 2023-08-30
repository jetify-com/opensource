package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"go.jetpack.io/auth"
	"golang.org/x/oauth2"
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
	tok, err := auth.Login(issuer, clientID)
	if err != nil {
		return err
	}

	err = printToken(tok)
	if err != nil {
		return err
	}
	return nil
}

func printToken(tok *oauth2.Token) error {
	data, err := json.MarshalIndent(tok, "", "  ")
	if err != nil {
		return err
	}
	err = printJSON(data)
	if err != nil {
		return err
	}

	return nil
}

func printJSON(bytes []byte) error {
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
