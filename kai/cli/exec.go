package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"go.jetpack.io/kai"
)

func ExecCmd() *cobra.Command {
	command := &cobra.Command{
		Use:           "exec",
		Short:         "Use AI to execute a shell command written in English",
		Args:          cobra.NoArgs,
		RunE:          execCmd,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	return command
}

func execCmd(cmd *cobra.Command, args []string) error {
	input := ""
	prompt := &survey.Input{
		Message: "What would you like to do?",
		Help:    "Type a query to run",
	}
	err := survey.AskOne(prompt, &input)

	if err != nil {
		return err
	}

	results, err := kai.Exec(input)
	if err != nil {
		return err
	}

	for _, result := range results {
		fmt.Println(result)
	}

	return nil
}
