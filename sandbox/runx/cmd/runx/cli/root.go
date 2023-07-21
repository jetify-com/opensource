package cli

import (
	"context"
	"fmt"
	"os"

	// "github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"go.jetpack.io/runx"
)

func Help() {
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("runx")
	fmt.Println()
	fmt.Println("Usage: runx [+<org>/<repo>]... [<cmd>] [<args>]...")
}

func Execute(ctx context.Context, args []string) int {
	if len(args) == 0 {
		Help()
		return 0
	}

	err := runx.Run(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		return 1
	}
	return 0
}

func Main() {
	code := Execute(context.Background(), os.Args[1:])
	os.Exit(code)
}
