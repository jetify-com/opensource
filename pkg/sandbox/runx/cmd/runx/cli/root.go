package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"go.jetpack.io/pkg/sandbox/runx"
)

func Help() {
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("runx")
	fmt.Println()
	fmt.Println(
		"Usage: runx [+<org>/<repo>]... [<cmd>] [<args>]...",
		"Usage: runx --install [<org>/<repo>]...",
	)
}

func Execute(ctx context.Context, args []string) int {
	if len(args) == 0 {
		Help()
		return 0
	}

	install := flag.Bool("install", false, "install packages only")
	flag.Parse()

	if *install {
		paths, err := runx.Install(args[1:]...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
			return 1
		}
		fmt.Println("Installed paths:")
		for _, path := range paths {
			fmt.Printf("  %s\n", path)
		}
	} else {
		if err := runx.Run(args...); err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
			return 1
		}
	}

	return 0
}

func Main() {
	code := Execute(context.Background(), os.Args[1:])
	os.Exit(code)
}
