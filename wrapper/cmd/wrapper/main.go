package main

import (
	"fmt"
	"log"
	"os"

	"github.com/k0kubun/pp"
	"go.jetpack.io/wrapper"
)

func main() {
	Main(os.Args[1:])
}

func Main(args []string) {
	if len(args) <= 0 {
		fmt.Println("Usage: wrapper <file>")
		log.Fatal("No arguments provided")
	}

	if args[0] == "debug" {
		Debug(args)
		return
	}

	Run(args)
}

func Debug(args []string) {
	if len(args) != 2 {
		fmt.Println("Usage: wrapper debug <file>")
		log.Fatal("Expected 2 arguments, got ", len(args))
	}

	w, err := wrapper.FromPath(args[1])
	if err != nil {
		log.Fatal(err)
	}
	exe := wrapper.ToExecutable(w.Config)
	pp.Println(exe)
}

func Run(args []string) {
	w, err := wrapper.FromPath(args[0])
	if err != nil {
		log.Fatal(err)
	}
	err = w.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
