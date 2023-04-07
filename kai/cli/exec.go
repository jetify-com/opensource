package cli

import (
	"fmt"
	"os"
	"os/exec"

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
	fmt.Println("I'm Kai, the AI assistant for your terminal\n")
	fmt.Println("I can translate English into shell commands\n")

	input := ""
	q1 := &survey.Input{
		Message: "What command would you like to run?",
		Suggest: func(toComplete string) []string {
			return []string{
				"list all files in the current directory",
				"show me the contents of the file named 'input.txt'",
				"create a new directory named 'dir'",
				"delete the file named 'output.txt'",
				"make all files in current directory read only",
				"start nginx using docker, forward 443 and 80 port",
			}
		},
	}
	err := survey.AskOne(q1, &input, survey.WithValidator(survey.Required))

	if err != nil {
		return err
	}

	results, err := kai.Exec(input)
	if err != nil {
		return err
	}

	selection := ""
	q2 := &survey.Select{
		Message: "Choose a command:",
		Options: results,
	}
	survey.AskOne(q2, &selection)

	shouldExecute := false
	q3 := &survey.Confirm{
		Message: "Execute this command?",
		Default: true,
	}
	survey.AskOne(q3, &shouldExecute)

	if shouldExecute {
		fmt.Printf("> %s\n", selection)
		shellexec(selection)
	}

	return nil
}

func shellexec(script string) {
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		// Check exit code:
		if exitError, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// Here you can get the exit code
			os.Exit(exitError.ExitCode())
		}
		panic(err)
	}
}
