package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func WhoAmICmd() *cobra.Command {
	flags := &sharedFlags{}

	command := &cobra.Command{
		Use:  "whoami",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := buildClient(flags.issuer, flags.client)
			if err != nil {
				return err
			}
			sessions, err := client.GetSessions()
			if err != nil {
				return err
			}

			fmt.Printf("Found %d sessions\n", len(sessions))

			for _, session := range sessions {
				tok := session.Peek()
				fmt.Printf("* %s\n", tok.IDClaims().Subject)
			}

			return nil
		},
	}

	command.Flags().StringVar(&flags.client, "client", "", "client id")
	command.Flags().StringVar(&flags.issuer, "issuer", "", "issuer")

	return command
}
